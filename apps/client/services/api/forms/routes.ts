'use server'

import { CreateFormSchemaType } from '@/app/schemas/form.schemas'
import { GET, POST } from '@/lib/api.config'
import { formsRoutes } from '@/lib/constants/apiRoutes/forms.routes'
import { errorHandler } from '@/lib/errorHandler'
import { FormFilter, FormType } from '@/types/form.types'

export const getForms = async ({ query }: { query?: FormFilter }) => {
  try {
    const res = await GET<FormType[]>(formsRoutes.base, { params: query })
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<FormType[]>(e)
  }
}

export const getFormById = async ({ id }: { id: string }) => {
  try {
    const res = await GET<FormType>(`${formsRoutes.base}/${id}`)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<FormType>(e)
  }
}

export const createForm = async (data: CreateFormSchemaType) => {
  try {
    const res = await POST<FormType>(formsRoutes.base, data)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<FormType>(e)
  }
}
