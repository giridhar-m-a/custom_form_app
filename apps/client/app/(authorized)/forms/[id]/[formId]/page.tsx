import { FormFieldContent } from '@/components/forms/FormFields'
import { getFormById, getFormFields } from '@/services/api/forms/routes'
import { FormField, FormType } from '@/types/form.types'
import { redirect } from 'next/navigation'

interface FormFieldUpsertProps {
  params: { id: string; formId: string }
}

const FormFieldUpsert = async ({ params }: FormFieldUpsertProps) => {
  const { id, formId } = await params

  let formData: FormType | null = null
  let fields: FormField[] = []

  try {
    const resp = await getFormById({ id: formId })
    if (resp.status === 200 && resp.data) {
      formData = resp.data
    }
    const fieldsResp = await getFormFields({ id: formId })
    if (fieldsResp.status === 200 && fieldsResp.data) {
      fields = fieldsResp.data
    } else if (fieldsResp.status === 200 && !fieldsResp.data?.length) {
      redirect(`/forms/new/${formId}`)
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
          formId={formId}
          formTitle={id === 'new' ? 'Create New Form Fields' : 'Edit or Create Form Fields'}
          mode={fields.length > 0 ? 'edit' : 'create'}
          initialFields={fields}
        />
      </div>
    </div>
  )
}

export default FormFieldUpsert
