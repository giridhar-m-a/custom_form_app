'use server'

import { ApiResponse } from '@/types/api.types'
import { DashboardData } from '@/types/dashboard.types'

import { dashboardRoutes } from '@/lib/constants/apiRoutes/dashboard.routes'
import { GET } from '@/lib/api.config'
import { errorHandler } from '@/lib/errorHandler'

export const getDashboardData = async (): Promise<ApiResponse<DashboardData>> => {
  try {
    return await GET<DashboardData>(dashboardRoutes.base)
  } catch (error) {
    return errorHandler<DashboardData>(error)
  }
}
