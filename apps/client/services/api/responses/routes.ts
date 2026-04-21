'use server'

import { FormSubmission } from '@/app/schemas/response.schemas'
import { GET, POST } from '@/lib/api.config'
import { responsesRoutes } from '@/lib/constants/apiRoutes/responses.routes'
import { errorHandler } from '@/lib/errorHandler'
import { ApiResponse } from '@/types/api.types'
import { ResponseDetails, ResponseFilter, ResponseList } from '@/types/responses.types'

export const createResponse = async ({
  data,
  token
}: {
  data: FormSubmission
  token: string
}): Promise<ApiResponse<any>> => {
  try {
    console.dir(data, { depth: null })
    const res = await POST<any>(`${responsesRoutes.base}`, data, { headers: { Authorization: `Bearer ${token}` } })
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<any>(e)
  }
}

export const getResponses = async (formId: string, params: ResponseFilter): Promise<ApiResponse<ResponseList[]>> => {
  try {
    const res = await GET<ResponseList[]>(`${responsesRoutes.base}/${formId}`, { params })
    return res
  } catch (e) {
    return errorHandler<ResponseList[]>(e)
  }
}

export const getSingleResponse = async (submissionId: string): Promise<ApiResponse<ResponseDetails>> => {
  try {
    const res = await GET<ResponseDetails>(`${responsesRoutes.submission}/${submissionId}`)
    return res
  } catch (e) {
    return errorHandler<ResponseDetails>(e)
  }
}

export const getSignedUrl = async (filePath: string): Promise<ApiResponse<{ signedUrl: string }>> => {
  try {
    const res = await POST<{ signedUrl: string }>(`${responsesRoutes.getSignedUrl}`, { filePath })
    return res
  } catch (e) {
    return errorHandler<{ signedUrl: string }>(e)
  }
}
