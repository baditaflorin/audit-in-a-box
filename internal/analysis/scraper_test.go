package analysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractManifestsFromHTML(t *testing.T) {
	candidates, err := ExtractManifestsFromHTML(`<html><body><pre>django==3.2.0
requests==2.25.1</pre></body></html>`)
	require.NoError(t, err)
	require.NotEmpty(t, candidates)
	require.Equal(t, "requirements.txt", candidates[0].FileName)
}
