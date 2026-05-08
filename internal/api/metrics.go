package api

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "audit_in_a_box_http_requests_total", Help: "HTTP requests by route, method, and status."},
		[]string{"method", "route", "status"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{Name: "audit_in_a_box_http_request_duration_seconds", Help: "HTTP request duration.", Buckets: prometheus.DefBuckets},
		[]string{"method", "route"},
	)
	auditRequests = prometheus.NewCounter(
		prometheus.CounterOpts{Name: "audit_in_a_box_audit_requests_total", Help: "Audit requests handled."},
	)
	auditDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{Name: "audit_in_a_box_audit_duration_seconds", Help: "Audit processing duration.", Buckets: []float64{1, 2, 5, 10, 30, 60, 120}},
	)
	scannerFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{Name: "audit_in_a_box_scanner_failures_total", Help: "Scanner failure count."},
		[]string{"scanner"},
	)
	riskScores = prometheus.NewHistogram(
		prometheus.HistogramOpts{Name: "audit_in_a_box_risk_score", Help: "Returned report risk score.", Buckets: []float64{0, 15, 35, 60, 80, 100}},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpDuration, auditRequests, auditDuration, scannerFailures, riskScores)
}

type statusRecorder struct {
	status int
	size   int
}

func (r *statusRecorder) observe(method, route string, start time.Time) {
	status := r.status
	if status == 0 {
		status = 200
	}
	httpRequests.WithLabelValues(method, route, strconv.Itoa(status)).Inc()
	httpDuration.WithLabelValues(method, route).Observe(time.Since(start).Seconds())
}
