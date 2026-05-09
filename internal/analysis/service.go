package analysis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
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
		return models.AuditReport{}, models.NewDomainError(
			"manifest_required",
			"Upload, paste, or scrape a dependency manifest first.",
			"The audit needs package data before it can build an SBOM or risk report.",
			"Provide package.json, package-lock.json, pnpm-lock.yaml, go.mod, pyproject.toml, or requirements.txt.",
			true,
		)
	}
	if fileName == "" {
		fileName = sbom.DefaultFileName(sbom.DetectEcosystem("", content))
	}

	deps, parseWarnings, inference, err := sbom.AnalyzeManifest(fileName, content)
	if err != nil {
		return models.AuditReport{}, err
	}
	warnings = append(warnings, parseWarnings...)
	content = inference.NormalizedContent
	fileName = firstNonEmpty(fileName, inference.FileName)

	if inference.Ecosystem == "unknown" && len(deps) == 0 {
		return models.AuditReport{}, models.NewDomainError(
			"unsupported_manifest",
			"This input does not look like a supported dependency manifest.",
			strings.Join(inference.Reasons, " "),
			"Upload package.json, package-lock.json, pnpm-lock.yaml, go.mod, pyproject.toml, or requirements.txt.",
			true,
		)
	}

	ecosystem := inference.Ecosystem
	if req.Ecosystem != "" && req.Ecosystem != "unknown" {
		ecosystem = req.Ecosystem
	}

	id := stableAuditID(fileName, content)
	workDir := filepath.Join(s.Config.WorkDir, id)
	if err := os.MkdirAll(workDir, 0o700); err != nil {
		return models.AuditReport{}, fmt.Errorf("create audit workdir: %w", err)
	}
	manifestPath := filepath.Join(workDir, filepath.Base(fileName))
	if err := os.WriteFile(manifestPath, []byte(content), 0o600); err != nil {
		return models.AuditReport{}, fmt.Errorf("write manifest: %w", err)
	}

	toolStatus := tools.CheckTools(ctx, s.Runner)
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
	anomalies := LicenseEvidenceAnomalies(deps)
	anomalies = append(anomalies, scaleAnomalies(len(deps))...)
	health := s.Maintainer.Assess(ctx, deps)
	rollup := duckdb.BuildRollup(ctx, s.Runner, workDir, deps, vulns, licenseRisks)
	if rollup.UsedDuckDB {
		markToolUsed(toolStatus, "duckdb")
	}

	risk := ScoreReport(deps, vulns, licenseRisks, health)
	input := models.AuditInput{
		FileName:        fileName,
		Ecosystem:       ecosystem,
		Bytes:           len(content),
		NormalizedBytes: len(content),
		Kind:            inference.Kind,
		Parser:          inference.Parser,
		Confidence:      inference.Confidence,
		Reasons:         inference.Reasons,
		SourceHash:      inference.SourceHash,
	}
	buildVersion := models.BuildVersion{
		Version: version.Version,
		Commit:  version.Commit,
		Date:    version.Date,
	}
	report := models.AuditReport{
		ID:          id,
		GeneratedAt: time.Now().UTC(),
		Input:       input,
		Provenance: models.ReportProvenance{
			SchemaVersion: inference.SchemaVersion,
			AppVersion:    buildVersion,
			SourceHash:    inference.SourceHash,
			Input:         input,
			Parameters: map[string]any{
				"max_maintainer_packages": s.Config.MaxMaintainerPackages,
				"tool_timeout_seconds":    int(s.Config.ToolTimeout.Seconds()),
				"scanners":                []string{"syft", "grype", "trivy", "duckdb"},
			},
		},
		ToolStatus:        toolStatus,
		Dependencies:      deps,
		Vulnerabilities:   vulns,
		LicenseRisks:      licenseRisks,
		Anomalies:         anomalies,
		MaintainerHealth:  health,
		DuckDBRollup:      rollup,
		Risk:              risk,
		Warnings:          compactWarnings(warnings),
		ScrapedCandidates: scraped,
		Version:           buildVersion,
		ElapsedMillis:     time.Since(start).Milliseconds(),
		RawEvidence: map[string]models.EvidenceInfo{
			"workdir":     {Available: true, Path: workDir, Message: "runtime artifacts are stored on the backend host"},
			"source_hash": {Available: true, Message: inference.SourceHash},
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

func stableAuditID(fileName string, content string) string {
	sum := sha256.Sum256([]byte(fileName + "\x00" + content))
	return "audit-" + hex.EncodeToString(sum[:])[:16]
}

func scaleAnomalies(dependencyCount int) []models.ReportAnomaly {
	if dependencyCount < 500 {
		return nil
	}
	return []models.ReportAnomaly{{
		Code:       "large_dependency_surface",
		Severity:   "medium",
		Message:    "This manifest has a very large dependency surface.",
		Why:        "Large manifests are still supported, but scanner and maintainer evidence may arrive more slowly or be partial.",
		NextStep:   "Review the highest-severity vulnerabilities and confirmed license risks first; use the dependency list as SBOM inventory.",
		Confidence: 0.98,
	}}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
