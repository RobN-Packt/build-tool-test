#!/usr/bin/env node
const { createRequire } = require('module');
const { spawnSync } = require('child_process');

const requireFromHere = createRequire(__filename);
const cliPath = requireFromHere.resolve('next/dist/bin/next');

const result = spawnSync(process.execPath, [cliPath, ...process.argv.slice(2)], {
  stdio: 'inherit',
});

if (result.error) {
  throw result.error;
}

if (typeof result.status === 'number') {
  process.exit(result.status);
}

if (typeof result.signal === 'string') {
  process.kill(process.pid, result.signal);
}
