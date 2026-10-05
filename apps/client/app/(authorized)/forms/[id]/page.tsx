import { InvitationTable } from '@/components/forms/Invitations/InvitationTable'
import { getFormById } from '@/services/api/forms/routes'
import { FormType } from '@/types/form.types'

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
    console.error(error)
    return <>not found</>
  }

  return (
    <div>
      <h1>{formData?.title}</h1>
      <p className="text-gray-500">{formData?.description}</p>
      <div className="flex gap-4">
        <div className="basis-3/4"></div>
        <div className="rounded-2xl border basis-1/4 bg-background w-full p-6 mt-8 ">
          <InvitationTable formId={id} />
        </div>
      </div>
    </div>
  )
}

export default FormDataPage
