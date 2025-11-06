import { spawnSync } from "node:child_process";
import { mkdirSync, writeFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import SwaggerParser from "@apidevtools/swagger-parser";

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const repoRoot = resolve(__dirname, "../../");
const specPath = resolve(__dirname, "openapi.yaml");

async function main() {
  console.log("Validating OpenAPI spec...");
  await SwaggerParser.validate(specPath);
  console.log("Spec is valid.");

  generateGo();
  generateTypescript();
}

function run(cmd: string, args: string[], cwd = repoRoot) {
  const result = spawnSync(cmd, args, {
    cwd,
    stdio: "inherit",
    env: process.env,
  });

  if (result.status !== 0) {
    throw new Error(`Command failed: ${cmd} ${args.join(" ")}`);
  }
}

function generateGo() {
  const outputDir = resolve(repoRoot, "apps/api/openapi");
  mkdirSync(outputDir, { recursive: true });
  console.log("Generating Go types...");
  run("pnpm", ["exec", "oapi-codegen", "-generate", "types", "-package", "openapi", "-o", join(outputDir, "types.gen.go"), specPath]);
}

function generateTypescript() {
  const outputDir = resolve(repoRoot, "apps/web/lib/api");
  mkdirSync(outputDir, { recursive: true });
  console.log("Generating TypeScript types...");
  run("pnpm", ["exec", "openapi-typescript", specPath, "--output", join(outputDir, "types.gen.ts")], repoRoot);

  const clientPath = join(outputDir, "client.ts");
  const relativeTypesPath = "./types.gen";
  const clientSource = `import createClient from "openapi-fetch";
import type { paths } from "${relativeTypesPath}";

const baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8080";

export const apiClient = createClient<paths>({ baseUrl });

export type ApiClient = typeof apiClient;
`;
  writeFileSync(clientPath, clientSource, "utf-8");
  console.log(`Wrote client stub to ${clientPath}`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
