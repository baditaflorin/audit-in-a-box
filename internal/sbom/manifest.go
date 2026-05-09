package sbom

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
)

const schemaVersion = "audit-report.v2"

type ManifestInference struct {
	FileName          string
	Ecosystem         string
	Kind              string
	Parser            string
	Confidence        float64
	Reasons           []string
	NormalizedContent string
	SourceHash        string
	SchemaVersion     string
}

func DetectEcosystem(fileName string, content string) string {
	return InferManifest(fileName, content).Ecosystem
}

func DefaultFileName(ecosystem string) string {
	switch ecosystem {
	case "npm":
		return "package.json"
	case "go":
		return "go.mod"
	case "python":
		return "requirements.txt"
	default:
		return "manifest.txt"
	}
}

func NormalizeContent(content string) string {
	normalized := strings.TrimPrefix(content, "\ufeff")
	normalized = strings.ReplaceAll(normalized, "\r\n", "\n")
	normalized = strings.ReplaceAll(normalized, "\r", "\n")
	normalized = strings.ReplaceAll(normalized, "\u00a0", " ")
	normalized = strings.ReplaceAll(normalized, "\x00", "")
	return strings.TrimSpace(normalized)
}

func InferManifest(fileName string, content string) ManifestInference {
	normalized := NormalizeContent(content)
	name := strings.ToLower(filepath.Base(strings.TrimSpace(fileName)))
	trimmed := strings.TrimSpace(normalized)
	lower := strings.ToLower(trimmed)

	inference := ManifestInference{
		FileName:          fileName,
		Ecosystem:         "unknown",
		Kind:              "unknown",
		Parser:            "none",
		Confidence:        0.1,
		NormalizedContent: normalized,
		SourceHash:        sourceHash(normalized),
		SchemaVersion:     schemaVersion,
	}

	set := func(ecosystem, kind, parser string, confidence float64, reasons ...string) ManifestInference {
		inference.Ecosystem = ecosystem
		inference.Kind = kind
		inference.Parser = parser
		inference.Confidence = confidence
		inference.Reasons = append(inference.Reasons, reasons...)
		return inference
	}

	switch {
	case name == "package-lock.json":
		return set("npm", "package-lock", "package-lock.json", 0.99, "file name is package-lock.json")
	case strings.HasSuffix(name, "pnpm-lock.yaml") || strings.HasSuffix(name, "pnpm-lock.yml"):
		return set("npm", "pnpm-lock", "pnpm-lock.yaml", 0.99, "file name is pnpm-lock.yaml")
	case strings.HasSuffix(name, "package.json"):
		return set("npm", "package-json", "package.json", 0.98, "file name is package.json")
	case strings.HasSuffix(name, "go.mod"):
		return set("go", "go-mod", "go.mod", 0.99, "file name is go.mod")
	case name == "pyproject.toml":
		return set("python", "pyproject", "pyproject.toml", 0.99, "file name is pyproject.toml")
	case strings.Contains(name, "requirements") && strings.HasSuffix(name, ".txt"):
		return set("python", "requirements", "requirements.txt", 0.96, "file name looks like a Python requirements file")
	case json.Valid([]byte(trimmed)) && (strings.Contains(trimmed, `"dependencies"`) || strings.Contains(trimmed, `"devDependencies"`) || strings.Contains(trimmed, `"peerDependencies"`)):
		return set("npm", "package-json", "package.json", 0.9, "JSON contains npm dependency fields")
	case strings.HasPrefix(trimmed, "module ") || strings.Contains(trimmed, "\nrequire "):
		return set("go", "go-mod", "go.mod", 0.92, "content matches go.mod module/require shape")
	case strings.Contains(lower, "[tool.poetry.dependencies]") || strings.Contains(lower, "[project]"):
		return set("python", "pyproject", "pyproject.toml", 0.88, "content contains pyproject dependency sections")
	case strings.Contains(lower, "lockfileversion:") && (strings.Contains(lower, "\nimporters:") || strings.Contains(lower, "\npackages:")):
		return set("npm", "pnpm-lock", "pnpm-lock.yaml", 0.84, "YAML contains pnpm lockfile sections")
	case looksLikeRequirements(trimmed):
		return set("python", "requirements", "requirements.txt", 0.82, "lines look like Python requirement specifiers")
	default:
		inference.Reasons = []string{"input does not match package.json, package-lock.json, pnpm-lock.yaml, go.mod, pyproject.toml, or requirements.txt"}
		return inference
	}
}

func ParseManifest(fileName string, content string) ([]models.Dependency, []string, error) {
	deps, warnings, _, err := AnalyzeManifest(fileName, content)
	return deps, warnings, err
}

func AnalyzeManifest(fileName string, content string) ([]models.Dependency, []string, ManifestInference, error) {
	inference := InferManifest(fileName, content)
	content = inference.NormalizedContent

	var (
		deps     []models.Dependency
		warnings []string
		err      error
	)
	switch inference.Kind {
	case "package-json":
		deps, warnings, err = parsePackageJSON(content)
	case "package-lock":
		deps, warnings, err = parsePackageLock(content)
	case "pnpm-lock":
		deps, warnings, err = parsePNPMLock(content)
	case "go-mod":
		deps, warnings, err = parseGoMod(fileName, content)
	case "pyproject":
		deps, warnings, err = parsePyproject(content)
	case "requirements":
		deps, warnings = parseRequirements(content)
	default:
		warnings = append(warnings, "Unsupported manifest type. Upload package.json, package-lock.json, pnpm-lock.yaml, go.mod, pyproject.toml, or requirements.txt.")
	}
	if err != nil {
		return nil, warnings, inference, err
	}

	for i := range deps {
		if deps[i].Confidence == 0 {
			deps[i].Confidence = inference.Confidence
		}
		if len(deps[i].Reasons) == 0 {
			deps[i].Reasons = []string{"parsed from " + inference.Kind}
		}
	}
	sortDependencies(deps)
	if len(deps) > 100 {
		warnings = append(warnings, fmt.Sprintf("Large dependency surface detected: %d dependencies. The report is prioritized and keeps confidence on inferred data.", len(deps)))
	}
	return deps, warnings, inference, nil
}
