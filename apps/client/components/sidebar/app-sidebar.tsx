'use client'

import { Bot, FileQuestionMarkIcon, GalleryVerticalEnd, Home } from 'lucide-react'
import * as React from 'react'
import { NavMain } from '@/components/sidebar/nav-main'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
  SidebarTrigger
} from '@/components/ui/sidebar'
import { SidebarMenuItemType } from '@/types/sideBar.types'
import { useSelector } from 'react-redux'
import { getMe } from '@/store/slices/me.slice'
import { NavUser } from './nav-user'

type SideBarData = {
  details: {
    name: string
    logo: React.ElementType
    description: string
  }

  navMain: SidebarMenuItemType[]
}

const data: SideBarData = {
  details: {
    name: 'Custom Form App',
    logo: GalleryVerticalEnd,
    description: 'Custom Form App'
  },

  navMain: [
    {
      title: 'Dashboard',
      icon: Home,
      url: '/dashboard'
      // items: [
      //   {
      //     title: 'Dashboard',
      //     url: '/dashboard',
      //     icon: Home
      //   },
      //   {
      //     title: 'Ai Insights',
      //     url: '/dashboard/ai-insights',
      //     icon: Bot
      //   }
      // ]
    },
    {
      title: 'Forms',
      icon: FileQuestionMarkIcon,
      url: '/forms'
      // items: [
      //   {
      //     title: 'Forms',
      //     url: '/forms',
      //     icon: FileQuestionMarkIcon
      //   },
      //   {
      //     title: 'Form Report',
      //     url: '/forms/reports',
      //     icon: Bot
      //   }
      // ]
    }
  ]
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  const user = useSelector(getMe)
  return (
    <Sidebar collapsible="icon" {...props}>
      <SidebarHeader>
        <SidebarTrigger />
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
