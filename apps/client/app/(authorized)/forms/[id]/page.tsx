'use client'

import { Modal } from '@/components/common/Modal'
import { InvitationTable } from '@/components/forms/Invitations/InvitationTable'
import { ResponseTable } from '@/components/forms/Responses/ResponseTable'
import { UpsertForm } from '@/components/forms/UpsertForm'
import { Button } from '@/components/ui/button'
import { FormStatusBadge } from '@/components/forms/forms.config'
import { useGetFormById } from '@/hooks/queryHooks/useFormApp'
import { Pencil, SquarePen } from 'lucide-react'
import Link from 'next/link'
import { useParams } from 'next/navigation'

const FormDataPage = () => {
  const { id } = useParams<{ id: string }>()
  const { data: response, isLoading } = useGetFormById(id)
  const formData = response?.data

  if (id === 'new' || id === 'edit') return <>not found</>
  if (isLoading) return <div className="p-6">Loading form...</div>
  if (!formData) return <>not found</>

  return (
    <div>
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold">
            {formData.title}{' '}
            <span><FormStatusBadge status={formData.status || 'draft'} /></span>
          </h1>
          <p className="text-gray-500">{formData.description}</p>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <Modal
            title="Edit Form"
            description="Edit the form to get started."
            trigger={<Button><SquarePen /> Edit Form</Button>}>
            <UpsertForm formId={id} data={formData} />
          </Modal>
          <Link href={`/forms/edit/${id}`}>
            <Button><Pencil /> Edit Form Fields</Button>
          </Link>
        </div>
      </div>
      <div className="flex flex-col lg:flex-row gap-4 mt-8">
        <div className="lg:basis-3/4 rounded-2xl border bg-background w-full p-6">
          <ResponseTable formId={id} />
        </div>
        <div className="lg:basis-1/4 rounded-2xl border bg-background w-full p-6">
          <InvitationTable formId={id} status={formData.status || 'draft'} access={formData.access || 'restricted'} />
        </div>
      </div>
    </div>
  )
}

export default FormDataPage
