export interface ApiResponse<T> {
  status: number
  message: string
  data?: T
  pagination?: Pagination
}

export interface Pagination {
  totalRecords: number
  page: number
  limit: number
  totalPages: number
}
