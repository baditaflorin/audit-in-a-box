import { z } from "zod";

export const toolStatusSchema = z.object({
  name: z.string(),
  found: z.boolean(),
  path: z.string().optional(),
  version: z.string().optional(),
  error: z.string().optional(),
  required: z.boolean(),
  used: z.boolean(),
  available: z.boolean(),
});

export const dependencySchema = z.object({
  name: z.string(),
  version: z.string(),
  ecosystem: z.string(),
  scope: z.string(),
  licenses: z.array(z.string()).default([]),
  package_url: z.string().optional(),
  source: z.string(),
  confidence: z.number().optional(),
  reasons: z.array(z.string()).optional(),
});

export const vulnerabilitySchema = z.object({
  id: z.string(),
  package_name: z.string(),
  installed_version: z.string(),
  fixed_version: z.string().optional(),
  severity: z.string(),
  description: z.string().optional(),
  primary_url: z.string().optional(),
  source: z.string(),
  confidence: z.number().optional(),
});

export const licenseRiskSchema = z.object({
  package_name: z.string(),
  license: z.string(),
  severity: z.string(),
  reason: z.string(),
  confidence: z.number().optional(),
});

export const maintainerHealthSchema = z.object({
  package_name: z.string(),
  ecosystem: z.string(),
  repository: z.string().optional(),
  last_commit: z.string().datetime().nullable().optional(),
  bus_factor: z.number(),
  maintainer_count: z.number(),
  score: z.number(),
  signals: z.array(z.string()),
  source: z.string(),
  error: z.string().optional(),
});

export const auditReportSchema = z.object({
  id: z.string(),
  generated_at: z.string(),
  input: z.object({
    file_name: z.string(),
    ecosystem: z.string(),
    bytes: z.number(),
    normalized_bytes: z.number().optional(),
    kind: z.string().optional(),
    parser: z.string().optional(),
    confidence: z.number().optional(),
    reasons: z.array(z.string()).optional(),
    source_hash: z.string().optional(),
  }),
  provenance: z
    .object({
      schema_version: z.string(),
      app_version: z.object({
        version: z.string(),
        commit: z.string(),
        date: z.string(),
      }),
      source_hash: z.string(),
      input: z.unknown(),
      parameters: z.record(z.string(), z.unknown()),
    })
    .optional(),
  tool_status: z.record(z.string(), toolStatusSchema),
  dependencies: z.array(dependencySchema),
  vulnerabilities: z.array(vulnerabilitySchema),
  license_risks: z.array(licenseRiskSchema),
  anomalies: z
    .array(
      z.object({
        code: z.string(),
        severity: z.string(),
        message: z.string(),
        why: z.string(),
        next_step: z.string(),
        confidence: z.number(),
      }),
    )
    .default([]),
  maintainer_health: z.array(maintainerHealthSchema),
  duckdb_rollup: z.object({
    used_duckdb: z.boolean(),
    dependency_count: z.number(),
    vulnerability_count: z.number(),
    license_risk_count: z.number(),
    severity_counts: z.record(z.string(), z.number()),
    generated_artifact: z.string().optional(),
    diagnostic: z.string().optional(),
  }),
  risk: z.object({
    score: z.number(),
    grade: z.string(),
    confidence: z.number().optional(),
    counts: z.record(z.string(), z.number()),
    factors: z.array(z.string()),
  }),
  summary: z.string(),
  warnings: z.array(z.string()),
  version: z.object({
    version: z.string(),
    commit: z.string(),
    date: z.string(),
  }),
  elapsed_millis: z.number(),
});

export type ToolStatus = z.infer<typeof toolStatusSchema>;
export type AuditReport = z.infer<typeof auditReportSchema>;

export const apiErrorSchema = z.object({
  error: z.object({
    code: z.string().optional(),
    message: z.string().optional(),
    why: z.string().optional(),
    next_step: z.string().optional(),
    recoverable: z.boolean().optional(),
  }),
});

export const scrapeCandidateSchema = z.object({
  file_name: z.string(),
  ecosystem: z.string(),
  content: z.string(),
  score: z.number(),
});

export type ScrapeCandidate = z.infer<typeof scrapeCandidateSchema>;
