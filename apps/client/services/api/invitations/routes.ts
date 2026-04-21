'use server'

import { InvitationType } from '@/app/schemas/invitation.schemas'
import { DELETE, GET, POST, POST_FORM_DATA } from '@/lib/api.config'
import { invitationsRoutes } from '@/lib/constants/apiRoutes/invitations.routes'
import { errorHandler } from '@/lib/errorHandler'
import { ApiResponse } from '@/types/api.types'
import {
  BulkInvitationResponse,
  Invitation,
  InvitationFilter,
  VerifyInvitationResponse
} from '@/types/invitations.types'

interface GetInvitationsByFormIdParams {
  formId: string
  params: InvitationFilter
}

export const getInvitationsByFormId = async ({
  formId,
  params
}: GetInvitationsByFormIdParams): Promise<ApiResponse<Invitation[]>> => {
  try {
    const exclude = params.exclude?.length ? params.exclude : []
    const res = await GET<Invitation[]>(`${invitationsRoutes.base}`, {
      params: { ...params, formId, exclude }
    })
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<Invitation[]>(e)
  }
}

export const createInvitation = async ({ data }: { data: InvitationType }): Promise<ApiResponse<Invitation>> => {
  try {
    const res = await POST<Invitation>(`${invitationsRoutes.base}`, data)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<Invitation>(e)
  }
}

export const createBulkInvitation = async ({
  data,
  formId
}: {
  data: FormData
  formId: string
}): Promise<ApiResponse<BulkInvitationResponse>> => {
  try {
    const res = await POST_FORM_DATA<BulkInvitationResponse>(`${invitationsRoutes.base}/${formId}`, data)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<BulkInvitationResponse>(e)
  }
}

export const deleteInvitation = async ({ invitationId }: { invitationId: string }): Promise<ApiResponse<undefined>> => {
  try {
    const res = await DELETE<undefined>(`${invitationsRoutes.base}/${invitationId}`)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<undefined>(e)
  }
}

export const verifyInvitation = async ({
  token
}: {
  token: string
}): Promise<ApiResponse<VerifyInvitationResponse>> => {
  try {
    const res = await POST<VerifyInvitationResponse>(`${invitationsRoutes.verify}`, { token })
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<VerifyInvitationResponse>(e)
  }
}
