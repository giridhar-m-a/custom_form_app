import { Modal } from '@/components/common/Modal'
import { WarningModalContent } from '@/components/common/WarningModalContent'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { useDeleteInvitation } from '@/hooks/queryHooks/useInvitations'
import { Invitation } from '@/types/invitations.types'
import { ColumnDef } from '@tanstack/react-table'
import { useState } from 'react'

export const InvitationsColumn = ({ page, count }: { page: number; count: number }): ColumnDef<Invitation>[] => [
  {
    header: 'SI No.',
    accessorKey: 'invitationId',
    cell: ({ row }) => row.index + (page - 1) * count + 1
  },
  {
    header: 'Invited Email',
    accessorKey: 'invitedEmail'
  },
  {
    header: 'Invited Name',
    accessorKey: 'invitedName'
  },
  {
    header: 'Status',
    accessorKey: 'status',
    cell: ({ row }) => {
      switch (row.original.status) {
        case 'pending':
          return <Badge className="bg-gray-500 text-white">Pending</Badge>
        case 'delivered':
          return <Badge className="bg-blue-500 text-white">Delivered</Badge>
        case 'submitted':
          return <Badge className="bg-green-500 text-white">Submitted</Badge>
        case 'failed':
          return <Badge className="bg-red-500 text-white">Failed</Badge>
        case 'bounced':
          return <Badge className="bg-red-300 text-white">Bounced</Badge>
        case 'clicked':
          return <Badge className="bg-blue-300 text-white">Clicked</Badge>
        case 'complained':
          return <Badge className="bg-red-400 text-white">Complained</Badge>
        case 'delayed':
          return <Badge className="bg-red-200 text-white">Delayed</Badge>
        case 'opened':
          return <Badge className="bg-green-300 text-white">Opened</Badge>
        default:
          return <Badge variant="default">Unknown</Badge>
      }
    }
  },
  {
    header: 'Delete',
    cell: ({ row }) => {
      const { mutate, isPending } = useDeleteInvitation()
      const [open, setOpen] = useState<boolean>(false)

      return (
        <Modal
          open={open}
          onOpenChange={setOpen}
          trigger={
            <Button variant={'destructive'} size={'sm'}>
              Delete
            </Button>
          }
          title="Delete Invitation"
          description={`Delete invitation for ${row.original.invitedName}`}>
          <WarningModalContent
            buttonLabel="Delete"
            message={`Are you sure delete invitation for ${row.original.invitedName}?`}
            handleCancel={setOpen}
            isLoading={isPending}
            handleDelete={() =>
              mutate(
                { invitationId: row.original.invitationId },
                {
                  onSuccess: () => {
                    setOpen(false)
                  }
                }
              )
            }
          />
        </Modal>
      )
    }
  }
]
