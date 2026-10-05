import type { Metadata } from 'next'
import { Geist, Geist_Mono } from 'next/font/google'
import './globals.css'
import { TanstackQueryWrapper } from '@/components/wrapper/TanstackQueryWrapper'
import { Toaster } from 'react-hot-toast'
import { GoogleOAuthProvider } from '@react-oauth/google'

import StoreProvider from '@/components/wrapper/StoreProvider'
import { ThemeProvider } from '@/components/theme/theme-provider'

const geistSans = Geist({
  variable: '--font-geist-sans',
  subsets: ['latin']
})

const geistMono = Geist_Mono({
  variable: '--font-geist-mono',
  subsets: ['latin']
})

export const metadata: Metadata = {
  title: 'Form Genius',
  description: 'Form Genius is a powerful form builder and response management platform that allows users to create, customize, and analyze forms with ease. With Form Genius, you can design forms for surveys, feedback, registrations, and more, while efficiently managing responses and gaining valuable insights.',
}

export default function RootLayout({
  children
}: Readonly<{
  children: React.ReactNode
}>) {
  return (
    <html lang="en">
      <GoogleOAuthProvider clientId="1091902491223-a8e9ga2qdf32a9tuqdf32p5j95jus9tt.apps.googleusercontent.com">
        <body className={`${geistSans.variable} ${geistMono.variable} antialiased`}>
          <ThemeProvider attribute="class" defaultTheme="system" enableSystem disableTransitionOnChange>
            <Toaster position="top-center" reverseOrder={false} />
            <StoreProvider>
              <TanstackQueryWrapper>{children}</TanstackQueryWrapper>
            </StoreProvider>
          </ThemeProvider>
        </body>
      </GoogleOAuthProvider>
    </html>
  )
}
