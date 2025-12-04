import { spawn } from 'node:child_process';
import { cpSync, existsSync, rmSync, symlinkSync } from 'node:fs';
import { join } from 'node:path';
import { createRequire } from 'node:module';

const require = createRequire(import.meta.url);

async function runNextBuild() {
  await new Promise((resolve, reject) => {
    const child = spawn(process.execPath, [require.resolve('next/dist/bin/next'), 'build'], {
      stdio: 'inherit',
      env: process.env,
    });

    child.on('close', (code) => {
      if (code === 0) {
        resolve();
      } else {
        reject(new Error(`next build exited with code ${code}`));
      }
    });
    child.on('error', reject);
  });
}

function resolveNodeModulesRoot() {
  const runfilesDir = process.env.RUNFILES_DIR;
  const workspaceName = process.env.BAZEL_WORKSPACE || '_main';
  const runfileNodeModules = runfilesDir
    ? join(runfilesDir, workspaceName, 'apps/web/node_modules')
    : null;

  if (runfileNodeModules && existsSync(runfileNodeModules)) {
    return runfileNodeModules;
  }

  const fallback = join(process.cwd(), 'node_modules');
  if (existsSync(fallback)) {
    return fallback;
  }

  throw new Error('Unable to locate node_modules; ensure Bazel runfiles were provided.');
}

function copyNodeModulesIntoNext() {
  const projectRoot = process.cwd();
  const nextDir = join(projectRoot, '.next');
  const sourceNodeModules = resolveNodeModulesRoot();
  const targetNodeModules = join(nextDir, 'node_modules');

  if (existsSync(targetNodeModules)) {
    rmSync(targetNodeModules, { recursive: true, force: true });
  }

  cpSync(sourceNodeModules, targetNodeModules, { recursive: true, dereference: true });

  const standaloneNodeModules = join(nextDir, 'standalone', 'node_modules');
  if (existsSync(standaloneNodeModules)) {
    rmSync(standaloneNodeModules, { recursive: true, force: true });
  }

  // Point standalone/node_modules back to the local copy within .next.
  const relativeTarget = join('..', 'node_modules');
  try {
    symlinkSync(relativeTarget, standaloneNodeModules, 'dir');
  } catch (err) {
    console.warn('Failed to create symlink for standalone/node_modules, copying instead', err);
    cpSync(targetNodeModules, standaloneNodeModules, { recursive: true, dereference: true });
  }
}

async function main() {
  await runNextBuild();
  copyNodeModulesIntoNext();
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
