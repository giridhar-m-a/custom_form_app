'use client'

import { FormFieldContent } from '@/components/forms/FormFields'
import { CustomLoader } from '@/components/common/CustomLoader'
import { useGetFormById, useGetFormFields } from '@/hooks/queryHooks/useFormApp'
import { useParams } from 'next/navigation'

const FormFieldUpsert = () => {
  const { id, formId } = useParams<{ id: string; formId: string }>()
  const { data: formResponse, isLoading: isFormLoading } = useGetFormById(formId)
  const { data: fieldsResponse, isLoading: areFieldsLoading } = useGetFormFields(formId)
  const fields = fieldsResponse?.data || []

  if (isFormLoading || areFieldsLoading) {
    return <div className="flex min-h-64 items-center justify-center"><CustomLoader /></div>
  }

  const formData = formResponse?.data
  if (!formData) return <>not found</>

  return (
    <div>
      <h1>{formData.title}</h1>
      <p className="text-gray-500">{formData.description}</p>
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
