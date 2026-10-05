import { ResponsePage } from '@/components/forms/Responses/ResponsePage'

interface ReportsPageProps {
  params: Promise<{ id: string; submissionId: string }>
}

const ReportsPage = async ({ params }: ReportsPageProps) => {
  const { id, submissionId } = await params
  return (
    <div>
      <ResponsePage params={{ id, submissionId }} />
    </div>
  )
}

export default ReportsPage
