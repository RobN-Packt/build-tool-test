/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',
  experimental: {
    optimizePackageImports: ['@testing-library/react', '@testing-library/user-event']
  }
};

module.exports = nextConfig;
