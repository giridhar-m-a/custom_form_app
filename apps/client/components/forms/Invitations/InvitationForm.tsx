import { InvitationSchema, InvitationType } from '@/app/schemas/invitation.schemas'
import { SubmitButton } from '@/components/common/SubmitButton'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { useCreateInvitation } from '@/hooks/queryHooks/useInvitations'
import { zodResolver } from '@hookform/resolvers/zod'
import { useForm } from 'react-hook-form'

interface InvitationFormProps {
  formId: string
  setInviteOpen?: (open: boolean) => void
}

export const InvitationForm = ({ formId, setInviteOpen }: InvitationFormProps) => {
  const form = useForm<InvitationType>({
    resolver: zodResolver(InvitationSchema),
    reValidateMode: 'onChange',
    mode: 'onChange',
    defaultValues: {
      form_id: formId
    }
  })
  const { handleSubmit, control } = form

  const { mutate: createInvitation, isPending } = useCreateInvitation()

  const onSubmit = (data: InvitationType) => {
    createInvitation(
      { data },
      {
        onSuccess: () => {
          setInviteOpen?.(false)
        }
      }
    )
  }
  return (
    <Form {...form}>
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={control}
          name="email"
          render={({ field }) => (
            <FormItem className="flex flex-col gap-2">
              <FormLabel htmlFor="email">Email</FormLabel>
              <FormControl>
                <Input id="email" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={control}
          name="name"
          render={({ field }) => (
            <FormItem className="flex flex-col gap-2">
              <FormLabel htmlFor="name">Name</FormLabel>
              <FormControl>
                <Input id="name" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <SubmitButton disabled={isPending} isLoading={isPending} className="w-full">
          Invite
        </SubmitButton>
      </form>
    </Form>
  )
}
