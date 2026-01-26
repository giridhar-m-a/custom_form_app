import { InputGroup, InputGroupAddon, InputGroupInput, InputGroupText } from '@/components/ui/input-group'
import { SearchIcon } from 'lucide-react'

interface SearchProps {
  placeholder?: string
}

export const Search: React.FC<React.ComponentProps<'input'>> = props => {
  return (
    <InputGroup>
      <InputGroupAddon>
        <InputGroupText>
          <SearchIcon />
        </InputGroupText>
      </InputGroupAddon>
      <InputGroupInput {...props} />
    </InputGroup>
  )
}
