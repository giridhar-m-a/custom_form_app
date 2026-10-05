'use server'

import { ChangePasswordSchemaType, UserUpdateSchemaType } from '@/app/schemas/user.schemas'
import { GET, PATCH, PUT, PUT_FORM_DATA } from '@/lib/api.config'
import { usersRoutes } from '@/lib/constants/apiRoutes/users.routes'
import { errorHandler } from '@/lib/errorHandler'
import { ApiResponse } from '@/types/api.types'
import { User } from '@/types/user.types'

export const getMe = async () => {
  try {
    return await GET<User>(usersRoutes.me)
  } catch (e) {
    return errorHandler<User>(e)
  }
}

export const uploadProfilePic = async (data: FormData): Promise<ApiResponse<User>> => {
  try {
    return await PUT_FORM_DATA(usersRoutes.profilePic, data)
  } catch (e) {
    return errorHandler<User>(e)
  }
}

export const updateProfile = async (data: UserUpdateSchemaType): Promise<ApiResponse<User>> => {
  try {
    return await PATCH(usersRoutes.me, data)
  } catch (e) {
    return errorHandler<User>(e)
  }
}

export const updatePassword = async (data: ChangePasswordSchemaType): Promise<ApiResponse<User>> => {
  try {
    return await PUT(usersRoutes.updatePassword, data)
  } catch (e) {
    return errorHandler<User>(e)
  }
}
