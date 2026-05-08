package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/baditaflorin/audit-in-a-box/internal/analysis"
	"github.com/baditaflorin/audit-in-a-box/internal/api"
	"github.com/baditaflorin/audit-in-a-box/internal/config"
	"github.com/baditaflorin/audit-in-a-box/internal/utils"
)

func main() {
	cfg, err := config.Load()
	if utils.HandleErrorOrLogWithMessages(err, "load config", "") {
		os.Exit(1)
	}

	level := slog.LevelInfo
	if strings.EqualFold(cfg.LogLevel, "debug") {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	service := analysis.NewService(cfg)
	handler := api.NewRouter(cfg, service)
	server := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       2 * time.Minute,
		WriteTimeout:      3 * time.Minute,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("server_starting", "addr", cfg.ServerAddr)
		errCh <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		slog.Info("shutdown_signal", "signal", sig.String())
	case err := <-errCh:
		if err != nil && err != http.ErrServerClosed {
			slog.Error("server_failed", "error", err)
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		slog.Error("shutdown_failed", "error", err)
		os.Exit(1)
	}
	slog.Info("server_stopped")
}
