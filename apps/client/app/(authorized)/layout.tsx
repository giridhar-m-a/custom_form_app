'use client'
import Header from '@/components/sidebar/Header'
import { AppSidebar } from '@/components/sidebar/app-sidebar'
import { Separator } from '@/components/ui/separator'
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar'
import { useGetMe } from '@/hooks/queryHooks/useUsers'
import { setMe } from '@/store/slices/me.slice'
import { useEffect } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import { useReAuth } from '@/hooks/queryHooks/useAuth'
import { getTokens, setTokens } from '@/store/slices/auth.slice'
import Cookies from 'js-cookie'

export default function Layout({ children }: { children: React.ReactNode }) {
  const { mutate } = useReAuth()
  const { data } = useGetMe()
  const dispatch = useDispatch()
  const session = useSelector(getTokens)
  useEffect(() => {
    if (data?.data) {
      dispatch(setMe(data.data))
    }
  }, [data?.data])
  useEffect(() => {
    const interval = setInterval(
      () => {
        if (!session.accessToken || !session.refreshToken) {
          console.log('tokens missing')
          return
        }

        console.log('tokens are present')

        mutate({
          accessToken: session.accessToken,
          refreshToken: session.refreshToken
        })
      },
      60 * 60 * 1000
    )

    return () => clearInterval(interval)
  }, [])
  useEffect(() => {
    const accessToken = Cookies.get('accessToken')
    const refreshToken = Cookies.get('refreshToken')
    if (accessToken && refreshToken) {
      dispatch(setTokens({ accessToken, refreshToken }))
    }
  }, [])
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset className="w-full">
        <header className="flex h-16 shrink-0 items-center gap-2 border-b w-ful">
          <div className="flex items-center gap-2 px-3 w-full">
            <Separator orientation="vertical" className="mr-2 h-4" />
            <Header />
          </div>
        </header>
        <div className="m-14 rounded-lg p-8 min-h-[calc(100vh-11rem)] max-h-[calc(100vh-11rem)] bg-accent">
          {children}
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
