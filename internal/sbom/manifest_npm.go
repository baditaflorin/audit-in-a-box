package sbom

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"gopkg.in/yaml.v3"
)

type packageJSON struct {
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

func parsePackageJSON(content string) ([]models.Dependency, []string, error) {
	var manifest packageJSON
	if err := json.Unmarshal([]byte(content), &manifest); err != nil {
		return nil, nil, manifestParseError("package.json", err)
	}

	var deps []models.Dependency
	appendDeps := func(scope string, values map[string]string) {
		for name, version := range values {
			deps = append(deps, models.Dependency{
				Name:       name,
				Version:    version,
				Ecosystem:  "npm",
				Scope:      scope,
				Source:     "manifest",
				Confidence: 0.98,
				Reasons:    []string{"declared in " + scope + " dependencies"},
			})
		}
	}

	appendDeps("runtime", manifest.Dependencies)
	appendDeps("development", manifest.DevDependencies)
	appendDeps("peer", manifest.PeerDependencies)
	appendDeps("optional", manifest.OptionalDependencies)

	return deps, nil, nil
}

func parsePackageLock(content string) ([]models.Dependency, []string, error) {
	var manifest struct {
		Packages map[string]struct {
			Name     string          `json:"name"`
			Version  string          `json:"version"`
			Dev      bool            `json:"dev"`
			Optional bool            `json:"optional"`
			License  json.RawMessage `json:"license"`
		} `json:"packages"`
		Dependencies map[string]struct {
			Version string          `json:"version"`
			Dev     bool            `json:"dev"`
			License json.RawMessage `json:"license"`
		} `json:"dependencies"`
	}
	if err := json.Unmarshal([]byte(content), &manifest); err != nil {
		return nil, nil, manifestParseError("package-lock.json", err)
	}

	var deps []models.Dependency
	for path, pkg := range manifest.Packages {
		if path == "" {
			continue
		}
		name := pkg.Name
		if name == "" {
			name = npmNameFromNodeModulesPath(path)
		}
		if name == "" {
			continue
		}
		scope := "runtime"
		if pkg.Dev {
			scope = "development"
		}
		if pkg.Optional {
			scope = "optional"
		}
		deps = append(deps, models.Dependency{
			Name:       name,
			Version:    pkg.Version,
			Ecosystem:  "npm",
			Scope:      scope,
			Licenses:   parseJSONLicenses(pkg.License),
			Source:     "package-lock",
			Confidence: 0.96,
			Reasons:    []string{"resolved package entry in package-lock.json"},
		})
	}
	for name, pkg := range manifest.Dependencies {
		if name == "" {
			continue
		}
		scope := "runtime"
		if pkg.Dev {
			scope = "development"
		}
		deps = append(deps, models.Dependency{
			Name:       name,
			Version:    pkg.Version,
			Ecosystem:  "npm",
			Scope:      scope,
			Licenses:   parseJSONLicenses(pkg.License),
			Source:     "package-lock",
			Confidence: 0.92,
			Reasons:    []string{"resolved dependency entry in package-lock.json"},
		})
	}

	return dedupeByNameVersion(deps), nil, nil
}

func parsePNPMLock(content string) ([]models.Dependency, []string, error) {
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(content)))
	var deps []models.Dependency
	for {
		var raw map[string]any
		err := decoder.Decode(&raw)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, nil, manifestParseError("pnpm-lock.yaml", err)
		}
		deps = append(deps, depsFromPNPMLockDoc(raw)...)
	}

	return dedupeByNameVersion(deps), nil, nil
}

func depsFromPNPMLockDoc(raw map[string]any) []models.Dependency {
	var deps []models.Dependency
	importers := asMap(raw["importers"])
	for _, importerValue := range importers {
		importer := asMap(importerValue)
		for _, dependencyField := range []struct {
			name  string
			scope string
		}{
			{name: "dependencies", scope: "runtime"},
			{name: "devDependencies", scope: "development"},
			{name: "optionalDependencies", scope: "optional"},
		} {
			for name, spec := range asMap(importer[dependencyField.name]) {
				deps = append(deps, models.Dependency{
					Name:       name,
					Version:    pnpmSpecVersion(spec),
					Ecosystem:  "npm",
					Scope:      dependencyField.scope,
					Source:     "pnpm-lock",
					Confidence: 0.93,
					Reasons:    []string{"declared in pnpm importer " + dependencyField.name},
				})
			}
		}
	}
	for key, pkgValue := range asMap(raw["packages"]) {
		name, version := parsePNPMPackageKey(key)
		if name == "" {
			continue
		}
		pkg := asMap(pkgValue)
		scope := "resolved"
		if dev, ok := pkg["dev"].(bool); ok && dev {
			scope = "development"
		}
		deps = append(deps, models.Dependency{
			Name:       name,
			Version:    firstNonEmpty(stringValue(pkg["version"]), version),
			Ecosystem:  "npm",
			Scope:      scope,
			Source:     "pnpm-lock",
			Confidence: 0.91,
			Reasons:    []string{"resolved package entry in pnpm-lock.yaml"},
		})
	}

	return deps
}

func parseJSONLicenses(raw json.RawMessage) []string {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}
	var single string
	if err := json.Unmarshal(raw, &single); err == nil && single != "" {
		return []string{single}
	}
	var many []string
	if err := json.Unmarshal(raw, &many); err == nil {
		return many
	}
	return nil
}

func npmNameFromNodeModulesPath(path string) string {
	const marker = "node_modules/"
	index := strings.LastIndex(path, marker)
	if index < 0 {
		return ""
	}
	name := path[index+len(marker):]
	if strings.HasPrefix(name, "@") {
		parts := strings.Split(name, "/")
		if len(parts) >= 2 {
			return parts[0] + "/" + parts[1]
		}
	}
	if index := strings.Index(name, "/"); index >= 0 {
		name = name[:index]
	}
	return name
}

func parsePNPMPackageKey(key string) (string, string) {
	key = strings.TrimPrefix(strings.TrimSpace(key), "/")
	key = strings.Split(key, "(")[0]
	if key == "" {
		return "", ""
	}
	separator := strings.LastIndex(key, "@")
	if strings.HasPrefix(key, "@") {
		if next := strings.Index(key[1:], "@"); next >= 0 {
			separator = next + 1
		}
	}
	if separator <= 0 || separator >= len(key)-1 {
		return key, ""
	}
	return key[:separator], key[separator+1:]
}

func pnpmSpecVersion(value any) string {
	switch item := value.(type) {
	case string:
		return item
	case map[string]any:
		return firstNonEmpty(stringValue(item["version"]), stringValue(item["specifier"]))
	default:
		return ""
	}
}
