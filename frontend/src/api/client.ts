import { auditReportSchema, scrapeCandidateSchema, ToolStatus } from "./schema";

const DEFAULT_BACKEND = "http://localhost:25342";
const STORAGE_KEY = "audit-in-a-box.backend-url";

export function getStoredBackendURL() {
  return localStorage.getItem(STORAGE_KEY) || DEFAULT_BACKEND;
}

export function setStoredBackendURL(value: string) {
  localStorage.setItem(STORAGE_KEY, trimSlash(value));
}

export function trimSlash(value: string) {
  return value.trim().replace(/\/+$/, "");
}

export async function fetchTools(
  baseURL: string,
): Promise<Record<string, ToolStatus>> {
  const response = await fetch(`${trimSlash(baseURL)}/api/v1/tools`);
  if (!response.ok) {
    throw new Error(`Backend returned ${response.status}`);
  }
  return (await response.json()) as Record<string, ToolStatus>;
}

export async function scrapeHTML(baseURL: string, html: string) {
  const response = await fetch(`${trimSlash(baseURL)}/api/v1/scrape`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ html }),
  });
  if (!response.ok) {
    throw new Error(await errorMessage(response));
  }
  const payload = (await response.json()) as { candidates: unknown[] };
  return payload.candidates.map((candidate) =>
    scrapeCandidateSchema.parse(candidate),
  );
}

export async function auditManifest(input: {
  baseURL: string;
  fileName: string;
  content: string;
  pastedHTML?: string;
  ecosystem?: string;
}) {
  const response = await fetch(`${trimSlash(input.baseURL)}/api/v1/audits`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      file_name: input.fileName,
      content: input.content,
      pasted_html: input.pastedHTML,
      ecosystem: input.ecosystem || undefined,
    }),
  });
  if (!response.ok) {
    throw new Error(await errorMessage(response));
  }
  return auditReportSchema.parse(await response.json());
}

async function errorMessage(response: Response) {
  try {
    const payload = (await response.json()) as { error?: { message?: string } };
    return payload.error?.message || `Backend returned ${response.status}`;
  } catch {
    return `Backend returned ${response.status}`;
  }
}
