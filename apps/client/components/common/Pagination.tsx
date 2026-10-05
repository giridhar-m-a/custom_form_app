'use client'

import { Button } from '@/components/ui/button'
import { ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight } from 'lucide-react'
import { cn } from '@/lib/utils'

interface PaginationProps {
  currentPage: number
  totalPages: number
  onPageChange: (page: number) => void
  totalRecords?: number
  pageSize?: number
  showPageInfo?: boolean
}

export function Pagination({
  currentPage,
  totalPages,
  onPageChange,
  totalRecords,
  pageSize,
  showPageInfo = true
}: PaginationProps) {
  // Generate page numbers to display
  const getPageNumbers = () => {
    const pages: (number | 'ellipsis')[] = []
    const maxVisible = 7 // Maximum number of page buttons to show

    if (totalPages <= maxVisible) {
      // Show all pages if total is less than max visible
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i)
      }
    } else {
      // Always show first page
      pages.push(1)

      if (currentPage <= 3) {
        // Near the start
        for (let i = 2; i <= 5; i++) {
          pages.push(i)
        }
        pages.push('ellipsis')
        pages.push(totalPages)
      } else if (currentPage >= totalPages - 2) {
        // Near the end
        pages.push('ellipsis')
        for (let i = totalPages - 4; i <= totalPages; i++) {
          pages.push(i)
        }
      } else {
        // In the middle
        pages.push('ellipsis')
        for (let i = currentPage - 1; i <= currentPage + 1; i++) {
          pages.push(i)
        }
        pages.push('ellipsis')
        pages.push(totalPages)
      }
    }

    return pages
  }

  const pageNumbers = getPageNumbers()

  // Calculate display info
  const startRecord = totalRecords ? (currentPage - 1) * (pageSize || 10) + 1 : 0
  const endRecord = totalRecords ? Math.min(currentPage * (pageSize || 10), totalRecords) : 0

  return (
    <div className="flex items-center justify-between w-full gap-4 flex-wrap">
      {/* Page info */}
      {showPageInfo && totalRecords && (
        <div className="text-sm text-muted-foreground">
          Showing <span className="font-semibold text-foreground">{startRecord}</span> to{' '}
          <span className="font-semibold text-foreground">{endRecord}</span> of{' '}
          <span className="font-semibold text-foreground">{totalRecords}</span> results
        </div>
      )}

      {/* Pagination controls */}
      <div className="flex items-center gap-1">
        {/* First page button */}
        <Button
          variant="outline"
          size="icon"
          onClick={() => onPageChange(1)}
          disabled={currentPage === 1}
          className={cn('h-9 w-9 transition-all duration-200', currentPage === 1 && 'opacity-50 cursor-not-allowed')}
          aria-label="Go to first page">
          <ChevronsLeft className="h-4 w-4" />
        </Button>

        {/* Previous page button */}
        <Button
          variant="outline"
          size="icon"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage === 1}
          className={cn('h-9 w-9 transition-all duration-200', currentPage === 1 && 'opacity-50 cursor-not-allowed')}
          aria-label="Go to previous page">
          <ChevronLeft className="h-4 w-4" />
        </Button>

        {/* Page number buttons - only show if more than 1 page */}
        {totalPages > 1 && (
          <div className="flex items-center gap-1">
            {pageNumbers.map((page, index) => {
              if (page === 'ellipsis') {
                return (
                  <span key={`ellipsis-${index}`} className="px-2 text-muted-foreground">
                    ...
                  </span>
                )
              }

              return (
                <Button
                  key={page}
                  variant={currentPage === page ? 'default' : 'outline'}
                  size="icon"
                  onClick={() => onPageChange(page)}
                  className={cn(
                    'h-9 w-9 transition-all duration-200',
                    currentPage === page && 'bg-primary text-primary-foreground shadow-md hover:bg-primary/90'
                  )}
                  aria-label={`Go to page ${page}`}
                  aria-current={currentPage === page ? 'page' : undefined}>
                  {page}
                </Button>
              )
            })}
          </div>
        )}

        {/* Next page button */}
        <Button
          variant="outline"
          size="icon"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage === totalPages}
          className={cn(
            'h-9 w-9 transition-all duration-200',
            currentPage === totalPages && 'opacity-50 cursor-not-allowed'
          )}
          aria-label="Go to next page">
          <ChevronRight className="h-4 w-4" />
        </Button>

        {/* Last page button */}
        <Button
          variant="outline"
          size="icon"
          onClick={() => onPageChange(totalPages)}
          disabled={currentPage === totalPages}
          className={cn(
            'h-9 w-9 transition-all duration-200',
            currentPage === totalPages && 'opacity-50 cursor-not-allowed'
          )}
          aria-label="Go to last page">
          <ChevronsRight className="h-4 w-4" />
        </Button>
      </div>
    </div>
  )
}
