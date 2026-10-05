import { InputGroup, InputGroupAddon, InputGroupInput, InputGroupText } from '@/components/ui/input-group'
import { SearchIcon } from 'lucide-react'
import { cn } from '@/lib/utils'

interface SearchProps {
  placeholder?: string
}

export const Search: React.FC<React.ComponentProps<'input'>> = ({ className, ...props }) => {
  return (
    <InputGroup className={cn('min-w-0', className)}>
      <InputGroupAddon>
        <InputGroupText>
          <SearchIcon />
        </InputGroupText>
      </InputGroupAddon>
      <InputGroupInput {...props} />
    </InputGroup>
  )
}
