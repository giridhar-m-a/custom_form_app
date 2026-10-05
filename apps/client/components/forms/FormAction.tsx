import { FormType } from '@/types/form.types'
import { Modal } from '../common/Modal'
import { Button } from '../ui/button'
import { UpsertForm } from './UpsertForm'
import { HiEllipsisVertical, HiPencilSquare, HiTrash } from 'react-icons/hi2'
import { WarningModalContent } from '../common/WarningModalContent'
import { CommonDropdown } from '../common/CommonDropdown'
import { useState } from 'react'
import { useDeleteForm } from '@/hooks/queryHooks/useFormApp'
import Link from 'next/link'
import { MdOpenInNew } from 'react-icons/md'

interface FormActionProps {
  form: FormType
}

export const FormAction = ({ form }: FormActionProps) => {
  const [editOpen, setEditOpen] = useState(false)
  const [deleteOpen, setDeleteOpen] = useState(false)
  const { mutate: deleteMutate, isPending: deleteIsPending } = useDeleteForm()

  return (
    <>
      <CommonDropdown
        options={[
          {
            content: (
              <div className="flex items-center gap-2 cursor-pointer">
                <HiPencilSquare /> Edit
              </div>
            ),
            onSelect: e => {
              e.preventDefault()
              setEditOpen(true)
            }
          },
          {
            content: (
              <div className="flex items-center gap-2 text-red-500 cursor-pointer">
                <HiTrash /> Delete
              </div>
            ),
            onSelect: e => {
              e.preventDefault()
              setDeleteOpen(true)
            }
          },
          {
            content: (
              <Link href={`/forms/${form.id}`} className="flex items-center gap-2 cursor-pointer">
                <MdOpenInNew /> Open
              </Link>
            )
          }
        ]}
        trigger={
          <Button variant="ghost">
            <HiEllipsisVertical />
          </Button>
        }
      />

      <Modal
        title="Edit Form"
        description={`Edit form ${form.title}`}
        onOpenChange={setEditOpen}
        open={editOpen}
        trigger={<></>}>
        <UpsertForm data={form} formId={form.id} onOpenChange={setEditOpen} />
      </Modal>

      <Modal
        title="Delete Form"
        description={`Delete form ${form.title}`}
        onOpenChange={setDeleteOpen}
        open={deleteOpen}
        trigger={<></>}>
        <WarningModalContent
          message={`Are you sure you want to delete form ${form.title}? This action will delete all the responses associated with this form.`}
          handleDelete={() =>
            deleteMutate(
              { id: form.id },
              {
                onSuccess: () => {
                  setDeleteOpen(false)
                }
              }
            )
          }
          handleCancel={setDeleteOpen}
          buttonLabel="Delete"
          isLoading={deleteIsPending}
        />
      </Modal>
    </>
  )
}
