'use client'
import { CreateFormSchema, CreateFormSchemaType } from '@/app/schemas/form.schemas'
import { useCreateForm, useUpdateForm } from '@/hooks/queryHooks/useFormApp'
import { FormType } from '@/types/form.types'
import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect, useMemo } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { SubmitButton } from '../common/SubmitButton'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import { Input } from '../ui/input'
import { Switch } from '../ui/switch'
import { Textarea } from '../ui/textarea'
import { addHours } from 'date-fns'
import { toDateTimeLocal } from '@/lib/date.utils'

interface UpsertFormProps {
  formId?: string
  data?: FormType
  onOpenChange?: (open: boolean) => void
}

export const UpsertForm = ({ formId, data, onOpenChange }: UpsertFormProps) => {
  const form = useForm({
    resolver: zodResolver(CreateFormSchema),
    defaultValues: {
      title: data?.title || '',
      description: data?.description || '',
      isScheduled: data?.isScheduled ?? false,
      scheduledTime: toDateTimeLocal(data?.scheduledTime),
      closingTime: toDateTimeLocal(data?.closingTime)
    },
    mode: 'onChange',
    reValidateMode: 'onChange'
  })
  const { mutate, isPending: createIsPending } = useCreateForm()
  const { mutate: updateMutate, isPending: updateIsPending } = useUpdateForm()

  const isPending = useMemo(() => createIsPending || updateIsPending, [createIsPending, updateIsPending])

  const { control, handleSubmit, reset, setValue } = form
  const { isScheduled, closingTime, scheduledTime } = useWatch({ control })

  useEffect(() => {
    if (!isScheduled) {
      setValue('scheduledTime', undefined)
    }
  }, [isScheduled])

  const onSubmit = (formData: CreateFormSchemaType) => {
    if (formId) {
      updateMutate(
        { id: formId, data: formData },
        {
          onSuccess: () => {
            reset()
            onOpenChange?.(false)
          }
        }
      )
      return
    }
    mutate(
      { data: formData },
      {
        onSuccess: () => {
          reset()
          onOpenChange?.(false)
        }
      }
    )
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
          <FormField
            name="isScheduled"
            control={control}
            render={({ field }) => (
              <FormItem className="flex items-center justify-between">
                <FormLabel>Should be scheduled</FormLabel>
                <FormControl>
                  <Switch checked={field.value} onCheckedChange={field.onChange} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            name="scheduledTime"
            control={control}
            render={({ field }) => (
              <FormItem>
                <FormLabel htmlFor="scheduledTime">Scheduled At</FormLabel>
                <FormControl>
                  <Input type="datetime-local" id="scheduledTime" {...field} disabled={isPending || !isScheduled} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
            disabled={!isScheduled}
          />
          <FormField
            name="closingTime"
            control={control}
            render={({ field }) => (
              <FormItem>
                <FormLabel htmlFor="closingTime">Closes At</FormLabel>
                <FormControl>
                  <Input type="datetime-local" id="closingTime" {...field} disabled={isPending} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <SubmitButton type="submit" disabled={isPending} isLoading={isPending}>
            Submit
          </SubmitButton>
        </div>
      </form>
    </Form>
  )
}
