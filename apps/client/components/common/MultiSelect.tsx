import * as React from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Check } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Option } from './CommonSelect'

export interface MultiSelectProps {
  options: Option[]
  placeholder: string
  value?: string[]
  onChange: (value: string[]) => void
  className?: string
}

export function MultiSelect({ options, placeholder, value = [], onChange, className }: MultiSelectProps) {
  const toggleValue = (val: string) => {
    if (value.includes(val)) {
      onChange(value.filter(v => v !== val))
    } else {
      onChange([...value, val])
    }
  }

  const selectedLabels = options
    .filter(o => value.includes(o.value))
    .map(o => o.label)
    .join(', ')

  return (
    <Select>
      <SelectTrigger className={cn('w-full', className)}>
        <SelectValue placeholder={placeholder} asChild>
          <span className="truncate">{selectedLabels || placeholder}</span>
        </SelectValue>
      </SelectTrigger>

      <SelectContent>
        {options.map(option => {
          const isSelected = value.includes(option.value)

          return (
            <SelectItem
              key={option.value}
              value={option.value}
              onSelect={e => {
                e.preventDefault() // ⛔ prevent default single-select behavior
                toggleValue(option.value)
              }}>
              <div className="flex items-center gap-2">
                <Check className={cn('h-4 w-4', isSelected ? 'opacity-100' : 'opacity-0')} />
                <span>{option.label}</span>
              </div>
            </SelectItem>
          )
        })}
      </SelectContent>
    </Select>
  )
}
