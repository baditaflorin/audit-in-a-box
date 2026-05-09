package analysis

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"github.com/baditaflorin/audit-in-a-box/internal/sbom"
	"github.com/stretchr/testify/require"
)

type realDataExpectation struct {
	File                                       string   `json:"file"`
	Mode                                       string   `json:"mode"`
	FileName                                   string   `json:"file_name"`
	ExpectedEcosystem                          string   `json:"expected_ecosystem"`
	MinDependencies                            int      `json:"min_dependencies"`
	MaxUnknownLicenseRisks                     int      `json:"max_unknown_license_risks"`
	MustNotScoreCriticalWithoutVulnerabilities bool     `json:"must_not_score_critical_without_vulnerabilities"`
	ExpectedErrorCode                          string   `json:"expected_error_code"`
	ExpectedRecoverable                        bool     `json:"expected_recoverable"`
	ExpectedNextStepContains                   string   `json:"expected_next_step_contains"`
	ForbiddenWarningSubstrings                 []string `json:"forbidden_warning_substrings"`
}

func TestRealDataFixturesParseWithUsefulFirstGuess(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "test", "fixtures", "realdata")
	expectedFiles, err := filepath.Glob(filepath.Join(fixtureDir, "*.expected.json"))
	require.NoError(t, err)
	require.NotEmpty(t, expectedFiles)

	for _, expectedFile := range expectedFiles {
		expectedFile := expectedFile
		t.Run(filepath.Base(expectedFile), func(t *testing.T) {
			expected := readExpectation(t, expectedFile)
			content := readFixture(t, fixtureDir, expected.File)
			fileName := expected.FileName
			if fileName == "" {
				fileName = expected.File
			}

			if expected.Mode == "pasted_html" {
				candidates, err := ExtractManifestsFromHTML(content)
				require.NoError(t, err)
				require.NotEmpty(t, candidates)
				content = candidates[0].Content
				fileName = candidates[0].FileName
			}

			deps, warnings, inference, err := sbom.AnalyzeManifest(fileName, content)
			if expected.ExpectedErrorCode != "" {
				require.Error(t, err)
				var domainErr models.DomainError
				require.True(t, errors.As(err, &domainErr), "expected a domain error, got %T", err)
				require.Equal(t, expected.ExpectedErrorCode, domainErr.Code)
				require.Equal(t, expected.ExpectedRecoverable, domainErr.Recoverable)
				require.Contains(t, strings.ToLower(domainErr.NextStep), strings.ToLower(expected.ExpectedNextStepContains))
				return
			}

			require.NoError(t, err)
			require.Equal(t, expected.ExpectedEcosystem, inference.Ecosystem)
			require.GreaterOrEqual(t, len(deps), expected.MinDependencies)
			require.GreaterOrEqual(t, inference.Confidence, 0.8)

			for _, forbidden := range expected.ForbiddenWarningSubstrings {
				require.NotContains(t, strings.Join(warnings, "\n"), forbidden)
			}

			licenseRisks := AssessLicenseRisks(deps, nil)
			require.LessOrEqual(t, unknownLicenseRisks(licenseRisks), expected.MaxUnknownLicenseRisks)

			score := ScoreReport(deps, nil, licenseRisks, nil)
			if expected.MustNotScoreCriticalWithoutVulnerabilities {
				require.NotEqual(t, "critical", score.Grade)
				require.Less(t, score.Score, 80)
			}

			depsAgain, _, _, err := sbom.AnalyzeManifest(fileName, content)
			require.NoError(t, err)
			first, err := json.Marshal(deps)
			require.NoError(t, err)
			second, err := json.Marshal(depsAgain)
			require.NoError(t, err)
			require.JSONEq(t, string(first), string(second), "same fixture must parse deterministically")
		})
	}
}

func readExpectation(t *testing.T, path string) realDataExpectation {
	t.Helper()
	data, err := os.ReadFile(path) // #nosec G304 -- test fixture paths are discovered from the committed fixture directory.
	require.NoError(t, err)
	var expected realDataExpectation
	require.NoError(t, json.Unmarshal(data, &expected))
	return expected
}

func readFixture(t *testing.T, dir string, fileName string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(dir, fileName)) // #nosec G304 -- fixture names come from committed expected JSON files.
	require.NoError(t, err)
	return string(data)
}

func unknownLicenseRisks(items []models.LicenseRisk) int {
	count := 0
	for _, item := range items {
		if strings.EqualFold(item.License, "UNKNOWN") {
			count++
		}
	}
	return count
}
