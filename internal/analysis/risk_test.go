package analysis

import (
	"testing"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
	"github.com/stretchr/testify/require"
)

func TestScoreReportCriticalVulnerability(t *testing.T) {
	score := ScoreReport(
		[]models.Dependency{{Name: "lodash", Ecosystem: "npm"}},
		[]models.Vulnerability{{ID: "CVE-1", Severity: "CRITICAL", PackageName: "lodash"}},
		nil,
		nil,
	)

	require.Equal(t, 18, score.Score)
	require.Equal(t, "low", score.Grade)
	require.Equal(t, 1, score.Counts["critical"])
}

func TestAssessLicenseRisksFlagsUnknown(t *testing.T) {
	risks := AssessLicenseRisks([]models.Dependency{{Name: "mystery", Ecosystem: "npm"}}, nil)
	require.Len(t, risks, 1)
	require.Equal(t, "UNKNOWN", risks[0].License)
}
