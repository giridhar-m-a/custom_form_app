import { Button } from '@/components/ui/button'
import { ResponseList } from '@/types/responses.types'
import { ColumnDef } from '@tanstack/react-table'
import { format } from 'date-fns'
import Link from 'next/link'
import { MdOpenInNew } from 'react-icons/md'

export const ResponseColumn = ({
  page,
  count,
  formId
}: {
  page: number
  count: number
  formId: string
}): ColumnDef<ResponseList>[] => [
  {
    header: 'SI No.',
    accessorKey: 'invitationId',
    cell: ({ row }) => row.index + (page - 1) * count + 1
  },
  {
    accessorKey: 'invitedEmail',
    header: 'Email'
  },
  {
    accessorKey: 'invitedName',
    header: 'Name'
  },
  {
    accessorKey: 'submittedAt',
    header: 'Submitted At',
    cell: ({ row }) => format(new Date(row.original.submittedAt), 'dd/MM/yyyy')
  },
  {
    accessorKey: 'submissionId',
    header: 'Open',
    cell: ({ row }) => (
      <Link href={`/forms/${formId}/reports/${row.original.submissionId}`}>
        <Button variant={'link'}>
          <MdOpenInNew />
        </Button>
      </Link>
    )
  }
]
