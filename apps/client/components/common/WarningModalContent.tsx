import { Button } from '../ui/button'
import { SubmitButton } from './SubmitButton'

interface WarningModalContentProps {
  message: string
  handleDelete: () => void
  buttonLabel: string
  handleCancel?: (val: boolean) => void
  isLoading?: boolean
}

export const WarningModalContent = ({
  message,
  handleDelete,
  buttonLabel,
  handleCancel,
  isLoading
}: WarningModalContentProps) => {
  return (
    <div>
      <p>{message}</p>
      <div className="flex items-center justify-end gap-2">
        {handleCancel && (
          <Button variant="outline" onClick={() => handleCancel(false)}>
            Cancel
          </Button>
        )}
        <SubmitButton variant="destructive" onClick={() => handleDelete()} isLoading={isLoading} disabled={isLoading}>
          {buttonLabel}
        </SubmitButton>
      </div>
    </div>
  )
}
