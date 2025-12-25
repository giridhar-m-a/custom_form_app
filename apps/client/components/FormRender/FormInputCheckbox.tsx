import * as React from 'react'
import { Checkbox } from '@/components/ui/checkbox'
import { FieldOption } from '@/types/form.types'
import { Label } from '../ui/label'

interface FormInputCheckboxProps {
  options: FieldOption[]
  value?: string // selected optionIds
  onChange?: (value: string) => void
  template?: (params: { option: FieldOption; checked: boolean; toggle: () => void }) => React.ReactNode
}

export const FormInputCheckbox = ({ options, value = '', onChange, template }: FormInputCheckboxProps) => {
  const toggleValue = (optionId: string) => {
    onChange?.(optionId)
  }

  return (
    <div className="gap-2 flex items-center flex-wrap">
      {options.map(option => {
        const checked = value.includes(option.optionId)

        if (template) {
          return (
            <div key={option.optionId}>
              {template({
                option,
                checked,
                toggle: () => toggleValue(option.optionId)
              })}
            </div>
          )
        }

        return (
          <div key={option.optionId} className="flex items-center gap-2">
            <Checkbox checked={checked} onCheckedChange={() => toggleValue(option.optionId)} id={option.optionId} />
            <Label htmlFor={option.optionId} className="text-sm">
              {option.optionLabel}
            </Label>
          </div>
        )
      })}
    </div>
  )
}
