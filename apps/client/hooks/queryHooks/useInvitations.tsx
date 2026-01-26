import { InvitationType } from '@/app/schemas/invitation.schemas'
import { GET_INVITATIONS_BY_FORM_ID } from '@/lib/constants/queryKeys/invitations.keys'
import {
  createBulkInvitation,
  createInvitation,
  deleteInvitation,
  getInvitationsByFormId
} from '@/services/api/invitations/routes'
import { InvitationFilter } from '@/types/invitations.types'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import toast from 'react-hot-toast'

export const useInvitations = ({ formId, params }: { formId: string; params: InvitationFilter }) => {
  return useQuery({
    queryKey: GET_INVITATIONS_BY_FORM_ID(formId, params),
    queryFn: () => getInvitationsByFormId({ formId, params }),
    enabled: !!formId
  })
}

export const useCreateInvitation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['createInvitation'],
    mutationFn: ({ data }: { data: InvitationType }) => createInvitation({ data }),
    onSuccess: ({ message }) => {
      toast.success(message)
      queryClient.invalidateQueries({ queryKey: ['invitations'] })
    },
    onError: ({ message }) => {
      toast.error(message)
    }
  })
}

export const useBulkCreateInvitation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['bulkCreateInvitation'],
    mutationFn: ({ data, formId }: { data: FormData; formId: string }) => createBulkInvitation({ data, formId }),
    onSuccess: ({ data }) => {
      if (data?.failed_count === data?.total_rows) {
        toast.error('All invitations failed')
      } else if (data) {
        const message = `${data?.success_count} invitations sent successfully and ${data?.failed_count} invitations failed`
        toast.success(message)
        queryClient.invalidateQueries({ queryKey: ['invitations'] })
      }
    },
    onError: ({ message }) => {
      toast.error(message)
    }
  })
}

export const useDeleteInvitation = () => {
  const queryClient = useQueryClient()
  return useMutation({
    mutationKey: ['deleteInvitation'],
    mutationFn: ({ invitationId }: { invitationId: string }) => deleteInvitation({ invitationId }),
    onSuccess: ({ message }) => {
      toast.success(message)
      queryClient.invalidateQueries({ queryKey: ['invitations'] })
    },
    onError: ({ message }) => {
      toast.error(message)
    }
  })
}
