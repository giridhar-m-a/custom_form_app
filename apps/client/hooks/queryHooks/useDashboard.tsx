'use client'

import { useQuery } from '@tanstack/react-query'
import { getDashboardData } from '@/services/api/dashboard/routes'

export const useDashboard = () => {
  return useQuery({
    queryKey: ['dashboard'],
    queryFn: getDashboardData
  })
}
