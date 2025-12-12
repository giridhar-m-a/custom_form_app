'use client'
import { useMemo, useState } from 'react'
import { DataTable } from '../common/DataTable'
import { Search } from '../common/Search'
import { FormTableDef } from './forms.config'
import { FormFilter, access, sort, status } from '@/types/form.types'
import { CommonSelect } from '../common/CommonSelect'
import { Button } from '../ui/button'
import { RxReset } from 'react-icons/rx'
import { useGetForms } from '@/hooks/queryHooks/useFormApp'
import { Pagination } from '@/types/api.types'
import { FormItemCard } from './FormItemCard'

export const FormPage = () => {
  const [filters, setFilters] = useState<FormFilter>({
    search: '',
    sort: undefined,
    page: 1,
    limit: 15
  })
  const { data, isFetching } = useGetForms(filters)
  const forms = useMemo(() => data?.data || [], [data])
  const pagination = useMemo<Pagination>(
    () =>
      data?.pagination || {
        page: 1,
        limit: 15,
        totalRecords: 0,
        totalPages: 1
      },
    [data]
  )
  const sortOptions = [
    { value: '-updated', label: 'Oldest First' },
    { value: 'updated', label: 'Newest First' },
    { value: 'title', label: 'Title A to Z' },
    { value: '-title', label: 'Title Z to A' }
  ]

  const accessOptions = [
    { value: 'restricted', label: 'Restricted' },
    { value: 'public', label: 'Public' }
  ]

  const statusOptions = [
    { value: 'draft', label: 'Draft' },
    { value: 'published', label: 'Published' },
    { value: 'archived', label: 'Archived' },
    { value: 'closed', label: 'Closed' }
  ]

  const resetFilters = () => {
    setFilters({
      search: '',
      sort: undefined,
      page: 1,
      limit: 10
    })
  }
  return (
    <DataTable
      columns={FormTableDef(filters.page || 1, filters.limit || 15)}
      data={forms}
      currentPage={pagination.page}
      pageSize={pagination.limit}
      totalRecords={pagination.totalRecords}
      totalPage={pagination.totalPages}
      handlePageChange={page => setFilters({ ...filters, page })}
      handlePageSizeChange={size => setFilters({ ...filters, limit: size, page: 1 })}
      headerTemplate={
        <div className="flex items-center justify-evenly gap-4">
          <Search
            placeholder="Search"
            value={filters.search}
            onChange={e => setFilters({ ...filters, search: e.target.value })}
          />
          <CommonSelect
            options={sortOptions}
            placeholder="Sort"
            value={filters.sort || ''}
            onChange={value => setFilters({ ...filters, sort: value as sort })}
          />
          <CommonSelect
            options={accessOptions}
            placeholder="Access"
            value={filters.access || ''}
            onChange={value => setFilters({ ...filters, access: value as access })}
          />
          <CommonSelect
            options={statusOptions}
            placeholder="Status"
            value={filters.status || ''}
            onChange={value => setFilters({ ...filters, status: value as status })}
          />
          <Button onClick={resetFilters} variant="outline">
            <RxReset />
          </Button>
        </div>
      }
      gridContentTemplate={row => <FormItemCard form={row} />}
    />
  )
}
