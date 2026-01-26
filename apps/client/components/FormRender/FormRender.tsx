'use client'

import { FormField } from '@/types/form.types'
import { FormInputWrapper } from './FormInputwrapper'
import { useForm } from 'react-hook-form'
import { Form } from '../ui/form'
import { ScrollArea } from '../ui/scroll-area'

export const FormRender = ({ fields }: { fields: FormField[] }) => {
  const form = useForm()
  const { control } = form
  return (
    <div>
      <ScrollArea className="h-[calc(100vh-20rem)]">
        <Form {...form}>
          <form className="space-y-4 px-4">
            {fields.map(field => (
              <FormInputWrapper key={field.fieldId} formField={field} control={control} />
            ))}
          </form>
        </Form>
      </ScrollArea>
    </div>
  )
}
