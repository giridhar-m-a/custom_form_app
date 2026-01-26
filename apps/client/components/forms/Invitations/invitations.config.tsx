import { Badge } from '@/components/ui/badge'
import { Invitation } from '@/types/invitations.types'
import { ColumnDef } from '@tanstack/react-table'

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
        case 'invited':
          return <Badge className="bg-blue-500 text-white">Invited</Badge>
        case 'submitted':
          return <Badge className="bg-green-500 text-white">Submitted</Badge>
        case 'failed':
          return <Badge className="bg-red-500 text-white">Failed</Badge>
        default:
          return <Badge variant="default">Unknown</Badge>
      }
    }
  }
]
