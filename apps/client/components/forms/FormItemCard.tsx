import { FormType } from '@/types/form.types'
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '../ui/card'
import { format } from 'date-fns'

export const FormItemCard = ({ form }: { form: FormType }) => {
  return (
    <div className="w-[250px] h-fit border-2 max-h-[250px] p-4 rounded-[10px] space-y-4">
      <h2 className="text-2xl font-bold">{form.title}</h2>
      <div className="border-2 rounded-[10px] bg-primary/20">
        <p className="p-2 h-[100px] overflow-y-auto">{form.description}</p>
      </div>
      <p>created at {format(form.createdAt, 'dd/MM/yyyy')}</p>
    </div>
  )
}
