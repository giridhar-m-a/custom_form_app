import { FormRender } from '@/components/FormRender/FormRender'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { getFieldsForResponse, getFormForResponse } from '@/services/api/forms/routes'
import { verifyInvitation } from '@/services/api/invitations/routes'

interface ResponsePageProps {
  searchParams: Promise<{
    token: string
  }>
}

export default async function ResponsePage({ searchParams }: ResponsePageProps) {
  const { token } = await searchParams

  const form = await getFormForResponse({ token })
  const fields = await getFieldsForResponse({ token })
  const verify = await verifyInvitation({ token })

  return (
    <main className="p-24">
      <Card className="w-full">
        <CardHeader>
          <CardTitle>{form.data?.title}</CardTitle>
          <CardDescription>Fill the form to submit your response</CardDescription>
        </CardHeader>
        <CardContent>
          <FormRender
            fields={fields.data || []}
            formId={form.data?.id || ''}
            respondentId={verify.data?.invitationId || ''}
            token={token}
          />
        </CardContent>
      </Card>
    </main>
  )
}
