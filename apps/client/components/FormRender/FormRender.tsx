'use client'

import { FormField as FormFieldType } from '@/types/form.types'
import { FormInputWrapper } from './FormInputwrapper'
import { useForm } from 'react-hook-form'
import { Form } from '../ui/form'
import { ScrollArea } from '../ui/scroll-area'
import { useMemo } from 'react'
import { FormSubmission } from '@/app/schemas/response.schemas'
import { useCreateResponse } from '@/hooks/queryHooks/useResponses'
import { SubmitButton } from '../common/SubmitButton'

interface FormRenderProps {
  fields: FormFieldType[]
  formId: string
  respondentId: string
  onSubmit?: (data: FormSubmission) => void
  token: string
}

interface FormValues {
  formId: string
  respondentId: string
  responses: {
    formFieldId: string
    responseText: string
    responseOptions: { optionId: string }[]
    responseFiles: { fileName: string; filePath: string; fileSize: number; fileType: string }[]
  }[]
}

/**
 * Builds default values matching the FormSubmissionSchema shape.
 * Each response item starts with the correct formFieldId and empty values.
 */
function buildDefaultValues(fields: FormFieldType[], formId: string, respondentId: string): FormValues {
  return {
    formId,
    respondentId,
    responses: fields.map(field => ({
      formFieldId: field.fieldId,
      responseText: '',
      responseOptions: [] as { optionId: string }[],
      responseFiles: [] as { fileName: string; filePath: string; fileSize: number; fileType: string }[]
    }))
  }
}

export const FormRender = ({ fields, formId, respondentId, onSubmit, token }: FormRenderProps) => {
  const defaultValues = useMemo(() => buildDefaultValues(fields, formId, respondentId), [fields, formId, respondentId])

  const { mutate: createResponse, isPending } = useCreateResponse()

  const form = useForm<FormValues>({
    defaultValues
  })

  const {
    control,
    handleSubmit,
    formState: { errors },
    watch
  } = form

  const handleFormSubmit = handleSubmit(data => {
    const submission: FormSubmission = {
      formId: data.formId,
      respondentId: data.respondentId,
      responses: data.responses.map(item => ({
        formFieldId: item.formFieldId,
        responseText: item.responseText,
        responseOptions: item.responseOptions,
        responseFiles: item.responseFiles
      }))
    }
    if (onSubmit) {
      onSubmit(submission)
    } else {
      createResponse({ data: submission, token })
    }
  })

  console.log('form errors: ', errors)
  console.log('form watch: ', watch())

  return (
    <div>
      <ScrollArea className="h-[calc(100vh-20rem)]">
        <Form {...form}>
          <form className="space-y-4 px-4" onSubmit={handleFormSubmit}>
            {fields.map((field, index) => (
              <FormInputWrapper key={field.fieldId} formField={field} control={control as any} index={index} />
            ))}
            <SubmitButton type="submit" className="w-full" isLoading={isPending}>
              Submit
            </SubmitButton>
          </form>
        </Form>
      </ScrollArea>
    </div>
  )
}
