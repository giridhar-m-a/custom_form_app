import * as z from 'zod'

export const CreateFormSchema = z
  .object({
    title: z.string().min(3, 'Title must be at least 3 characters long'),
    description: z
      .string()
      .min(10, 'Description must be at least 10 characters long')
      .max(100, 'Description must be at most 100 characters long'),
    isScheduled: z.boolean(),
    scheduledTime: z.iso.datetime().optional(),
    closingTime: z.iso.datetime().optional()
  })
  .superRefine((data, ctx) => {
    if (data.isScheduled && !data.scheduledTime) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Scheduled time is required when form is scheduled',
        path: ['scheduledTime']
      })
    }
    if (!data.isScheduled && data.scheduledTime) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Scheduled time is not allowed when form is not scheduled',
        path: ['scheduledTime']
      })
    }
    if (data.scheduledTime && new Date(data.scheduledTime) <= new Date()) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Scheduled time must be in the future',
        path: ['scheduledTime']
      })
    }
    if (data.closingTime && new Date(data.closingTime) <= new Date()) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Closing time must be in the future',
        path: ['closingTime']
      })
    }
  })

export type CreateFormSchemaType = z.infer<typeof CreateFormSchema>

export const FormFieldOptionSchema = z.object({
  optionLabel: z.string().min(3, 'Option label must be at least 3 characters long'),
  ordering: z.number().default(0),
  isAnswer: z.boolean().default(false),
  fieldId: z.string().optional(),
  optionId: z.string().optional()
})

export type FormFieldOptionSchemaType = z.infer<typeof FormFieldOptionSchema>

export const FormFieldCreateSchema = z
  .object({
    fieldLabel: z.string().min(3, 'Field label must be at least 3 characters long'),
    fieldType: z
      .enum([
        'text',
        'number',
        'date',
        'time',
        'datetime',
        'email',
        'phone',
        'url',
        'file',
        'image',
        'video',
        'audio',
        'checkbox',
        'radio',
        'dropdown',
        'multiselect',
        'rating',
        'slider',
        'color',
        'textArea'
      ])
      .default('text'),
    isRequired: z.boolean().default(false),
    ordering: z.number().default(0),
    options: z.array(FormFieldOptionSchema).default([]),
    fieldId: z.string().optional()
  })
  .superRefine((data, ctx) => {
    if (
      data.fieldType === 'checkbox' ||
      data.fieldType === 'radio' ||
      data.fieldType === 'dropdown' ||
      data.fieldType === 'multiselect'
    ) {
      if (data.options.length === 0) {
        ctx.addIssue({
          code: 'custom',
          message: 'Options are required for this field type',
          path: ['options']
        })
      }
    }

    // Validate correct answer count
    if (data.fieldType === 'radio' || data.fieldType === 'dropdown' || data.fieldType === 'checkbox') {
      const answerCount = data.options.filter(opt => opt.isAnswer).length
      if (answerCount > 1) {
        ctx.addIssue({
          code: 'custom',
          message: 'Only one correct answer is allowed for this field type',
          path: ['root', 'options']
        })
      }
    }
  })

export type FormFieldCreateSchemaType = z.infer<typeof FormFieldCreateSchema>

export const FormFieldSchema = z.object({
  formId: z.string(),
  formFields: z
    .array(FormFieldCreateSchema)
    .default([])
    .refine(data => data.length > 0, {
      message: 'At least one field is required'
    }),
  removedFields: z.array(z.string()).default([]),
  removedFieldOptions: z.array(z.string()).default([])
})

export type FormFieldSchemaType = z.infer<typeof FormFieldSchema>
