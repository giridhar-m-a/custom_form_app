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

export const RequestPasswordResetSchema = z.object({
  email: z.email('Invalid email address').min(1, 'Email is required')
})

export type RequestPasswordResetSchemaType = z.infer<typeof RequestPasswordResetSchema>

export const ResetPasswordSchema = z
  .object({
    token: z.string().optional(),
    newPassword: z
      .string()
      .min(8, 'Password must be at least 8 characters long')
      .regex(
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
        'Password must contain at least one uppercase letter, one lowercase letter, one number and one special character'
      ),
    confirmPassword: z
      .string()
      .min(1, 'Please confirm your password')
      .regex(
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$/,
        'Password must contain at least one uppercase letter, one lowercase letter, one number and one special character'
      )
  })
  .superRefine((data, ctx) => {
    if (data.newPassword !== data.confirmPassword) {
      ctx.addIssue({
        code: z.ZodIssueCode.custom,
        message: 'Passwords do not match',
        path: ['confirmPassword']
      })
    }
  })

export type ResetPasswordSchemaType = z.infer<typeof ResetPasswordSchema>

export const TempUserSchema = z.object({
  name: z.string().min(1, 'Name is required')
})

export type TempUserSchemaType = z.infer<typeof TempUserSchema>
