import { FormType } from '@/types/form.types'
import { format } from 'date-fns'
import { FormAction } from './FormAction'
import { FormStatusBadge } from './forms.config'

export const FormItemCard = ({ form }: { form: FormType }) => {
  return (
    <div className="w-[250px] h-fit border-2 max-h-[250px] p-4 rounded-[10px] space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">{form.title}</h2>
        <FormAction form={form} />
      </div>
      <div className="border-2 rounded-[10px] bg-primary/20">
        <p className="p-2 h-[100px] overflow-y-auto">{form.description}</p>
      </div>
      <div className="flex items-center justify-between">
        <p>{format(form.updatedAt, 'dd/MM/yyyy')}</p>
        <FormStatusBadge status={form.status} />
      </div>
    </div>
  )
}
