import { ApiResponse } from '@/types/api.types'
import { AxiosError } from 'axios'

export const errorHandler = <T>(error: any): ApiResponse<T> => {
  console.error(error)
  if (error instanceof AxiosError) {
    return {
      status: error.response?.status || 500,
      message: error.response?.data?.message || error.message || 'An error occurred',
      data: undefined,
      pagination: undefined
    }
  }
  return {
    status: 500,
    message: error.message || 'An error occurred',
    data: undefined,
    pagination: undefined
  }
}
