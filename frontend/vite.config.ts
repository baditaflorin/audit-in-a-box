import react from "@vitejs/plugin-react";
import { execSync } from "node:child_process";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { defineConfig } from "vitest/config";

const pkg = JSON.parse(
  readFileSync(new URL("./package.json", import.meta.url), "utf8"),
) as {
  version: string;
};

function gitCommit() {
  try {
    return execSync("git rev-parse --short=12 HEAD", {
      stdio: ["ignore", "pipe", "ignore"],
    })
      .toString()
      .trim();
  } catch {
    return "unknown";
  }
}

export default defineConfig({
  base: "/audit-in-a-box/",
  plugins: [react()],
  build: {
    outDir: resolve(__dirname, "../docs"),
    emptyOutDir: false,
    sourcemap: false,
    rollupOptions: {
      output: {
        assetFileNames: "assets/[name]-[hash][extname]",
        chunkFileNames: "assets/[name]-[hash].js",
        entryFileNames: "assets/[name]-[hash].js",
      },
    },
  },
  define: {
    __APP_VERSION__: JSON.stringify(
      process.env.VITE_APP_VERSION || pkg.version,
    ),
    __GIT_COMMIT__: JSON.stringify(process.env.VITE_GIT_COMMIT || gitCommit()),
  },
  test: {
    environment: "node",
    globals: true,
  },
});
