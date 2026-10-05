'use client'

import { ColumnDef, flexRender, getCoreRowModel, useReactTable } from '@tanstack/react-table'

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table'
import { useState } from 'react'
import { Button } from '../ui/button'
import { Grid2X2, List } from 'lucide-react'
import { CommonSelect } from './CommonSelect'
import { Pagination } from './Pagination'
import { ScrollArea } from '../ui/scroll-area'
import { CustomLoader } from './CustomLoader'

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[]
  data: TData[]
  headerTemplate?: React.ReactNode
  footerTemplate?: React.ReactNode
  gridContentTemplate?: (row: TData) => React.ReactNode
  pageSize?: number
  currentPage?: number
  handlePageChange?: (page: number) => void
  handlePageSizeChange?: (size: number) => void
  totalPage?: number
  totalRecords?: number
  isLoading?: boolean
  maxHeight?: string
  minHeight?: string
}

export function DataTable<TData, TValue>({
  columns,
  data,
  headerTemplate,
  gridContentTemplate,
  footerTemplate,
  currentPage,
  handlePageChange,
  handlePageSizeChange,
  pageSize,
  totalPage,
  totalRecords,
  isLoading,
  maxHeight,
  minHeight
}: DataTableProps<TData, TValue>) {
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel()
  })

  const [isGrid, setGrid] = useState(false)

  const pageSizeOptions = [
    {
      label: '15',
      value: '15'
    },
    {
      label: '20',
      value: '20'
    },
    {
      label: '30',
      value: '30'
    },
    {
      label: '40',
      value: '40'
    },
    {
      label: '50',
      value: '50'
    }
  ]

  return (
    <div className="overflow-hidden rounded-md border bg-background w-full">
      {(gridContentTemplate || headerTemplate) && (
        <div className="p-6 flex items-center">
          {gridContentTemplate && (
            <div className="basis-[4%]">
              <Button onClick={() => setGrid(!isGrid)}>{!isGrid ? <Grid2X2 /> : <List />}</Button>
            </div>
          )}
          <div className="basis-[96%]">{headerTemplate}</div>
        </div>
      )}
      {!isLoading ? (
        <>
          <div
            className="relative overflow-auto"
            style={{ maxHeight: maxHeight || 'calc(100vh - 50vh)', minHeight: minHeight || 'calc(100vh - 50vh)' }}>
            {!isGrid && (
              <Table>
                <TableHeader className="sticky top-0 z-10 bg-accent border-background border-2">
                  {table.getHeaderGroups().map(headerGroup => (
                    <TableRow key={headerGroup.id}>
                      {headerGroup.headers.map(header => {
                        return (
                          <TableHead key={header.id} className="px-6">
                            {header.isPlaceholder
                              ? null
                              : flexRender(header.column.columnDef.header, header.getContext())}
                          </TableHead>
                        )
                      })}
                    </TableRow>
                  ))}
                </TableHeader>
                <TableBody className="w-full">
                  {table.getRowModel().rows?.length ? (
                    table.getRowModel().rows.map((row, i) => (
                      <TableRow
                        key={row.id}
                        className={i % 2 !== 0 ? 'bg-accent/30' : ''}
                        data-state={row.getIsSelected() && 'selected'}>
                        {row.getVisibleCells().map(cell => (
                          <TableCell className="px-6" key={cell.id}>
                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </TableCell>
                        ))}
                      </TableRow>
                    ))
                  ) : (
                    <TableRow>
                      <TableCell colSpan={columns.length} className="h-24 text-center">
                        No results.
                      </TableCell>
                    </TableRow>
                  )}
                </TableBody>
              </Table>
            )}
            {isGrid && gridContentTemplate && (
              <div className="w-full h-full flex items-center justify-start gap-6 flex-wrap p-12">
                {table.getRowModel().rows?.length ? (
                  table.getRowModel().rows.map(row => gridContentTemplate(row.original))
                ) : (
                  <div className="p-6 flex items-center justify-center w-full h-full">No results.</div>
                )}
              </div>
            )}
          </div>
          {!footerTemplate && (
            <div className="p-6 flex items-center justify-between gap-4">
              {/* Page size selector */}
              <div className="flex items-center gap-2">
                <span className="text-sm text-muted-foreground whitespace-nowrap">Rows per page:</span>
                <CommonSelect
                  value={pageSize?.toString() || '15'}
                  options={pageSizeOptions}
                  onChange={value => {
                    if (!isNaN(Number(value))) handlePageSizeChange?.(Number(value))
                  }}
                  placeholder="Select Page Size"
                />
              </div>

              {/* Pagination controls */}
              {totalPage && totalPage > 1 && (
                <Pagination
                  currentPage={currentPage || 1}
                  totalPages={totalPage}
                  onPageChange={page => handlePageChange?.(page)}
                  totalRecords={totalRecords}
                  pageSize={pageSize}
                />
              )}
            </div>
          )}
          {footerTemplate && footerTemplate}
        </>
      ) : (
        <div className="flex items-center justify-center min-h-[calc(100vh-60vh)] h-full w-full">
          <CustomLoader />
        </div>
      )}
    </div>
  )
}
