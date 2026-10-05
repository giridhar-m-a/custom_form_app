'use client'

import { Modal } from '@/components/common/Modal'
import { SubmitButton } from '@/components/common/SubmitButton'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useCreateAnonymousInvitation } from '@/hooks/queryHooks/useInvitations'
import { format } from 'date-fns'
import { Copy, Link } from 'lucide-react'
import { useState } from 'react'
import copy from 'copy-to-clipboard'
import toast from 'react-hot-toast'

export const AnonymousInvitation = ({ formId }: { formId: string }) => {
  const { mutate: createAnonymousInvitation, isPending: createAnonymousInvitationIsPending } =
    useCreateAnonymousInvitation()
  const [open, setOpen] = useState(false)
  const [url, setUrl] = useState<string>()
  const [expiresIn, setExpiresIn] = useState<string>()

  const copyToClipboard = (text: string) => {
    copy(text)
    toast.success('Link copied to clipboard')
  }

  const handleGenerateLink = () => {
    createAnonymousInvitation(
      { data: { formId } },
      {
        onSuccess: data => {
          console.log(data)
          setUrl(`${process.env.NEXT_PUBLIC_FRONTEND_URL}/user-response?token=${data.data?.token}`)
          setExpiresIn(data.data?.expiresIn || '')
        }
      }
    )
  }

  console.log(expiresIn)

  return (
    <Modal
      title="Anonymous Invitation"
      description="Get Anonymous Invitation Link"
      trigger={
        <Button variant={'outline'}>
          <Link />
        </Button>
      }
      open={open}
      onOpenChange={setOpen}>
      <div>
        {!url && (
          <SubmitButton isLoading={createAnonymousInvitationIsPending} onClick={handleGenerateLink} className="w-full">
            Generate Link
          </SubmitButton>
        )}
        {url && (
          <>
            <div className="flex items-center gap-2">
              <Input type="text" value={url} readOnly disabled />
              <Button onClick={() => copyToClipboard(url || '')} variant={'outline'}>
                <Copy />
              </Button>
            </div>
            <p className="text-xs text-muted-foreground">
              {' '}
              Invitation Expores at {format(parseGoTime(expiresIn || ''), 'dd-MMM-yyyy HH:mm')}
            </p>
          </>
        )}
      </div>
    </Modal>
  )
}

export function parseGoTime(goTime: string): Date {
  if (!goTime) {
    throw new Error('Invalid input')
  }

  // Remove monotonic clock part (m=+...)
  const cleaned = goTime.split(' m=')[0]

  // Convert to ISO-compatible string
  const isoString = cleaned
    .replace(' ', 'T') // "YYYY-MM-DD HH:mm:ss" → "YYYY-MM-DDTHH:mm:ss"
    .replace(' +0000 UTC', 'Z') // UTC → Z

  return new Date(isoString)
}
