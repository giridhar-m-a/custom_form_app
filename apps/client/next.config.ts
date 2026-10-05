import type { NextConfig } from 'next'

const isDev = process.env.NODE_ENV === 'development'

const protocol = isDev ? 'http' : 'https'
const hostname = process.env.MINIO_SERVER || 'rustfs-me-obs.giridhar.dev'

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
