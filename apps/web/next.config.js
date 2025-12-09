/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  experimental: {
    optimizePackageImports: ['@testing-library/react', '@testing-library/user-event']
  },
  eslint: {
    // Next's built-in lint runner isn't yet compatible with ESLint 9 flat config.
    ignoreDuringBuilds: true
  }
};

module.exports = nextConfig;
