import { FormRender } from '@/components/FormRender/FormRender'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { FormFields } from '@/services/dummydata/formfields'

export default function ResponsePage() {
  const data = FormFields

  return (
    <main className="flex items-center justify-center w-screen h-screen">
      <Card className="w-full max-w-6xl">
        <CardHeader>
          <CardTitle>Response</CardTitle>
          <CardDescription>Fill the form to submit your response</CardDescription>
        </CardHeader>
        <CardContent>
          <FormRender fields={data} />
        </CardContent>
      </Card>
    </main>
  )
}
