'use client'

import { useQuery } from '@tanstack/react-query'
import { getMe } from '@/services/api/users/routes'

export const useGetMe = () => {
  return useQuery({
    queryKey: ['users', 'me'],
    queryFn: getMe
  })
}
