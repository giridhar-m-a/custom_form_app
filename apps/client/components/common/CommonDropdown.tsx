import { useState } from 'react'
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '../ui/dropdown-menu'

interface DropdownOption {
  content: React.ReactNode
  onSelect?: (e: Event) => void
}

interface CommonDropdownProps {
  options: DropdownOption[]
  trigger: React.ReactNode | string
}

export const CommonDropdown = ({ options, trigger }: CommonDropdownProps) => {
  const [open, setOpen] = useState(false)
  return (
    <DropdownMenu open={open} onOpenChange={setOpen}>
      <DropdownMenuTrigger>{trigger}</DropdownMenuTrigger>
      <DropdownMenuContent>
        {options.map((option, index) => (
          <DropdownMenuItem
            key={index}
            onSelect={e => {
              option.onSelect?.(e)
              setOpen(false)
            }}>
            {option.content}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
