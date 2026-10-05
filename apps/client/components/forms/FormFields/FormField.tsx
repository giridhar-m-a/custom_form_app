'use client'

import { FormFieldCreateSchemaType, FormFieldOptionSchemaType } from '@/app/schemas/form.schemas'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { FieldType } from '@/types/form.types'
import { Plus, Trash } from 'lucide-react'
import {
  SortableContext,
  verticalListSortingStrategy,
  useSortable,
  arrayMove,
  sortableKeyboardCoordinates
} from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import {
  DndContext,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
  closestCorners,
  DragEndEvent
} from '@dnd-kit/core'
import { MdDragIndicator } from 'react-icons/md'
import { FieldErrors } from 'react-hook-form'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { FormFieldOptions } from './FormFieldOptions'
import { CommonSelect } from '@/components/common/CommonSelect'
import { FIELD_TYPE_OPTIONS } from './formFields.config'
import { FormStarRating } from '@/components/FormRender/FormStarRating'

interface FormFieldProps {
  field: FormFieldCreateSchemaType
  onLabelChange: (label: string) => void
  onRemove: (index: number) => void
  onOptionsChange: (options: FormFieldOptionSchemaType[]) => void
  isEdit: boolean
  fieldType: FieldType
  addField: (position: number) => void
  errors?: FieldErrors<FormFieldOptionSchemaType>[]
  labelError?: string
  onChangeFieldType?: (type: FieldType) => void
  onRequiredChange?: (required: boolean) => void
  rootError?: string
  // Form Actions
  onRemoveField: () => void
  onDuplicateField: () => void
  onAddField: (type: FieldType) => void
  onSave?: () => void
  onCancel?: () => void
}

export const FormField = ({
  field,
  onLabelChange,
  onRemove,
  onOptionsChange,
  isEdit,
  fieldType,
  addField,
  errors,
  labelError,
  onChangeFieldType,
  onRequiredChange,
  rootError,
  onRemoveField,
  onDuplicateField,
  onAddField,
  onSave,
  onCancel
}: FormFieldProps) => {
  const sensors = useSensors(
    useSensor(PointerSensor),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates
    })
  )

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event
    if (active.id !== over?.id) {
      const oldIndex = field.options.findIndex(option => option.optionId === active.id)
      const newIndex = field.options.findIndex(option => option.optionId === over?.id)

      // Fallback to ordering search if optionId lookup fails (though it shouldn't)
      const finalOldIndex = oldIndex !== -1 ? oldIndex : field.options.findIndex(o => o.ordering === active.id)
      const finalNewIndex = newIndex !== -1 ? newIndex : field.options.findIndex(o => o.ordering === over?.id)

      if (finalOldIndex !== -1 && finalNewIndex !== -1) {
        const reorderedOptions = arrayMove(field.options, finalOldIndex, finalNewIndex)
        // Update ordering values to match the new positions (starting from 0)
        const newOptions = reorderedOptions.map((option, index) => ({
          ...option,
          ordering: index
        }))
        onOptionsChange(newOptions)
      }
    }
  }

  return (
    <div className="w-full relative">
      {/* Edit Mode */}
      {isEdit && (
        <div className="flex flex-col gap-4">
          <div className="flex items-start gap-4">
            <div className="flex-1 space-y-2">
              <Textarea
                value={field.fieldLabel}
                onChange={e => onLabelChange(e.target.value)}
                className="w-full resize-none text-lg font-semibold border-none shadow-none focus-visible:ring-0 placeholder:text-muted-foreground/50 min-h-10"
                placeholder="Question / Field Title"
                rows={1}
                onInput={e => {
                  const target = e.target as HTMLTextAreaElement
                  target.style.height = 'auto'
                  target.style.height = `${target.scrollHeight}px`
                }}
              />
              {labelError && <p className="text-[0.8rem] font-medium text-destructive">{labelError}</p>}
            </div>

            {/* Field Type Selector */}
            {onChangeFieldType && (
              <div className="w-40 shrink-0">
                <CommonSelect
                  options={FIELD_TYPE_OPTIONS}
                  placeholder="Type"
                  value={fieldType}
                  onChange={val => onChangeFieldType(val as FieldType)}
                />
              </div>
            )}

            {onRequiredChange && (
              <div className="flex h-8 items-center gap-2 border-l pl-4 border-border">
                <Switch
                  checked={field.isRequired}
                  onCheckedChange={onRequiredChange}
                  className="scale-75 data-[state=checked]:bg-primary"
                />
                <span className="text-sm font-medium text-muted-foreground">Required</span>
              </div>
            )}

            <FormFieldOptions
              onDuplicateField={onDuplicateField}
              onRemoveField={onRemoveField}
              isEdit={isEdit}
              setEdit={onCancel || (() => {})}
              handleSave={onSave || (() => {})}
              onAddField={onAddField}
            />
          </div>

          {fieldType === 'checkbox' ||
          fieldType === 'radio' ||
          fieldType === 'dropdown' ||
          fieldType === 'multiselect' ? (
            <div className="space-y-4">
              {rootError && <p className="text-sm font-medium text-destructive">{rootError}</p>}
              {field.options.length < 1 && (
                <FieldOption
                  option={{
                    isAnswer: false,
                    ordering: 0,
                    optionLabel: ''
                  }}
                  index={0}
                  onOptionsChange={onOptionsChange}
                  onRemove={onRemove}
                  options={field.options}
                  addField={addField}
                  fieldType={fieldType}
                  // Dummy ID for empty state if needed
                  sortId="empty-state"
                />
              )}
              <DndContext sensors={sensors} collisionDetection={closestCorners} onDragEnd={handleDragEnd}>
                <div className="touch-none">
                  <SortableContext
                    items={field.options.map(option => ({ ...option, id: option.optionId || option.ordering }))}
                    strategy={verticalListSortingStrategy}>
                    <div className="flex flex-col gap-2">
                      {field.options.map((option, index) => {
                        const uniqueId = option.optionId || option.ordering
                        return (
                          <div key={uniqueId} className="space-y-1">
                            <FieldOption
                              option={option}
                              index={index}
                              onOptionsChange={onOptionsChange}
                              onRemove={onRemove}
                              options={field.options}
                              addField={addField}
                              fieldType={fieldType}
                              sortId={uniqueId}
                            />
                            {errors?.[index]?.optionLabel?.message && (
                              <p className="px-2 text-[0.8rem] font-medium text-destructive">
                                {errors[index]?.optionLabel?.message}
                              </p>
                            )}
                          </div>
                        )
                      })}
                    </div>
                  </SortableContext>
                </div>
              </DndContext>
            </div>
          ) : (
            <div className="text-sm text-muted-foreground italic px-1">User input field (no options to configure)</div>
          )}
        </div>
      )}

      {/* View Mode */}
      {!isEdit && (
        <div className="space-y-3 px-1 group">
          <div className="flex items-start justify-between">
            <div className="flex items-center gap-2">
              <p className="text-base font-semibold text-foreground">
                {field.fieldLabel}
                {field.isRequired && <span className="text-destructive ml-1">*</span>}
              </p>
              <span className="text-[10px] font-medium px-1.5 py-0.5 rounded bg-muted text-muted-foreground uppercase tracking-wider">
                {fieldType}
              </span>
            </div>

            <div className="opacity-0 transition-opacity duration-200 group-hover:opacity-100">
              <FormFieldOptions
                onAddField={onAddField}
                onDuplicateField={onDuplicateField}
                onRemoveField={onRemoveField}
                isEdit={isEdit}
                setEdit={onCancel || (() => {})}
                handleSave={onSave || (() => {})}
              />
            </div>
          </div>

          {(fieldType === 'text' ||
            fieldType === 'email' ||
            fieldType === 'url' ||
            fieldType === 'phone' ||
            fieldType === 'number') && (
            <div className="h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm text-muted-foreground/50">
              Short answer text
            </div>
          )}

          {fieldType === 'textArea' && (
            <div className="h-20 w-full rounded-md border border-input bg-background px-3 py-2 text-sm text-muted-foreground/50">
              Long answer text
            </div>
          )}

          {(fieldType === 'checkbox' || fieldType === 'multiselect') && (
            <div className="space-y-2">
              {field.options.map((option, index) => (
                <div key={index} className="flex items-center gap-2">
                  <div
                    className={`h-4 w-4 rounded border flex items-center justify-center ${
                      option.isAnswer ? 'bg-primary border-primary' : 'border-primary/50'
                    }`}>
                    {option.isAnswer && <div className="h-2 w-2 bg-primary-foreground rounded-[1px]" />}
                  </div>
                  <span className="text-sm">{option.optionLabel}</span>
                  {option.isAnswer && <span className="text-[10px] text-primary font-medium">(Correct)</span>}
                </div>
              ))}
            </div>
          )}

          {fieldType === 'radio' && (
            <div className="space-y-2">
              {field.options.map((option, index) => (
                <div key={index} className="flex items-center gap-2">
                  <div
                    className={`h-4 w-4 rounded-full border flex items-center justify-center ${
                      option.isAnswer ? 'bg-primary border-primary' : 'border-primary/50'
                    }`}>
                    {option.isAnswer && <div className="h-2 w-2 bg-primary-foreground rounded-full" />}
                  </div>
                  <span className="text-sm">{option.optionLabel}</span>
                  {option.isAnswer && <span className="text-[10px] text-primary font-medium">(Correct)</span>}
                </div>
              ))}
            </div>
          )}

          {fieldType === 'dropdown' && (
            <div className="h-10 w-full max-w-[200px] flex items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm text-muted-foreground">
              <span>Select visible option</span>
              <div className="h-4 w-4 border-l border-b border-muted-foreground/50 -rotate-45 origin-center -translate-y-0.5 scale-50" />
            </div>
          )}

          {(fieldType === 'date' || fieldType === 'time' || fieldType === 'datetime') && (
            <div className="h-10 w-full flex items-center gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm text-muted-foreground/50">
              <span>Select {fieldType}</span>
            </div>
          )}

          {!!(fieldType === 'file' || fieldType === 'image' || fieldType === 'video' || fieldType === 'audio') && (
            <div className="h-20 w-full flex items-center justify-center rounded-md border-2 border-dashed border-input bg-background/50 text-sm text-muted-foreground">
              Click to upload file
            </div>
          )}
          {fieldType === 'rating' && (
            <div className="h-10 w-full flex items-center gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm text-muted-foreground/50">
              <FormStarRating disabled />
            </div>
          )}
          {fieldType === 'slider' && (
            <div className="h-10 w-full flex items-center gap-2 rounded-md border border-input bg-background px-3 py-2 text-sm text-muted-foreground/50">
              <Input type="range" disabled />
            </div>
          )}
        </div>
      )}
    </div>
  )
}

const FieldOption = ({
  option,
  index,
  onOptionsChange,
  onRemove,
  options,
  addField,
  fieldType,
  sortId
}: {
  option: FormFieldOptionSchemaType
  index: number
  onOptionsChange: (options: FormFieldOptionSchemaType[]) => void
  onRemove: (index: number) => void
  options: FormFieldOptionSchemaType[]
  addField: (index: number) => void
  fieldType: FieldType
  sortId?: string | number
}) => {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
    id: sortId ?? option.ordering,
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

  // Is this field type capable of having 'correct' answers?
  const isMcqType =
    fieldType === 'radio' || fieldType === 'checkbox' || fieldType === 'dropdown' || fieldType === 'multiselect'

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={`group flex items-start gap-2 rounded-md border border-input bg-background p-2 transition-colors hover:border-primary/50 ${
        isDragging ? 'shadow-md ring-1 ring-primary/20' : ''
      }`}>
      <div
        {...attributes}
        {...listeners}
        className="cursor-grab active:cursor-grabbing text-muted-foreground hover:text-foreground mt-2">
        <MdDragIndicator size={20} />
      </div>

      {/* Option Answer Toggle */}
      {isMcqType && (
        <div className="pt-2">
          <Switch
            checked={option.isAnswer}
            onCheckedChange={checked => {
              const newOptions = [...options]
              // If single choice (radio/dropdown/checkbox per schema logic), uncheck others if checking this one
              if (checked && fieldType !== 'multiselect') {
                newOptions.forEach((opt, idx) => {
                  if (idx !== index) opt.isAnswer = false
                })
              }
              newOptions[index] = { ...option, isAnswer: checked }
              onOptionsChange(newOptions)
            }}
            className="scale-75 data-[state=checked]:bg-green-500"
            title="Mark as correct answer"
          />
        </div>
      )}

      <div className="flex-1">
        <Input
          value={option.optionLabel}
          className="h-9 border-0 bg-transparent shadow-none focus-visible:ring-0 focus-visible:ring-offset-0 px-0 font-medium"
          placeholder="Option Label"
          onChange={e => {
            const newOptions = [...options]
            newOptions[index] = { ...option, optionLabel: e.target.value || `option ${index + 1}` }
            onOptionsChange(newOptions)
          }}
        />
      </div>
      <div className="flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100 focus-within:opacity-100 mt-1">
        {options.length > 1 && (
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className="h-8 w-8 text-muted-foreground hover:text-destructive"
            onClick={() => onRemove(index)}>
            <Trash size={14} />
          </Button>
        )}
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="h-8 w-8 text-muted-foreground hover:text-primary"
          onClick={() => addField(index + 1)}>
          <Plus size={14} />
        </Button>
      </div>
    </div>
  )
}
