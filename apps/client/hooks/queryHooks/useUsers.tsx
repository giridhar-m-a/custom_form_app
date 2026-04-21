'use client'

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { getMe, updatePassword, updateProfile, uploadProfilePic } from '@/services/api/users/routes'
import toast from 'react-hot-toast'
import { ChangePasswordSchemaType, UserUpdateSchemaType } from '@/app/schemas/user.schemas'
import { useDispatch } from 'react-redux'
import { setMyName } from '@/store/slices/me.slice'

export const useGetMe = () => {
  return useQuery({
    queryKey: ['users', 'me'],
    queryFn: getMe
  })
}

export const useUploadProfilePic = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['users', 'upload-profile-pic'],
    mutationFn: (data: FormData) => uploadProfilePic(data),
    onSuccess: ({ message, status }) => {
      if (status === 200) {
        toast.success(message)
        queryClient.invalidateQueries({ predicate: query => query.queryKey[0] === 'users' })
      } else {
        throw new Error(message)
      }
    },
    onError: error => {
      toast.error(error.message)
    }
  })
}

export const useUpdatePassword = () => {
  return useMutation({
    mutationKey: ['users', 'update-password'],
    mutationFn: (data: ChangePasswordSchemaType) => updatePassword(data),
    onSuccess: ({ message, status }) => {
      if (status === 200) {
        toast.success(message)
      } else {
        throw new Error(message)
      }
    },
    onError: error => {
      toast.error(error.message)
    }
  })
}

export const useUpdateUser = () => {
  const queryClient = useQueryClient()
  const dispatch = useDispatch()
  return useMutation({
    mutationKey: ['users', 'update-user'],
    mutationFn: (data: UserUpdateSchemaType) => updateProfile(data),
    onSuccess: ({ message, status, data }) => {
      if (status === 200 && data) {
        dispatch(setMyName(data.fullName))
        queryClient.invalidateQueries({ predicate: query => query.queryKey.includes('users') })
        toast.success(message)
      } else {
        throw new Error(message)
      }
    },
    onError: error => {
      toast.error(error.message)
    }
  })
}
