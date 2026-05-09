package analysis

import (
	"sort"
	"strings"
	"time"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
)

func ScoreReport(deps []models.Dependency, vulns []models.Vulnerability, licenses []models.LicenseRisk, health []models.MaintainerHealth) models.RiskScore {
	score := 0
	counts := map[string]int{
		"dependencies":     len(deps),
		"vulnerabilities":  len(vulns),
		"license_risks":    len(licenses),
		"maintainer_risks": 0,
	}
	var factors []string

	for _, vuln := range vulns {
		switch strings.ToUpper(vuln.Severity) {
		case "CRITICAL":
			score += 18
			counts["critical"]++
		case "HIGH":
			score += 10
			counts["high"]++
		case "MEDIUM":
			score += 5
			counts["medium"]++
		case "LOW":
			score += 2
			counts["low"]++
		default:
			score++
			counts["unknown"]++
		}
	}

	if counts["critical"] > 0 {
		factors = append(factors, "critical vulnerabilities are present")
	}
	if counts["high"] > 0 {
		factors = append(factors, "high-severity vulnerabilities are present")
	}

	for _, item := range licenses {
		switch strings.ToUpper(item.Severity) {
		case "HIGH", "CRITICAL":
			score += 8
		case "MEDIUM":
			score += 4
		default:
			score += 1
		}
	}
	if len(licenses) > 0 {
		factors = append(factors, "license review is required")
	}

	now := time.Now().UTC()
	for _, item := range health {
		stale := item.LastCommit != nil && now.Sub(*item.LastCommit) > 540*24*time.Hour
		fragile := item.BusFactor > 0 && item.BusFactor <= 1
		if stale || fragile || item.Score < 50 {
			counts["maintainer_risks"]++
			score += 4
		}
	}
	if counts["maintainer_risks"] > 0 {
		factors = append(factors, "some packages have stale or fragile maintainer signals")
	}

	if len(deps) > 100 {
		score += 5
		factors = append(factors, "large dependency surface")
	}

	if score > 100 {
		score = 100
	}
	sort.Strings(factors)

	return models.RiskScore{
		Score:      score,
		Grade:      grade(score),
		Confidence: scoreConfidence(deps, vulns, licenses),
		Counts:     counts,
		Factors:    factors,
	}
}

func grade(score int) string {
	switch {
	case score >= 80:
		return "critical"
	case score >= 60:
		return "high"
	case score >= 35:
		return "medium"
	case score >= 15:
		return "low"
	default:
		return "clean"
	}
}

func DeduplicateDependencies(items []models.Dependency) []models.Dependency {
	seen := map[string]models.Dependency{}
	for _, item := range items {
		if item.Name == "" {
			continue
		}
		key := item.Ecosystem + "\x00" + strings.ToLower(item.Name)
		existing, ok := seen[key]
		if !ok || (len(existing.Licenses) == 0 && len(item.Licenses) > 0) || existing.Version == "" {
			seen[key] = item
		}
	}

	out := make([]models.Dependency, 0, len(seen))
	for _, item := range seen {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Ecosystem != out[j].Ecosystem {
			return out[i].Ecosystem < out[j].Ecosystem
		}
		leftName := strings.ToLower(out[i].Name)
		rightName := strings.ToLower(out[j].Name)
		if leftName != rightName {
			return leftName < rightName
		}
		return out[i].Version < out[j].Version
	})
	return out
}

func DeduplicateVulnerabilities(items []models.Vulnerability) []models.Vulnerability {
	seen := map[string]bool{}
	out := make([]models.Vulnerability, 0, len(items))
	for _, item := range items {
		key := item.ID + "\x00" + strings.ToLower(item.PackageName) + "\x00" + item.Source
		if item.ID == "" || seen[key] {
			continue
		}
		seen[key] = true
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity {
			return severityRank(out[i].Severity) > severityRank(out[j].Severity)
		}
		if out[i].PackageName != out[j].PackageName {
			return out[i].PackageName < out[j].PackageName
		}
		return out[i].ID < out[j].ID
	})
	return out
}

func scoreConfidence(deps []models.Dependency, vulns []models.Vulnerability, licenses []models.LicenseRisk) float64 {
	if len(deps) == 0 {
		return 0.35
	}
	total := 0.0
	count := 0.0
	for _, dep := range deps {
		if dep.Confidence > 0 {
			total += dep.Confidence
			count++
		}
	}
	for _, vuln := range vulns {
		if vuln.Confidence > 0 {
			total += vuln.Confidence
			count++
		}
	}
	for _, license := range licenses {
		if license.Confidence > 0 {
			total += license.Confidence
			count++
		}
	}
	if count == 0 {
		return 0.6
	}
	return total / count
}

func severityRank(value string) int {
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
