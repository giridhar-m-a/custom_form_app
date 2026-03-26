import { Textarea } from '../ui/textarea'

interface FormLongInputProps extends React.ComponentProps<'textarea'> {
  template?: (props: React.ComponentProps<'textarea'>) => React.ReactNode
}

export const FormLongInput = ({ template, ...props }: FormLongInputProps) => {
  return template ? template(props) : <Textarea {...props} />
}
