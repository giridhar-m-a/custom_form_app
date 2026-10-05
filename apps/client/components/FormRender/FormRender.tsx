'use client'

import { FormField as FormFieldType } from '@/types/form.types'
import { FormInputWrapper } from './FormInputwrapper'
import { useForm } from 'react-hook-form'
import { Form } from '../ui/form'
import { useMemo } from 'react'
import { FormSubmission } from '@/app/schemas/response.schemas'
import { useCreateResponse } from '@/hooks/queryHooks/useResponses'
import { SubmitButton } from '../common/SubmitButton'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card'
import { CheckCircle2 } from 'lucide-react'

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

  const { mutate: createResponse, isPending, data: submissionResult } = useCreateResponse()
  const isSubmitted = submissionResult?.status === 200

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

  return (
    <div className="w-full min-w-0">
      {isSubmitted ? (
        <Card className="border-emerald-500/30 bg-emerald-500/5 shadow-none">
          <CardHeader className="justify-items-center text-center">
            <div className="mb-2 flex size-14 items-center justify-center rounded-full bg-emerald-500/10 text-emerald-600 dark:text-emerald-400">
              <CheckCircle2 className="size-8" aria-hidden="true" />
            </div>
            <CardTitle className="text-xl">Response submitted</CardTitle>
            <CardDescription className="text-base">
              Thank you for completing this form. Your response has been received.
            </CardDescription>
          </CardHeader>
          {submissionResult.message && (
            <CardContent className="text-center text-sm text-muted-foreground">
              {submissionResult.message}
            </CardContent>
          )}
        </Card>
      ) : (
      <Form {...form}>
        <form className="w-full min-w-0 space-y-4" onSubmit={handleFormSubmit} noValidate>
          {fields.map((field, index) => (
            <FormInputWrapper key={field.fieldId} formField={field} control={control as any} index={index} />
          ))}
          {fields.length > 0 && (
            <SubmitButton type="submit" className="w-full" isLoading={isPending}>
              Submit
            </SubmitButton>
          )}
        </form>
      </Form>
      )}
    </div>
  )
}
