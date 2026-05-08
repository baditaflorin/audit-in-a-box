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

func manifestCandidate(text string) (models.ScrapedManifest, bool) {
	normalized := strings.TrimSpace(text)
	score := 0
	ecosystem := sbom.DetectEcosystem("", normalized)
	switch ecosystem {
	case "npm":
		score += 60
		if strings.Contains(normalized, `"dependencies"`) {
			score += 30
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

	return models.ScrapedManifest{
		FileName:  sbom.DefaultFileName(ecosystem),
		Ecosystem: ecosystem,
		Content:   compact,
		Score:     score,
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
