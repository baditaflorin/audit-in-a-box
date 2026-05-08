package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/baditaflorin/audit-in-a-box/internal/analysis"
	"github.com/baditaflorin/audit-in-a-box/internal/config"
	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"github.com/baditaflorin/audit-in-a-box/internal/tools"
	"github.com/baditaflorin/audit-in-a-box/pkg/version"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	Config   config.Config
	Service  analysis.Service
	Validate *validator.Validate
}

func NewRouter(cfg config.Config, service analysis.Service) http.Handler {
	server := Server{
		Config:   cfg,
		Service:  service,
		Validate: validator.New(),
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(server.logMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/healthz", server.health)
	r.Get("/readyz", server.ready)
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/tools", server.tools)
		api.Post("/scrape", server.scrape)
		api.Post("/audits", server.audit)
	})

	return r
}

func (s Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "version": version.Version, "commit": version.Commit})
}

func (s Server) ready(w http.ResponseWriter, r *http.Request) {
	status := tools.CheckTools(r.Context(), s.Service.Runner)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "tools": status})
}

func (s Server) tools(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, tools.CheckTools(r.Context(), s.Service.Runner))
}

func (s Server) scrape(w http.ResponseWriter, r *http.Request) {
	var request struct {
		HTML string `json:"html" validate:"required"`
	}
	if err := decodeJSON(r, s.Config.MaxUploadBytes, &request); err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err)
		return
	}
	if err := s.Validate.Struct(request); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", err)
		return
	}
	candidates, err := analysis.ExtractManifestsFromHTML(request.HTML)
	if err != nil {
		writeError(w, http.StatusBadRequest, "scrape_failed", err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"candidates": candidates})
}

func (s Server) audit(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	auditRequests.Inc()

	request, err := s.decodeAudit(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "bad_request", err)
		return
	}
	if err := s.Validate.Struct(request); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", err)
		return
	}

	report, err := s.Service.Audit(r.Context(), request)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "audit_failed", err)
		return
	}
	for _, warning := range report.Warnings {
		for _, scanner := range []string{"trivy", "syft", "grype", "duckdb"} {
			if strings.Contains(strings.ToLower(warning), scanner) {
				scannerFailures.WithLabelValues(scanner).Inc()
			}
		}
	}

	auditDuration.Observe(time.Since(start).Seconds())
	riskScores.Observe(float64(report.Risk.Score))
	writeJSON(w, http.StatusOK, report)
}

func (s Server) decodeAudit(r *http.Request) (models.AuditRequest, error) {
	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(s.Config.MaxUploadBytes); err != nil {
			return models.AuditRequest{}, err
		}
		var request models.AuditRequest
		request.FileName = r.FormValue("file_name")
		request.Ecosystem = r.FormValue("ecosystem")
		request.PastedHTML = r.FormValue("pasted_html")
		file, header, err := r.FormFile("file")
		if err == nil {
			defer file.Close()
			data, readErr := io.ReadAll(io.LimitReader(file, s.Config.MaxUploadBytes+1))
			if readErr != nil {
				return models.AuditRequest{}, readErr
			}
			if int64(len(data)) > s.Config.MaxUploadBytes {
				return models.AuditRequest{}, fmt.Errorf("file exceeds %d bytes", s.Config.MaxUploadBytes)
			}
			if request.FileName == "" {
				request.FileName = header.Filename
			}
			request.Content = string(data)
		}
		return request, nil
	}

	var request models.AuditRequest
	if err := decodeJSON(r, s.Config.MaxUploadBytes, &request); err != nil {
		return models.AuditRequest{}, err
	}
	return request, nil
}

func (s Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w, status: 200}
		next.ServeHTTP(recorder, r)
		route := routePattern(r)
		(&statusRecorder{status: recorder.status, size: recorder.size}).observe(r.Method, route, start)
		slog.Info("http_request",
			"method", r.Method,
			"path", r.URL.Path,
			"route", route,
			"status", recorder.status,
			"bytes", recorder.size,
			"duration_ms", time.Since(start).Milliseconds(),
			"trace_id", middleware.GetReqID(r.Context()),
		)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(data []byte) (int, error) {
	n, err := r.ResponseWriter.Write(data)
	r.size += n
	return n, err
}

func decodeJSON(r *http.Request, maxBytes int64, target any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, maxBytes+1))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, code string, err error) {
	if err == nil {
		err = errors.New(http.StatusText(status))
	}
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": err.Error(),
		},
	})
}

func routePattern(r *http.Request) string {
	if ctx := chi.RouteContext(r.Context()); ctx != nil {
		if pattern := ctx.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
}
