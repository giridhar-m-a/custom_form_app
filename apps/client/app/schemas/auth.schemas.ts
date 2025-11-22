import * as z from 'zod'

export const SignInSchema = z.object({
  email: z.email('Invalid email address').min(1, 'Email is required'),
  password: z.string().min(1, 'Password is required')
})

export type SignInSchemaType = z.infer<typeof SignInSchema>

export const SignUpSchema = z
  .object({
    email: z.email('Invalid email address').min(1, 'Email is required'),
    password: z.string().min(8, 'Password must be at least 8 characters long'),
    confirmPassword: z.string().min(1, 'Please confirm your password'),
    name: z.string().min(1, 'Full name is required')
  })
  .refine(data => data.password === data.confirmPassword, {
    message: 'Passwords do not match',
    path: ['confirmPassword']
  })

export type SignUpSchemaType = z.infer<typeof SignUpSchema>
