import openapiTS from 'openapi-typescript';
import { mkdir, writeFile } from 'node:fs/promises';
import { spawn } from 'node:child_process';
import path from 'node:path';
import process from 'node:process';

const ROOT = path.resolve(process.cwd());
const SPEC_PATH = path.join(ROOT, 'apps/api/openapi/openapi.yaml');
const TS_OUTPUT = path.join(ROOT, 'apps/web/lib/api/types.ts');
const GO_OUTPUT = path.join(ROOT, 'apps/api/openapi/gen.models.go');

async function run() {
  console.log('Generating TypeScript types...');
  await mkdir(path.dirname(TS_OUTPUT), { recursive: true });
  const tsResult = await openapiTS(SPEC_PATH);
  await writeFile(TS_OUTPUT, tsResult, 'utf8');

  console.log('Generating Go models...');
  await mkdir(path.dirname(GO_OUTPUT), { recursive: true });
  await runCommand('go', [
    'run',
    'github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1',
    '-generate',
    'types',
    '-package',
    'openapi',
    '-o',
    GO_OUTPUT,
    SPEC_PATH
  ]);

  console.log('Code generation completed.');
}

function runCommand(command: string, args: string[]) {
  return new Promise<void>((resolve, reject) => {
    const child = spawn(command, args, { stdio: 'inherit' });
    child.on('close', (code) => {
      if (code !== 0) {
        reject(new Error(`${command} ${args.join(' ')} exited with code ${code}`));
        return;
      }
      resolve();
    });
    child.on('error', reject);
  });
}

run().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});
