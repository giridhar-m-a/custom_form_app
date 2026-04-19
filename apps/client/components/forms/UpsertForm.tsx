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
import { DEFAULT_INVITATION_SCHEDULE_GAP } from '@/lib/constants/constants'

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
      scheduledTime: data?.scheduledTime ? new Date(data?.scheduledTime).toISOString() : '',
      closingTime: data?.closingTime ? new Date(data?.closingTime).toISOString() : '',
      invitationScheduleGap: data?.invitationScheduleGap
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
                  <Switch
                    checked={field.value}
                    onCheckedChange={value => {
                      field.onChange(value)
                      setValue('scheduledTime', undefined)
                      setValue('invitationScheduleGap', undefined)
                    }}
                    disabled={isPending || data?.status === 'published'}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          {isScheduled && data?.status !== 'published' && (
            <div className="flex flex-col gap-4">
              <FormField
                name="scheduledTime"
                control={control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel htmlFor="scheduledTime">Scheduled At</FormLabel>
                    <FormControl>
                      <Input
                        type="datetime-local"
                        id="scheduledTime"
                        {...field}
                        value={field.value ? toDateTimeLocal(new Date(field.value)) : undefined}
                        onChange={e => {
                          field.onChange(new Date(e.target.value).toISOString())
                          setValue('invitationScheduleGap', DEFAULT_INVITATION_SCHEDULE_GAP)
                        }}
                        disabled={isPending || !isScheduled}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
                disabled={!isScheduled}
              />
              <FormField
                name="invitationScheduleGap"
                control={control}
                render={({ field }) => (
                  <FormItem>
                    <FormLabel htmlFor="invitationScheduleGap">Invitation Schedule Gap (in minutes)</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        id="invitationScheduleGap"
                        {...field}
                        onChange={e => field.onChange(Number(e.target.value))}
                        disabled={isPending || !isScheduled}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
                disabled={!isScheduled}
              />
            </div>
          )}
          <FormField
            name="closingTime"
            control={control}
            render={({ field }) => (
              <FormItem>
                <FormLabel htmlFor="closingTime">Closes At</FormLabel>
                <FormControl>
                  <Input
                    type="datetime-local"
                    id="closingTime"
                    {...field}
                    value={field.value ? toDateTimeLocal(new Date(field.value)) : undefined}
                    onChange={e => field.onChange(new Date(e.target.value).toISOString())}
                    disabled={isPending}
                  />
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
