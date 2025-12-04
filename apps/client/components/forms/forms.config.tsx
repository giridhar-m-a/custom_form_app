import { FormType } from '@/types/form.types'
import { ColumnDef } from '@tanstack/react-table'

export const FormTableDef: ColumnDef<FormType>[] = [
  {
    header: 'S No.',
    cell: ({ row }) => row.index + 1
  },
  {
    header: 'Title',
    accessorKey: 'title'
  },
  {
    header: 'Updated At',
    accessorKey: 'updatedAt'
  },
  {
    header: 'Status',
    accessorKey: 'status'
  },
  {
    header: 'Access',
    accessorKey: 'access'
  }
]
