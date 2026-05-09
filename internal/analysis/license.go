package analysis

import (
	"sort"
	"strings"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
)

func AssessLicenseRisks(deps []models.Dependency, scannerRisks []models.LicenseRisk) []models.LicenseRisk {
	risks := append([]models.LicenseRisk{}, scannerRisks...)
	for _, dep := range deps {
		if dep.Name == "" {
			continue
		}
		for _, license := range dep.Licenses {
			normalized := strings.ToUpper(license)
			switch {
			case strings.Contains(normalized, "AGPL"):
				risks = append(risks, models.LicenseRisk{PackageName: dep.Name, License: license, Severity: "HIGH", Reason: "strong network copyleft license", Confidence: 0.92})
			case strings.Contains(normalized, "GPL") && !strings.Contains(normalized, "LGPL"):
				risks = append(risks, models.LicenseRisk{PackageName: dep.Name, License: license, Severity: "HIGH", Reason: "strong copyleft license", Confidence: 0.9})
			case strings.Contains(normalized, "LGPL"), strings.Contains(normalized, "MPL"), strings.Contains(normalized, "CDDL"):
				risks = append(risks, models.LicenseRisk{PackageName: dep.Name, License: license, Severity: "MEDIUM", Reason: "weak copyleft license", Confidence: 0.88})
			case strings.Contains(normalized, "BUSL"), strings.Contains(normalized, "SSPL"), strings.Contains(normalized, "PROPRIETARY"):
				risks = append(risks, models.LicenseRisk{PackageName: dep.Name, License: license, Severity: "HIGH", Reason: "source-available or proprietary-style license", Confidence: 0.9})
			}
		}
	}

	return dedupeLicenseRisks(risks)
}

func LicenseEvidenceAnomalies(deps []models.Dependency) []models.ReportAnomaly {
	missing := 0
	for _, dep := range deps {
		if dep.Name != "" && len(dep.Licenses) == 0 {
			missing++
		}
	}
	if missing == 0 {
		return nil
	}

	severity := "info"
	nextStep := "Install Syft/Trivy or audit with a lockfile/SBOM so license evidence can be confirmed."
	if missing > len(deps)/2 {
		severity = "medium"
		nextStep = "Treat the license section as incomplete until scanner evidence is available; do not read missing metadata as confirmed license risk."
	}
	return []models.ReportAnomaly{{
		Code:       "license_evidence_missing",
		Severity:   severity,
		Message:    "License evidence is incomplete for some dependencies.",
		Why:        "The manifest lists dependencies, but it does not carry authoritative package license metadata.",
		NextStep:   nextStep,
		Confidence: 0.95,
	}}
}

func dedupeLicenseRisks(items []models.LicenseRisk) []models.LicenseRisk {
	seen := map[string]bool{}
	out := make([]models.LicenseRisk, 0, len(items))
	for _, item := range items {
		key := item.PackageName + "\x00" + item.License + "\x00" + item.Severity
		if seen[key] {
			continue
		}
		seen[key] = true
		if item.Confidence == 0 {
			item.Confidence = 0.85
		}
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity {
			return licenseSeverityRank(out[i].Severity) > licenseSeverityRank(out[j].Severity)
		}
		leftName := strings.ToLower(out[i].PackageName)
		rightName := strings.ToLower(out[j].PackageName)
		if leftName != rightName {
			return leftName < rightName
		}
		return out[i].License < out[j].License
	})
	return out
}

func licenseSeverityRank(value string) int {
	switch strings.ToUpper(value) {
	case "CRITICAL":
		return 5
	case "HIGH":
		return 4
	case "MEDIUM":
		return 3
	case "LOW":
		return 2
	default:
		return 1
	}
}
