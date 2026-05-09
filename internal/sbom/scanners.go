package sbom

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"github.com/baditaflorin/audit-in-a-box/internal/tools"
)

type Scanner struct {
	Runner tools.Runner
}

func (s Scanner) RunSyft(ctx context.Context, workDir string) ([]models.Dependency, []string) {
	stdout, _, err := s.Runner.Command(ctx, "syft", "dir:"+workDir, "-o", "json", "--quiet")
	if err != nil {
		return nil, []string{err.Error()}
	}

	var payload struct {
		Artifacts []struct {
			Name     string          `json:"name"`
			Version  string          `json:"version"`
			Type     string          `json:"type"`
			PURL     string          `json:"purl"`
			Licenses json.RawMessage `json:"licenses"`
		} `json:"artifacts"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		return nil, []string{"parse syft output: " + err.Error()}
	}

	deps := make([]models.Dependency, 0, len(payload.Artifacts))
	for _, artifact := range payload.Artifacts {
		ecosystem := ecosystemFromType(artifact.Type)
		if artifact.Name == "" || ecosystem == "unknown" {
			continue
		}
		deps = append(deps, models.Dependency{
			Name:       artifact.Name,
			Version:    artifact.Version,
			Ecosystem:  ecosystem,
			Scope:      "detected",
			Licenses:   parseSyftLicenses(artifact.Licenses),
			PackageURL: artifact.PURL,
			Source:     "syft",
			Confidence: 0.96,
			Reasons:    []string{"detected by Syft SBOM scan"},
		})
	}

	return deps, nil
}

func (s Scanner) RunGrype(ctx context.Context, workDir string) ([]models.Vulnerability, []string) {
	stdout, _, err := s.Runner.Command(ctx, "grype", "dir:"+workDir, "-o", "json", "--quiet")
	if err != nil {
		return nil, []string{err.Error()}
	}

	var payload struct {
		Matches []struct {
			Vulnerability struct {
				ID          string `json:"id"`
				Severity    string `json:"severity"`
				Description string `json:"description"`
				DataSource  string `json:"dataSource"`
				Fix         struct {
					Versions []string `json:"versions"`
				} `json:"fix"`
				URLs []string `json:"urls"`
			} `json:"vulnerability"`
			Artifact struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"artifact"`
		} `json:"matches"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		return nil, []string{"parse grype output: " + err.Error()}
	}

	vulns := make([]models.Vulnerability, 0, len(payload.Matches))
	for _, match := range payload.Matches {
		vuln := models.Vulnerability{
			ID:               match.Vulnerability.ID,
			PackageName:      match.Artifact.Name,
			InstalledVersion: match.Artifact.Version,
			FixedVersion:     strings.Join(match.Vulnerability.Fix.Versions, ", "),
			Severity:         strings.ToUpper(match.Vulnerability.Severity),
			Description:      match.Vulnerability.Description,
			Source:           "grype",
			Confidence:       0.95,
		}
		if len(match.Vulnerability.URLs) > 0 {
			vuln.PrimaryURL = match.Vulnerability.URLs[0]
		} else {
			vuln.PrimaryURL = match.Vulnerability.DataSource
		}
		vulns = append(vulns, vuln)
	}

	return vulns, nil
}

func (s Scanner) RunTrivy(ctx context.Context, workDir string) ([]models.Vulnerability, []models.LicenseRisk, []string) {
	stdout, _, err := s.Runner.Command(ctx, "trivy", "fs", "--quiet", "--format", "json", "--scanners", "vuln,license", workDir)
	if err != nil {
		return nil, nil, []string{err.Error()}
	}

	var payload struct {
		Results []struct {
			Vulnerabilities []struct {
				VulnerabilityID  string   `json:"VulnerabilityID"`
				PkgName          string   `json:"PkgName"`
				InstalledVersion string   `json:"InstalledVersion"`
				FixedVersion     string   `json:"FixedVersion"`
				Severity         string   `json:"Severity"`
				Description      string   `json:"Description"`
				PrimaryURL       string   `json:"PrimaryURL"`
				References       []string `json:"References"`
			} `json:"Vulnerabilities"`
			LicenseFindings []struct {
				PkgName    string `json:"PkgName"`
				Name       string `json:"Name"`
				Severity   string `json:"Severity"`
				Category   string `json:"Category"`
				Confidence any    `json:"Confidence"`
			} `json:"Licenses"`
		} `json:"Results"`
	}
	if err := json.Unmarshal([]byte(stdout), &payload); err != nil {
		return nil, nil, []string{"parse trivy output: " + err.Error()}
	}

	var vulns []models.Vulnerability
	var licenses []models.LicenseRisk
	for _, result := range payload.Results {
		for _, item := range result.Vulnerabilities {
			primaryURL := item.PrimaryURL
			if primaryURL == "" && len(item.References) > 0 {
				primaryURL = item.References[0]
			}
			vulns = append(vulns, models.Vulnerability{
				ID:               item.VulnerabilityID,
				PackageName:      item.PkgName,
				InstalledVersion: item.InstalledVersion,
				FixedVersion:     item.FixedVersion,
				Severity:         strings.ToUpper(item.Severity),
				Description:      item.Description,
				PrimaryURL:       primaryURL,
				Source:           "trivy",
				Confidence:       0.94,
			})
		}
		for _, item := range result.LicenseFindings {
			if item.Severity == "" || strings.EqualFold(item.Severity, "UNKNOWN") {
				continue
			}
			licenses = append(licenses, models.LicenseRisk{
				PackageName: item.PkgName,
				License:     item.Name,
				Severity:    strings.ToUpper(item.Severity),
				Reason:      strings.TrimSpace(item.Category + " " + confidenceText(item.Confidence)),
				Confidence:  0.9,
			})
		}
	}

	return vulns, licenses, nil
}

func confidenceText(value any) string {
	switch item := value.(type) {
	case string:
		return item
	case float64:
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", item), "0"), ".")
	default:
		raw, err := json.Marshal(item)
		if err != nil || string(raw) == "null" {
			return ""
		}
		return string(raw)
	}
}

func parseSyftLicenses(raw json.RawMessage) []string {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var stringsOnly []string
	if err := json.Unmarshal(raw, &stringsOnly); err == nil {
		return stringsOnly
	}

	var objects []struct {
		Value          string `json:"value"`
		SPDXExpression string `json:"spdxExpression"`
	}
	if err := json.Unmarshal(raw, &objects); err != nil {
		return nil
	}

	licenses := make([]string, 0, len(objects))
	for _, item := range objects {
		switch {
		case item.SPDXExpression != "":
			licenses = append(licenses, item.SPDXExpression)
		case item.Value != "":
			licenses = append(licenses, item.Value)
		}
	}

	return licenses
}

func ecosystemFromType(value string) string {
	switch strings.ToLower(value) {
	case "npm":
		return "npm"
	case "go-module":
		return "go"
	case "python", "python-package":
		return "python"
	default:
		return "unknown"
	}
}
