import { defineConfig } from 'vitest/config';
import react from '@vitejs/plugin-react';
import { resolve } from 'node:path';
import { createRequire } from 'node:module';

const require = createRequire(import.meta.url);
const jestDomVitest = require.resolve('@testing-library/jest-dom/vitest');
const jestDomMatchers = require.resolve('@testing-library/jest-dom/matchers');
const reactJsxRuntime = require.resolve('react/jsx-runtime');
const reactJsxDevRuntime = require.resolve('react/jsx-dev-runtime');
const testingLibraryReact = require.resolve('@testing-library/react');
const testingLibraryUserEvent = require.resolve('@testing-library/user-event');
const reactEntry = require.resolve('react');
const reactDomEntry = require.resolve('react-dom');
const nextLink = require.resolve('next/link');

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': resolve(process.cwd(), '.'),
      '@testing-library/jest-dom/vitest': jestDomVitest,
      '@testing-library/jest-dom/matchers': jestDomMatchers,
      '@testing-library/react': testingLibraryReact,
      '@testing-library/user-event': testingLibraryUserEvent,
      react: reactEntry,
      'react-dom': reactDomEntry,
      'next/link': nextLink,
      'react/jsx-runtime': reactJsxRuntime,
      'react/jsx-dev-runtime': reactJsxDevRuntime,
    }
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts'],
    globals: true,
    css: false
  }
});
