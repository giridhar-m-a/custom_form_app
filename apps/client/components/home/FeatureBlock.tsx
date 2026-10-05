import { FC } from 'react'

interface FeatureBlockProps {
  icon: FC<{ className?: string }>
  title: string
  description: string
}

export const FeatureBlock: FC<FeatureBlockProps> = ({ icon: Icon, title, description }) => (
  <div className="flex items-start space-x-4 transition-transform hover:translate-x-1">
    <div className="shrink-0 p-3 rounded-full bg-indigo-500/20 text-indigo-400 shadow-inner">
      <Icon className="w-6 h-6" />
    </div>
    <div>
      <h3 className="text-lg font-semibold text-white">{title}</h3>
      <p className="mt-1 text-gray-400 leading-relaxed">{description}</p>
    </div>
  </div>
)
