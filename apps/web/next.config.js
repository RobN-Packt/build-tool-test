/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  eslint: {
    // Linting is handled via Bazel (see //apps/web:lint)
    ignoreDuringBuilds: true,
  },
  experimental: {
    optimizePackageImports: ['@testing-library/react', '@testing-library/user-event'],
  },
};

module.exports = nextConfig;
