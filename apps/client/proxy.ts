import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import { unAuthorizedPages } from './lib/constants/constants'
import { verifyRefreshToken, verifyToken } from './services/api/login/route'

export default async function proxy(req: NextRequest) {
  const pathname = req.nextUrl.pathname

  const accessToken = req.cookies.get('accessToken')?.value
  const refreshToken = req.cookies.get('refreshToken')?.value

  const isPublic = unAuthorizedPages.has(pathname)

  // If no token and a protected route → redirect to login/home
  if (!accessToken && !isPublic) {
    return NextResponse.redirect(new URL('/', req.url))
  }

  // Token exists AND it's a GET request AND page is protected
  if (accessToken && req.method === 'GET' && !isPublic) {
    try {
      const verifyResp = await verifyToken(accessToken)

      // If token is valid → allow
      if (verifyResp.status === 200 && verifyResp.data?.userID) {
        return NextResponse.next()
      }

      // If access token expired but refresh exists
      if (verifyResp.status === 401 && refreshToken) {
        const refreshResp = await verifyRefreshToken(refreshToken)

        if (refreshResp.status === 200 && refreshResp.data?.accessToken && refreshResp.data?.refreshToken) {
          const res = NextResponse.next()
          res.cookies.set('accessToken', refreshResp.data.accessToken)
          res.cookies.set('refreshToken', refreshResp.data.refreshToken)
          return res
        }

        return NextResponse.redirect(new URL('/', req.url))
      }

      return NextResponse.redirect(new URL('/', req.url))
    } catch (err) {
      console.error('Token verification failed', err)
      return NextResponse.redirect(new URL('/', req.url))
    }
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    // Skip Next.js internals and all static files, unless found in search params
    '/((?!_next|[^?]*\\.(?:html?|css|js(?!on)|jpe?g|webp|png|gif|svg|ttf|woff2?|ico|csv|docx?|xlsx?|zip|webmanifest)).*)',
    // Always run for API routes
    '/(api|trpc)(.*)'
  ]
}
