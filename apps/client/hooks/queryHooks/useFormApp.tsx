'use client'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { createForm, getFormById, getForms } from '@/services/api/forms/routes'
import { CreateFormSchemaType } from '@/app/schemas/form.schemas'
import toast from 'react-hot-toast'
import { useRouter } from 'next/navigation'
import { formsKeys } from '@/lib/constants/queryKeys/forms.keys'
import { FormFilter } from '@/types/form.types'

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
        // router.push(`forms/new/${data.id}`)
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
