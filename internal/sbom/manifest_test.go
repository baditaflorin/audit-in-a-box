package sbom

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePackageJSON(t *testing.T) {
	deps, warnings, err := ParseManifest("package.json", `{"dependencies":{"lodash":"4.17.20"},"devDependencies":{"jest":"26.6.0"}}`)
	require.NoError(t, err)
	require.Empty(t, warnings)
	require.Len(t, deps, 2)
	require.Equal(t, "npm", deps[0].Ecosystem)
}

func TestParseGoMod(t *testing.T) {
	deps, _, err := ParseManifest("go.mod", "module example.com/demo\n\ngo 1.22\n\nrequire github.com/go-chi/chi/v5 v5.0.10\n")
	require.NoError(t, err)
	require.Len(t, deps, 1)
	require.Equal(t, "github.com/go-chi/chi/v5", deps[0].Name)
}

func TestParseRequirements(t *testing.T) {
	deps, _, err := ParseManifest("requirements.txt", "django==3.2.0\n# comment\nrequests>=2.25\n")
	require.NoError(t, err)
	require.Len(t, deps, 2)
	require.Equal(t, "python", deps[0].Ecosystem)
}
