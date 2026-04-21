'use client'
import { TempUserSchema, TempUserSchemaType } from '@/app/schemas/auth.schemas'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useTempAuth } from '@/hooks/queryHooks/useAuth'
import { zodResolver } from '@hookform/resolvers/zod'
import { Loader2, Mail } from 'lucide-react'
import { useForm } from 'react-hook-form'
import { Modal } from '../common/Modal'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'

const TempUser = () => {
  const { mutate: tempLogin, isPending } = useTempAuth()

  const form = useForm<TempUserSchemaType>({
    resolver: zodResolver(TempUserSchema),
    defaultValues: {
      name: ''
    }
  })

  const { handleSubmit, control } = form

  const handleFormLogin = (data: TempUserSchemaType) => {
    tempLogin(data)
  }

  return (
    <div className="w-full flex justify-center mt-4">
      <Modal
        trigger={<Button className="w-full bg-indigo-500">Try as a Temp User</Button>}
        title="Sign Up as a Temp User"
        description="All the data related to this user will be lost after 1 hour">
        <Form {...form}>
          <form onSubmit={handleSubmit(handleFormLogin)} className="space-y-4">
            <div>
              <FormField
                control={control}
                name="name"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel htmlFor="name" className="text-black">
                      Name
                    </FormLabel>
                    <FormControl>
                      <div className="relative mt-1">
                        <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                        <Input
                          id="name"
                          type="text"
                          placeholder="Full Name"
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
              {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : 'Sign In'}
            </Button>
          </form>
        </Form>
      </Modal>
    </div>
  )
}

export default TempUser
