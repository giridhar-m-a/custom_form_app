import Link from 'next/link'
import { ArrowLeft, FileQuestion } from 'lucide-react'
import { Button } from '@/components/ui/button'

export default function NotFound() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-background px-6 py-16 text-foreground">
      <div className="w-full max-w-lg text-center">
        <div className="mx-auto mb-6 flex size-16 items-center justify-center rounded-2xl bg-primary/10 text-primary">
          <FileQuestion className="size-8" aria-hidden="true" />
        </div>
        <p className="text-sm font-semibold uppercase tracking-[0.2em] text-primary">404 · Page not found</p>
        <h1 className="mt-4 text-3xl font-bold tracking-tight sm:text-4xl">We can’t find that page</h1>
        <p className="mx-auto mt-4 max-w-md text-muted-foreground">
          The link may be incorrect, expired, or the page may have been moved.
        </p>
        <Button asChild className="mt-8">
          <Link href="/">
            <ArrowLeft aria-hidden="true" />
            Back to home
          </Link>
        </Button>
      </div>
    </main>
  )
}
