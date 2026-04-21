'use client'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import {
  createForm,
  createFormField,
  deleteForm,
  getFormById,
  getFormFields,
  getForms,
  updateForm,
  updateFormField
} from '@/services/api/forms/routes'
import { CreateFormSchemaType, FormFieldSchemaType } from '@/app/schemas/form.schemas'
import toast from 'react-hot-toast'
import { useRouter } from 'next/navigation'
import { formsKeys } from '@/lib/constants/queryKeys/forms.keys'
import { FormField, FormFilter, FormUpdateType } from '@/types/form.types'

export const useCreateForm = () => {
  const queryClient = useQueryClient()
  const router = useRouter()
  return useMutation({
    mutationKey: ['forms', 'create'],
    mutationFn: async ({ data }: { data: CreateFormSchemaType }) => {
      const res = await createForm(data)
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message, data }) => {
      queryClient.invalidateQueries({ predicate: ({ queryKey }) => queryKey.includes('form') })
      if (message && data?.id) {
        toast.success(message)
        router.push(`forms/new/${data.id}`)
      }
    },
    onError: error => {
      toast.error(error.message)
    }
  })
}

export const useUpdateFormField = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['forms', 'update-field'],
    mutationFn: async ({ data }: { data: FormFieldSchemaType }) => {
      const res = await updateFormField({ data })
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message, data }) => {
      queryClient.invalidateQueries({ predicate: ({ queryKey }) => queryKey.includes('form') })
      if (message && data) {
        toast.success(message)
      }
    },
    onError: error => {
      toast.error(error.message)
    }
  })
}

export const useGetForms = (params?: FormFilter) => {
  return useQuery({
    queryKey: formsKeys.list({ query: params }),
    queryFn: () => getForms({ query: params })
  })
}

export const useGetFormById = (id: string) => {
  return useQuery({
    queryKey: formsKeys.detail(id),
    queryFn: () => getFormById({ id })
  })
}

export const useGetFormFields = (id: string, initialData?: FormField[]) => {
  return useQuery({
    queryKey: formsKeys.fields(id),
    queryFn: () => getFormFields({ id }),
    initialData: initialData && {
      data: initialData,
      status: 200,      message: 'Form fields loaded successfully'
    }
  })
}

export const useUpdateForm = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['forms', 'update'],
    mutationFn: async ({ id, data }: { id: string; data: FormUpdateType }) => {
      const res = await updateForm({ id, data })
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message, data }) => {
      queryClient.invalidateQueries({ predicate: ({ queryKey }) => queryKey.includes('form') })
      if (message && data?.id) {
        toast.success(message)
        // router.push(`forms/new/${data.id}`)
      }
    },
    onError: error => {
      toast.error(error.message)
    }
  })
}

export const useDeleteForm = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['forms', 'delete'],
    mutationFn: async ({ id }: { id: string }) => {
      const res = await deleteForm({ id })
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message, data }) => {
      queryClient.invalidateQueries({ predicate: ({ queryKey }) => queryKey.includes('form') })
      if (message && data?.id) {
        toast.success(message)
      }
    },
    onError: error => {
      toast.error(error.message)
    }
  })
}

export const useCreateFormField = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['forms', 'create-fields'],
    mutationFn: async ({ data }: { data: FormFieldSchemaType }) => {
      const res = await createFormField({ data })
      if (res.status === 200 || res.status === 201) {
        return res
      }
      throw new Error(res.message)
    },
    onSuccess: ({ message, data }) => {
      queryClient.invalidateQueries({ predicate: ({ queryKey }) => queryKey.includes('form') })
      if (message && data?.length) {
        toast.success(message)
      }
    },
    onError: error => {
      toast.error(error.message)
    }
  })
}
