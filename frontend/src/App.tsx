import { useMutation, useQuery } from "@tanstack/react-query";
import {
  Box,
  CheckCircle2,
  Code2,
  Database,
  FileUp,
  Github,
  HeartHandshake,
  Loader2,
  Play,
  Server,
  ShieldCheck,
  Sparkles,
} from "lucide-react";
import { FormEvent, useMemo, useState } from "react";
import {
  auditManifest,
  fetchTools,
  getStoredBackendURL,
  scrapeHTML,
  setStoredBackendURL,
} from "./api/client";
import { AuditReport } from "./api/schema";
import { ReportView } from "./features/audit/ReportView";
import {
  sampleGoMod,
  samplePackageJSON,
  sampleRequirements,
} from "./features/audit/samples";

const reportStorageKey = "audit-in-a-box.last-report";

export function App() {
  const [backendURL, setBackendURL] = useState(getStoredBackendURL);
  const [fileName, setFileName] = useState("package.json");
  const [content, setContent] = useState(samplePackageJSON);
  const [htmlPaste, setHtmlPaste] = useState("");
  const [lastReport, setLastReport] = useState<AuditReport | null>(() => {
    const stored = localStorage.getItem(reportStorageKey);
    if (!stored) return null;
    try {
      return JSON.parse(stored) as AuditReport;
    } catch {
      return null;
    }
  });
  const [error, setError] = useState<string | null>(null);

  const tools = useQuery({
    queryKey: ["tools", backendURL],
    queryFn: () => fetchTools(backendURL),
  });

  const audit = useMutation({
    mutationFn: auditManifest,
    onSuccess: (report) => {
      setLastReport(report);
      localStorage.setItem(reportStorageKey, JSON.stringify(report));
      setError(null);
    },
    onError: (err) => {
      setError(err instanceof Error ? err.message : "Audit failed");
    },
  });

  const scrape = useMutation({
    mutationFn: () => scrapeHTML(backendURL, htmlPaste),
    onSuccess: (candidates) => {
      if (candidates.length === 0) {
        setError("No manifest-like content was found in that HTML paste.");
        return;
      }
      setFileName(candidates[0].file_name);
      setContent(candidates[0].content);
      setError(null);
    },
    onError: (err) => {
      setError(err instanceof Error ? err.message : "HTML scrape failed");
    },
  });

  const toolItems = useMemo(
    () => Object.values(tools.data ?? {}),
    [tools.data],
  );

  function submit(event: FormEvent) {
    event.preventDefault();
    setStoredBackendURL(backendURL);
    audit.mutate({ baseURL: backendURL, fileName, content });
  }

  function loadSample(name: string) {
    if (name === "go.mod") {
      setFileName("go.mod");
      setContent(sampleGoMod);
      return;
    }
    if (name === "requirements.txt") {
      setFileName("requirements.txt");
      setContent(sampleRequirements);
      return;
    }
    setFileName("package.json");
    setContent(samplePackageJSON);
  }

  async function onFile(file: File | undefined) {
    if (!file) return;
    setFileName(file.name);
    setContent(await file.text());
  }

  return (
    <main className="min-h-screen bg-paper text-ink">
      <section className="border-b border-ink/10 bg-white">
        <div className="mx-auto grid min-h-[96vh] w-full max-w-7xl gap-8 px-4 py-6 sm:px-6 lg:grid-cols-[minmax(0,1fr)_420px] lg:px-8">
          <div className="flex flex-col justify-between gap-10">
            <header className="flex flex-wrap items-center justify-between gap-3">
              <div className="flex items-center gap-3">
                <div className="grid h-11 w-11 place-items-center rounded-lg bg-moss text-white">
                  <ShieldCheck size={24} aria-hidden />
                </div>
                <div>
                  <div className="text-xl font-black">audit-in-a-box</div>
                  <div className="text-sm text-ink/60">
                    OSS dependency risk, local-first
                  </div>
                </div>
              </div>
              <nav className="flex flex-wrap gap-2 text-sm font-semibold">
                <a
                  className="inline-flex items-center gap-2 rounded border border-ink/15 px-3 py-2 text-ink no-underline hover:border-moss"
                  href="https://github.com/baditaflorin/audit-in-a-box"
                >
                  <Github size={16} />
                  Star on GitHub
                </a>
                <a
                  className="inline-flex items-center gap-2 rounded border border-ink/15 px-3 py-2 text-ink no-underline hover:border-moss"
                  href="https://www.paypal.com/paypalme/florinbadita"
                >
                  <HeartHandshake size={16} />
                  PayPal
                </a>
              </nav>
            </header>

            <div>
              <p className="mb-3 inline-flex items-center gap-2 rounded bg-sky/10 px-3 py-1 text-sm font-bold text-sky">
                <Sparkles size={16} />
                Trivy + Syft + Grype + DuckDB + local LLM
              </p>
              <h1 className="max-w-4xl text-5xl font-black leading-[0.98] tracking-normal sm:text-7xl">
                Drop a manifest. Get an audit you can actually read.
              </h1>
              <p className="mt-6 max-w-2xl text-lg leading-8 text-ink/70">
                Upload or paste `package.json`, `go.mod`, or `requirements.txt`
                and run a local Docker analyzer for SBOM, vulnerabilities,
                license risks, maintainer health, and a plain-English summary.
              </p>
              <div className="mt-6 grid gap-3 sm:grid-cols-3">
                <Signal icon={<Box size={18} />} label="SBOM" />
                <Signal icon={<Database size={18} />} label="DuckDB rollups" />
                <Signal icon={<Server size={18} />} label="Local backend" />
              </div>
            </div>

            <footer className="flex flex-wrap items-center gap-3 text-xs text-ink/55">
              <span>Version {__APP_VERSION__}</span>
              <span>Commit {__GIT_COMMIT__}</span>
              <span>
                Repository https://github.com/baditaflorin/audit-in-a-box
              </span>
            </footer>
          </div>

          <form
            className="self-center rounded-lg border border-ink/15 bg-paper p-4 shadow-soft"
            onSubmit={submit}
            aria-label="Run dependency audit"
          >
            <label className="text-sm font-bold" htmlFor="backend-url">
              Backend URL
            </label>
            <div className="mt-2 flex gap-2">
              <input
                className="min-w-0 flex-1 rounded border border-ink/20 bg-white px-3 py-2 text-sm"
                id="backend-url"
                value={backendURL}
                onChange={(event) => setBackendURL(event.target.value)}
              />
              <button
                className="grid h-10 w-10 place-items-center rounded bg-ink text-white"
                type="button"
                title="Refresh backend tool status"
                onClick={() => void tools.refetch()}
              >
                {tools.isFetching ? (
                  <Loader2 className="animate-spin" size={18} />
                ) : (
                  <Server size={18} />
                )}
              </button>
            </div>

            <div className="mt-3 grid grid-cols-2 gap-2 sm:grid-cols-4">
              {["trivy", "syft", "grype", "duckdb"].map((name) => {
                const item = toolItems.find((tool) => tool.name === name);
                return (
                  <div
                    className="rounded border border-ink/10 bg-white px-2 py-2 text-xs"
                    key={name}
                  >
                    <div className="flex items-center justify-between gap-2">
                      <span className="font-bold">{name}</span>
                      {item?.available ? (
                        <CheckCircle2 className="text-moss" size={15} />
                      ) : (
                        <span className="h-2 w-2 rounded-full bg-coral" />
                      )}
                    </div>
                  </div>
                );
              })}
            </div>

            <div className="mt-4 flex flex-wrap gap-2">
              {["package.json", "go.mod", "requirements.txt"].map((sample) => (
                <button
                  className="rounded border border-ink/15 bg-white px-3 py-2 text-sm font-semibold hover:border-moss"
                  key={sample}
                  type="button"
                  onClick={() => loadSample(sample)}
                >
                  {sample}
                </button>
              ))}
            </div>

            <label className="mt-4 block text-sm font-bold" htmlFor="file-name">
              File name
            </label>
            <input
              className="mt-2 w-full rounded border border-ink/20 bg-white px-3 py-2 text-sm"
              id="file-name"
              value={fileName}
              onChange={(event) => setFileName(event.target.value)}
            />

            <label
              className="mt-4 flex cursor-pointer items-center justify-center gap-2 rounded border border-dashed border-ink/30 bg-white px-3 py-3 text-sm font-bold hover:border-moss"
              htmlFor="manifest-file"
            >
              <FileUp size={18} />
              Upload manifest
            </label>
            <input
              className="sr-only"
              id="manifest-file"
              type="file"
              accept=".json,.mod,.txt"
              onChange={(event) => void onFile(event.target.files?.[0])}
            />

            <label
              className="mt-4 block text-sm font-bold"
              htmlFor="manifest-content"
            >
              Manifest content
            </label>
            <textarea
              className="mt-2 h-64 w-full resize-y rounded border border-ink/20 bg-white p-3 font-mono text-xs leading-5"
              id="manifest-content"
              value={content}
              onChange={(event) => setContent(event.target.value)}
              spellCheck={false}
            />

            <div className="mt-4 rounded-lg border border-ink/10 bg-white p-3">
              <label className="block text-sm font-bold" htmlFor="html-paste">
                Paste HTML scraper
              </label>
              <textarea
                className="mt-2 h-24 w-full resize-y rounded border border-ink/20 bg-paper p-3 font-mono text-xs leading-5"
                id="html-paste"
                value={htmlPaste}
                onChange={(event) => setHtmlPaste(event.target.value)}
                spellCheck={false}
              />
              <button
                className="mt-2 inline-flex items-center gap-2 rounded bg-ink px-3 py-2 text-sm font-bold text-white disabled:opacity-60"
                disabled={scrape.isPending || htmlPaste.trim() === ""}
                type="button"
                onClick={() => scrape.mutate()}
              >
                {scrape.isPending ? (
                  <Loader2 className="animate-spin" size={16} />
                ) : (
                  <Code2 size={16} />
                )}
                Extract manifest
              </button>
            </div>

            {error ? (
              <div className="mt-3 rounded bg-coral/10 p-3 text-sm font-semibold text-coral">
                {error}
              </div>
            ) : null}

            <button
              className="mt-4 inline-flex w-full items-center justify-center gap-2 rounded bg-moss px-4 py-3 text-base font-black text-white disabled:cursor-not-allowed disabled:opacity-60"
              disabled={audit.isPending}
              type="submit"
            >
              {audit.isPending ? (
                <Loader2 className="animate-spin" size={18} />
              ) : (
                <Play size={18} />
              )}
              Run audit
            </button>
          </form>
        </div>
      </section>

      <section className="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
        {lastReport ? (
          <ReportView report={lastReport} />
        ) : (
          <div className="rounded-lg border border-dashed border-ink/20 bg-white p-8 text-center">
            <Code2 className="mx-auto text-moss" size={34} />
            <h2 className="mt-3 text-2xl font-black">No report yet</h2>
            <p className="mx-auto mt-2 max-w-xl text-ink/65">
              Start the Docker backend, run a sample audit, and the SBOM,
              vulnerabilities, licenses, maintainer health, version, and commit
              metadata will appear here.
            </p>
          </div>
        )}
      </section>
    </main>
  );
}

function Signal({ icon, label }: { icon: React.ReactNode; label: string }) {
  return (
    <div className="flex items-center gap-2 rounded-lg border border-ink/10 bg-white px-3 py-3 text-sm font-bold shadow-soft">
      {icon}
      {label}
    </div>
  );
}
