// app/api/auth/refresh/route.ts
import { NextResponse } from 'next/server'
import { verifyRefreshToken } from '@/services/api/auth/route'

export async function GET(request: Request) {
  const searchParams = new URL(request.url).searchParams
  const refreshToken = searchParams.get('refreshToken')

  console.log('[api/auth/refresh] refresh token: ', refreshToken)

  if (!refreshToken) {
    return NextResponse.json({ error: 'No refresh token' }, { status: 401 })
  }

  try {
    console.log('refresh token: ', refreshToken)
    const refreshResp = await verifyRefreshToken(refreshToken)
    console.log('refresh resp: ', refreshResp)

    const cookieExpiry = new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)

    const response = NextResponse.json({
      accessToken: refreshResp.data?.accessToken,
      refreshToken: refreshResp.data?.refreshToken
    })

    // Next.js sets these on the browser via Set-Cookie headers
    if (refreshResp.data) {
      response.cookies.set('accessToken', refreshResp.data?.accessToken, { expires: cookieExpiry, path: '/' })
      response.cookies.set('refreshToken', refreshResp.data?.refreshToken, { expires: cookieExpiry, path: '/' })
    }
    return response
  } catch (error) {
    console.error('Refresh failed:', error)
    const res = NextResponse.redirect('/')
    return res
  }
}
