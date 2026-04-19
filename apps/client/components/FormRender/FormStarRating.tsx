import * as React from 'react'
import { Star } from 'lucide-react'
import { cn } from '@/lib/utils'

interface FormStarRatingProps {
  value?: number
  max?: number
  onChange?: (value: number) => void
  disabled?: boolean
  template?: (props: {
    value: number
    max: number
    setValue: (value: number) => void
    disabled?: boolean
  }) => React.ReactNode
}

export const FormStarRating = ({ value = 0, max = 5, onChange, disabled, template }: FormStarRatingProps) => {
  const setValue = (val: number) => {
    if (!disabled) {
      onChange?.(val)
    }
  }

  if (template) {
    return template({
      value,
      max,
      setValue,
      disabled
    })
  }

  return (
    <div className="flex gap-1">
      {Array.from({ length: max }).map((_, index) => {
        const starValue = index + 1
        const isActive = starValue <= value

        return (
          <button
            key={starValue}
            type="button"
            onClick={() => setValue(starValue)}
            disabled={disabled}
            aria-label={`Rate ${starValue} star`}
            className={cn('focus:outline-none', disabled && 'cursor-not-allowed opacity-50')}>
            <Star
              className={cn(
                'h-6 w-6 transition-colors',
                isActive ? 'fill-yellow-400 text-yellow-400' : 'text-muted-foreground'
              )}
            />
          </button>
        )
      })}
    </div>
  )
}
