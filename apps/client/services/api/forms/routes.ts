'use server'

import { CreateFormSchemaType, FormFieldCreateSchemaType, FormFieldSchemaType } from '@/app/schemas/form.schemas'
import { DELETE, GET, PATCH, POST } from '@/lib/api.config'
import { formsRoutes } from '@/lib/constants/apiRoutes/forms.routes'
import { errorHandler } from '@/lib/errorHandler'
import { FormField, FormFilter, FormType, FormUpdateType } from '@/types/form.types'

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

export const updateForm = async ({ id, data }: { id: string; data: FormUpdateType }) => {
  try {
    const res = await PATCH<FormType>(`${formsRoutes.base}/${id}`, data)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<FormType>(e)
  }
}

export const deleteForm = async ({ id }: { id: string }) => {
  try {
    const res = await DELETE<FormType>(`${formsRoutes.base}/${id}`)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<FormType>(e)
  }
}

export const createFormField = async ({ data }: { data: FormFieldSchemaType }) => {
  try {
    const payload = {
      ...data,
      removedFields: undefined,
      removedFieldOptions: undefined
    }
    const res = await POST<FormField[]>(`${formsRoutes.fields}`, payload)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<FormField[]>(e)
  }
}

export const updateFormField = async ({ data }: { data: FormFieldSchemaType }) => {
  try {
    const res = await PATCH<FormField[]>(`${formsRoutes.fields}`, data)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<FormField[]>(e)
  }
}

export const getFormFields = async ({ id }: { id: string }) => {
  try {
    const res = await GET<FormField[]>(`${formsRoutes.fields}/${id}`)
    return res
  } catch (e) {
    console.error(e)
    return errorHandler<FormField[]>(e)
  }
}
