'use client'

import { DataTable } from '@/components/common/DataTable'
import { Search } from '@/components/common/Search'
import { Button } from '@/components/ui/button'
import { useResponses } from '@/hooks/queryHooks/useResponses'
import { Pagination } from '@/types/api.types'
import { ResponseFilter } from '@/types/responses.types'
import { useMemo, useState } from 'react'
import { RxReset } from 'react-icons/rx'
import { ResponseColumn } from './Response.config'

interface ResponseTableProps {
  formId: string
}

export const ResponseTable = ({ formId }: ResponseTableProps) => {
  const [params, setParams] = useState<ResponseFilter>({
    page: 1,
    limit: 15,
    search: ''
  })
  const { data, isFetching } = useResponses({ formId, params })

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

  const resetFilters = () => {
    setParams({ page: pagination.page, limit: pagination.limit, search: '' })
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h1>Responses</h1>
      </div>
      <DataTable
        isLoading={isFetching}
        columns={ResponseColumn({ page: pagination.page, count: pagination.limit, formId })}
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
            <Button onClick={resetFilters} variant="outline">
              <RxReset />
            </Button>
          </div>
        }
      />
    </div>
  )
}
