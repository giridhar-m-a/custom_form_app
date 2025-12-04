import * as z from 'zod'

export const CreateFormSchema = z.object({
  title: z.string().min(3, 'Title must be at least 3 characters long'),
  description: z
    .string()
    .min(10, 'Description must be at least 10 characters long')
    .max(100, 'Description must be at most 100 characters long')
})

export type CreateFormSchemaType = z.infer<typeof CreateFormSchema>
