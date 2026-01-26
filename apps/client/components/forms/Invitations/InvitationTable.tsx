'use client'

import { DataTable } from '@/components/common/DataTable'
import { useInvitations } from '@/hooks/queryHooks/useInvitations'
import { Invitation, InvitationFilter, InvitationStatus } from '@/types/invitations.types'
import { useMemo, useState } from 'react'
import { InvitationsColumn } from './invitations.config'
import { Pagination } from '@/types/api.types'
import { Search } from '@/components/common/Search'
import { CommonSelect } from '@/components/common/CommonSelect'
import { Button } from '@/components/ui/button'
import { RxReset } from 'react-icons/rx'
import { Modal } from '@/components/common/Modal'
import { InvitationForm } from './InvitationForm'
import { Mail } from 'lucide-react'
import BulkInvitationForm from './BulkInvitationForm'

interface InvitationTableProps {
  formId: string
}

export const InvitationTable = ({ formId }: InvitationTableProps) => {
  const [params, setParams] = useState<InvitationFilter>({
    page: 1,
    limit: 15,
    search: ''
  })
  const [inviteOpen, setInviteOpen] = useState(false)
  const [bulkInviteOpen, setBulkInviteOpen] = useState(false)
  const { data, isFetching } = useInvitations({ formId, params })

  const pagination = useMemo<Pagination>(
    () => ({
      page: data?.pagination?.page || 1,
      limit: !!(data?.pagination?.limit && data?.pagination?.limit < 15)
        ? 15
        : !data?.pagination?.limit
          ? 15
          : data?.pagination?.limit,
      totalRecords: data?.pagination?.totalRecords || 0,
      totalPages: data?.pagination?.totalPages || 1
    }),
    [data]
  )

  const statusOptions: { value: InvitationStatus; label: string }[] = [
    { value: 'invited', label: 'Invited' },
    { value: 'pending', label: 'Pending' },
    { value: 'submitted', label: 'Submitted' },
    { value: 'failed', label: 'Failed' }
  ]

  const resetFilters = () => {
    setParams({})
  }

  const handleStatusChange = (value: InvitationStatus) => {
    setParams({ ...params, status: value, page: 1 })
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1>Invitations</h1>
        <div className="flex items-center gap-2 justify-end">
          <Modal
            description="Invite new users to fill the form"
            title="Invite Users"
            open={inviteOpen}
            onOpenChange={setInviteOpen}
            trigger={
              <Button variant={'outline'} size={'icon'}>
                <Mail />
              </Button>
            }>
            <InvitationForm formId={formId} setInviteOpen={setInviteOpen} />
          </Modal>
          <Modal
            description="Invite new users to fill the form"
            title="Invite Users"
            open={bulkInviteOpen}
            onOpenChange={setBulkInviteOpen}
            trigger={<Button variant={'outline'}>Bulk Invite</Button>}>
            <BulkInvitationForm formId={formId} setInviteOpen={setBulkInviteOpen} />
          </Modal>
        </div>
      </div>
      <DataTable
        isLoading={isFetching}
        columns={InvitationsColumn({ page: pagination.page, count: pagination.limit })}
        data={data?.data || []}
        currentPage={pagination.page}
        pageSize={pagination.limit}
        totalPage={pagination.totalPages}
        totalRecords={pagination.totalRecords}
        handlePageChange={page => setParams({ ...params, page })}
        handlePageSizeChange={size => setParams({ ...params, limit: size })}
        maxHeight="calc(100vh - 65vh)"
        minHeight="calc(100vh - 65vh)"
        headerTemplate={
          <div className="flex items-center justify-evenly gap-4">
            <Search
              placeholder="Search"
              value={params.search}
              onChange={e => setParams({ ...params, search: e.target.value, page: 1 })}
            />
            <CommonSelect
              options={statusOptions}
              placeholder="Status"
              value={params.status || ''}
              onChange={value => handleStatusChange(value as InvitationStatus)}
            />
            <Button onClick={resetFilters} variant="outline">
              <RxReset />
            </Button>
          </div>
        }
      />
    </div>
  )
}
