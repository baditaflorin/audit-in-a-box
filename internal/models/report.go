package models

import "time"

type AuditRequest struct {
	FileName   string `json:"file_name" validate:"omitempty,max=256"`
	Content    string `json:"content" validate:"omitempty"`
	PastedHTML string `json:"pasted_html" validate:"omitempty"`
	Ecosystem  string `json:"ecosystem" validate:"omitempty,oneof=npm go python unknown"`
}

type AuditReport struct {
	ID                string                  `json:"id"`
	GeneratedAt       time.Time               `json:"generated_at"`
	Input             AuditInput              `json:"input"`
	ToolStatus        map[string]ToolStatus   `json:"tool_status"`
	Dependencies      []Dependency            `json:"dependencies"`
	Vulnerabilities   []Vulnerability         `json:"vulnerabilities"`
	LicenseRisks      []LicenseRisk           `json:"license_risks"`
	MaintainerHealth  []MaintainerHealth      `json:"maintainer_health"`
	DuckDBRollup      DuckDBRollup            `json:"duckdb_rollup"`
	Risk              RiskScore               `json:"risk"`
	Summary           string                  `json:"summary"`
	Warnings          []string                `json:"warnings"`
	ScrapedCandidates []ScrapedManifest       `json:"scraped_candidates,omitempty"`
	Version           BuildVersion            `json:"version"`
	ElapsedMillis     int64                   `json:"elapsed_millis"`
	RawEvidence       map[string]EvidenceInfo `json:"raw_evidence"`
}

type BuildVersion struct {
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Date    string `json:"date"`
}

type AuditInput struct {
	FileName  string `json:"file_name"`
	Ecosystem string `json:"ecosystem"`
	Bytes     int    `json:"bytes"`
}

type ToolStatus struct {
	Name      string `json:"name"`
	Found     bool   `json:"found"`
	Path      string `json:"path,omitempty"`
	Version   string `json:"version,omitempty"`
	Error     string `json:"error,omitempty"`
	Required  bool   `json:"required"`
	Used      bool   `json:"used"`
	Available bool   `json:"available"`
}

type Dependency struct {
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Ecosystem  string   `json:"ecosystem"`
	Scope      string   `json:"scope"`
	Licenses   []string `json:"licenses"`
	PackageURL string   `json:"package_url,omitempty"`
	Source     string   `json:"source"`
}

type Vulnerability struct {
	ID               string `json:"id"`
	PackageName      string `json:"package_name"`
	InstalledVersion string `json:"installed_version"`
	FixedVersion     string `json:"fixed_version,omitempty"`
	Severity         string `json:"severity"`
	Description      string `json:"description,omitempty"`
	PrimaryURL       string `json:"primary_url,omitempty"`
	Source           string `json:"source"`
}

type LicenseRisk struct {
	PackageName string `json:"package_name"`
	License     string `json:"license"`
	Severity    string `json:"severity"`
	Reason      string `json:"reason"`
}

type MaintainerHealth struct {
	PackageName     string     `json:"package_name"`
	Ecosystem       string     `json:"ecosystem"`
	Repository      string     `json:"repository,omitempty"`
	LastCommit      *time.Time `json:"last_commit,omitempty"`
	BusFactor       int        `json:"bus_factor"`
	MaintainerCount int        `json:"maintainer_count"`
	Score           int        `json:"score"`
	Signals         []string   `json:"signals"`
	Source          string     `json:"source"`
	Error           string     `json:"error,omitempty"`
}

type DuckDBRollup struct {
	UsedDuckDB         bool           `json:"used_duckdb"`
	DependencyCount    int            `json:"dependency_count"`
	VulnerabilityCount int            `json:"vulnerability_count"`
	LicenseRiskCount   int            `json:"license_risk_count"`
	SeverityCounts     map[string]int `json:"severity_counts"`
	GeneratedArtifact  string         `json:"generated_artifact,omitempty"`
	Diagnostic         string         `json:"diagnostic,omitempty"`
}

type RiskScore struct {
	Score   int            `json:"score"`
	Grade   string         `json:"grade"`
	Counts  map[string]int `json:"counts"`
	Factors []string       `json:"factors"`
}

type ScrapedManifest struct {
	FileName  string `json:"file_name"`
	Ecosystem string `json:"ecosystem"`
	Content   string `json:"content"`
	Score     int    `json:"score"`
}

type EvidenceInfo struct {
	Available bool   `json:"available"`
	Path      string `json:"path,omitempty"`
	Message   string `json:"message,omitempty"`
}
