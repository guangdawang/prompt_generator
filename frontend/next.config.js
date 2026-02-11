/** @type {import('next').NextConfig} */
const nextConfig = {
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api',
  },
  output: 'standalone',
  // React 19 默认启用 Strict Mode
}

module.exports = nextConfig
