package sbom

import (
	"regexp"
	"strings"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"github.com/pelletier/go-toml/v2"
)

var (
	requirementName = regexp.MustCompile(`^([A-Za-z0-9_.-]+)(?:\[[^\]]+\])?\s*(.*)$`)
	requirementOp   = regexp.MustCompile(`^(==|!=|<=|>=|~=|<|>|===)`)
)

func parsePyproject(content string) ([]models.Dependency, []string, error) {
	var raw map[string]any
	if err := toml.Unmarshal([]byte(content), &raw); err != nil {
		return nil, nil, manifestParseError("pyproject.toml", err)
	}

	var deps []models.Dependency
	project := asMap(raw["project"])
	for _, item := range asStringSlice(project["dependencies"]) {
		if dep, ok := dependencyFromRequirement(item, "runtime", "project.dependencies in pyproject.toml"); ok {
			deps = append(deps, dep)
		}
	}
	for group, values := range asMap(project["optional-dependencies"]) {
		for _, item := range asStringSlice(values) {
			if dep, ok := dependencyFromRequirement(item, "optional:"+group, "project.optional-dependencies in pyproject.toml"); ok {
				deps = append(deps, dep)
			}
		}
	}

	tool := asMap(raw["tool"])
	poetry := asMap(tool["poetry"])
	for name, value := range asMap(poetry["dependencies"]) {
		if dep, ok := dependencyFromPoetryEntry(name, value, "runtime", "tool.poetry.dependencies in pyproject.toml"); ok {
			deps = append(deps, dep)
		}
	}
	for name, value := range asMap(poetry["dev-dependencies"]) {
		if dep, ok := dependencyFromPoetryEntry(name, value, "development", "tool.poetry.dev-dependencies in pyproject.toml"); ok {
			deps = append(deps, dep)
		}
	}
	for groupName, groupValue := range asMap(poetry["group"]) {
		dependencySection := asMap(asMap(groupValue)["dependencies"])
		for name, value := range dependencySection {
			if dep, ok := dependencyFromPoetryEntry(name, value, "group:"+groupName, "tool.poetry.group dependencies in pyproject.toml"); ok {
				deps = append(deps, dep)
			}
		}
	}

	return dedupeByNameVersion(deps), nil, nil
}

func parseRequirements(content string) ([]models.Dependency, []string) {
	var (
		deps     []models.Dependency
		warnings []string
	)
	for _, rawLine := range strings.Split(content, "\n") {
		line := normalizeRequirementLine(rawLine)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "-r ") || strings.HasPrefix(line, "--requirement") {
			warnings = append(warnings, "Nested requirements file reference skipped: "+line)
			continue
		}
		if strings.HasPrefix(line, "-") {
			continue
		}
		if dep, ok := dependencyFromRequirement(line, "runtime", "line in requirements file"); ok {
			deps = append(deps, dep)
		}
	}

	return dedupeByNameVersion(deps), warnings
}

func dependencyFromRequirement(value string, scope string, reason string) (models.Dependency, bool) {
	value = normalizeRequirementLine(value)
	if value == "" || strings.HasPrefix(value, "-") {
		return models.Dependency{}, false
	}
	if strings.Contains(value, " @ ") {
		parts := strings.SplitN(value, " @ ", 2)
		name := normalizedPythonName(parts[0])
		if name == "" {
			return models.Dependency{}, false
		}
		return models.Dependency{Name: name, Version: "@ " + strings.TrimSpace(parts[1]), Ecosystem: "python", Scope: scope, Source: "manifest", Confidence: 0.9, Reasons: []string{reason}}, true
	}
	match := requirementName.FindStringSubmatch(value)
	if len(match) < 2 {
		return models.Dependency{}, false
	}
	name := normalizedPythonName(match[1])
	if name == "" {
		return models.Dependency{}, false
	}
	version := strings.TrimSpace(match[2])
	if !requirementOp.MatchString(version) {
		version = strings.TrimSpace(strings.TrimPrefix(version, " "))
	}
	return models.Dependency{Name: name, Version: version, Ecosystem: "python", Scope: scope, Source: "manifest", Confidence: 0.94, Reasons: []string{reason}}, true
}

func dependencyFromPoetryEntry(name string, value any, scope string, reason string) (models.Dependency, bool) {
	name = normalizedPythonName(name)
	if name == "" || name == "python" {
		return models.Dependency{}, false
	}
	version := stringValue(value)
	if version == "" {
		version = stringValue(asMap(value)["version"])
	}
	return models.Dependency{Name: name, Version: version, Ecosystem: "python", Scope: scope, Source: "manifest", Confidence: 0.95, Reasons: []string{reason}}, true
}

func normalizeRequirementLine(line string) string {
	line = strings.TrimSpace(strings.TrimPrefix(line, "\ufeff"))
	if line == "" || strings.HasPrefix(line, "#") {
		return ""
	}
	if index := strings.Index(line, " --hash="); index >= 0 {
		line = line[:index]
	}
	if index := strings.Index(line, ";"); index >= 0 {
		line = line[:index]
	}
	if index := strings.Index(line, " #"); index >= 0 {
		line = line[:index]
	}
	return strings.TrimSpace(line)
}

func looksLikeRequirements(text string) bool {
	lines := strings.Split(text, "\n")
	matches := 0
	for _, line := range lines {
		line = normalizeRequirementLine(line)
		if line == "" {
			continue
		}
		if _, ok := dependencyFromRequirement(line, "runtime", "requirements-like line"); ok {
			matches++
		}
	}
	return matches >= 2
}

func normalizedPythonName(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	name = strings.ReplaceAll(name, "_", "-")
	return name
}
