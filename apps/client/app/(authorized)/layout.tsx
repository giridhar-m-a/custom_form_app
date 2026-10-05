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
        <div className="flex-1 m-4 md:m-8 lg:m-14 rounded-lg p-4 md:p-8 bg-accent overflow-y-auto min-h-0">
          {children}
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
