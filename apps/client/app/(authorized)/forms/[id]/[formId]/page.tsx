interface FormFieldUpsertProps {
  params: { id: string; formId: string }
}

const FormFieldUpsert = async ({ params }: FormFieldUpsertProps) => {
  return (
    <div>
      <h1>
        Form Field Upsert {params.id} {params.formId}
      </h1>
    </div>
  )
}

export default FormFieldUpsert
