'use client'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { CreateFormSchema, CreateFormSchemaType } from '@/app/schemas/form.schemas'
import { Input } from '../ui/input'
import { Textarea } from '../ui/textarea'
import { Button } from '../ui/button'
import { FormType } from '@/types/form.types'
import { useCreateForm } from '@/hooks/queryHooks/useFormApp'

interface UpsertFormProps {
  formId?: string
  data?: FormType
}

export const UpsertForm = ({ formId, data }: UpsertFormProps) => {
  const form = useForm<CreateFormSchemaType>({
    resolver: zodResolver(CreateFormSchema),
    defaultValues: {
      title: data?.title || '',
      description: data?.description || ''
    },
    mode: 'onChange',
    reValidateMode: 'onChange'
  })
  const { mutate, isPending } = useCreateForm()

  const { control, handleSubmit, reset } = form

  const onSubmit = (formData: CreateFormSchemaType) => {
    if (formId) {
      return
    }
    mutate({ data: formData })
  }

  return (
    <Form {...form}>
      <form onSubmit={handleSubmit(onSubmit)} action="">
        <div className="flex flex-col gap-4 w-full">
          <FormField
            name="title"
            control={control}
            render={({ field }) => (
              <FormItem>
                <FormLabel>Title</FormLabel>
                <FormControl>
                  <Input {...field} disabled={isPending} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            name="description"
            control={control}
            render={({ field }) => (
              <FormItem>
                <FormLabel>Description</FormLabel>
                <FormControl>
                  <Textarea className="h-36" {...field} disabled={isPending} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit" disabled={isPending}>
            Submit
          </Button>
        </div>
      </form>
    </Form>
  )
}
