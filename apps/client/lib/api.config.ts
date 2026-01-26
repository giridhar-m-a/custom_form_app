import { ApiResponse } from '@/types/api.types'
import axios, { AxiosInstance, AxiosResponse, AxiosRequestConfig, InternalAxiosRequestConfig } from 'axios'
import { errorHandler } from './errorHandler'
import { getStore } from '../store/store'
import Cookies from 'js-cookie'

// Extend InternalAxiosRequestConfig to support explicit token
interface CustomAxiosRequestConfig extends InternalAxiosRequestConfig {
  token?: string
}

class ApiConfig {
  private axiosInstance: AxiosInstance

  constructor() {
    const baseURL: string = `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/`
    this.axiosInstance = axios.create({
      baseURL,
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    // Request interceptor
    this.axiosInstance.interceptors.request.use(
      async (config: InternalAxiosRequestConfig): Promise<InternalAxiosRequestConfig> => {
        try {
          let token: string | null | undefined = null
          const customConfig = config as CustomAxiosRequestConfig

          // Priority 1: Explicit token in config (for server-side calls)
          if (customConfig.token) {
            token = customConfig.token
          }
          // Priority 2: Check if we're on client-side
          else if (typeof window !== 'undefined') {
            // Client-side: Try Redux store first
            try {
              const store = getStore()
              token = store?.getState().auth.accessToken
            } catch (error) {
              if (process.env.NODE_ENV === 'development') {
                console.warn('[API Config] Failed to get token from Redux store:', error)
              }
            }

            // Client-side fallback: Direct cookie access
            if (!token) {
              try {
                token = Cookies.get('accessToken')
              } catch (error) {
                if (process.env.NODE_ENV === 'development') {
                  console.warn('[API Config] Failed to get token from cookies:', error)
                }
              }
            }
          }
          // Priority 3: Server-side - use Next.js cookies()
          else {
            try {
              const { cookies } = await import('next/headers')
              const cookieStore = await cookies()
              token = cookieStore.get('accessToken')?.value
            } catch (error) {
              if (process.env.NODE_ENV === 'development') {
                console.warn('[API Config] Failed to get token from Next.js cookies:', error)
              }
            }
          }

          // Set Authorization header if token is available
          if (token) {
            config.headers.Authorization = `Bearer ${token}`
          }
        } catch (error) {
          // Catch any unexpected errors and log them in development
          if (process.env.NODE_ENV === 'development') {
            console.error('[API Config] Unexpected error in request interceptor:', error)
          }
          // Don't fail the request - proceed without token
        }

        return config
      },
      error => {
        return Promise.reject(error)
      }
    )

    // Response interceptor
    this.axiosInstance.interceptors.response.use(
      (response: AxiosResponse) => {
        return response
      },
      error => {
        // Handle common errors
        if (error.response?.status === 401) {
          // Handle unauthorized access
          //   localStorage.removeItem('authToken')
          // Redirect to login or refresh token
        }
        return Promise.reject(error)
      }
    )
  }

  async GET<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<ApiResponse<T>> = await this.axiosInstance.get(url, config)
      return {
        status: response.status,
        message: response.data.message,
        data: response.data.data,
        pagination: response.data.pagination
      }
    } catch (error: any) {
      return errorHandler(error)
    }
  }

  async POST<T, U = any>(url: string, data?: U, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<ApiResponse<T>> = await this.axiosInstance.post(url, data, config)
      return {
        status: response.status,
        message: response.data.message,
        data: response.data.data,
        pagination: response.data.pagination
      }
    } catch (error: any) {
      console.log('error in post method')
      return errorHandler(error)
    }
  }

  async PUT<T, U = any>(url: string, data?: U, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<ApiResponse<T>> = await this.axiosInstance.put(url, data, config)
      return {
        status: response.status,
        message: response.data.message,
        data: response.data.data,
        pagination: response.data.pagination
      }
    } catch (error: any) {
      return errorHandler(error)
    }
  }

  async PATCH<T, U = any>(url: string, data?: U, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<ApiResponse<T>> = await this.axiosInstance.patch(url, data, config)
      return {
        status: response.status,
        message: response.data.message,
        data: response.data.data,
        pagination: response.data.pagination
      }
    } catch (error: any) {
      return errorHandler(error)
    }
  }

  async DELETE<T>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<ApiResponse<T>> = await this.axiosInstance.delete(url, config)
      return {
        status: response.status,
        message: response.data.message,
        data: response.data.data,
        pagination: response.data.pagination
      }
    } catch (error: any) {
      return errorHandler(error)
    }
  }

  // New method specifically for file uploads
  async UPLOAD_FILE<T>(
    url: string,
    file: File,
    fieldName: string = 'file',
    additionalData?: Record<string, any>,
    config?: AxiosRequestConfig
  ): Promise<ApiResponse<T>> {
    try {
      const formData = new FormData()
      formData.append(fieldName, file)

      // Add any additional form fields
      if (additionalData) {
        Object.entries(additionalData).forEach(([key, value]) => {
          formData.append(key, value)
        })
      }

      const response: AxiosResponse<ApiResponse<T>> = await this.axiosInstance.post(url, formData, {
        ...config,
        headers: {
          ...config?.headers,
          'Content-Type': 'multipart/form-data'
        }
      })

      return {
        status: response.status,
        message: response.data.message,
        data: response.data.data
      }
    } catch (error: any) {
      return errorHandler(error)
    }
  }

  // Alternative method for uploading with FormData directly
  async POST_FORM_DATA<T>(url: string, formData: FormData, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<ApiResponse<T>> = await this.axiosInstance.post(url, formData, {
        ...config,
        headers: {
          ...config?.headers,
          'Content-Type': 'multipart/form-data'
        }
      })

      return {
        status: response.status,
        message: response.data.message,
        data: response.data.data
      }
    } catch (error: any) {
      return errorHandler(error)
    }
  }

  async PUT_FORM_DATA<T>(url: string, formData: FormData, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    try {
      const response: AxiosResponse<ApiResponse<T>> = await this.axiosInstance.put(url, formData, {
        ...config,
        headers: {
          ...config?.headers,
          'Content-Type': 'multipart/form-data'
        }
      })

      return {
        status: response.status,
        message: response.data.message,
        data: response.data.data
      }
    } catch (error: any) {
      return errorHandler(error)
    }
  }
}

// Create and export a singleton instance
const apiConfig = new ApiConfig()

// Export the methods bound to the instance
export const GET = apiConfig.GET.bind(apiConfig)
export const POST = apiConfig.POST.bind(apiConfig)
export const PUT = apiConfig.PUT.bind(apiConfig)
export const PATCH = apiConfig.PATCH.bind(apiConfig)
export const DELETE = apiConfig.DELETE.bind(apiConfig)
export const UPLOAD_FILE = apiConfig.UPLOAD_FILE.bind(apiConfig)
export const POST_FORM_DATA = apiConfig.POST_FORM_DATA.bind(apiConfig)
export const PUT_FORM_DATA = apiConfig.PUT_FORM_DATA.bind(apiConfig)
export default apiConfig
