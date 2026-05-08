package analysis

import (
	"strings"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
)

func AssessLicenseRisks(deps []models.Dependency, scannerRisks []models.LicenseRisk) []models.LicenseRisk {
	risks := append([]models.LicenseRisk{}, scannerRisks...)
	for _, dep := range deps {
		if dep.Name == "" {
			continue
		}
		if len(dep.Licenses) == 0 {
			risks = append(risks, models.LicenseRisk{
				PackageName: dep.Name,
				License:     "UNKNOWN",
				Severity:    "MEDIUM",
				Reason:      "no license detected in SBOM output",
			})
			continue
		}
		for _, license := range dep.Licenses {
			normalized := strings.ToUpper(license)
			switch {
			case strings.Contains(normalized, "AGPL"):
				risks = append(risks, models.LicenseRisk{PackageName: dep.Name, License: license, Severity: "HIGH", Reason: "strong network copyleft license"})
			case strings.Contains(normalized, "GPL") && !strings.Contains(normalized, "LGPL"):
				risks = append(risks, models.LicenseRisk{PackageName: dep.Name, License: license, Severity: "HIGH", Reason: "strong copyleft license"})
			case strings.Contains(normalized, "LGPL"), strings.Contains(normalized, "MPL"), strings.Contains(normalized, "CDDL"):
				risks = append(risks, models.LicenseRisk{PackageName: dep.Name, License: license, Severity: "MEDIUM", Reason: "weak copyleft license"})
			case strings.Contains(normalized, "BUSL"), strings.Contains(normalized, "SSPL"), strings.Contains(normalized, "PROPRIETARY"):
				risks = append(risks, models.LicenseRisk{PackageName: dep.Name, License: license, Severity: "HIGH", Reason: "source-available or proprietary-style license"})
			}
		}
	}

	return dedupeLicenseRisks(risks)
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
		out = append(out, item)
	}
	return out
}
