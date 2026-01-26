'use client'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger
} from '@/components/ui/dropdown-menu'
import { FieldType } from '@/types/form.types'
import { MdAdd, MdContentCopy, MdDelete, MdEdit, MdSave, MdShortText } from 'react-icons/md'
import { FIELD_TYPE_OPTIONS } from './formFields.config'

interface FormFieldOptionsProps {
  onAddField: (type: FieldType) => void
  onDuplicateField: () => void
  onRemoveField: () => void
  isEdit: boolean
  setEdit: () => void
  handleSave: () => void
  className?: string
}

export const FormFieldOptions = ({
  onAddField,
  onDuplicateField,
  onRemoveField,
  isEdit,
  setEdit,
  handleSave,
  className
}: FormFieldOptionsProps) => {
  return (
    <div className={`flex items-center ${isEdit ? 'gap-2' : 'gap-4'} ${className || ''}`}>
      {!isEdit ? (
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="h-8 w-8 text-muted-foreground hover:text-foreground"
          onClick={setEdit}>
          <MdEdit size={16} />
        </Button>
      ) : (
        <Button type="button" variant="default" size="icon" className="h-8 w-8" onClick={handleSave}>
          <MdSave size={16} />
        </Button>
      )}

      <div className="mx-1 h-4 w-px bg-border" />

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            type="button"
            variant="ghost"
            size="icon"
            className="h-8 w-8 text-muted-foreground hover:text-primary">
            <MdAdd size={18} />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-56 max-h-96 overflow-y-auto">
          <DropdownMenuLabel>Add New Field</DropdownMenuLabel>
          <DropdownMenuSeparator />
          {FIELD_TYPE_OPTIONS.map(option => (
            <DropdownMenuItem onClick={() => onAddField(option.value)}>
              {option.icon && <option.icon className="mr-2 h-4 w-4" />}
              <span>{option.label}</span>
            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>

      <Button
        type="button"
        variant="ghost"
        size="icon"
        className="h-8 w-8 text-muted-foreground hover:text-primary"
        onClick={onDuplicateField}>
        <MdContentCopy size={16} />
      </Button>
      <Button
        type="button"
        variant="ghost"
        size="icon"
        className="h-8 w-8 text-muted-foreground hover:text-destructive"
        onClick={onRemoveField}>
        <MdDelete size={18} />
      </Button>
    </div>
  )
}
