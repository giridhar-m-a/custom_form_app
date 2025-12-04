export interface FormField {
  fieldId: string
  fieldLabel: string
  fieldType: fieldType
  isRequired: boolean
  ordering: number
  formId: string
  options: FieldOption[]
}

export interface FieldOption {
  optionId: string
  optionLabel: string
  ordering: number
  fieldId: string
}

export interface FormFieldOptionRecord {
  [key: string]: FieldOption
}

type fieldType =
  | 'text'
  | 'number'
  | 'date'
  | 'time'
  | 'datetime'
  | 'email'
  | 'phone'
  | 'url'
  | 'file'
  | 'image'
  | 'video'
  | 'audio'
  | 'checkbox'
  | 'radio'
  | 'dropdown'
  | 'multiselect'
  | 'rating'
  | 'slider'
  | 'color'
  | 'textArea'

export interface FormType {
  access: string
  createdAt: string
  createdBy: string
  description: string
  id: string
  status: string
  title: string
  updatedAt: string
}

export interface FormFilter {
  search?: string
  sort?: sort
  status?: status
  access?: access
  page?: number
  limit?: number
}

export type sort = 'updated' | '-updated' | 'title' | '-title'

export type status = 'draft' | 'published' | 'archived' | 'closed'

export type access = 'public' | 'restricted'
