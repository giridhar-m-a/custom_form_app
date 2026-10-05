import { LucideIcon } from 'lucide-react'

export type SidebarMenuItemType = {
  title: string
  url?: string
  icon?: LucideIcon
  isActive?: boolean
  items?: {
    title: string
    url: string
    icon?: LucideIcon
  }[]
}
