package sbom

import (
	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"golang.org/x/mod/modfile"
)

func parseGoMod(fileName string, content string) ([]models.Dependency, []string, error) {
	parsed, err := modfile.Parse(fileName, []byte(content), nil)
	if err != nil {
		return nil, nil, manifestParseError("go.mod", err)
	}

	deps := make([]models.Dependency, 0, len(parsed.Require))
	for _, req := range parsed.Require {
		scope := "runtime"
		if req.Indirect {
			scope = "indirect"
		}
		deps = append(deps, models.Dependency{
			Name:       req.Mod.Path,
			Version:    req.Mod.Version,
			Ecosystem:  "go",
			Scope:      scope,
			Source:     "manifest",
			Confidence: 0.99,
			Reasons:    []string{"require directive in go.mod"},
		})
	}

	return deps, nil, nil
}
