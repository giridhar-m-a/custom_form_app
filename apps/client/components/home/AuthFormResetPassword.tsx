'use client'

import { ResetPasswordSchema, ResetPasswordSchemaType } from '@/app/schemas/auth.schemas'
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useResetPassword } from '@/hooks/queryHooks/useAuth'
import { zodResolver } from '@hookform/resolvers/zod'
import { Mail } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { SubmitButton } from '../common/SubmitButton'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import { useEffect } from 'react'
import toast from 'react-hot-toast'

const AuthFormResetPassword = ({ token }: { token?: string }) => {
  const { mutate: resetPassword, isPending: isRegisterLoading } = useResetPassword()

  const form = useForm<ResetPasswordSchemaType>({
    resolver: zodResolver(ResetPasswordSchema),
    defaultValues: {
      token,
      newPassword: '',
      confirmPassword: ''
    }
  })

  const onSubmit = (data: ResetPasswordSchemaType) => {
    if (!data.token) {
      toast.error('Token is missing')
      return
    }
    resetPassword(data)
  }

  const { control, handleSubmit } = form

  return (
    <Card className="w-full max-w-md border-none shadow-xl bg-white/95 backdrop-blur-sm">
      <CardHeader className="text-center pb-4">
        <h2 className="text-3xl font-bold text-gray-800">Reset Password</h2>
        <p className="text-sm text-gray-500">Enter your new password and confirm it.</p>
      </CardHeader>

      <CardContent>
        <Form {...form}>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={control}
              name="newPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-black" htmlFor="password">
                    Password
                  </FormLabel>
                  <FormControl>
                    <div className="relative mt-1">
                      <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                      <Input
                        id="password"
                        type="password"
                        placeholder="••••••••"
                        {...field}
                        required
                        className="pl-10 text-gray-700 border! border-gray-300! bg-white! hover:bg-gray-50!"
                        disabled={isRegisterLoading}
                      />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={control}
              name="confirmPassword"
              render={({ field }) => (
                <FormItem>
                  <FormLabel className="text-black" htmlFor="confirmPassword">
                    Confirm Password
                  </FormLabel>
                  <FormControl>
                    <div className="relative mt-1">
                      <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                      <Input
                        id="confirmPassword"
                        type="password"
                        placeholder="••••••••"
                        {...field}
                        required
                        className="pl-10 text-gray-700 border! border-gray-300! bg-white! hover:bg-gray-50!"
                        disabled={isRegisterLoading}
                      />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </form>
        </Form>
      </CardContent>
      <CardFooter className="flex flex-col space-y-4">
        <SubmitButton
          className="w-full"
          disabled={isRegisterLoading}
          isLoading={isRegisterLoading}
          onClick={handleSubmit(onSubmit)}>
          Reset Password
        </SubmitButton>
      </CardFooter>
    </Card>
  )
}

export default AuthFormResetPassword
