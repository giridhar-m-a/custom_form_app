import PasswordChange from '@/components/account/PasswordChange'
import { Profile } from '@/components/account/Profile'

export default async function AccountPage() {
  return (
    <div>
      <div className="mb-8 text-center sm:text-left space-y-4">
        <h1 className="text-3xl font-bold tracking-tight">Account Settings</h1>
        <p className="mt-2">Manage your profile information and security settings.</p>
        <div className="flex flex-2 gap-8 pt-10">
          <Profile />
          <PasswordChange />
        </div>
      </div>
    </div>
  )
}
