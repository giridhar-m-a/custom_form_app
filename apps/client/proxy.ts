import { NextResponse } from 'next/server'
import { NextRequest } from 'next/server'
import { unAuthorizedPages } from './lib/constants/constants'
import { verifyRefreshToken, verifyToken } from './services/api/auth/route'

async function isAuthenticated(
  accessToken: string,
  refreshToken: string | undefined,
  req: NextRequest
): Promise<NextResponse | null> {
  try {
    const verifyResp = await verifyToken(accessToken)

    // Token valid
    if (verifyResp.status === 200 && verifyResp.data?.userID) {
      console.log('Token is valid')
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

        return res
      } else {
        return null
      }
    }

    // Token invalid
    return null
  } catch (err) {
    console.error('Token verification error:', err)
    return null
  }
}

export default async function proxy(req: NextRequest) {
  const pathname = req.nextUrl.pathname

  if (pathname === '/api/auth/refresh') {
    return NextResponse.next()
  }

  const accessToken = req.cookies.get('accessToken')?.value
  const refreshToken = req.cookies.get('refreshToken')?.value

  const isPublic = unAuthorizedPages.has(pathname)

  // ⛔ No token + protected route → send to login
  if (!accessToken && !isPublic) {
    return NextResponse.redirect(new URL('/', req.url))
  }

  // ✔ Logged-in user on a public page → verify first, then redirect to dashboard
  if (accessToken && isPublic) {
    const authResult = await isAuthenticated(accessToken, refreshToken, req)
    if (authResult) {
      // Token is valid (or refreshed) → redirect to dashboard
      return NextResponse.redirect(new URL('/dashboard', req.url))
    } else {
      const res = NextResponse.next()
      res.cookies.delete('accessToken')
      res.cookies.delete('refreshToken')
      return res
    }
  }

  // ✔ Token check only on first (full) page load, not on client-side route changes.
  // Next.js App Router sets the `RSC: 1` header on all client-side navigations.
  // A true full page load has no RSC header.
  const isFullPageLoad = req.headers.get('RSC') !== '1'

  if (accessToken && req.method === 'GET' && isFullPageLoad) {
    const authResult = await isAuthenticated(accessToken, refreshToken, req)
    if (authResult) {
      return authResult
    }
    // Token invalid/expired → send to login

    const res = NextResponse.redirect(new URL('/', req.url))
    res.cookies.delete('accessToken')
    res.cookies.delete('refreshToken')
    return res
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    '/((?!_next|[^?]*\\.(?:html?|css|js(?!on)|jpe?g|webp|png|gif|svg|ttf|woff2?|ico|csv|docx?|xlsx?|zip|webmanifest)).*)',
    '/(api|trpc)(.*)'
  ]
}
