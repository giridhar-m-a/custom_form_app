'use client'
import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Globe, Loader2, Lock, Mail } from 'lucide-react'
import AuthFormSignUp from './AuthFormSignUp'

const AuthFormLogin = () => {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [isGoogleLoading, setIsGoogleLoading] = useState(false)
  const [isFormLoading, setIsFormLoading] = useState(false)
  const [isSignUp, setIsSignUp] = useState(false)

  const handleGoogleLogin = () => {
    setIsGoogleLoading(true)
    setTimeout(() => setIsGoogleLoading(false), 1500)
  }

  const handleFormLogin = (e: React.FormEvent) => {
    e.preventDefault()
    setIsFormLoading(true)
    setTimeout(() => setIsFormLoading(false), 1500)
  }

  return (
    <Card className="w-full max-w-md border-none shadow-xl bg-white/95 backdrop-blur-sm">
      <CardHeader className="text-center pb-4">
        <h2 className="text-3xl font-bold text-gray-800">{isSignUp ? 'Sign Up' : 'Welcome Back'}</h2>
        <p className="text-sm text-gray-500">{isSignUp ? 'Create' : 'Sign in to'} your FormGenius account</p>
      </CardHeader>

      <CardContent>
        <Button
          onClick={handleGoogleLogin}
          className="w-full mb-6 text-gray-700 border border-gray-300 bg-white hover:bg-gray-50 transition-all"
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

        <div className="relative my-4">
          <Separator />
          <span className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-white px-2 text-xs font-semibold text-gray-400 uppercase">
            or
          </span>
        </div>

        {!isSignUp && (
          <form onSubmit={handleFormLogin} className="space-y-4">
            <div>
              <Label htmlFor="email">Email</Label>
              <div className="relative mt-1">
                <Mail className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                <Input
                  id="email"
                  type="email"
                  placeholder="you@example.com"
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                  required
                  className="pl-10"
                  disabled={isFormLoading || isGoogleLoading}
                />
              </div>
            </div>

            <div>
              <Label htmlFor="password">Password</Label>
              <div className="relative mt-1">
                <Lock className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400" />
                <Input
                  id="password"
                  type="password"
                  placeholder="••••••••"
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  required
                  className="pl-10"
                  disabled={isFormLoading || isGoogleLoading}
                />
              </div>
            </div>

            <Button
              type="submit"
              className="w-full mt-6 bg-indigo-600 hover:bg-indigo-700 transition-all"
              size="lg"
              disabled={isFormLoading || isGoogleLoading}>
              {isFormLoading ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : 'Sign In'}
            </Button>
          </form>
        )}
        {isSignUp && <AuthFormSignUp />}
        <p className="mt-4 text-center text-sm text-gray-500">
          {isSignUp ? 'Already have an account? ' : "Don't have an account? "}
          <button
            type="button"
            className="text-indigo-600 hover:underline font-medium"
            onClick={() => setIsSignUp(!isSignUp)}>
            {isSignUp ? 'Log in' : 'Sign up'}
          </button>
        </p>
      </CardContent>
    </Card>
  )
}

export default AuthFormLogin
