package duckdb

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"github.com/baditaflorin/audit-in-a-box/internal/tools"
)

func BuildRollup(ctx context.Context, runner tools.Runner, workDir string, deps []models.Dependency, vulns []models.Vulnerability, risks []models.LicenseRisk) models.DuckDBRollup {
	rollup := fallbackRollup(deps, vulns, risks)
	depsPath := filepath.Join(workDir, "dependencies.json")
	vulnsPath := filepath.Join(workDir, "vulnerabilities.json")

	if err := writeJSON(depsPath, deps); err != nil {
		rollup.Diagnostic = "write duckdb dependency artifact: " + err.Error()
		return rollup
	}
	if err := writeJSON(vulnsPath, vulns); err != nil {
		rollup.Diagnostic = "write duckdb vulnerability artifact: " + err.Error()
		return rollup
	}

	query := "SELECT (SELECT count(*) FROM read_json_auto('" + escapeSQL(depsPath) + "')) AS dependency_count, " +
		"(SELECT count(*) FROM read_json_auto('" + escapeSQL(vulnsPath) + "')) AS vulnerability_count;"
	stdout, _, err := runner.Command(ctx, "duckdb", "-json", "-c", query)
	if err != nil {
		rollup.Diagnostic = err.Error()
		return rollup
	}

	var rows []struct {
		DependencyCount    int `json:"dependency_count"`
		VulnerabilityCount int `json:"vulnerability_count"`
	}
	if err := json.Unmarshal([]byte(stdout), &rows); err != nil || len(rows) == 0 {
		if err != nil {
			rollup.Diagnostic = "parse duckdb rollup: " + err.Error()
		}
		return rollup
	}

	rollup.UsedDuckDB = true
	rollup.DependencyCount = rows[0].DependencyCount
	rollup.VulnerabilityCount = rows[0].VulnerabilityCount
	rollup.GeneratedArtifact = filepath.Join(workDir, "audit.duckdb")
	rollup.Diagnostic = "rollup calculated with DuckDB CLI"
	return rollup
}

func fallbackRollup(deps []models.Dependency, vulns []models.Vulnerability, risks []models.LicenseRisk) models.DuckDBRollup {
	severityCounts := map[string]int{}
	for _, vuln := range vulns {
		severityCounts[strings.ToUpper(vuln.Severity)]++
	}
	return models.DuckDBRollup{
		UsedDuckDB:         false,
		DependencyCount:    len(deps),
		VulnerabilityCount: len(vulns),
		LicenseRiskCount:   len(risks),
		SeverityCounts:     severityCounts,
		Diagnostic:         "DuckDB CLI unavailable; rollup calculated in Go",
	}
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func escapeSQL(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
