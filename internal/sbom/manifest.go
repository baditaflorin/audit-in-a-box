package sbom

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"golang.org/x/mod/modfile"
)

var requirementName = regexp.MustCompile(`^([A-Za-z0-9_.-]+)`)

func DetectEcosystem(fileName string, content string) string {
	name := strings.ToLower(fileName)
	switch {
	case strings.HasSuffix(name, "package.json") || strings.Contains(content, `"dependencies"`):
		return "npm"
	case strings.HasSuffix(name, "go.mod") || strings.HasPrefix(strings.TrimSpace(content), "module "):
		return "go"
	case strings.HasSuffix(name, "requirements.txt"):
		return "python"
	default:
		return "unknown"
	}
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

func ParseManifest(fileName string, content string) ([]models.Dependency, []string, error) {
	ecosystem := DetectEcosystem(fileName, content)
	switch ecosystem {
	case "npm":
		return parsePackageJSON(content)
	case "go":
		return parseGoMod(fileName, content)
	case "python":
		return parseRequirements(content), nil, nil
	default:
		return nil, []string{"unknown manifest type; scanner tools may still produce results"}, nil
	}
}

type packageJSON struct {
	Dependencies         map[string]string `json:"dependencies"`
	DevDependencies      map[string]string `json:"devDependencies"`
	PeerDependencies     map[string]string `json:"peerDependencies"`
	OptionalDependencies map[string]string `json:"optionalDependencies"`
}

func parsePackageJSON(content string) ([]models.Dependency, []string, error) {
	var manifest packageJSON
	if err := json.Unmarshal([]byte(content), &manifest); err != nil {
		return nil, nil, fmt.Errorf("parse package.json: %w", err)
	}

	var deps []models.Dependency
	appendDeps := func(scope string, values map[string]string) {
		for name, version := range values {
			deps = append(deps, models.Dependency{
				Name:      name,
				Version:   version,
				Ecosystem: "npm",
				Scope:     scope,
				Source:    "manifest",
			})
		}
	}

	appendDeps("runtime", manifest.Dependencies)
	appendDeps("development", manifest.DevDependencies)
	appendDeps("peer", manifest.PeerDependencies)
	appendDeps("optional", manifest.OptionalDependencies)
	sortDependencies(deps)

	return deps, nil, nil
}

func parseGoMod(fileName string, content string) ([]models.Dependency, []string, error) {
	parsed, err := modfile.Parse(fileName, []byte(content), nil)
	if err != nil {
		return nil, nil, fmt.Errorf("parse go.mod: %w", err)
	}

	deps := make([]models.Dependency, 0, len(parsed.Require))
	for _, req := range parsed.Require {
		scope := "runtime"
		if req.Indirect {
			scope = "indirect"
		}
		deps = append(deps, models.Dependency{
			Name:      req.Mod.Path,
			Version:   req.Mod.Version,
			Ecosystem: "go",
			Scope:     scope,
			Source:    "manifest",
		})
	}
	sortDependencies(deps)

	return deps, nil, nil
}

func parseRequirements(content string) []models.Dependency {
	var deps []models.Dependency
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") {
			continue
		}
		line = strings.Split(line, "#")[0]
		match := requirementName.FindStringSubmatch(line)
		if len(match) < 2 {
			continue
		}
		name := match[1]
		version := strings.TrimSpace(strings.TrimPrefix(line, name))
		deps = append(deps, models.Dependency{
			Name:      name,
			Version:   version,
			Ecosystem: "python",
			Scope:     "runtime",
			Source:    "manifest",
		})
	}
	sortDependencies(deps)

	return deps
}

func sortDependencies(deps []models.Dependency) {
	sort.Slice(deps, func(i, j int) bool {
		if deps[i].Ecosystem == deps[j].Ecosystem {
			return deps[i].Name < deps[j].Name
		}
		return deps[i].Ecosystem < deps[j].Ecosystem
	})
}
