'use client'

import { FormSubmission } from '@/app/schemas/response.schemas'
import { createResponse, getResponses, getSignedUrl, getSingleResponse } from '@/services/api/responses/routes'
import { ResponseFilter } from '@/types/responses.types'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'

export const useCreateResponse = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['createResponse'],
    mutationFn: ({ data, token }: { data: FormSubmission; token: string }) => createResponse({ data, token }),
    onSuccess: ({ message, status }) => {
      if (status === 200) {
        toast.success(message)
        queryClient.invalidateQueries({ queryKey: ['responses'] })
      } else {
        toast.error(message)
      }
    },
    onError: ({ message }) => {
      toast.error(message)
    }
  })
}

export const useResponses = ({ formId, params }: { formId: string; params: ResponseFilter }) => {
  return useQuery({
    queryKey: ['responses', formId, params],
    queryFn: () => getResponses(formId, params)
  })
}

export const useSingleResponse = (submissionId: string) => {
  return useQuery({
    queryKey: ['singleResponse', submissionId],
    queryFn: () => getSingleResponse(submissionId)
  })
}

export const useDownloadResponseFiles = () => {
  return useMutation({
    mutationKey: ['downloadResponseFiles'],
    mutationFn: ({ filePath }: { filePath: string; fileName: string }) => getSignedUrl(filePath),
    onSuccess: ({ message, status, data }, variables) => {
      if (status === 200 && data?.signedUrl) {
        downloadFile(data.signedUrl, variables.fileName)
      } else {
        throw new Error(message)
      }
    },
    onError: ({ message }) => {
      toast.error(message)
    }
  })
}

const downloadFile = async (url: string, filename: string) => {
  try {
    const response = await fetch(url)
    const blob = await response.blob()
    const objectUrl = URL.createObjectURL(blob)

    const link = document.createElement('a')
    link.href = objectUrl
    link.download = filename
    link.click()

    URL.revokeObjectURL(objectUrl)
  } catch (e) {
    throw new Error('Failed to download file')
  }
}
