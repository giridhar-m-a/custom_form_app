import AuthFormLogin from '@/components/home/AuthFormLogin'
import { FeatureBlock } from '@/components/home/FeatureBlock'
import { Check, Rocket, Zap } from 'lucide-react'

const HomeContent = ({ children }: { children: React.ReactNode }) => {
  const features = [
    {
      icon: Zap,
      title: 'Lightning Fast Setup',
      description: 'Design complex forms in minutes with our intuitive drag-and-drop builder.'
    },
    {
      icon: Check,
      title: 'Automated Workflows',
      description: 'Trigger webhooks, send notifications, and connect with 3,000+ integrations.'
    },
    {
      icon: Rocket,
      title: 'Scale with Confidence',
      description: 'Handle millions of submissions effortlessly. Enterprise-grade reliability.'
    }
  ]

  return (
    <div className="min-h-screen font-sans bg-linear-to-br from-indigo-50 to-white">
      <div className="grid min-h-screen md:grid-cols-2">
        {/* LEFT PANEL - Marketing */}
        <div className="flex flex-col justify-center p-10 md:p-16 lg:p-24 bg-gray-900 text-white space-y-10 relative overflow-hidden">
          {/* background glow */}
          <div className="absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(99,102,241,0.25),transparent_70%)]" />
          <div className="relative space-y-6">
            <h1 className="text-5xl md:text-6xl font-extrabold tracking-tight leading-tight">
              Build <span className="text-indigo-400">Forms That Work</span>
              <br /> Instantly.
            </h1>
            <p className="text-lg md:text-xl text-gray-300 max-w-md">
              The only platform you need to capture data, automate processes, and connect with your audience.
            </p>
          </div>

          <div className="relative grid gap-8">
            {features.map((feature, i) => (
              <FeatureBlock key={i} {...feature} />
            ))}
          </div>

          <p className="relative text-sm text-gray-500 pt-6">© 2025 FormGenius. All rights reserved.</p>
        </div>

        {/* RIGHT PANEL - Auth */}
        <div className="flex items-center justify-center bg-white p-10 md:p-16">{children}</div>
      </div>
    </div>
  )
}

export default HomeContent
