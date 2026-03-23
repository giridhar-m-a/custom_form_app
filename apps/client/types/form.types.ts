export interface FormField {
  fieldId: string
  fieldLabel: string
  fieldType: FieldType
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
  isAnswer?: boolean
}

export interface FormFieldOptionRecord {
  [key: string]: FieldOption
}

export interface FormUpdateType {
  title?: string
  access?: access
  description?: string
  status?: status
}

export type FieldType =
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
  access: access
  createdAt: string
  createdBy: string
  description: string
  id: string
  status: status
  title: string
  updatedAt: string
  closingTime?: Date
  scheduledTime?: Date
  isScheduled: boolean
  invitationScheduleGap?: number
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
