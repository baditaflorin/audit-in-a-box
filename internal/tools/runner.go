package tools

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
)

type Runner struct {
	Timeout time.Duration
}

func (r Runner) Command(ctx context.Context, name string, args ...string) (string, string, error) {
	if r.Timeout <= 0 {
		r.Timeout = 60 * time.Second
	}

	cmdCtx, cancel := context.WithTimeout(ctx, r.Timeout)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, name, args...) // #nosec G204 -- scanner command names are selected by backend code.
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if errors.Is(cmdCtx.Err(), context.DeadlineExceeded) {
		return stdout.String(), stderr.String(), fmt.Errorf("%s timed out after %s", name, r.Timeout)
	}
	if err != nil {
		return stdout.String(), stderr.String(), fmt.Errorf("%s failed: %w: %s", name, err, strings.TrimSpace(stderr.String()))
	}

	return stdout.String(), stderr.String(), nil
}

func CheckTools(ctx context.Context, runner Runner) map[string]models.ToolStatus {
	tools := []struct {
		name     string
		required bool
		args     []string
	}{
		{name: "trivy", required: true, args: []string{"--version"}},
		{name: "syft", required: true, args: []string{"version"}},
		{name: "grype", required: true, args: []string{"version"}},
		{name: "duckdb", required: true, args: []string{"--version"}},
	}

	status := make(map[string]models.ToolStatus, len(tools))
	for _, candidate := range tools {
		path, err := exec.LookPath(candidate.name)
		item := models.ToolStatus{
			Name:      candidate.name,
			Required:  candidate.required,
			Found:     err == nil,
			Available: err == nil,
			Path:      path,
		}
		if err != nil {
			item.Error = err.Error()
			status[candidate.name] = item
			continue
		}

		out, _, versionErr := runner.Command(ctx, candidate.name, candidate.args...)
		if versionErr != nil {
			item.Error = versionErr.Error()
		} else {
			item.Version = firstLines(out, 3)
		}
		status[candidate.name] = item
	}

	return status
}

func firstLines(value string, maxLines int) string {
	lines := strings.Split(strings.TrimSpace(value), "\n")
	if len(lines) > maxLines {
		lines = lines[:maxLines]
	}
	return strings.Join(lines, "\n")
}
