//go:build integration

package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/baditaflorin/audit-in-a-box/internal/analysis"
	"github.com/baditaflorin/audit-in-a-box/internal/api"
	"github.com/baditaflorin/audit-in-a-box/internal/config"
	"github.com/stretchr/testify/require"
)

func TestHealthz(t *testing.T) {
	cfg := config.Config{
		ServerAddr:            ":0",
		AllowedOrigins:        []string{"http://localhost:4173"},
		WorkDir:               t.TempDir(),
		ToolTimeout:           1,
		MaxUploadBytes:        1024,
		MaxMaintainerPackages: 1,
	}
	router := api.NewRouter(cfg, analysis.NewService(cfg))
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
}
