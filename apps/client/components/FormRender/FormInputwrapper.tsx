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

interface FormInputWrapperProps {
  formField: FormFieldType
  control: Control<FieldValues>
}

export const FormInputWrapper = ({ formField, control }: FormInputWrapperProps) => {
  return (
    <FormField
      name={formField.fieldId} // ✅ stable RHF name
      control={control}
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
}

const RenderFormInput = ({ fieldType, field, id, ariaLabel, formField }: RenderFormInputProps) => {
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
      return <FormInputCheckbox options={formField?.options} {...field} aria-label={ariaLabel} />

    case 'dropdown':
      return (
        <FormInputSelect
          options={formField?.options.map(option => ({ value: option.optionId, label: option.optionLabel }))}
          onChange={field.onChange}
          value={field.value}
          aria-label={ariaLabel}
          placeholder={'Select an option'}
        />
      )

    case 'multiselect':
      return (
        <MultiSelect
          options={formField?.options.map(option => ({ value: option.optionId, label: option.optionLabel }))}
          onChange={field.onChange}
          value={field.value}
          aria-label={ariaLabel}
          placeholder={'Select an option'}
        />
      )

    case 'radio':
      return (
        <RadioGroup
          options={formField?.options.map(option => ({ value: option.optionId, label: option.optionLabel }))}
          onChange={field.onChange}
          value={field.value}
          aria-label={ariaLabel}
        />
      )

    case 'file':
      return <FormShortInput id={id} type="file" {...field} aria-label={ariaLabel} />

    case 'image':
      return <FormShortInput id={id} type="file" {...field} aria-label={ariaLabel} accept="image/*" />

    case 'video':
      return <FormShortInput id={id} type="file" {...field} aria-label={ariaLabel} accept="video/*" />

    case 'audio':
      return <FormShortInput id={id} type="file" {...field} aria-label={ariaLabel} accept="audio/*" />

    case 'slider':
      return <FormShortInput id={id} type="range" {...field} aria-label={ariaLabel} />

    case 'rating':
      return <FormStarRating {...field} aria-label={ariaLabel} />

    default:
      return <FormShortInput id={id} type="text" {...field} aria-label={ariaLabel} />
  }
}
