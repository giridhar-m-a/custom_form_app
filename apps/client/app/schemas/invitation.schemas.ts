import * as z from 'zod'

export const InvitationSchema = z.object({
  email: z.email({ error: 'Email is required' }),
  form_id: z.uuid(),
  name: z
    .string({ error: 'Name is required' })
    .min(3, 'Name must be at least 3 characters long')
    .max(100, 'Name must be at most 100 characters long')
})

export type InvitationType = z.infer<typeof InvitationSchema>
