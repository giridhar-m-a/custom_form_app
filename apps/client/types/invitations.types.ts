export interface Invitation {
  invitationId: string
  formId: string
  invitedEmail: string
  status: InvitationStatus
  invitedBy: string
  invitedName: string
}

export type InvitationStatus = 'pending' | 'submitted' | 'invited' | 'failed'

export interface InvitationFilter {
  search?: string
  status?: InvitationStatus
  excludeStatus?: InvitationStatus
  page?: number
  limit?: number
}

export interface BulkInvitationResponse {
  failed_count: number
  success_count: number
  total_rows: number
}
