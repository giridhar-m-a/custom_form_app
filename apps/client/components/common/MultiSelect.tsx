import * as React from 'react'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Check, X } from 'lucide-react'
import { cn } from '@/lib/utils'
import { Option } from './CommonSelect'
import { Badge } from '@/components/ui/badge'

export interface MultiSelectProps {
  options: Option[]
  placeholder: string
  value?: string[]
  onChange: (value: string[]) => void
  className?: string
}

export function MultiSelect({ options, placeholder, value = [], onChange, className }: MultiSelectProps) {
  const toggleValue = (val: string) => {
    const nextValue = value.includes(val) ? value.filter(v => v !== val) : [...value, val]
    onChange(nextValue)
  }

  const selectedOptions = options.filter(o => value.includes(o.value))

  return (
    <Select>
      <SelectTrigger className={cn('w-full h-auto min-h-10 py-2 flex-wrap gap-1', className)}>
        <div className="flex flex-wrap gap-1">
          {selectedOptions.length > 0 ? (
            selectedOptions.map(option => (
              <Badge
                key={option.value}
                variant="secondary"
                className="flex items-center gap-1 px-1 py-0 text-xs"
                onClick={e => {
                  e.stopPropagation()
                  toggleValue(option.value)
                }}>
                {option.label}
                <X className="h-3 w-3" />
              </Badge>
            ))
          ) : (
            <span className="text-muted-foreground">{placeholder}</span>
          )}
        </div>
      </SelectTrigger>

      <SelectContent>
        {options.map(option => {
          const isSelected = value.includes(option.value)

          return (
            <div
              key={option.value}
              className={cn(
                'relative flex w-full cursor-default select-none items-center rounded-sm py-1.5 pl-8 pr-2 text-sm outline-none hover:bg-accent hover:text-accent-foreground data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
                isSelected && 'bg-accent/50'
              )}
              onClick={e => {
                e.preventDefault()
                e.stopPropagation()
                toggleValue(option.value)
              }}>
              <span className="absolute left-2 flex h-3.5 w-3.5 items-center justify-center">
                {isSelected && <Check className="h-4 w-4" />}
              </span>
              <span>{option.label}</span>
            </div>
          )
        })}
      </SelectContent>
    </Select>
  )
}
