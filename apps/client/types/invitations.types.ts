export interface Invitation {
  invitationId: string
  formId: string
  invitedEmail: string
  status: InvitationStatus
  invitedBy: string
  invitedName: string
}

export type InvitationStatus =
  | 'pending'
  | 'bounced'
  | 'clicked'
  | 'opened'
  | 'delivered'
  | 'complained'
  | 'delayed'
  | 'failed'
  | 'submitted'

export interface InvitationFilter {
  search?: string
  status?: InvitationStatus
  exclude?: InvitationStatus[]
  page?: number
  limit?: number
}

export interface BulkInvitationResponse {
  failed_count: number
  success_count: number
  total_rows: number
}

export interface VerifyInvitationResponse {
  formId: string
  invitationId?: string
}

export type AnonymousInvitationPayload = {
  formId: string
}

export type AnonymousInvitationResponse = {
  token: string
  expiresIn: string
}
