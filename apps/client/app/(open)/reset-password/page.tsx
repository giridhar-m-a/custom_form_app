import AuthFormResetPassword from '@/components/home/AuthFormResetPassword'
import HomeContent from '@/components/home/HomeContent'

interface ResetPasswordPageProps {
  searchParams: Promise<{
    token?: string
  }>
}

export default async function ResetPasswordPage({ searchParams }: ResetPasswordPageProps) {
  const { token } = await searchParams

  return (
    <HomeContent>
      <AuthFormResetPassword token={token} />
    </HomeContent>
  )
}
