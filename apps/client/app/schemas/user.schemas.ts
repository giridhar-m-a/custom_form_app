import * as z from 'zod'
import { ImageSchema } from './image.schemas'

export const UserUpdateSchema = z.object({
  userFullName: z
    .string()
    .min(2, { message: 'Name must be at least 2 characters long' })
    .max(100, { message: 'Name must be at most 100 characters long' })
})

export const UserProfileSchema = z.object({
  file: ImageSchema
})

export const ChangePasswordSchema = z
  .object({
    oldPassword: z.string().min(8, { message: 'Password must be at least 8 characters long' }),
    userPassword: z.string().min(8, { message: 'Password must be at least 8 characters long' }),
    userVerifyPassword: z.string()
  })
  .refine(data => data.userPassword === data.userVerifyPassword, {
    message: 'Passwords do not match',
    path: ['userVerifyPassword']
  })

export type ChangePasswordSchemaType = z.infer<typeof ChangePasswordSchema>
export type UserProfileSchemaType = z.infer<typeof UserProfileSchema>
export type UserUpdateSchemaType = z.infer<typeof UserUpdateSchema>
