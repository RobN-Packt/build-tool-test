const path = require('path');

/** @type {import('next').NextConfig} */
const nextConfig = {
  ...(process.env.DOCKER_BUILD === 'true' ? { output: 'standalone' } : {}),

  experimental: {
    optimizePackageImports: [
      '@testing-library/react',
      '@testing-library/user-event',
    ],
  },

  eslint: {
    ignoreDuringBuilds: true,
  },
  outputFileTracing: true,
  outputFileTracingRoot: path.join(__dirname, '../../'),
};

module.exports = nextConfig;
