'use client'

import { FormFieldCreateSchema, FormFieldCreateSchemaType, FormFieldOptionSchemaType } from '@/app/schemas/form.schemas'
import { zodResolver } from '@hookform/resolvers/zod'
import { useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { FormField } from './FormField'
import toast from 'react-hot-toast'

interface EditFieldFormProps {
  initialField: FormFieldCreateSchemaType
  index: number
  onSave: (data: FormFieldCreateSchemaType) => void
  onCancel: () => void
  onRemove: () => void
  onDuplicate: () => void
}

export const EditFieldForm = ({ initialField, index, onSave, onCancel, onRemove, onDuplicate }: EditFieldFormProps) => {
  const schema = FormFieldCreateSchema

  // Helper to ensure at least one option exists for MCQ types
  const ensureDefaultOption = (field: FormFieldCreateSchemaType): FormFieldOptionSchemaType[] => {
    const mcqTypes = ['radio', 'checkbox', 'dropdown', 'multiselect']
    if (mcqTypes.includes(field.fieldType) && (!field.options || field.options.length === 0)) {
      return [
        {
          optionLabel: '',
          ordering: 0,
          isAnswer: false,
          optionId: undefined,
          fieldId: field.fieldId
        }
      ]
    }
    return field.options.map(opt => ({
      ...opt,
      optionId: opt.optionId
    }))
  }

  const {
    handleSubmit,
    setValue,
    watch,
    formState: { errors }
  } = useForm<FormFieldCreateSchemaType>({
    resolver: zodResolver(schema) as any,
    defaultValues: {
      ...initialField,
      options: ensureDefaultOption(initialField)
    },
    mode: 'onChange'
  })

  const field = watch() as FormFieldCreateSchemaType

  // Ensure options never drops to zero for MCQ types (Fixes focus loss on delete-all)
  useEffect(() => {
    const mcqTypes = ['radio', 'checkbox', 'dropdown', 'multiselect']
    if (mcqTypes.includes(field.fieldType) && field.options.length === 0) {
      // Auto-add an empty option if all are deleted, to maintain stable UI
      // Use standard handleAddOption logic or direct set, but we need unique ID
      const newOption = {
        optionLabel: '',
        ordering: 0,
        isAnswer: false,
        optionId: undefined,
        fieldId: field.fieldId
      }
      setValue('options', [newOption], { shouldValidate: true })
    }
  }, [field.options?.length, field.fieldType])

  // Update form if initialField changes
  useEffect(() => {
    // We could sync external changes here if needed, but usually edit mode is isolated
  }, [initialField])

  const handleLabelChange = (label: string) => {
    setValue('fieldLabel', label, { shouldValidate: true, shouldDirty: true })
  }

  const handleOptionsChange = (options: FormFieldOptionSchemaType[]) => {
    setValue('options', options, { shouldValidate: true, shouldDirty: true })
  }

  const handleRemoveOption = (optionIndex: number) => {
    const newOptions = field.options.filter((_, idx) => idx !== optionIndex)
    setValue('options', newOptions, { shouldValidate: true, shouldDirty: true })
  }

  const handleAddOption = (position: number) => {
    const newOption: FormFieldOptionSchemaType = {
      optionLabel: `Option ${field.options.length + 1}`,
      ordering: position,
      isAnswer: false,
      optionId: undefined, // Generate stable ID
      fieldId: field.fieldId
    }
    const newOptions = [...field.options]
    newOptions.splice(position, 0, newOption)

    // Recalc ordering but KEEP optionId
    const reordered = newOptions.map((opt, idx) => ({ ...opt, ordering: idx }))

    setValue('options', reordered, { shouldValidate: true, shouldDirty: true })
  }

  const handleRequiredChange = (checked: boolean) => {
    setValue('isRequired', checked, { shouldValidate: true, shouldDirty: true })
  }

  const onChangeFieldType = (type: any) => {
    setValue('fieldType', type, { shouldValidate: true })
  }

  const onSubmit = (data: FormFieldCreateSchemaType) => {
    console.log('Form Submitted (Valid):', data)
    onSave(data)
  }

  const onError = (errors: any) => {
    console.error('Form Validation Errors:', errors)
  }

  const handleSaveClick = () => {
    handleSubmit(onSubmit, onError)()
  }

  const handleDuplicateClick = () => {
    if (Object.keys(errors).length > 0) {
      console.warn('Cannot duplicate field with errors', errors)
      toast.error('Cannot duplicate field with errors')
      // Optionally show a toast here
      return
    }
    onDuplicate()
  }

  // Check for root errors on the options array (e.g., min length, or custom refinements)
  const optionsError = Array.isArray(errors.options) ? undefined : (errors.options as any)?.message

  return (
    <div className="flex w-full flex-col gap-4">
      <FormField
        field={field}
        onLabelChange={handleLabelChange}
        onRemove={handleRemoveOption}
        onOptionsChange={handleOptionsChange}
        isEdit={true}
        fieldType={field.fieldType}
        addField={handleAddOption}
        // We pass the local form errors
        errors={Array.isArray(errors.options) ? (errors.options as any) : []}
        labelError={errors.fieldLabel?.message}
        onChangeFieldType={onChangeFieldType}
        onRequiredChange={handleRequiredChange}
        rootError={optionsError || errors?.root?.options?.message}
        // Form Actions
        onDuplicateField={handleDuplicateClick}
        onRemoveField={onRemove}
        onAddField={() => {}} // No-op for edit mode, or could pass a handler if we wanted to allow adding from here
        onSave={handleSaveClick}
        onCancel={onCancel}
      />
      {/* Show root errors (like answer count) */}
      {/* Root errors are now handled inside FormField via rootError prop */}
    </div>
  )
}
