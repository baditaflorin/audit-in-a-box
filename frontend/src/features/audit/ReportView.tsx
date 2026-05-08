import {
  AlertTriangle,
  BadgeCheck,
  CircleGauge,
  Clock3,
  Database,
  ExternalLink,
  FileCode2,
  GitCommit,
  HeartPulse,
  ShieldAlert,
} from "lucide-react";
import { AuditReport } from "../../api/schema";

type ReportViewProps = {
  report: AuditReport;
};

export function ReportView({ report }: ReportViewProps) {
  const severity = report.risk.grade;

  return (
    <section className="space-y-5" aria-label="Audit report">
      <div className="grid gap-4 lg:grid-cols-[1fr_280px]">
        <div className="rounded-lg border border-ink/15 bg-white p-5 shadow-soft">
          <div className="flex flex-wrap items-start justify-between gap-4">
            <div>
              <p className="text-sm font-semibold uppercase tracking-wide text-moss">
                Plain-English risk summary
              </p>
              <h2 className="mt-1 text-2xl font-bold text-ink">
                {report.input.file_name}
              </h2>
            </div>
            <span className={badgeClass(severity)}>{severity}</span>
          </div>
          <p className="mt-4 text-base leading-7 text-ink/80">
            {report.summary}
          </p>
          {report.warnings.length > 0 ? (
            <div className="mt-4 rounded-lg border border-gold/40 bg-gold/10 p-3 text-sm text-ink">
              <div className="flex items-center gap-2 font-semibold">
                <AlertTriangle size={17} />
                Scanner notes
              </div>
              <ul className="mt-2 list-disc space-y-1 pl-5">
                {report.warnings.slice(0, 5).map((warning) => (
                  <li key={warning}>{warning}</li>
                ))}
              </ul>
            </div>
          ) : null}
        </div>

        <div className="rounded-lg border border-ink/15 bg-ink p-5 text-white shadow-soft">
          <div className="flex items-center gap-2 text-sm font-semibold text-white/70">
            <CircleGauge size={18} />
            Risk score
          </div>
          <div className="mt-4 flex items-end gap-2">
            <span className="text-6xl font-black leading-none">
              {report.risk.score}
            </span>
            <span className="pb-2 text-lg text-white/70">/100</span>
          </div>
          <dl className="mt-5 grid grid-cols-2 gap-3 text-sm">
            <Metric label="Deps" value={report.dependencies.length} />
            <Metric label="Vulns" value={report.vulnerabilities.length} />
            <Metric label="Licenses" value={report.license_risks.length} />
            <Metric
              label="Maintainers"
              value={report.risk.counts.maintainer_risks ?? 0}
            />
          </dl>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <InfoCard
          icon={<FileCode2 size={20} />}
          label="SBOM packages"
          value={String(
            report.duckdb_rollup.dependency_count || report.dependencies.length,
          )}
          detail={
            report.dependencies
              .slice(0, 3)
              .map((item) => item.name)
              .join(", ") || "No packages detected"
          }
        />
        <InfoCard
          icon={<ShieldAlert size={20} />}
          label="Vulnerabilities"
          value={String(report.vulnerabilities.length)}
          detail={topSeverity(report)}
        />
        <InfoCard
          icon={<Database size={20} />}
          label="DuckDB"
          value={report.duckdb_rollup.used_duckdb ? "used" : "fallback"}
          detail={report.duckdb_rollup.diagnostic || "Rollup complete"}
        />
        <InfoCard
          icon={<Clock3 size={20} />}
          label="Elapsed"
          value={`${Math.max(1, Math.round(report.elapsed_millis / 1000))}s`}
          detail={`Generated ${new Date(report.generated_at).toLocaleString()}`}
        />
      </div>

      <div className="grid gap-4 xl:grid-cols-2">
        <TablePanel title="Vulnerabilities">
          {report.vulnerabilities.length === 0 ? (
            <EmptyState text="No vulnerabilities returned by Trivy or Grype." />
          ) : (
            <table className="w-full text-left text-sm">
              <thead className="text-xs uppercase text-ink/55">
                <tr>
                  <th className="pb-2">Severity</th>
                  <th className="pb-2">Package</th>
                  <th className="pb-2">ID</th>
                  <th className="pb-2">Fix</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-ink/10">
                {report.vulnerabilities.slice(0, 12).map((vuln) => (
                  <tr key={`${vuln.source}-${vuln.id}-${vuln.package_name}`}>
                    <td className="py-2">
                      <span className={severityClass(vuln.severity)}>
                        {vuln.severity || "UNKNOWN"}
                      </span>
                    </td>
                    <td className="py-2 font-medium">{vuln.package_name}</td>
                    <td className="py-2">
                      {vuln.primary_url ? (
                        <a
                          className="inline-flex items-center gap-1 text-moss underline"
                          href={vuln.primary_url}
                        >
                          {vuln.id}
                          <ExternalLink size={13} />
                        </a>
                      ) : (
                        vuln.id
                      )}
                    </td>
                    <td className="py-2 text-ink/70">
                      {vuln.fixed_version || "unknown"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </TablePanel>

        <TablePanel title="Maintainer health">
          {report.maintainer_health.length === 0 ? (
            <EmptyState text="No maintainer metadata was available." />
          ) : (
            <table className="w-full text-left text-sm">
              <thead className="text-xs uppercase text-ink/55">
                <tr>
                  <th className="pb-2">Package</th>
                  <th className="pb-2">Score</th>
                  <th className="pb-2">Bus</th>
                  <th className="pb-2">Last commit</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-ink/10">
                {report.maintainer_health.slice(0, 12).map((item) => (
                  <tr key={`${item.ecosystem}-${item.package_name}`}>
                    <td className="py-2">
                      <div className="font-medium">{item.package_name}</div>
                      <div className="text-xs text-ink/55">
                        {item.signals.join("; ")}
                      </div>
                    </td>
                    <td className="py-2">{item.score}</td>
                    <td className="py-2">{item.bus_factor || "n/a"}</td>
                    <td className="py-2 text-ink/70">
                      {item.last_commit
                        ? new Date(item.last_commit).toLocaleDateString()
                        : "unknown"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </TablePanel>
      </div>

      <TablePanel title="License risks">
        {report.license_risks.length === 0 ? (
          <EmptyState text="No license risks were found in the returned evidence." />
        ) : (
          <div className="grid gap-2 md:grid-cols-2">
            {report.license_risks.slice(0, 16).map((risk) => (
              <div
                className="rounded-lg border border-ink/10 bg-paper p-3 text-sm"
                key={`${risk.package_name}-${risk.license}-${risk.severity}`}
              >
                <div className="flex items-center justify-between gap-2">
                  <span className="font-semibold">{risk.package_name}</span>
                  <span className={severityClass(risk.severity)}>
                    {risk.severity}
                  </span>
                </div>
                <div className="mt-1 text-ink/70">{risk.license}</div>
                <div className="mt-1 text-xs text-ink/55">{risk.reason}</div>
              </div>
            ))}
          </div>
        )}
      </TablePanel>

      <div className="rounded-lg border border-ink/15 bg-white p-4 text-sm text-ink/70">
        <div className="flex flex-wrap items-center gap-3">
          <span className="inline-flex items-center gap-2">
            <BadgeCheck size={16} />
            Backend version {report.version.version}
          </span>
          <span className="inline-flex items-center gap-2">
            <GitCommit size={16} />
            {report.version.commit}
          </span>
          <span className="inline-flex items-center gap-2">
            <HeartPulse size={16} />
            {report.tool_status.trivy?.used ? "Trivy used" : "Trivy not used"}
          </span>
        </div>
      </div>
    </section>
  );
}

function Metric({ label, value }: { label: string; value: number }) {
  return (
    <div>
      <dt className="text-white/55">{label}</dt>
      <dd className="text-lg font-bold">{value}</dd>
    </div>
  );
}

function InfoCard({
  icon,
  label,
  value,
  detail,
}: {
  icon: React.ReactNode;
  label: string;
  value: string;
  detail: string;
}) {
  return (
    <div className="rounded-lg border border-ink/15 bg-white p-4 shadow-soft">
      <div className="flex items-center gap-2 text-sm font-semibold text-ink/60">
        {icon}
        {label}
      </div>
      <div className="mt-3 text-3xl font-black text-ink">{value}</div>
      <div className="mt-2 line-clamp-2 text-sm leading-5 text-ink/60">
        {detail}
      </div>
    </div>
  );
}

function TablePanel({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <section className="rounded-lg border border-ink/15 bg-white p-4 shadow-soft">
      <h3 className="mb-4 text-lg font-bold text-ink">{title}</h3>
      <div className="overflow-x-auto">{children}</div>
    </section>
  );
}

function EmptyState({ text }: { text: string }) {
  return (
    <div className="rounded-lg border border-dashed border-ink/20 bg-paper p-4 text-sm text-ink/60">
      {text}
    </div>
  );
}

function topSeverity(report: AuditReport) {
  const counts = report.duckdb_rollup.severity_counts;
  const order = ["CRITICAL", "HIGH", "MEDIUM", "LOW", "UNKNOWN"];
  const found = order.find((severity) => counts?.[severity]);
  return found
    ? `${counts[found]} ${found.toLowerCase()} finding(s)`
    : "No severity counts";
}

function badgeClass(grade: string) {
  const base = "rounded px-3 py-1 text-sm font-black uppercase";
  if (grade === "critical" || grade === "high")
    return `${base} bg-coral text-white`;
  if (grade === "medium") return `${base} bg-gold text-ink`;
  return `${base} bg-moss text-white`;
}

function severityClass(severity: string) {
  const base = "rounded px-2 py-1 text-xs font-bold uppercase";
  const normalized = severity.toUpperCase();
  if (normalized === "CRITICAL" || normalized === "HIGH")
    return `${base} bg-coral/15 text-coral`;
  if (normalized === "MEDIUM") return `${base} bg-gold/20 text-ink`;
  return `${base} bg-moss/10 text-moss`;
}
