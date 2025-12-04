import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'

interface CommonSelectProps {
  options: Option[]
  placeholder: string
  value?: string
  onChange: (value: string) => void
}

export interface Option {
  value: string
  label: string
}

export function CommonSelect({ options, placeholder, value, onChange }: CommonSelectProps) {
  return (
    <Select onValueChange={onChange} value={value}>
      <SelectTrigger className="w-[180px]">
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
