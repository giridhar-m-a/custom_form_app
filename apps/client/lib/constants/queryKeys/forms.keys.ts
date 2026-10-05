import { FormFilter } from '@/types/form.types'

export const formsKeys = {
  list: ({ query }: { query?: FormFilter }) => ['form', 'list', JSON.stringify(query)],
  detail: (id: string) => ['form', 'detail', id],
  fields: (id: string) => ['form', 'fields', id]
}
