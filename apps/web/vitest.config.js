const { defineConfig } = require('vitest/config');
const react = require('@vitejs/plugin-react');
const path = require('node:path');

module.exports = defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(process.cwd(), '.')
    }
  },
  test: {
    environment: 'jsdom',
    setupFiles: ['./tests/setup.ts'],
    globals: true,
    css: false
  }
});
