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
    <div className="min-w-0">
      <div className="flex flex-col gap-4 xl:flex-row xl:justify-between xl:items-center">
        <div>
          <h1 className="text-2xl font-bold">
            {formData.title}{' '}
            <span><FormStatusBadge status={formData.status || 'draft'} /></span>
          </h1>
          <p className="text-gray-500">{formData.description}</p>
        </div>
        <div className="flex flex-wrap gap-3">
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
      <div className="mt-8 grid min-w-0 grid-cols-1 gap-4 2xl:grid-cols-2">
        <div className="min-w-0 rounded-2xl border bg-background p-4 md:p-6">
          <ResponseTable formId={id} />
        </div>
        <div className="min-w-0 rounded-2xl border bg-background p-4 md:p-6">
          <InvitationTable formId={id} status={formData.status || 'draft'} access={formData.access || 'restricted'} />
        </div>
      </div>
    </div>
  )
}

export default FormDataPage
