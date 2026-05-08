package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/baditaflorin/audit-in-a-box/internal/config"
	"github.com/baditaflorin/audit-in-a-box/internal/models"
)

type Summarizer struct {
	Config config.Config
	Client *http.Client
}

func (s Summarizer) Summarize(ctx context.Context, report models.AuditReport) string {
	prompt := buildPrompt(report)
	if s.Config.OllamaBaseURL != "" {
		if summary, err := s.ollama(ctx, prompt); err == nil && strings.TrimSpace(summary) != "" {
			return strings.TrimSpace(summary)
		}
	}

	if s.Config.LLMCommand != "" {
		if summary, err := s.command(ctx, prompt); err == nil && strings.TrimSpace(summary) != "" {
			return strings.TrimSpace(summary)
		}
	}

	return fallback(report)
}

func (s Summarizer) ollama(ctx context.Context, prompt string) (string, error) {
	client := s.Client
	if client == nil {
		client = &http.Client{Timeout: 45 * time.Second}
	}

	body, _ := json.Marshal(map[string]any{
		"model":  s.Config.OllamaModel,
		"prompt": prompt,
		"stream": false,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(s.Config.OllamaBaseURL, "/")+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("ollama returned %s", resp.Status)
	}

	var payload struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	return payload.Response, nil
}

func (s Summarizer) command(ctx context.Context, prompt string) (string, error) {
	parts := strings.Fields(s.Config.LLMCommand)
	if len(parts) == 0 {
		return "", fmt.Errorf("empty llm command")
	}

	cmdCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	cmd := exec.CommandContext(cmdCtx, parts[0], parts[1:]...)
	cmd.Stdin = strings.NewReader(prompt)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("llm command failed: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return stdout.String(), nil
}

func buildPrompt(report models.AuditReport) string {
	return fmt.Sprintf(`Write a concise dependency risk summary for a developer.
Score: %d/100 (%s)
Dependencies: %d
Vulnerabilities: %d
License risks: %d
Maintainer risks: %d
Key factors: %s
Use plain English. Do not invent exact remediation versions unless present in the evidence.`,
		report.Risk.Score,
		report.Risk.Grade,
		len(report.Dependencies),
		len(report.Vulnerabilities),
		len(report.LicenseRisks),
		report.Risk.Counts["maintainer_risks"],
		strings.Join(report.Risk.Factors, "; "),
	)
}

func fallback(report models.AuditReport) string {
	if report.Risk.Score == 0 {
		return "No obvious dependency risk was found from the available evidence. Keep scanner databases fresh and review unknown licenses before release."
	}

	parts := []string{
		fmt.Sprintf("This manifest currently scores %d/100 (%s risk).", report.Risk.Score, report.Risk.Grade),
		fmt.Sprintf("The audit found %d dependencies, %d vulnerabilities, and %d license-risk flags.", len(report.Dependencies), len(report.Vulnerabilities), len(report.LicenseRisks)),
	}
	if len(report.Risk.Factors) > 0 {
		parts = append(parts, "Main drivers: "+strings.Join(report.Risk.Factors, "; ")+".")
	}
	parts = append(parts, "Prioritize critical/high vulnerabilities, packages with restrictive licenses, and dependencies with stale or fragile maintainer signals.")
	return strings.Join(parts, " ")
}
