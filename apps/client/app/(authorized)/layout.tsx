'use client'
import Header from '@/components/sidebar/Header'
import { AppSidebar } from '@/components/sidebar/app-sidebar'
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator
} from '@/components/ui/breadcrumb'
import { Separator } from '@/components/ui/separator'
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { useGetMe } from '@/hooks/queryHooks/useUsers'
import { setMe } from '@/store/slices/me.slice'
import { useEffect } from 'react'
import { useDispatch } from 'react-redux'

export default function Layout({ children }: { children: React.ReactNode }) {
  const { data } = useGetMe()
  const dispatch = useDispatch()
  useEffect(() => {
    if (data?.data) {
      dispatch(setMe(data.data))
    }
  }, [data?.data])
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
