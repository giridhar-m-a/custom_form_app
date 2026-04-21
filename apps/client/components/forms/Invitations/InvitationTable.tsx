'use client'

import { CommonSelect } from '@/components/common/CommonSelect'
import { DataTable } from '@/components/common/DataTable'
import { Modal } from '@/components/common/Modal'
import { Search } from '@/components/common/Search'
import { Button } from '@/components/ui/button'
import { useInvitations } from '@/hooks/queryHooks/useInvitations'
import { Pagination } from '@/types/api.types'
import { status } from '@/types/form.types'
import { InvitationFilter, InvitationStatus } from '@/types/invitations.types'
import { Mail } from 'lucide-react'
import { useMemo, useState } from 'react'
import { RxReset } from 'react-icons/rx'
import BulkInvitationForm from './BulkInvitationForm'
import { InvitationForm } from './InvitationForm'
import { InvitationsColumn } from './invitations.config'

interface InvitationTableProps {
  formId: string
  status: status
}

export const InvitationTable = ({ formId, status }: InvitationTableProps) => {
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
    { value: 'pending', label: 'Pending' },
    { value: 'delivered', label: 'Delivered' },
    { value: 'submitted', label: 'Submitted' },
    { value: 'bounced', label: 'Bounced' },
    { value: 'clicked', label: 'Clicked' },
    { value: 'complained', label: 'Complained' },
    { value: 'delayed', label: 'Delayed' },
    { value: 'opened', label: 'Opened' },
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
          {status !== 'closed' && (
            <Modal
              description="Invite new users to fill the form"
              title="Invite Users"
              open={bulkInviteOpen}
              onOpenChange={setBulkInviteOpen}
              trigger={<Button variant={'outline'}>Bulk Invite</Button>}>
              <BulkInvitationForm formId={formId} setInviteOpen={setBulkInviteOpen} />
            </Modal>
          )}
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
