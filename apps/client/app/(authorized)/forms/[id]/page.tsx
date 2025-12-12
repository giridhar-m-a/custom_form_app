import { FormFieldContent } from '@/components/forms/FormFields'
import { FormFieldWrapper } from '@/components/forms/FormFields/FormFieldWrapper'
import { getFormById, getFormFields } from '@/services/api/forms/routes'
import { FormField, FormType } from '@/types/form.types'

interface FormDataPageProps {
  params: { id: string }
}

const FormDataPage: React.FC<FormDataPageProps> = async ({ params }) => {
  const { id } = await params
  if (id === 'new' || id === 'edit') {
    return <>not found</>
  }

  let formData: FormType | null = null
  let fields: FormField[] = []

  try {
    const resp = await getFormById({ id })
    if (resp.status === 200 && resp.data) {
      formData = resp.data
    }
    const fieldsResp = await getFormFields({ id })
    if (fieldsResp.status === 200 && fieldsResp.data) {
      fields = fieldsResp.data
    }
  } catch (error) {
    console.error(error)
    return <>not found</>
  }

  return (
    <div>
      <h1>{formData?.title}</h1>
      <p className="text-gray-500">{formData?.description}</p>
      <div className="rounded-2xl border bg-background w-full p-6 mt-8">
        <FormFieldContent
          formId={id}
          formTitle={id === 'new' ? 'Create New Form Fields' : 'Edit or Create Form Fields'}
          mode={fields.length > 0 ? 'edit' : 'create'}
          initialFields={fields}
        />
      </div>
    </div>
  )
}

export default FormDataPage
