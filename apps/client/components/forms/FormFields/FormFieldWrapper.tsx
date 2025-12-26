'use client'

import { FormFieldCreateSchemaType, FormFieldOptionSchemaType } from '@/app/schemas/form.schemas'
import { FieldType } from '@/types/form.types'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import { useEffect, useState } from 'react'
import { FieldErrors } from 'react-hook-form'
import { MdDragIndicator } from 'react-icons/md'
import { EditFieldForm } from './EditFieldForm'
import { FormField } from './FormField'

interface FormFieldWrapperProps {
  formField: FormFieldCreateSchemaType
  handleRemoveField: (index: number) => void
  handleDuplicateField: (index: number) => void
  handleAddField: (index: number, type?: FieldType) => void
  index: number
  handleSubmitField: (field: FormFieldCreateSchemaType, index: number) => void
  errors?: FieldErrors<FormFieldCreateSchemaType>
  isEdit: boolean
  setEdit: () => void
}

export const FormFieldWrapper = ({
  formField,
  handleRemoveField,
  handleDuplicateField,
  handleAddField,
  index,
  handleSubmitField,
  errors,
  isEdit,
  setEdit
}: FormFieldWrapperProps) => {
  const [field, setField] = useState<FormFieldCreateSchemaType>(formField)

  // Setup sortable
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: formField.ordering ?? 0,
    transition: {
      duration: 200,
      easing: 'cubic-bezier(0.25, 1, 0.5, 1)'
    }
  })

  const style = {
    transition: transition || 'transform 200ms cubic-bezier(0.25, 1, 0.5, 1)',
    transform: CSS.Transform.toString(transform),
    opacity: isDragging ? 0.5 : 1,
    zIndex: isDragging ? 1000 : 'auto'
  }

  // Sync with parent state when formField changes
  useEffect(() => {
    setField(formField)
  }, [formField])

  // Auto-revert to saved state when isEdit becomes false
  useEffect(() => {
    if (!isEdit) {
      setField(formField)
    }
  }, [isEdit, formField])

  const handleLabelChange = (label: string) => {
    setField(prev => ({ ...prev, fieldLabel: label }))
  }

  const handleSave = () => {
    // Block save if there are errors
    if (errors && Object.keys(errors).length > 0) {
      return
    }
    handleSubmitField(field, index)
  }

  // Handler to remove an option by index
  const handleRemoveOption = (index: number) => {
    console.log('Removing option at index:', index)
    setField(prev => {
      const newOptions = prev.options.filter((_, i) => i !== index)
      // Recalculate ordering after removal
      const options = newOptions.map((option, i) => ({
        ...option,
        ordering: i
      }))
      return { ...prev, options }
    })
  }

  // Handler to add a new option at a specific position
  const handleAddOption = (position: number) => {
    const newOption: FormFieldOptionSchemaType = {
      optionLabel: `New Option ${field.options.length + 1}`,
      ordering: position,
      isAnswer: false
    }

    setField(prev => {
      const newOptions = [...prev.options]
      newOptions.splice(position, 0, newOption)
      // Recalculate ordering after insertion
      const options = newOptions.map((option, i) => ({
        ...option,
        ordering: i
      }))
      return { ...prev, options }
    })
  }

  const handleOptionsChange = (options: FormFieldOptionSchemaType[]) => {
    setField(prev => ({ ...prev, options }))
  }

  if (isEdit) {
    return (
      <div
        ref={setNodeRef}
        style={style}
        className={`group relative flex w-full rounded-xl border border-primary/50 bg-card p-4 pl-14 pr-4 shadow-md ring-2 ring-primary/20 transition-all ${
          isDragging ? 'shadow-lg bg-accent/50 z-50' : ''
        }`}>
        <div
          {...attributes}
          {...listeners}
          className="absolute left-0 top-0 bottom-0 flex w-10 cursor-grab active:cursor-grabbing items-center justify-center rounded-l-xl border-r border-border bg-muted/30 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground">
          <MdDragIndicator size={20} />
        </div>
        <EditFieldForm
          initialField={field}
          index={index}
          onSave={data => {
            setField(data)
            handleSubmitField(data, index)
          }}
          onCancel={() => setEdit()} // Toggling edit off triggers the auto-revert from useEffect
          onRemove={() => handleRemoveField(index)}
          onDuplicate={() => handleDuplicateField(index)}
        />
      </div>
    )
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`group relative flex w-full rounded-xl border border-border bg-card p-4 pl-14 pr-4 shadow-sm transition-all hover:shadow-md ${
        isDragging ? 'shadow-lg ring-2 ring-primary/20 bg-accent/50' : ''
      }`}>
      <div
        {...attributes}
        {...listeners}
        className="absolute left-0 top-0 bottom-0 flex w-10 cursor-grab active:cursor-grabbing items-center justify-center rounded-l-xl border-r border-border bg-muted/30 text-muted-foreground transition-colors hover:bg-muted hover:text-foreground">
        <MdDragIndicator size={20} />
      </div>

      <FormField
        field={field}
        onLabelChange={handleLabelChange}
        onRemove={handleRemoveOption}
        onOptionsChange={handleOptionsChange}
        isEdit={false}
        fieldType={field.fieldType}
        addField={handleAddOption}
        errors={errors?.options as unknown as FieldErrors<FormFieldOptionSchemaType>[]}
        labelError={errors?.fieldLabel?.message}
        // Form Actions
        onDuplicateField={() => handleDuplicateField(index)}
        onRemoveField={() => handleRemoveField(index)}
        onAddField={type => handleAddField(index, type)}
        onCancel={setEdit} // Acts as setEdit for view mode
      />
    </div>
  )
}
