import { CommonSelect, CommonSelectProps } from '../common/CommonSelect'

interface FormInputSelectProps extends CommonSelectProps {
  template?: (params: CommonSelectProps) => React.ReactNode
}

export const FormInputSelect = ({ template, options, ...props }: FormInputSelectProps) => {
  return template ? (
    template({ options, ...props })
  ) : (
    <CommonSelect onChange={value => props.onChange?.(value)} options={options} placeholder={props.placeholder} />
  )
}
