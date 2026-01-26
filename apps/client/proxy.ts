import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import { unAuthorizedPages } from './lib/constants/constants'
import { verifyRefreshToken, verifyToken } from './services/api/login/route'

export default async function proxy(req: NextRequest) {
  const pathname = req.nextUrl.pathname

  const accessToken = req.cookies.get('accessToken')?.value
  const refreshToken = req.cookies.get('refreshToken')?.value

  const isPublic = unAuthorizedPages.has(pathname)
  const isDashboard = pathname.startsWith('/dashboard')

  // ⛔ No token + protected route → send to login
  if (!accessToken && !isPublic) {
    return NextResponse.redirect(new URL('/', req.url))
  }

  // ✔ Token check only during page rendering (GET requests)
  if (accessToken && !isPublic && req.method === 'GET') {
    try {
      const verifyResp = await verifyToken(accessToken)

      // Token valid
      if (verifyResp.status === 200 && verifyResp.data?.userID) {
        // ⛔ Do NOT redirect if already on dashboard
        if (isPublic) {
          return NextResponse.redirect(new URL('/dashboard', req.url))
        }

        return NextResponse.next()
      }

      // Access token expired → try refresh
      if (verifyResp.status === 401 && refreshToken) {
        const refreshResp = await verifyRefreshToken(refreshToken)

        if (refreshResp.status === 200 && refreshResp.data?.accessToken && refreshResp.data?.refreshToken) {
          const res = NextResponse.next()

          res.cookies.set('accessToken', refreshResp.data.accessToken, {
            httpOnly: true,
            path: '/',
            expires: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)
          })

          res.cookies.set('refreshToken', refreshResp.data.refreshToken, {
            httpOnly: true,
            path: '/',
            expires: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)
          })

          // After refresh, redirect only if on public page
          if (isPublic) {
            return NextResponse.redirect(new URL('/dashboard', req.url))
          }

          return res
        }

        return NextResponse.redirect(new URL('/', req.url))
      }

      // Anything else → redirect to login
      return NextResponse.redirect(new URL('/', req.url))
    } catch (err) {
      console.error('Token verification error:', err)
      return NextResponse.redirect(new URL('/', req.url))
    }
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    '/((?!_next|[^?]*\\.(?:html?|css|js(?!on)|jpe?g|webp|png|gif|svg|ttf|woff2?|ico|csv|docx?|xlsx?|zip|webmanifest)).*)',
    '/(api|trpc)(.*)'
  ]
}
