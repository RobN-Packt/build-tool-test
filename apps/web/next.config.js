/** @type {import('next').NextConfig} */
const nextConfig = {
  // Only use standalone output for Docker builds, not Vercel
  ...(process.env.DOCKER_BUILD === 'true' ? { output: 'standalone' } : {}),
  experimental: {
    optimizePackageImports: ['@testing-library/react', '@testing-library/user-event']
  },
  eslint: {
    // Next's built-in lint runner isn't yet compatible with ESLint 9 flat config.
    ignoreDuringBuilds: true
  },
    outputFileTracing: false,
};

module.exports = nextConfig;
