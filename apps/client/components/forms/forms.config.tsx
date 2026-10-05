'use client'
import { FormType, access, status } from '@/types/form.types'
import { ColumnDef } from '@tanstack/react-table'
import { FormAction } from './FormAction'
import { format } from 'date-fns'
import { Badge } from '../ui/badge'

export const FormTableDef = (page: number, count: number): ColumnDef<FormType>[] => [
  {
    header: 'S No.',
    cell: ({ row }) => row.index + (page - 1) * count + 1
  },
  {
    header: 'Title',
    accessorKey: 'title'
  },
  {
    header: 'Updated At',
    accessorKey: 'updatedAt',
    cell: ({ row }) => format(row.original.updatedAt, 'dd/MM/yyyy')
  },
  {
    header: 'Status',
    accessorKey: 'status',
    cell: ({ row }) => FormStatusBadge({ status: row.original.status })
  },
  {
    header: 'Access',
    accessorKey: 'access',
    cell: ({ row }) => FormAccessBadge({ access: row.original.access })
  },
  {
    header: 'Action',
    cell: ({ row }) => <FormAction form={row.original} />
  }
]

export const FormStatusBadge = ({ status }: { status: status }) => {
  switch (status) {
    case 'draft':
      return <Badge className="bg-gray-500 text-white">Draft</Badge>
    case 'published':
      return <Badge className="bg-green-500 text-white">Published</Badge>
    case 'archived':
      return <Badge className="bg-red-400 text-white">Archived</Badge>
    case 'closed':
      return <Badge className="bg-red-500 text-white">Closed</Badge>
    default:
      return <Badge className="bg-red-300 text-white">Unknown</Badge>
  }
}

export const FormAccessBadge = ({ access }: { access: access }) => {
  switch (access) {
    case 'restricted':
      return <Badge variant="destructive">Restricted</Badge>
    case 'public':
      return <Badge className="bg-green-500 text-white">Public</Badge>
    default:
      return <Badge className="bg-red-300 text-white">Unknown</Badge>
  }
}
