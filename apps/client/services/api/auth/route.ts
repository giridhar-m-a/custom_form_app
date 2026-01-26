'use server'

import {
  RequestPasswordResetSchemaType,
  ResetPasswordSchemaType,
  SignInSchema,
  SignInSchemaType,
  SignUpSchema,
  SignUpSchemaType
} from '@/app/schemas/auth.schemas'
import { GET, POST } from '@/lib/api.config'
import { AUTH_ROUTES } from '@/lib/constants/apiRoutes/auth.routes'
import { errorHandler } from '@/lib/errorHandler'
import { AuthResponse, RefreshTokenResponse, VerifyTokenResponse } from '@/types/auth.types'

export const loginWithCredentials = async (data: SignInSchemaType) => {
  try {
    const parsed = SignInSchema.safeParse(data)
    if (!parsed.success) {
      throw new Error(JSON.stringify(parsed.error))
    }
    const response = await POST<AuthResponse>(AUTH_ROUTES.login.credentials, parsed.data)
    return response
  } catch (e) {
    console.error(e)
    return errorHandler<AuthResponse>(e)
  }
}

export const loginWithGoogle = async (code: string) => {
  try {
    const response = await GET<AuthResponse>(`${AUTH_ROUTES.login.google}?code=${code}`)
    return response
  } catch (e) {
    return errorHandler<AuthResponse>(e)
  }
}

export const register = async (data: SignUpSchemaType) => {
  try {
    const parsed = SignUpSchema.safeParse(data)
    if (!parsed.success) {
      throw new Error(JSON.stringify(parsed.error))
    }
    const response = await POST<AuthResponse>(AUTH_ROUTES.register, { ...parsed.data, confirmPassword: undefined })
    return response
  } catch (e) {
    return errorHandler<AuthResponse>(e)
  }
}

export const verifyToken = async (token: string) => {
  try {
    const response = await GET<VerifyTokenResponse>(`${AUTH_ROUTES.verify}?token=${token}`)
    return response
  } catch (e) {
    return errorHandler<VerifyTokenResponse>(e)
  }
}

export const verifyRefreshToken = async (token: string) => {
  try {
    const response = await GET<RefreshTokenResponse>(`${AUTH_ROUTES.refreshToken}?token=${token}`)
    return response
  } catch (e) {
    return errorHandler<RefreshTokenResponse>(e)
  }
}

export const requestPasswordReset = async (data: RequestPasswordResetSchemaType) => {
  try {
    const response = await POST<AuthResponse>(AUTH_ROUTES.resetRequest, data)
    return response
  } catch (e) {
    return errorHandler<AuthResponse>(e)
  }
}

export const resetPassword = async (data: ResetPasswordSchemaType) => {
  try {
    const response = await POST<AuthResponse>(AUTH_ROUTES.resetPassword, data)
    return response
  } catch (e) {
    return errorHandler<AuthResponse>(e)
  }
}
