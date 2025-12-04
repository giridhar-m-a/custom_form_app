'use server'

import { GET } from '@/lib/api.config'
import { usersRoutes } from '@/lib/constants/apiRoutes/users.routes'
import { errorHandler } from '@/lib/errorHandler'
import { User } from '@/types/user.types'

export const getMe = async () => {
  try {
    return await GET<User>(usersRoutes.me)
  } catch (e) {
    return errorHandler<User>(e)
  }
}
