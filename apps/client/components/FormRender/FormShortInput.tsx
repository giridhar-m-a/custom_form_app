import React from 'react'
import { Input } from '../ui/input'

interface FormInputProps extends React.ComponentProps<'input'> {
  template?: (props: React.ComponentProps<'input'>) => React.ReactNode
}

export const FormShortInput = ({ template, ...props }: FormInputProps) => {
  return template ? template(props) : <Input {...props} />
}
