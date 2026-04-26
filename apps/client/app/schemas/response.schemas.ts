import { FieldType } from '@/types/form.types'
import * as z from 'zod'

// -------------------------
// FILE SCHEMA
// -------------------------
export const ResponseFileSchema = z.object({
  fileName: z.string().min(1, 'File name is required'),
  filePath: z.string().min(1, 'File path is required'),
  fileSize: z.number().min(1, 'File size is required'),
  fileType: z.string().min(1, 'File type is required')
})

// -------------------------
// OPTION SCHEMA
// -------------------------
export const ResponseOptionSchema = z.object({
  optionId: z.string().min(1, 'Option ID is required')
})

// -------------------------
// PER-FIELD RESPONSE SCHEMA (dynamic, based on field config)
// -------------------------
export const ResponseItemSchema = ({
  type,
  isRequired,
  isMultiple
}: {
  type: FieldType
  isRequired: boolean
  isMultiple: boolean
}) => {
  console.log('isRequired', isRequired, 'type: ', type, 'isMultiple: ', isMultiple)
  let responseText: z.ZodTypeAny = z.string()
  let responseOptions: z.ZodTypeAny = z.array(ResponseOptionSchema)
  let responseFiles: z.ZodTypeAny = z.array(ResponseFileSchema)

  // TEXT / TEXTAREA
  if (['text', 'textArea'].includes(type)) {
    responseText = isRequired ? z.string().min(1, 'Response is required') : z.string().optional()
  }

  // NUMBER / SLIDER
  if (type === 'number' || type === 'slider') {
    const base = z.string().refine(val => !isNaN(Number(val)), 'Must be a valid number')
    responseText = isRequired ? base.min(1, 'Response is required') : base.optional()
  }

  // EMAIL
  if (type === 'email') {
    const base = z.string().email('Invalid email')
    responseText = isRequired ? base : base.optional()
  }

  // PHONE
  if (type === 'phone') {
    const base = z.string().regex(/^[0-9]{10,15}$/, 'Invalid phone number')
    responseText = isRequired ? base : base.optional()
  }

  // URL
  if (type === 'url') {
    const base = z.string().url('Invalid URL')
    responseText = isRequired ? base : base.optional()
  }

  // DATE / DATETIME
  if (['date', 'datetime'].includes(type)) {
    const base = z.string().refine(val => !isNaN(Date.parse(val)), 'Invalid date')
    responseText = isRequired ? base : base.optional()
  }

  // TIME
  if (type === 'time') {
    const base = z.string().regex(/^([01]\d|2[0-3]):([0-5]\d)$/, 'Invalid time format')
    responseText = isRequired ? base : base.optional()
  }

  // COLOR
  if (type === 'color') {
    const base = z.string().regex(/^#([0-9A-Fa-f]{3}){1,2}$/, 'Invalid color')
    responseText = isRequired ? base : base.optional()
  }

  // RATING
  if (type === 'rating') {
    const base = z.string().refine(val => {
      const num = Number(val)
      return num >= 1 && num <= 5
    }, 'Rating must be between 1 and 5')
    responseText = isRequired ? base : base.optional()
  }

  // OPTIONS (single select)
  if (['radio', 'dropdown'].includes(type)) {
    const base = z.array(ResponseOptionSchema)
    responseOptions = isRequired
      ? base.min(1, 'Selection required').max(1, 'Only one option allowed')
      : base.max(1).optional()
  }

  // OPTIONS (multi select)
  if (['checkbox', 'multiselect'].includes(type)) {
    const base = z.array(ResponseOptionSchema)
    responseOptions = isRequired ? base.min(1, 'Select at least one option') : base.optional()
  }

  // FILES / MEDIA
  if (['file', 'image', 'video', 'audio'].includes(type)) {
    let base = z.array(ResponseFileSchema)

    if (isRequired) {
      base = base.min(1, 'At least one file required')
    }

    if (!isMultiple) {
      base = base.max(1, 'Only one file allowed')
    }
    responseFiles = isRequired ? base : base.optional()
  }

  return z.object({
    formFieldId: z.string().min(1, 'Field ID is required'),
    responseText,
    responseOptions,
    responseFiles
  })
}

// -------------------------
// BASE RESPONSE ITEM SHAPE
// Used for structural validation of the POST body before per-field validation
// -------------------------
export const BaseResponseItemSchema = z.object({
  formFieldId: z.string().min(1, 'Field ID is required'),
  responseText: z.string().optional(),
  responseOptions: z.array(ResponseOptionSchema).optional(),
  responseFiles: z.array(ResponseFileSchema).optional()
})

// -------------------------
// FULL FORM SUBMISSION SCHEMA
// Validates the entire POST request body shape
// -------------------------
export const FormSubmissionSchema = z.object({
  formId: z.string().min(1, 'Form ID is required'),
  respondentId: z.string().min(1, 'Respondent ID is required'),
  responses: z.array(BaseResponseItemSchema).min(1, 'At least one response is required')
})

// -------------------------
// INFERRED TYPES
// -------------------------
export type ResponseFile = z.infer<typeof ResponseFileSchema>
export type ResponseOption = z.infer<typeof ResponseOptionSchema>
export type BaseResponseItem = z.infer<typeof BaseResponseItemSchema>
export type FormSubmission = z.infer<typeof FormSubmissionSchema>
