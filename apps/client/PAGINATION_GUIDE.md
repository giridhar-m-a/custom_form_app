# Pagination Implementation Guide

## Overview
This document describes the pagination implementation for the DataTable component in the custom form application.

## Components Created

### 1. Pagination Component (`/components/common/Pagination.tsx`)

A modern, accessible pagination component with the following features:

#### Features
- **Smart Page Number Display**: Shows a maximum of 7 page buttons with ellipsis for long page lists
- **Navigation Buttons**: 
  - First page (double chevron left)
  - Previous page (single chevron left)
  - Page numbers (with intelligent display logic)
  - Next page (single chevron right)
  - Last page (double chevron right)
- **Page Information**: Displays "Showing X to Y of Z results"
- **Accessibility**: Full ARIA labels and keyboard navigation support
- **Disabled States**: Appropriate button disabling for edge cases
- **Modern Design**: Smooth transitions and hover effects

#### Props
```typescript
interface PaginationProps {
  currentPage: number         // Current active page (1-indexed)
  totalPages: number          // Total number of pages
  onPageChange: (page: number) => void  // Callback when page changes
  totalRecords?: number       // Optional: Total number of records
  pageSize?: number           // Optional: Number of items per page
  showPageInfo?: boolean      // Optional: Show/hide page info (default: true)
}
```

#### Page Number Logic
The component intelligently displays page numbers:
- **7 or fewer pages**: Shows all pages
- **Near the start** (pages 1-3): Shows pages 1, 2, 3, 4, 5, ..., last
- **Near the end**: Shows 1, ..., (last-4), (last-3), (last-2), (last-1), last
- **In the middle**: Shows 1, ..., (current-1), current, (current+1), ..., last

## Updated Components

### 2. DataTable Component (`/components/common/DataTable.tsx`)

#### Changes Made
1. **Added Pagination Import**: Imported the new Pagination component
2. **Updated Footer Layout**: 
   - Changed from flex with basis percentages to `justify-between` layout
   - Added label "Rows per page:" next to the page size selector
   - Integrated Pagination component
3. **Added Props**: Now accepts and uses `totalPage` and `totalRecords` props

#### Updated Footer Structure
```tsx
{!footerTemplate && (
  <div className="p-6 flex items-center justify-between gap-4">
    {/* Page size selector */}
    <div className="flex items-center gap-2">
      <span className="text-sm text-muted-foreground whitespace-nowrap">
        Rows per page:
      </span>
      <CommonSelect ... />
    </div>

    {/* Pagination controls */}
    {totalPage && totalPage > 1 && (
      <Pagination ... />
    )}
  </div>
)}
```

#### Props Interface
```typescript
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
}
```

### 3. FormPage Component (`/components/forms/FormPage.tsx`)

#### Changes Made
1. **Added Pagination Props**: Passes pagination-related props to DataTable
2. **Page Change Handler**: Updates filter state when page changes
3. **Page Size Change Handler**: Updates filter state and resets to page 1 when page size changes
4. **TODO Comments**: Added placeholders for actual API data

#### Updated DataTable Usage
```tsx
<DataTable
  columns={FormTableDef}
  data={[]}
  currentPage={filters.page}
  totalPage={10} // TODO: Replace with actual total pages from API
  totalRecords={100} // TODO: Replace with actual total records from API
  pageSize={filters.limit}
  handlePageChange={page => setFilters({ ...filters, page })}
  handlePageSizeChange={size => setFilters({ ...filters, limit: size, page: 1 })}
  // ... other props
/>
```

## Integration with API

When you integrate with your backend API, you'll need to:

1. **Update the API call** to accept pagination parameters:
   ```typescript
   const response = await api.get('/forms', {
     params: {
       page: filters.page,
       limit: filters.limit,
       search: filters.search,
       sort: filters.sort,
       // ... other filters
     }
   })
   ```

2. **Update the component** to use real data from API response:
   ```tsx
   const [data, setData] = useState([])
   const [totalRecords, setTotalRecords] = useState(0)
   const [totalPages, setTotalPages] = useState(0)

   useEffect(() => {
     fetchData()
   }, [filters])

   const fetchData = async () => {
     const response = await api.get('/forms', { params: filters })
     setData(response.data.items)
     setTotalRecords(response.data.total)
     setTotalPages(response.data.totalPages)
   }
   ```

3. **Replace the hardcoded values** in FormPage:
   ```tsx
   <DataTable
     data={data}
     totalPage={totalPages}
     totalRecords={totalRecords}
     // ... other props
   />
   ```

## Design Features

### Visual Excellence
- **Modern Aesthetics**: Clean, professional design with proper spacing
- **Interactive Elements**: Hover effects and smooth transitions
- **Responsive Layout**: Adapts to different screen sizes with flex-wrap
- **Accessible Design**: Proper ARIA labels and keyboard navigation
- **Visual Feedback**: Clear active state for current page
- **Disabled States**: Grayed-out buttons when navigation is not available

### User Experience
- **Smart Navigation**: Quick jump to first/last page
- **Page Info Display**: Always know where you are in the dataset
- **Flexible Page Size**: Users can choose how many rows to display
- **Automatic Reset**: Changing page size resets to page 1
- **Error Prevention**: Disabled buttons prevent invalid navigation

## Testing Checklist

- [ ] Pagination appears when `totalPage > 1`
- [ ] First/Last buttons work correctly
- [ ] Previous/Next buttons work correctly
- [ ] Direct page number buttons work
- [ ] Ellipsis appears correctly for long page lists
- [ ] Current page is highlighted
- [ ] Disabled states work correctly (first/last page)
- [ ] Page info displays correct record range
- [ ] Page size change resets to page 1
- [ ] Keyboard navigation works
- [ ] Screen reader accessibility works

## Browser Compatibility
The pagination component uses modern React and Lucide icons, which are compatible with:
- Chrome/Edge 90+
- Firefox 88+
- Safari 14+
- Mobile browsers (iOS Safari, Chrome Mobile)

## Dependencies
- `@tanstack/react-table`: For table management
- `lucide-react`: For icons (ChevronLeft, ChevronRight, ChevronsLeft, ChevronsRight)
- `clsx` and `tailwind-merge`: For className utilities (cn function)
- UI components: Button, Table components from shadcn/ui

## Future Enhancements

Potential improvements for future development:
1. **Keyboard Shortcuts**: Add shortcuts for quick navigation (e.g., Ctrl+← for previous page)
2. **Jump to Page**: Input field to jump directly to a specific page
3. **Loading States**: Show skeleton/spinner during data fetch
4. **URL Persistence**: Sync pagination state with URL query parameters
5. **Infinite Scroll**: Alternative pagination mode for mobile/touch devices
6. **Page Size Persistence**: Remember user's preferred page size in localStorage
7. **Customizable Page Ranges**: Allow custom page size options per table
8. **Export Functionality**: Export current page or all pages to CSV/Excel
