'use client'
import {
  RequestPasswordResetSchemaType,
  ResetPasswordSchemaType,
  SignInSchemaType,
  SignUpSchemaType
} from '@/app/schemas/auth.schemas'
import { LOGIN_KEYS } from '@/lib/constants/queryKeys/login.keys'
import {
  loginWithCredentials,
  loginWithGoogle,
  register,
  requestPasswordReset,
  resetPassword,
  verifyRefreshToken,
  verifyToken
} from '@/services/api/auth/route'
import { clearTokens, setTokens } from '@/store/slices/auth.slice'
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

export const useResetPassword = () => {
  const { push } = useRouter()
  return useMutation({
    mutationKey: LOGIN_KEYS.resetPassword,
    mutationFn: async (data: ResetPasswordSchemaType) => {
      const res = await resetPassword(data)
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message }) => {
      toast.success(message)
      push('/')
    },
    onError: ({ message }) => {
      toast.error(message)
    }
  })
}

export const useRequestPasswordReset = () => {
  return useMutation({
    mutationKey: LOGIN_KEYS.requestPasswordReset,
    mutationFn: async (data: RequestPasswordResetSchemaType) => {
      const res = await requestPasswordReset(data)
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message }) => {
      toast.success(message)
    },
    onError: ({ message }) => {
      toast.error(message)
    }
  })
}

export const useReAuth = () => {
  const router = useRouter()
  const dispatch = useDispatch()
  return useMutation({
    mutationKey: LOGIN_KEYS.reAuth,
    mutationFn: async (data: { accessToken: string; refreshToken: string }) => {
      const res = await verifyToken(data.accessToken)
      if (res.status === 200 || res.status === 201) {
        return true
      } else if (res.status === 401 && data.refreshToken) {
        const refreshRes = await verifyRefreshToken(data.refreshToken)
        if ((refreshRes.status === 200 || refreshRes.status === 201) && refreshRes.data && refreshRes.data) {
          dispatch(setTokens({ accessToken: refreshRes.data.accessToken, refreshToken: refreshRes.data.refreshToken }))
          return true
        }
      }
      throw new Error(res.message)
    },
    onError: () => {
      dispatch(clearTokens())
      router.push('/')
    }
  })
}
