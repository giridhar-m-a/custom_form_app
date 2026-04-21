import { Modal } from '@/components/common/Modal'
import { InvitationTable } from '@/components/forms/Invitations/InvitationTable'
import { ResponseTable } from '@/components/forms/Responses/ResponseTable'
import { UpsertForm } from '@/components/forms/UpsertForm'
import { Button } from '@/components/ui/button'
import { getFormById } from '@/services/api/forms/routes'
import { FormType } from '@/types/form.types'
import { PiNotePencil } from 'react-icons/pi'
import Link from 'next/link'
import { Pencil } from 'lucide-react'
import { Badge } from '@/components/ui/badge'
import { FormStatusBadge } from '@/components/forms/forms.config'

interface FormDataPageProps {
  params: { id: string }
}

const FormDataPage: React.FC<FormDataPageProps> = async ({ params }) => {
  const { id } = await params
  if (id === 'new' || id === 'edit') {
    return <>not found</>
  }

  let formData: FormType | null = null

  try {
    const resp = await getFormById({ id })
    if (resp.status === 200 && resp.data) {
      formData = resp.data
    }
  } catch (error) {
    console.log(error)
    return <>not found</>
  }

  return (
    <div>
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold">
            {formData?.title}{' '}
            <span>
              <FormStatusBadge status={formData?.status || 'draft'} />
            </span>
          </h1>
          <p className="text-gray-500">{formData?.description}</p>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <Modal
            title="Edit Form"
            description="Edit the form to get started."
            trigger={
              <Button>
                {' '}
                <PiNotePencil /> Edit Form
              </Button>
            }>
            <UpsertForm formId={id} data={formData || undefined} />
          </Modal>
          <Link href={`/forms/edit/${id}`}>
            <Button>
              <Pencil />
              Edit Form Fields
            </Button>
          </Link>
        </div>
      </div>
      <div className="flex gap-4">
        <div className="basis-3/4 rounded-2xl border bg-background w-full p-6 mt-8 ">
          <ResponseTable formId={id} />
        </div>
        <div className="rounded-2xl border basis-1/4 bg-background w-full p-6 mt-8 ">
          <InvitationTable formId={id} status={formData?.status || 'draft'} />
        </div>
      </div>
    </div>
  )
}

export default FormDataPage
