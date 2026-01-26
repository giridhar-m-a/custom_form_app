'use client'
import { ChevronRight, Lock } from 'lucide-react'
import { useForm } from 'react-hook-form'

import { zodResolver } from '@hookform/resolvers/zod'
import { ChangePasswordSchema, ChangePasswordSchemaType } from '@/app/schemas/user.schemas'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import { Input } from '../ui/input'
import { SubmitButton } from '../common/SubmitButton'
import { Card } from '../ui/card'
import { Separator } from '../ui/separator'
import { HookFormDevTool } from '../wrapper/HookFormDevTool'
import { useUpdatePassword } from '@/hooks/queryHooks/useUsers'

export default function PasswordChange() {
  const { isPending, mutate } = useUpdatePassword()
  const form = useForm<ChangePasswordSchemaType>({
    resolver: zodResolver(ChangePasswordSchema),
    defaultValues: {
      oldPassword: '',
      userPassword: '',
      userVerifyPassword: ''
    },
    mode: 'onChange'
  })
  const { handleSubmit, control, reset } = form
  const onSubmit = handleSubmit(async data => {
    mutate(data, {
      onSuccess: () => {
        reset()
      }
    })
  })

  return (
    <Card className="rounded-2xl py-0 shadow-sm border overflow-hidden w-full basis-1/2">
      <div className="p-6 sm:p-8">
        <div className="flex items-center gap-3 mb-6">
          <div className="p-2 rounded-lg">
            <Lock size={20} />
          </div>
          <h2 className="text-xl font-bold">Change Password</h2>
        </div>

        <Form {...form}>
          <HookFormDevTool control={control} />
          <form className="space-y-4">
            <FormField
              control={control}
              name="oldPassword"
              render={({ field }) => (
                <div>
                  <FormItem>
                    <FormLabel htmlFor="oldPassword" className="block text-sm font-semibold mb-1">
                      Current Password
                    </FormLabel>
                    <FormControl>
                      <Input type="password" id="oldPassword" className="h-10" placeholder="••••••••" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                </div>
              )}
            />

            <FormField
              control={control}
              name="userPassword"
              render={({ field }) => (
                <div>
                  <FormItem>
                    <FormLabel htmlFor="userPassword" className="block text-sm font-semibold mb-1">
                      New Password
                    </FormLabel>
                    <FormControl>
                      <Input type="password" id="userPassword" className="h-10" placeholder="••••••••" {...field} />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                </div>
              )}
            />

            <FormField
              control={control}
              name="userVerifyPassword"
              render={({ field }) => (
                <div>
                  <FormItem>
                    <FormLabel htmlFor="userVerifyPassword" className="block text-sm font-semibold mb-1">
                      Confirm New Password
                    </FormLabel>
                    <FormControl>
                      <Input
                        type="password"
                        id="userVerifyPassword"
                        className="h-10"
                        placeholder="••••••••"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                </div>
              )}
            />
            <Separator className="mt-9" />
            <div className="pt-4 flex justify-end">
              <SubmitButton
                type="submit"
                onClick={onSubmit}
                className="h-12"
                disabled={isPending}
                isLoading={isPending}>
                Update Password
              </SubmitButton>
            </div>
          </form>
        </Form>
      </div>
    </Card>
  )
}
