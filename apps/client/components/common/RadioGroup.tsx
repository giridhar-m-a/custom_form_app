import { Label } from '@/components/ui/label'
import { RadioGroup as RadioGroupComponent, RadioGroupItem } from '@/components/ui/radio-group'
import { Option } from './CommonSelect'

export interface RadioGroupProps {
  options: Option[]
  value?: string
  onChange: (value: string) => void
}

export function RadioGroup({ options, value, onChange }: RadioGroupProps) {
  return (
    <RadioGroupComponent className="flex flex-col gap-2" value={value} onValueChange={onChange}>
      {options.map(option => (
        <div className="flex items-center gap-3">
          <RadioGroupItem value={option.value} id={option.value} />
          <Label htmlFor={option.value}>{option.label}</Label>
        </div>
      ))}
    </RadioGroupComponent>
  )
}
