/** @type {import('next').NextConfig} */
const nextConfig = {
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api',
  },
  output: 'standalone',
  // React 19 默认启用 Strict Mode
}

// Rewrite /favicon.ico to the svg in public so browsers requesting
// /favicon.ico receive a valid asset even if no .ico file exists.
nextConfig.rewrites = async () => {
  return [
    {
      source: '/favicon.ico',
      destination: '/favicon.svg',
    },
  ]
}

module.exports = nextConfig
