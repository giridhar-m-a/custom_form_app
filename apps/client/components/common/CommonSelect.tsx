import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

export interface CommonSelectProps {
  options: Option[]
  placeholder: string
  value?: string
  onChange: (value: string) => void
  className?: string
}

export interface Option {
  value: string
  label: string
}

export function CommonSelect({ options, placeholder, value, onChange, className }: CommonSelectProps) {
  return (
    <Select onValueChange={onChange} value={value}>
      <SelectTrigger className={className || 'w-[180px]'}>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {options.map(option => (
          <SelectItem key={option.value} value={option.value}>
            {option.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}
