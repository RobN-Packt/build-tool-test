import { execSync } from 'node:child_process';
import { promises as fs } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const repoRoot = path.resolve(__dirname, '..');
const openapiDir = path.join(repoRoot, 'apps/api/openapi');
const openapiSpecPath = path.join(openapiDir, 'openapi.yaml');
const tsTypesPath = path.join(repoRoot, 'apps/web/lib/api/types.ts');
const tsClientPath = path.join(repoRoot, 'apps/web/lib/api/client.ts');
const goModelsPath = path.join(openapiDir, 'gen.models.go');

async function generateTypescript() {
  await fs.mkdir(path.dirname(tsTypesPath), { recursive: true });
  execSync(`pnpm exec openapi-typescript ${openapiSpecPath} --output ${tsTypesPath}`, {
    cwd: repoRoot,
    stdio: 'inherit'
  });
  const clientSource = `import createClient from 'openapi-fetch';\nimport type { paths } from './types';\n\nconst baseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080';\n\nexport const apiClient = createClient<paths>({ baseUrl });\n`;
  await fs.writeFile(tsClientPath, clientSource, 'utf8');
  console.log('Generated TypeScript types and client.');
}

function generateGoModels() {
  const goPath = execSync('go env GOPATH', { encoding: 'utf8' }).trim();
  const binary = path.join(goPath, 'bin', 'oapi-codegen');
  const cmd = `${binary} -package openapi -generate types -o ${goModelsPath} ${openapiSpecPath}`;
  execSync(cmd, { stdio: 'inherit' });
  console.log('Generated Go models.');
}

async function main() {
  try {
    await generateTypescript();
    generateGoModels();
  } catch (err) {
    console.error(err);
    process.exitCode = 1;
  }
}

main();
