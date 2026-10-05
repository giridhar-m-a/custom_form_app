'use client'
import { SignInSchema, SignInSchemaType } from '@/app/schemas/auth.schemas'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Separator } from '@/components/ui/separator'
import { useCredentialAuth, useGoogleAuth } from '@/hooks/queryHooks/useAuth'
import { zodResolver } from '@hookform/resolvers/zod'
import { Globe, Loader2, Mail } from 'lucide-react'
import { useMemo, useState } from 'react'
import { useForm } from 'react-hook-form'
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '../ui/form'
import AuthFormSignUp from './AuthFormSignUp'
import { useGoogleLogin } from '@react-oauth/google'
import AuthFormReset from './AuthFormReset'
import TempUser from './TempUser'

const AuthFormLogin = () => {
  const [isSignUp, setIsSignUp] = useState<'login' | 'signup' | 'forget'>('login')
  const { mutate: credentialLogin, isPending: isFormLoading } = useCredentialAuth()
  const { mutate: googleLogin, isPending: isGoogleLoading } = useGoogleAuth()

  const isPending = useMemo(() => isFormLoading || isGoogleLoading, [isFormLoading, isGoogleLoading])

  const form = useForm<SignInSchemaType>({
    resolver: zodResolver(SignInSchema),
    defaultValues: {
      email: '',
      password: ''
    }
  })

  const { handleSubmit, control } = form

  const handleGoogleLogin = useGoogleLogin({
    onSuccess: async (codeResponse: any) => {
      googleLogin(codeResponse.code)
    },
    flow: 'auth-code'
  })

  const handleFormLogin = (data: SignInSchemaType) => {
    credentialLogin(data)
  }

  return (
    <Card className="w-full max-w-md border-none shadow-xl bg-white/95 backdrop-blur-sm">
      {isSignUp !== 'forget' && (
        <CardHeader className="text-center pb-4">
          <h2 className="text-3xl font-bold text-gray-800">{isSignUp === 'signup' ? 'Sign Up' : 'Welcome Back'}</h2>
          <p className="text-sm text-gray-500">
            {isSignUp === 'signup' ? 'Create' : 'Sign in to'} your FormGenius account
          </p>
        </CardHeader>
      )}

      <CardContent>
        <Button
          onClick={handleGoogleLogin}
          className="w-full hidden mb-6 text-gray-700 hover:text-gray-700! border! border-gray-300! bg-white! hover:bg-gray-50! transition-all"
          variant="outline"
          size="lg"
          disabled={isGoogleLoading || isFormLoading}>
          {isGoogleLoading ? (
            <Loader2 className="mr-2 h-5 w-5 animate-spin" />
          ) : (
            <Globe className="mr-3 h-6 w-6 text-indigo-600" />
          )}
          Continue with Google
        </Button>

        <div className="relative hidden my-4">
          <Separator />
          <span className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-white px-2 text-xs font-semibold text-gray-400 uppercase">
            or
          </span>
        </div>

        {isSignUp === 'login' && (
          <Form {...form}>
            <form onSubmit={handleSubmit(handleFormLogin)} className="space-y-4">
              <div>
                <FormField
                  control={control}
                  name="email"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel htmlFor="email" className="text-black">
                        Email
                      </FormLabel>
                      <FormControl>
                        <div className="relative mt-1">
                          <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                          <Input
                            id="email"
                            type="email"
                            placeholder="you@example.com"
                            {...field}
                            required
                            className="pl-10 text-gray-700 border! border-gray-300! bg-white! hover:bg-gray-50!"
                            disabled={isFormLoading || isGoogleLoading}
                          />
                        </div>
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>

              <div>
                <FormField
                  control={control}
                  name="password"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel htmlFor="password" className="text-black">
                        Password
                      </FormLabel>
                      <FormControl>
                        <div className="relative mt-1">
                          <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                          <Input
                            id="password"
                            type="password"
                            placeholder="••••••••"
                            {...field}
                            required
                            className="pl-10 text-gray-700 border! border-gray-300! bg-white! hover:bg-gray-50!"
                            disabled={isFormLoading || isGoogleLoading}
                          />
                        </div>
                      </FormControl>
                      <FormMessage />
                      <p
                        className="underline text-right text-primary hover:text-indigo-700"
                        onClick={() => setIsSignUp('forget')}>
                        Forgot Password?
                      </p>
                    </FormItem>
                  )}
                />
              </div>

              <Button
                type="submit"
                className="w-full mt-6 bg-indigo-600 hover:bg-indigo-700 transition-all"
                size="lg"
                disabled={isPending}>
                {isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : 'Sign In'}
              </Button>
            </form>
          </Form>
        )}
        {isSignUp === 'signup' && <AuthFormSignUp />}
        {isSignUp === 'forget' && <AuthFormReset />}
        <TempUser />
        <p className="mt-4 text-center text-sm text-gray-500">
          {isSignUp === 'login' || isSignUp === 'forget' ? 'Already have an account? ' : "Don't have an account? "}
          <button
            type="button"
            className="text-indigo-600 hover:underline font-medium"
            onClick={() => setIsSignUp(prev => (prev === 'login' ? 'signup' : 'login'))}>
            {isSignUp === 'login' ? 'Sign Up' : 'Login'}
          </button>
        </p>
      </CardContent>
    </Card>
  )
}

export default AuthFormLogin
