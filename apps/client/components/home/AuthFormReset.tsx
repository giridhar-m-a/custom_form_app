'use client'
import { RequestPasswordResetSchema, RequestPasswordResetSchemaType } from '@/app/schemas/auth.schemas'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { useRequestPasswordReset } from '@/hooks/queryHooks/useAuth'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2, Mail } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'

const AuthFormReset = () => {
  const { mutate: requestPasswordReset, isPending } = useRequestPasswordReset()

  const form = useForm<RequestPasswordResetSchemaType>({
    resolver: zodResolver(RequestPasswordResetSchema),
    defaultValues: {
      email: ''
    }
  })

  const { handleSubmit, control } = form

  const handleFormReset = (data: RequestPasswordResetSchemaType) => {
    requestPasswordReset(data)
  }

  return (
    <Card className="w-full max-w-md border-none shadow-xl bg-white/95 backdrop-blur-sm">
      <CardHeader className="text-center pb-4">
        <h2 className="text-3xl font-bold text-gray-800">Reset Password</h2>
        <p className="text-sm text-gray-500">Enter your email to reset your password</p>
      </CardHeader>

      <CardContent>
        <Form {...form}>
          <form onSubmit={handleSubmit(handleFormReset)} className="space-y-4">
            <div>
              <FormField
                control={control}
                name="email"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel className="text-black" htmlFor="email">
                      Email
                    </FormLabel>
                    <FormControl>
                      <div className="relative mt-1">
                        <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <Input
                          id="email"
                          type="email"
                          placeholder="you@example.com"
                          {...field}
                          required
                          className="pl-10 text-gray-700 border! border-gray-300! bg-white! hover:bg-gray-50!"
                          disabled={isPending}
                        />
                      </div>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </div>
            <Button
              type="submit"
              className="w-full mt-6 bg-indigo-600 hover:bg-indigo-700 transition-all"
              size="lg"
              disabled={isPending}>
              {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : 'Submit'}
            </Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  )
}

export default AuthFormReset
