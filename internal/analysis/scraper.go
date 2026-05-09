package analysis

import (
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"github.com/baditaflorin/audit-in-a-box/internal/sbom"
)

var whitespace = regexp.MustCompile(`\s+`)

func ExtractManifestsFromHTML(html string) ([]models.ScrapedManifest, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	var candidates []models.ScrapedManifest
	if candidate, ok := githubBlobCandidate(doc); ok {
		candidates = append(candidates, candidate)
	}
	doc.Find("pre, code, textarea").Each(func(_ int, selection *goquery.Selection) {
		text := strings.TrimSpace(selection.Text())
		if text == "" {
			return
		}
		if candidate, ok := manifestCandidate(text); ok {
			candidates = append(candidates, candidate)
		}
	})

	if len(candidates) == 0 {
		text := strings.TrimSpace(doc.Text())
		if candidate, ok := manifestCandidate(text); ok {
			candidates = append(candidates, candidate)
		}
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Score > candidates[j].Score
	})

	return candidates, nil
}

func githubBlobCandidate(doc *goquery.Document) (models.ScrapedManifest, bool) {
	var lines []string
	doc.Find(".blob-code, td.blob-code, table.highlight td:nth-child(2)").Each(func(_ int, selection *goquery.Selection) {
		text := strings.TrimRight(selection.Text(), "\n")
		if strings.TrimSpace(text) == "" {
			lines = append(lines, "")
			return
		}
		lines = append(lines, text)
	})
	if len(lines) == 0 {
		return models.ScrapedManifest{}, false
	}
	text := strings.TrimSpace(strings.Join(lines, "\n"))
	candidate, ok := manifestCandidate(text)
	if !ok {
		return models.ScrapedManifest{}, false
	}
	candidate.Score += 35
	candidate.Reason = "GitHub blob code table"
	return candidate, true
}

func manifestCandidate(text string) (models.ScrapedManifest, bool) {
	normalized := strings.TrimSpace(text)
	score := 0
	ecosystem := sbom.DetectEcosystem("", normalized)
	kind := sbom.InferManifest("", normalized).Kind
	switch ecosystem {
	case "npm":
		score += 60
		if strings.Contains(normalized, `"dependencies"`) {
			score += 30
		}
		if kind == "pnpm-lock" || kind == "package-lock" {
			score += 20
		}
	case "go":
		score += 70
		if strings.Contains(normalized, "\nrequire ") || strings.Contains(normalized, "\nrequire (") {
			score += 20
		}
	case "python":
		score += 50
		for _, line := range strings.Split(normalized, "\n") {
			if strings.Contains(line, "==") || strings.Contains(line, ">=") {
				score += 5
			}
		}
	default:
		if looksLikeRequirements(normalized) {
			ecosystem = "python"
			score = 55
		}
	}

	if score < 50 {
		return models.ScrapedManifest{}, false
	}

	compact := normalized
	if len(compact) > 20000 {
		compact = compact[:20000]
	}

	fileName := sbom.DefaultFileName(ecosystem)
	switch kind {
	case "package-lock":
		fileName = "package-lock.json"
	case "pnpm-lock":
		fileName = "pnpm-lock.yaml"
	case "pyproject":
		fileName = "pyproject.toml"
	}

	return models.ScrapedManifest{
		FileName:  fileName,
		Ecosystem: ecosystem,
		Content:   compact,
		Score:     score,
		Reason:    "manifest-like text block",
	}, true
}

func looksLikeRequirements(text string) bool {
	lines := strings.Split(text, "\n")
	matches := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = whitespace.ReplaceAllString(line, "")
		if strings.Contains(line, "==") || strings.Contains(line, ">=") || strings.Contains(line, "~=") {
			matches++
		}
	}
	return matches >= 2
}
