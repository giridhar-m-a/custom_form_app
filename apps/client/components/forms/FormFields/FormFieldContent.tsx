'use client'

import { FormFieldCreateSchemaType, FormFieldSchema, FormFieldSchemaType } from '@/app/schemas/form.schemas'
import { CustomLoader } from '@/components/common/CustomLoader'
import { SubmitButton } from '@/components/common/SubmitButton'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useCreateFormField, useGetFormFields, useUpdateFormField } from '@/hooks/queryHooks/useFormApp'
import { FieldType, FormField } from '@/types/form.types'
import {
  DndContext,
  DragEndEvent,
  KeyboardSensor,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors
} from '@dnd-kit/core'
import { SortableContext, arrayMove, sortableKeyboardCoordinates, verticalListSortingStrategy } from '@dnd-kit/sortable'
import { zodResolver } from '@hookform/resolvers/zod'
import { useState } from 'react'
import { useForm } from 'react-hook-form'
import toast from 'react-hot-toast'
import { MdAdd } from 'react-icons/md'
import { FormFieldWrapper } from './FormFieldWrapper'
import { FIELD_TYPE_OPTIONS } from './formFields.config'
import { v4 as uuidv4 } from 'uuid'

// Augment the schema type locally to include tempId for UI stability
type FormFieldWithId = FormFieldCreateSchemaType & { tempId: string }

interface FormFieldContentProps {
  formId: string
  formTitle: string
  initialFields?: FormField[]
  mode?: 'create' | 'edit'
}

export const FormFieldContent = ({ formId, initialFields = [], formTitle, mode = 'create' }: FormFieldContentProps) => {
  const [editingField, setEditingField] = useState<string | null>(null)
  const { mutateAsync: createFormField, isPending: isCreatingFormField } = useCreateFormField()
  const { mutateAsync: updateFormField, isPending: isUpdatingFormField } = useUpdateFormField()
  const { data: FieldRes, isLoading } = useGetFormFields(formId, initialFields)

  const isPending = isCreatingFormField || isUpdatingFormField

  const schema = FormFieldSchema

  const {
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
    reset
  } = useForm({
    resolver: zodResolver(schema),
    defaultValues: {
      formId,
      formFields:
        FieldRes?.data?.map(f => ({
          ...f,
          tempId: uuidv4()
        })) || [],
      removedFields: [],
      removedFieldOptions: []
    }
  })

  // We cast watch result to our augmented type
  const formFields = (watch('formFields') || []) as FormFieldWithId[]
  const removedFields = watch('removedFields') || []
  const removedFieldOptions = watch('removedFieldOptions') || []

  // Setup drag and drop sensors
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates
    })
  )

  // Handle drag end for field reordering
  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (active.id !== over?.id) {
      const oldIndex = formFields.findIndex(field => field.tempId === active.id)
      const newIndex = formFields.findIndex(field => field.tempId === over?.id)
      const reorderedFields = arrayMove(formFields, oldIndex, newIndex)

      // Update ordering values to match the new positions
      const updatedFields = reorderedFields.map((field, index) => ({
        ...field,
        ordering: index
      }))

      setValue('formFields', updatedFields as any, { shouldValidate: true, shouldDirty: true })
      // No need to manually update editingField since it now tracks ID, which didn't change
    }
  }

  // Toggle edit mode for a specific field
  const toggleEdit = (id: string) => {
    setEditingField(prev => (prev === id ? null : id))
  }

  // Add a new field
  const handleAddField = (afterIndex: number, fieldType: FieldType = 'text') => {
    const newField: FormFieldWithId = {
      fieldLabel: 'New Field',
      fieldType: fieldType,
      isRequired: false,
      ordering: afterIndex + 1,
      options: [],
      tempId: uuidv4()
    }

    const updatedFields = [...formFields]
    updatedFields.splice(afterIndex + 1, 0, newField)

    // Recalculate ordering
    const reorderedFields = updatedFields.map((field, idx) => ({
      ...field,
      ordering: idx
    }))

    setValue('formFields', reorderedFields as any, { shouldValidate: true, shouldDirty: true })
    // Auto-edit the new field
    setEditingField(newField.tempId)
  }

  // Duplicate a field
  const handleDuplicateField = (index: number) => {
    const fieldToDuplicate = formFields[index]
    if (errors?.formFields?.[index] || !fieldToDuplicate) {
      toast.error('Errors found in the field')
      return
    }
    handleSubmitField(fieldToDuplicate as FormFieldCreateSchemaType, index)
    const duplicatedField: FormFieldWithId = {
      fieldLabel: `${fieldToDuplicate.fieldLabel} (Copy)`,
      fieldType: fieldToDuplicate.fieldType ?? 'text',
      isRequired: fieldToDuplicate.isRequired ?? false,
      ordering: index + 1,
      fieldId: undefined, // Remove fieldId for new field
      options: fieldToDuplicate?.options?.map(opt => ({
        optionLabel: opt.optionLabel,
        ordering: opt.ordering ?? 0,
        isAnswer: opt.isAnswer ?? false,
        optionId: undefined,
        fieldId: undefined
      })),
      tempId: uuidv4()
    } as FormFieldWithId

    const updatedFields = [...formFields]
    updatedFields.splice(index + 1, 0, duplicatedField)

    // Recalculate ordering
    const reorderedFields = updatedFields.map((field, idx) => ({
      ...field,
      ordering: idx
    }))

    setValue('formFields', reorderedFields as any, { shouldValidate: true, shouldDirty: true })
    setEditingField(duplicatedField.tempId)
  }

  // Remove a field
  const handleRemoveField = (index: number) => {
    const fieldToRemove = formFields[index]

    // If editing an existing field, track it for deletion
    if (fieldToRemove.fieldId) {
      setValue('removedFields', [...removedFields, fieldToRemove.fieldId], { shouldValidate: true })
      // Also track removed options
      const optionIds = fieldToRemove?.options?.filter(opt => opt.optionId).map(opt => opt.optionId!) || []
      setValue('removedFieldOptions', [...removedFieldOptions, ...optionIds], { shouldValidate: true })
    }

    // Remove from formFields array and recalculate ordering
    const updatedFields = formFields.filter((_, idx) => idx !== index)
    const reorderedFields = updatedFields.map((field, idx) => ({
      ...field,
      ordering: idx
    }))

    setValue('formFields', reorderedFields as any, { shouldValidate: true, shouldDirty: true })

    // Adjust editingField
    if (editingField === fieldToRemove.tempId) {
      setEditingField(null)
    }
  }

  // Submit a single field (save changes)
  const handleSubmitField = (field: FormFieldCreateSchemaType, index: number) => {
    const updatedFields = [...formFields]
    // Preserve the tempId when updating the field data from edit form
    updatedFields[index] = { ...field, tempId: formFields[index].tempId } as FormFieldWithId
    setValue('formFields', updatedFields as any, { shouldValidate: true, shouldDirty: true })

    // Exit edit mode for this field if valid
    // TODO: Adding validation check here would be good, but react-hook-form handles validation on submit
    setEditingField(null)
  }

  // Form submission
  const onFormSubmit = async (data: FormFieldSchemaType) => {
    if (mode === 'create') {
      createFormField(
        { data },
        {
          onSuccess: data => {
            if (data.data?.length) {
              reset({
                formFields: data.data,
                removedFields: [],
                removedFieldOptions: []
              })
            }
          }
        }
      )
    } else {
      updateFormField(
        { data },
        {
          onSuccess: ({ data }) => {
            if (data?.length) {
              // Re-map with new tempIds to ensure consistency, though we could try to preserve them if we matched fieldIds
              const fieldsWithIds = data.map(f => ({
                ...f,
                tempId: uuidv4()
              }))
              setValue('formFields', fieldsWithIds)
              setValue('removedFields', [])
              setValue('removedFieldOptions', [])
            }
          }
        }
      )
    }
  }

  return (
    <>
      {!isLoading && (
        <div className="space-y-6">
          {/* Top Action Bar */}
          <div className="flex items-center justify-between rounded-xl border border-border bg-card p-4 shadow-sm">
            <h2 className="text-lg font-semibold">{formTitle}</h2>
            <div className="flex items-center gap-2">
              <DropdownMenu>
                <DropdownMenuTrigger asChild>
                  <Button variant="outline" className="gap-2">
                    <MdAdd size={18} />
                    Add Field
                  </Button>
                </DropdownMenuTrigger>
                <DropdownMenuContent align="end" className="w-56 max-h-96 overflow-y-auto">
                  <DropdownMenuLabel>Select Field Type</DropdownMenuLabel>
                  <DropdownMenuSeparator />
                  {FIELD_TYPE_OPTIONS.map(ft => (
                    <DropdownMenuItem key={ft.value} onClick={() => handleAddField(formFields.length - 1, ft.value)}>
                      {ft.icon && <ft.icon className="mr-2 h-4 w-4" />}
                      <span>{ft.label}</span>
                    </DropdownMenuItem>
                  ))}
                </DropdownMenuContent>
              </DropdownMenu>
              <SubmitButton onClick={handleSubmit(onFormSubmit as any)} disabled={isPending} isLoading={isPending}>
                Save
              </SubmitButton>
            </div>
          </div>
          <ScrollArea className="max-h-[65vh] min-h-[300px] p-4">
            <form onSubmit={handleSubmit(onFormSubmit as any)} className="space-y-6">
              <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
                <SortableContext items={formFields.map(field => field.tempId)} strategy={verticalListSortingStrategy}>
                  <div className="space-y-4">
                    {formFields.map((field, index) => (
                      <FormFieldWrapper
                        key={field.tempId}
                        formField={field as FormFieldCreateSchemaType}
                        tempId={field.tempId}
                        handleRemoveField={index => handleRemoveField(index)}
                        handleDuplicateField={index => handleDuplicateField(index)}
                        handleAddField={(index, type) => handleAddField(index, type)}
                        index={index}
                        handleSubmitField={(f, i) => handleSubmitField(f, i)}
                        errors={errors.formFields?.[index]}
                        isEdit={editingField === field.tempId}
                        setEdit={() => toggleEdit(field.tempId)}
                      />
                    ))}
                  </div>
                </SortableContext>
              </DndContext>

              {formFields.length === 0 && (
                <div className="flex flex-col items-center justify-center py-12 text-center border-2 border-dashed border-muted rounded-xl bg-muted/10">
                  <p className="text-muted-foreground mb-4">No fields added yet.</p>
                  <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                      <Button variant="outline" className="gap-2">
                        <MdAdd size={18} />
                        Add Your First Field
                      </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="center" className="w-56 overflow-y-auto max-h-96">
                      {FIELD_TYPE_OPTIONS.map(ft => (
                        <DropdownMenuItem key={ft.value} onClick={() => handleAddField(-1, ft.value)}>
                          {ft.icon && <ft.icon className="mr-2 h-4 w-4" />}
                          <span>{ft.label}</span>
                        </DropdownMenuItem>
                      ))}
                    </DropdownMenuContent>
                  </DropdownMenu>
                </div>
              )}

              {/* Display form-level errors */}
              {errors.root && <p className="text-destructive text-sm">{errors.root.message}</p>}
            </form>
          </ScrollArea>
        </div>
      )}
      {isLoading && (
        <div className="flex items-center justify-center h-full w-full">
          <CustomLoader />
        </div>
      )}
    </>
  )
}
