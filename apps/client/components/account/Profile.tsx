'use client'

import { UserUpdateSchema, UserUpdateSchemaType } from '@/app/schemas/user.schemas'
import { useUpdateUser } from '@/hooks/queryHooks/useUsers'
import { getMe } from '@/store/slices/me.slice'
import { zodResolver } from '@hookform/resolvers/zod'
import { Mail, Save, User } from 'lucide-react'
import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { useSelector } from 'react-redux'
import { SubmitButton } from '../common/SubmitButton'
import { Badge } from '../ui/badge'
import { Card } from '../ui/card'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import { Input } from '../ui/input'
import { HookFormDevTool } from '../wrapper/HookFormDevTool'
import { ProfilePic } from './ProfilePic'

export const Profile = () => {
  const profile = useSelector(getMe)
  const { mutate, isPending } = useUpdateUser()

  const form = useForm<UserUpdateSchemaType>({
    resolver: zodResolver(UserUpdateSchema),
    mode: 'onChange',
    defaultValues: {
      userFullName: profile.fullName
    }
  })

  const { handleSubmit, control, setValue } = form

  useEffect(() => {
    setValue('userFullName', profile.fullName)
  }, [profile.fullName, setValue])

  const onSubmit = handleSubmit(async data => {
    mutate(data)
  })

  return (
    <Card className="rounded-2xl py-0 shadow-sm border overflow-hidden basis-1/2">
      <div className="p-6 sm:p-8">
        <div className="flex flex-col items-center gap-8">
          {/* Row 1: Profile Picture Centered */}
          <ProfilePic user={profile} />

          {/* Row 2: Name and Email in separate rows */}
          <Form {...form}>
            <HookFormDevTool control={control} />
            <div className="w-full space-y-5">
              <FormField
                control={control}
                name="userFullName"
                render={({ field }) => (
                  <div>
                    <FormItem>
                      <FormLabel htmlFor="oldPassword" className="block text-sm font-semibold mb-1">
                        Full Name
                      </FormLabel>
                      <FormControl>
                        <div className="relative">
                          <User className="absolute left-3 top-1/2 -translate-y-1/2" size={18} />
                          <Input
                            type="text"
                            {...field}
                            className="w-full pl-10 pr-4 py-2.5 "
                            placeholder="Enter your name"
                          />
                        </div>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  </div>
                )}
              />

              <div>
                <label className=" text-sm font-semibold mb-1 flex items-center gap-2">
                  Email Address
                  <Badge className="text-[10px] uppercase tracking-wider px-1.5 py-0.5 rounded">Read Only</Badge>
                </label>
                <div className="relative">
                  <Mail
                    className="absolute left-3 top-1/2 -translate-y-1/2"
                    size={18}
                    style={{ color: 'oklch(0.556 0 0)' }}
                  />
                  <Input
                    type="email"
                    value={profile.email}
                    disabled
                    className="w-full pl-10 pr-4 py-2.5 cursor-not-allowed"
                  />
                </div>
              </div>
            </div>
          </Form>
        </div>

        <div className="mt-8 pt-6 border-t flex justify-end">
          <SubmitButton
            className="flex items-center gap-2 px-6 py-2.5 font-semibold h-12"
            disabled={isPending}
            onClick={onSubmit}
            isLoading={isPending}>
            <Save size={18} />
            Save Profile
          </SubmitButton>
        </div>
      </div>
    </Card>
  )
}
