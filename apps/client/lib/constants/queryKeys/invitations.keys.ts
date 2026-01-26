import { InvitationFilter } from '@/types/invitations.types'

const GET_INVITATIONS_BY_FORM_ID = (formId: string, params: InvitationFilter) => ['invitations', formId, params]

export { GET_INVITATIONS_BY_FORM_ID }
