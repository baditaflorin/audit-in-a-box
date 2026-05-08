package analysis

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/baditaflorin/audit-in-a-box/internal/config"
	"github.com/baditaflorin/audit-in-a-box/internal/duckdb"
	"github.com/baditaflorin/audit-in-a-box/internal/llm"
	"github.com/baditaflorin/audit-in-a-box/internal/maintainer"
	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"github.com/baditaflorin/audit-in-a-box/internal/sbom"
	"github.com/baditaflorin/audit-in-a-box/internal/tools"
	"github.com/baditaflorin/audit-in-a-box/pkg/version"
	"github.com/google/uuid"
)

type Service struct {
	Config     config.Config
	Runner     tools.Runner
	Scanner    sbom.Scanner
	Maintainer maintainer.Client
	Summarizer llm.Summarizer
}

func NewService(cfg config.Config) Service {
	runner := tools.Runner{Timeout: cfg.ToolTimeout}
	return Service{
		Config:     cfg,
		Runner:     runner,
		Scanner:    sbom.Scanner{Runner: runner},
		Maintainer: maintainer.Client{Limit: cfg.MaxMaintainerPackages},
		Summarizer: llm.Summarizer{Config: cfg},
	}
}

func (s Service) Audit(ctx context.Context, req models.AuditRequest) (models.AuditReport, error) {
	start := time.Now()
	id := uuid.NewString()
	warnings := []string{}

	content := strings.TrimSpace(req.Content)
	fileName := strings.TrimSpace(req.FileName)
	var scraped []models.ScrapedManifest

	if content == "" && strings.TrimSpace(req.PastedHTML) != "" {
		candidates, err := ExtractManifestsFromHTML(req.PastedHTML)
		if err != nil {
			return models.AuditReport{}, fmt.Errorf("scrape pasted html: %w", err)
		}
		scraped = candidates
		if len(candidates) == 0 {
			return models.AuditReport{}, fmt.Errorf("no manifest-like content found in pasted html")
		}
		content = candidates[0].Content
		fileName = candidates[0].FileName
	}

	if content == "" {
		return models.AuditReport{}, fmt.Errorf("content or pasted_html is required")
	}
	if fileName == "" {
		fileName = sbom.DefaultFileName(sbom.DetectEcosystem("", content))
	}

	ecosystem := req.Ecosystem
	if ecosystem == "" || ecosystem == "unknown" {
		ecosystem = sbom.DetectEcosystem(fileName, content)
	}

	workDir := filepath.Join(s.Config.WorkDir, id)
	if err := os.MkdirAll(workDir, 0o700); err != nil {
		return models.AuditReport{}, fmt.Errorf("create audit workdir: %w", err)
	}
	manifestPath := filepath.Join(workDir, filepath.Base(fileName))
	if err := os.WriteFile(manifestPath, []byte(content), 0o600); err != nil {
		return models.AuditReport{}, fmt.Errorf("write manifest: %w", err)
	}

	toolStatus := tools.CheckTools(ctx, s.Runner)
	deps, parseWarnings, err := sbom.ParseManifest(fileName, content)
	if err != nil {
		return models.AuditReport{}, err
	}
	warnings = append(warnings, parseWarnings...)

	if syftDeps, syftWarnings := s.Scanner.RunSyft(ctx, workDir); len(syftWarnings) > 0 {
		warnings = append(warnings, syftWarnings...)
	} else {
		deps = append(deps, syftDeps...)
		markToolUsed(toolStatus, "syft")
	}
	deps = DeduplicateDependencies(deps)

	var vulns []models.Vulnerability
	if grypeVulns, grypeWarnings := s.Scanner.RunGrype(ctx, workDir); len(grypeWarnings) > 0 {
		warnings = append(warnings, grypeWarnings...)
	} else {
		vulns = append(vulns, grypeVulns...)
		markToolUsed(toolStatus, "grype")
	}

	var scannerLicenseRisks []models.LicenseRisk
	if trivyVulns, trivyLicenses, trivyWarnings := s.Scanner.RunTrivy(ctx, workDir); len(trivyWarnings) > 0 {
		warnings = append(warnings, trivyWarnings...)
	} else {
		vulns = append(vulns, trivyVulns...)
		scannerLicenseRisks = append(scannerLicenseRisks, trivyLicenses...)
		markToolUsed(toolStatus, "trivy")
	}
	vulns = DeduplicateVulnerabilities(vulns)

	licenseRisks := AssessLicenseRisks(deps, scannerLicenseRisks)
	health := s.Maintainer.Assess(ctx, deps)
	rollup := duckdb.BuildRollup(ctx, s.Runner, workDir, deps, vulns, licenseRisks)
	if rollup.UsedDuckDB {
		markToolUsed(toolStatus, "duckdb")
	}

	risk := ScoreReport(deps, vulns, licenseRisks, health)
	report := models.AuditReport{
		ID:                id,
		GeneratedAt:       time.Now().UTC(),
		Input:             models.AuditInput{FileName: fileName, Ecosystem: ecosystem, Bytes: len(content)},
		ToolStatus:        toolStatus,
		Dependencies:      deps,
		Vulnerabilities:   vulns,
		LicenseRisks:      licenseRisks,
		MaintainerHealth:  health,
		DuckDBRollup:      rollup,
		Risk:              risk,
		Warnings:          compactWarnings(warnings),
		ScrapedCandidates: scraped,
		Version: models.BuildVersion{
			Version: version.Version,
			Commit:  version.Commit,
			Date:    version.Date,
		},
		ElapsedMillis: time.Since(start).Milliseconds(),
		RawEvidence: map[string]models.EvidenceInfo{
			"workdir": {Available: true, Path: workDir, Message: "runtime artifacts are stored on the backend host"},
		},
	}
	report.Summary = s.Summarizer.Summarize(ctx, report)

	return report, nil
}

func markToolUsed(status map[string]models.ToolStatus, name string) {
	item := status[name]
	item.Used = true
	status[name] = item
}

func compactWarnings(warnings []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(warnings))
	for _, warning := range warnings {
		warning = strings.TrimSpace(warning)
		if warning == "" || seen[warning] {
			continue
		}
		seen[warning] = true
		out = append(out, warning)
	}
	return out
}
