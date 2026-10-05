import { Modal } from '@/components/common/Modal'
import { FormPage } from '@/components/forms/FormPage'
import { UpsertForm } from '@/components/forms/UpsertForm'
import { Button } from '@/components/ui/button'
import { CirclePlusIcon } from 'lucide-react'

const FormsPage = async () => {
  return (
    <div className="space-y-6">
      <Modal
        title="Create New Form"
        description="Add a new form to get started."
        trigger={
          <Button>
            {' '}
            <CirclePlusIcon /> Add New Form
          </Button>
        }>
        <UpsertForm />
      </Modal>
      <FormPage />
    </div>
  )
}

export default FormsPage
