import { VariantProps } from 'class-variance-authority'
import { Button } from '../ui/button'
import { CustomLoader } from './CustomLoader'

type size = 'default' | 'sm' | 'lg' | 'icon' | 'icon-sm' | 'icon-lg' | null | undefined

type variant = 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link' | null | undefined

export const SubmitButton: React.FC<
  React.ComponentProps<'button'> & { isLoading?: boolean } & { size?: size } & { variant?: variant }
> = ({ disabled, isLoading, children, variant, size, ...props }) => {
  return (
    <Button type="submit" {...props} disabled={disabled || isLoading} variant={variant} size={size}>
      {isLoading && <CustomLoader />}
      {children}
    </Button>
  )
}
