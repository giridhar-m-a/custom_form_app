import { FormRender } from '@/components/FormRender/FormRender'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { getFieldsForResponse, getFormForResponse } from '@/services/api/forms/routes'
import { verifyInvitation } from '@/services/api/invitations/routes'
import { notFound } from 'next/navigation'

interface ResponsePageProps {
  searchParams: Promise<{
    token?: string
  }>
}

export default async function ResponsePage({ searchParams }: ResponsePageProps) {
  const { token } = await searchParams
  if (!token) notFound()

  const [form, fields, verify] = await Promise.all([
    getFormForResponse({ token }),
    getFieldsForResponse({ token }),
    verifyInvitation({ token })
  ])

  if (!form.data || !verify.data?.invitationId) notFound()

  return (
    <main className="mx-auto w-full max-w-4xl px-4 py-8 sm:px-6 sm:py-12">
      <Card className="w-full min-w-0 overflow-hidden">
        <CardHeader>
          <CardTitle>{form.data?.title}</CardTitle>
          <CardDescription>Fill the form to submit your response</CardDescription>
        </CardHeader>
        <CardContent className="min-w-0 px-4 sm:px-6">
          <FormRender
            fields={fields.data || []}
            formId={form.data.id}
            respondentId={verify.data.invitationId}
            token={token}
          />
        </CardContent>
      </Card>
    </main>
  )
}
