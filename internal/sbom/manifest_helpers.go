package sbom

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
)

func manifestParseError(kind string, err error) error {
	message := strings.ToLower(err.Error())
	if strings.Contains(message, "unexpected end") || strings.Contains(message, "unexpected eof") || strings.Contains(message, "eof") {
		return models.WrapDomainError(
			"manifest_truncated",
			"The manifest looks incomplete.",
			"The parser reached the end of the file before the manifest structure was closed.",
			"Paste or upload the full "+kind+" file, then run the audit again.",
			true,
			err,
		)
	}
	return models.WrapDomainError(
		"manifest_invalid",
		"The manifest could not be parsed.",
		"The file was recognized as "+kind+", but its syntax is not valid for that format.",
		"Check the highlighted file type, fix the syntax, or upload the original manifest file.",
		true,
		err,
	)
}

func sortDependencies(deps []models.Dependency) {
	sort.Slice(deps, func(i, j int) bool {
		if deps[i].Ecosystem != deps[j].Ecosystem {
			return deps[i].Ecosystem < deps[j].Ecosystem
		}
		leftName := strings.ToLower(deps[i].Name)
		rightName := strings.ToLower(deps[j].Name)
		if leftName != rightName {
			return leftName < rightName
		}
		return deps[i].Version < deps[j].Version
	})
}

func dedupeByNameVersion(items []models.Dependency) []models.Dependency {
	seen := map[string]models.Dependency{}
	for _, item := range items {
		if item.Name == "" {
			continue
		}
		key := item.Ecosystem + "\x00" + strings.ToLower(item.Name) + "\x00" + item.Version
		existing, ok := seen[key]
		if !ok || item.Confidence > existing.Confidence || (item.Confidence == existing.Confidence && scopePriority(item.Scope) > scopePriority(existing.Scope)) {
			seen[key] = item
		}
	}
	out := make([]models.Dependency, 0, len(seen))
	for _, item := range seen {
		out = append(out, item)
	}
	sortDependencies(out)
	return out
}

func scopePriority(scope string) int {
	switch scope {
	case "runtime":
		return 5
	case "development":
		return 4
	case "optional":
		return 3
	case "peer", "indirect":
		return 2
	default:
		return 1
	}
}

func asMap(value any) map[string]any {
	switch item := value.(type) {
	case map[string]any:
		return item
	case map[any]any:
		out := make(map[string]any, len(item))
		for key, val := range item {
			out[fmt.Sprint(key)] = val
		}
		return out
	default:
		return nil
	}
}

func asStringSlice(value any) []string {
	switch item := value.(type) {
	case []string:
		return item
	case []any:
		out := make([]string, 0, len(item))
		for _, value := range item {
			if text := stringValue(value); text != "" {
				out = append(out, text)
			}
		}
		return out
	default:
		return nil
	}
}

func stringValue(value any) string {
	switch item := value.(type) {
	case string:
		return item
	default:
		return ""
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func sourceHash(content string) string {
	sum := sha256.Sum256([]byte(content))
	return hex.EncodeToString(sum[:])
}
