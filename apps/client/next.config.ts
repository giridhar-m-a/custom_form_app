import type { NextConfig } from 'next'

const isDev = process.env.NODE_ENV === 'development'

const protocol = isDev ? 'http' : 'https'
const hostname = process.env.MINIO_SERVER || 'minio.custom-form-app.home'

const nextConfig: NextConfig = {
  images: {
    unoptimized: isDev,
    remotePatterns: [
      {
        protocol,
        hostname
      }
    ]
  },
  experimental: {
    serverActions: {
      bodySizeLimit: '100mb'
    }
  }
}

export default nextConfig
