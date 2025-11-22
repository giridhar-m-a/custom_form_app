import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'
import { unAuthorizedPages } from './lib/constants/constants'

export function middleware(req: NextRequest) {
  // Read token from cookies
  const token = req.cookies.get('accessToken')?.value
  const refreshToken = req.cookies.get('refreshToken')?.value

  // If token missing → redirect to home page
  if (!token && !unAuthorizedPages.has(req.nextUrl.pathname)) {
    return NextResponse.redirect(new URL('/', req.url))
  }

  if (token && unAuthorizedPages.has(req.nextUrl.pathname)) {
    return NextResponse.redirect(new URL('/dashboard', req.url))
  }

  // Otherwise allow request through
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
