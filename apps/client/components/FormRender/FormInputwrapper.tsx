import * as React from 'react'
import { Control, ControllerRenderProps, FieldValues } from 'react-hook-form'

import { FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'

import { FieldType, FormField as FormFieldType } from '@/types/form.types'
import { FormShortInput } from './FormShortInput'
import { FormLongInput } from './FormLongInput'
import { FormInputCheckbox } from './FormInputCheckbox'
import { FormInputSelect } from './FormInputSelect'
import { MultiSelect } from '../common/MultiSelect'
import { RadioGroup } from '../common/RadioGroup'
import { FormStarRating } from './FormStarRating'
import FileUpload from './FileUpload'
import { useSearchParams } from 'next/navigation'
import { ResponseItemSchema } from '@/app/schemas/response.schemas'

interface FormInputWrapperProps {
  formField: FormFieldType
  control: Control<FieldValues>
  index: number
}

/**
 * Returns the appropriate ResponseItemSchema field name based on field type.
 * - Text-like fields → responseText
 * - Option-like fields → responseOptions
 * - File-like fields → responseFiles
 */
function getResponseFieldName(
  fieldType: FieldType
): 'responseText' | 'responseOptions' | 'responseFiles' {
  if (['checkbox', 'radio', 'dropdown', 'multiselect'].includes(fieldType)) {
    return 'responseOptions'
  }
  if (['file', 'image', 'video', 'audio'].includes(fieldType)) {
    return 'responseFiles'
  }
  return 'responseText'
}

/**
 * Builds a validate function for this field using ResponseItemSchema.
 * Constructs a full response item from the current field value,
 * runs safeParse, and returns the error message for the specific field.
 */
function buildValidate(formField: FormFieldType, responseField: 'responseText' | 'responseOptions' | 'responseFiles') {
  const itemSchema = ResponseItemSchema({
    type: formField.fieldType,
    isRequired: formField.isRequired,
    isMultiple: false
  })

  return (value: any) => {
    // Build a complete response item with this field's current value
    const responseItem = {
      formFieldId: formField.fieldId,
      responseText: responseField === 'responseText' ? value : '',
      responseOptions: responseField === 'responseOptions' ? value : [],
      responseFiles: responseField === 'responseFiles' ? value : []
    }

    const result = itemSchema.safeParse(responseItem)
    if (result.success) return true

    // Find the first error matching this specific field
    const fieldError = result.error.issues.find(issue =>
      issue.path.includes(responseField)
    )

    return fieldError?.message || true
  }
}

export const FormInputWrapper = ({ formField, control, index }: FormInputWrapperProps) => {
  const searchParams = useSearchParams()
  const token = searchParams.get('token')

  const responseField = getResponseFieldName(formField.fieldType)
  const fieldName = `responses.${index}.${responseField}`

  const validate = React.useMemo(
    () => buildValidate(formField, responseField),
    [formField, responseField]
  )

  return (
    <FormField
      name={fieldName}
      control={control}
      rules={{ validate }}
      render={({ field }) => (
        <FormItem>
          <FormLabel htmlFor={formField.fieldId}>
            {formField.fieldLabel}
            {formField.isRequired && <span className="ml-1 text-red-500">*</span>}
          </FormLabel>

          <FormControl>
            <RenderFormInput
              fieldType={formField.fieldType}
              field={field}
              id={formField.fieldId}
              ariaLabel={formField.fieldLabel}
              formField={formField}
              token={`${token}`}
            />
          </FormControl>

          <FormMessage />
        </FormItem>
      )}
    />
  )
}

/* ------------------------------------------------------------------ */
/* Internal input renderer (single responsibility)                     */
/* ------------------------------------------------------------------ */

interface RenderFormInputProps {
  fieldType: FieldType
  field: ControllerRenderProps<FieldValues, string>
  id: string
  ariaLabel?: string
  formField: FormFieldType
  token: string
}

const RenderFormInput = ({ fieldType, field, id, ariaLabel, formField, token }: RenderFormInputProps) => {
  switch (fieldType) {
    case 'text':
      return <FormShortInput id={id} type="text" {...field} aria-label={ariaLabel} />

    case 'textArea':
      return <FormLongInput id={id} {...field} aria-label={ariaLabel} />

    case 'number':
      return <FormShortInput id={id} type="number" {...field} aria-label={ariaLabel} />

    case 'email':
      return <FormShortInput id={id} type="email" {...field} aria-label={ariaLabel} />

    case 'phone':
      return <FormShortInput id={id} type="tel" {...field} aria-label={ariaLabel} />

    case 'url':
      return <FormShortInput id={id} type="url" {...field} aria-label={ariaLabel} />

    case 'date':
      return <FormShortInput id={id} type="date" {...field} aria-label={ariaLabel} />

    case 'time':
      return <FormShortInput id={id} type="time" {...field} aria-label={ariaLabel} />

    case 'datetime':
      return <FormShortInput id={id} type="datetime-local" {...field} aria-label={ariaLabel} />

    case 'checkbox':
      return (
        <FormInputCheckbox
          options={formField?.options}
          value={
            Array.isArray(field.value)
              ? field.value.map((o: { optionId: string }) => o.optionId).join(',')
              : ''
          }
          onChange={(optionId: string) => {
            const current: { optionId: string }[] = Array.isArray(field.value) ? field.value : []
            const exists = current.some(o => o.optionId === optionId)
            const next = exists
              ? current.filter(o => o.optionId !== optionId)
              : [...current, { optionId }]
            field.onChange(next)
          }}
          aria-label={ariaLabel}
        />
      )

    case 'dropdown':
      return (
        <FormInputSelect
          options={formField?.options.map(option => ({ value: option.optionId, label: option.optionLabel }))}
          onChange={(value: string) => {
            field.onChange([{ optionId: value }])
          }}
          value={
            Array.isArray(field.value) && field.value.length > 0
              ? field.value[0].optionId
              : ''
          }
          aria-label={ariaLabel}
          placeholder={'Select an option'}
        />
      )

    case 'multiselect':
      return (
        <MultiSelect
          options={formField?.options.map(option => ({ value: option.optionId, label: option.optionLabel }))}
          onChange={(values: string[]) => {
            field.onChange(values.map(v => ({ optionId: v })))
          }}
          value={
            Array.isArray(field.value)
              ? field.value.map((o: { optionId: string }) => o.optionId)
              : []
          }
          aria-label={ariaLabel}
          placeholder={'Select an option'}
        />
      )

    case 'radio':
      return (
        <RadioGroup
          options={formField?.options.map(option => ({ value: option.optionId, label: option.optionLabel }))}
          onChange={(value: string) => {
            field.onChange([{ optionId: value }])
          }}
          value={
            Array.isArray(field.value) && field.value.length > 0
              ? field.value[0].optionId
              : ''
          }
          aria-label={ariaLabel}
        />
      )

    case 'file':
      return (
        <FileUpload
          autoUpload={true}
          uploadPath={`forms/${formField.formId}/${formField.fieldId}`}
          handleResponse={data => {
            const current: any[] = Array.isArray(field.value) ? field.value : []
            field.onChange([
              ...current,
              {
                fileName: data.fileName,
                filePath: data.filePath,
                fileSize: data.fileSize,
                fileType: data.fileType
              }
            ])
          }}
          token={token as string}
        />
      )

    case 'image':
      return (
        <FileUpload
          autoUpload={true}
          uploadPath={`forms/${formField.formId}/${formField.fieldId}`}
          accept="image/*"
          handleResponse={data => {
            const current: any[] = Array.isArray(field.value) ? field.value : []
            field.onChange([
              ...current,
              {
                fileName: data.fileName,
                filePath: data.filePath,
                fileSize: data.fileSize,
                fileType: data.fileType
              }
            ])
          }}
          token={token as string}
        />
      )

    case 'video':
      return (
        <FileUpload
          autoUpload={true}
          uploadPath={`forms/${formField.formId}/${formField.fieldId}`}
          accept="video/*"
          handleResponse={data => {
            const current: any[] = Array.isArray(field.value) ? field.value : []
            field.onChange([
              ...current,
              {
                fileName: data.fileName,
                filePath: data.filePath,
                fileSize: data.fileSize,
                fileType: data.fileType
              }
            ])
          }}
          token={token as string}
        />
      )

    case 'audio':
      return (
        <FileUpload
          autoUpload={true}
          uploadPath={`forms/${formField.formId}/${formField.fieldId}`}
          accept="audio/*"
          handleResponse={data => {
            const current: any[] = Array.isArray(field.value) ? field.value : []
            field.onChange([
              ...current,
              {
                fileName: data.fileName,
                filePath: data.filePath,
                fileSize: data.fileSize,
                fileType: data.fileType
              }
            ])
          }}
          token={token as string}
        />
      )

    case 'slider':
      return <FormShortInput id={id} type="range" {...field} aria-label={ariaLabel} />

    case 'color':
      return <FormShortInput id={id} type="color" {...field} aria-label={ariaLabel} />

    case 'rating':
      return (
        <FormStarRating
          value={field.value ? Number(field.value) : 0}
          onChange={(val: number) => field.onChange(String(val))}
          aria-label={ariaLabel}
        />
      )

    default:
      return <FormShortInput id={id} type="text" {...field} aria-label={ariaLabel} />
  }
}
