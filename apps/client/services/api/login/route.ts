'use server'

import { SignInSchema, SignInSchemaType, SignUpSchema, SignUpSchemaType } from '@/app/schemas/auth.schemas'
import { GET, POST } from '@/lib/api.config'
import { AUTH_ROUTES } from '@/lib/constants/apiRoutes/auth.routes'
import { errorHandler } from '@/lib/errorHandler'
import { AuthResponse } from '@/types/auth.types'

export const loginWithCredentials = async (data: SignInSchemaType) => {
  try {
    const parsed = await SignInSchema.safeParse(data)
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
    const parsed = await SignUpSchema.safeParse(data)
    if (!parsed.success) {
      throw new Error(JSON.stringify(parsed.error))
    }
    const response = await POST<AuthResponse>(AUTH_ROUTES.register, parsed.data)
    return response
  } catch (e) {
    return errorHandler<AuthResponse>(e)
  }
}
