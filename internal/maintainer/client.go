package maintainer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/baditaflorin/audit-in-a-box/internal/models"
)

type Client struct {
	HTTPClient *http.Client
	UserAgent  string
	Limit      int
}

func (c Client) Assess(ctx context.Context, deps []models.Dependency) []models.MaintainerHealth {
	limit := c.Limit
	if limit <= 0 {
		limit = 40
	}
	if len(deps) < limit {
		limit = len(deps)
	}

	results := make([]models.MaintainerHealth, 0, limit)
	for _, dep := range deps[:limit] {
		health := c.assessOne(ctx, dep)
		results = append(results, health)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score < results[j].Score
	})

	return results
}

func (c Client) assessOne(ctx context.Context, dep models.Dependency) models.MaintainerHealth {
	base := models.MaintainerHealth{
		PackageName: dep.Name,
		Ecosystem:   dep.Ecosystem,
		Score:       70,
		Source:      "registry",
	}

	var repo string
	var maintainerCount int
	var err error
	switch dep.Ecosystem {
	case "npm":
		repo, maintainerCount, err = c.npm(ctx, dep.Name)
	case "python":
		repo, maintainerCount, err = c.pypi(ctx, dep.Name)
	case "go":
		repo = repoFromGoModule(dep.Name)
	default:
		err = fmt.Errorf("unsupported ecosystem")
	}

	if maintainerCount > 0 {
		base.MaintainerCount = maintainerCount
	}
	if err != nil && repo == "" {
		base.Error = err.Error()
		base.Score = 45
		base.Signals = append(base.Signals, "registry metadata unavailable")
		return base
	}

	base.Repository = repo
	if repo == "" {
		base.Score = 50
		base.Signals = append(base.Signals, "repository URL unavailable")
		return base
	}

	gh, err := c.github(ctx, repo)
	if err != nil {
		base.Error = err.Error()
		base.Signals = append(base.Signals, "GitHub metadata unavailable")
		return base
	}

	base.Source = "github"
	base.LastCommit = gh.LastCommit
	base.BusFactor = gh.BusFactor
	if base.MaintainerCount == 0 {
		base.MaintainerCount = gh.Contributors
	}
	base.Score = healthScore(base.LastCommit, base.BusFactor, base.MaintainerCount, gh.Archived)
	base.Signals = healthSignals(base.LastCommit, base.BusFactor, base.MaintainerCount, gh.Archived)
	return base
}

func (c Client) npm(ctx context.Context, name string) (string, int, error) {
	var payload struct {
		Repository struct {
			URL string `json:"url"`
		} `json:"repository"`
		Maintainers []any `json:"maintainers"`
	}
	if err := c.getJSON(ctx, "https://registry.npmjs.org/"+url.PathEscape(name), &payload); err != nil {
		return "", 0, err
	}
	return normalizeGitHubURL(payload.Repository.URL), len(payload.Maintainers), nil
}

func (c Client) pypi(ctx context.Context, name string) (string, int, error) {
	var payload struct {
		Info struct {
			HomePage    string            `json:"home_page"`
			Author      string            `json:"author"`
			ProjectURLs map[string]string `json:"project_urls"`
		} `json:"info"`
	}
	if err := c.getJSON(ctx, "https://pypi.org/pypi/"+url.PathEscape(name)+"/json", &payload); err != nil {
		return "", 0, err
	}

	repo := ""
	for _, key := range []string{"Source", "Source Code", "Repository", "Homepage", "Home"} {
		if payload.Info.ProjectURLs != nil && payload.Info.ProjectURLs[key] != "" {
			repo = payload.Info.ProjectURLs[key]
			break
		}
	}
	if repo == "" {
		repo = payload.Info.HomePage
	}
	count := 0
	if payload.Info.Author != "" {
		count = 1
	}
	return normalizeGitHubURL(repo), count, nil
}

type githubInfo struct {
	LastCommit   *time.Time
	BusFactor    int
	Contributors int
	Archived     bool
}

func (c Client) github(ctx context.Context, repoURL string) (githubInfo, error) {
	owner, repo, ok := githubOwnerRepo(repoURL)
	if !ok {
		return githubInfo{}, fmt.Errorf("not a GitHub repository: %s", repoURL)
	}

	var repoPayload struct {
		Archived bool `json:"archived"`
	}
	if err := c.getJSON(ctx, "https://api.github.com/repos/"+owner+"/"+repo, &repoPayload); err != nil {
		return githubInfo{}, err
	}

	var commits []struct {
		Commit struct {
			Committer struct {
				Date time.Time `json:"date"`
			} `json:"committer"`
		} `json:"commit"`
	}
	_ = c.getJSON(ctx, "https://api.github.com/repos/"+owner+"/"+repo+"/commits?per_page=1", &commits)

	var contributors []struct {
		Login         string `json:"login"`
		Contributions int    `json:"contributions"`
	}
	_ = c.getJSON(ctx, "https://api.github.com/repos/"+owner+"/"+repo+"/contributors?per_page=25", &contributors)

	var lastCommit *time.Time
	if len(commits) > 0 {
		value := commits[0].Commit.Committer.Date
		lastCommit = &value
	}

	return githubInfo{
		LastCommit:   lastCommit,
		BusFactor:    busFactor(contributors),
		Contributors: len(contributors),
		Archived:     repoPayload.Archived,
	}, nil
}

func (c Client) getJSON(ctx context.Context, endpoint string, target any) error {
	client := c.HTTPClient
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	userAgent := c.UserAgent
	if userAgent == "" {
		userAgent = "audit-in-a-box"
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			return
		}
	}()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("%s returned %s", endpoint, resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func repoFromGoModule(module string) string {
	parts := strings.Split(module, "/")
	if len(parts) >= 3 && parts[0] == "github.com" {
		return "https://github.com/" + parts[1] + "/" + parts[2]
	}
	return ""
}

func normalizeGitHubURL(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "git+")
	value = strings.TrimSuffix(value, ".git")
	value = strings.ReplaceAll(value, "git://github.com/", "https://github.com/")
	value = strings.ReplaceAll(value, "git@github.com:", "https://github.com/")
	if strings.Contains(value, "github.com/") {
		if strings.HasPrefix(value, "http://") {
			value = strings.Replace(value, "http://", "https://", 1)
		}
		return value
	}
	return ""
}

func githubOwnerRepo(repoURL string) (string, string, bool) {
	parsed, err := url.Parse(repoURL)
	if err != nil || !strings.EqualFold(parsed.Host, "github.com") {
		return "", "", false
	}
	parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", false
	}
	return parts[0], strings.TrimSuffix(parts[1], ".git"), true
}

func busFactor(contributors []struct {
	Login         string `json:"login"`
	Contributions int    `json:"contributions"`
}) int {
	if len(contributors) == 0 {
		return 0
	}
	total := 0
	for _, contributor := range contributors {
		total += contributor.Contributions
	}
	if total == 0 {
		return len(contributors)
	}
	covered := 0
	for i, contributor := range contributors {
		covered += contributor.Contributions
		if float64(covered)/float64(total) >= 0.5 {
			return i + 1
		}
	}
	return len(contributors)
}

func healthScore(lastCommit *time.Time, busFactor int, maintainerCount int, archived bool) int {
	score := 85
	if archived {
		score -= 45
	}
	if lastCommit == nil {
		score -= 15
	} else {
		age := time.Since(*lastCommit)
		if age > 730*24*time.Hour {
			score -= 35
		} else if age > 365*24*time.Hour {
			score -= 20
		} else if age > 180*24*time.Hour {
			score -= 10
		}
	}
	if busFactor <= 1 {
		score -= 25
	} else if busFactor == 2 {
		score -= 10
	}
	if maintainerCount <= 1 {
		score -= 10
	}
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func healthSignals(lastCommit *time.Time, busFactor int, maintainerCount int, archived bool) []string {
	var signals []string
	if archived {
		signals = append(signals, "repository archived")
	}
	if lastCommit == nil {
		signals = append(signals, "last commit unknown")
	} else if time.Since(*lastCommit) > 365*24*time.Hour {
		signals = append(signals, "last commit older than one year")
	}
	if busFactor <= 1 {
		signals = append(signals, "bus factor appears to be 1")
	}
	if maintainerCount <= 1 {
		signals = append(signals, "single maintainer signal")
	}
	if len(signals) == 0 {
		signals = append(signals, "maintainer signals look healthy")
	}
	return signals
}
