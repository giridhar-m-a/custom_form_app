'use client'
import { SignInSchemaType, SignUpSchemaType } from '@/app/schemas/auth.schemas'
import { LOGIN_KEYS } from '@/lib/constants/queryKeys/login.keys'
import { loginWithCredentials, loginWithGoogle, register } from '@/services/api/login/route'
import { setTokens } from '@/store/slices/auth.slice'
import { useMutation } from '@tanstack/react-query'
import { useRouter } from 'next/navigation'
import toast from 'react-hot-toast'
import { useDispatch } from 'react-redux'

export const useCredentialAuth = () => {
  const router = useRouter()
  const dispatch = useDispatch()
  return useMutation({
    mutationKey: LOGIN_KEYS.loginWithCredential,
    mutationFn: async (data: SignInSchemaType) => {
      const res = await loginWithCredentials(data)
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message, data }) => {
      if (data) {
        dispatch(setTokens({ accessToken: data.accessToken, refreshToken: data.refreshToken }))
      }
      toast.success(message)
      router.push('/dashboard')
    },
    onError: ({ message }) => {
      toast.error(message)
    }
  })
}

export const useGoogleAuth = () => {
  const router = useRouter()
  const dispatch = useDispatch()
  return useMutation({
    mutationKey: LOGIN_KEYS.loginWithGoogle,
    mutationFn: async (code: string) => {
      const res = await loginWithGoogle(code)
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message, data }) => {
      if (data) {
        dispatch(setTokens({ accessToken: data.accessToken, refreshToken: data.refreshToken }))
      }
      toast.success(message)
      router.push('/dashboard')
    },
    onError: ({ message }) => {
      console.log('error: ', message)
      toast.error(message)
    }
  })
}

export const useRegister = () => {
  const router = useRouter()
  const dispatch = useDispatch()
  return useMutation({
    mutationKey: LOGIN_KEYS.register,
    mutationFn: async (data: SignUpSchemaType) => {
      const res = await register(data)
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message, data }) => {
      if (data) {
        dispatch(setTokens({ accessToken: data.accessToken, refreshToken: data.refreshToken }))
      }
      toast.success(message)
      router.push('/dashboard')
    },
    onError: ({ message }) => {
      toast.error(message)
    }
  })
}
