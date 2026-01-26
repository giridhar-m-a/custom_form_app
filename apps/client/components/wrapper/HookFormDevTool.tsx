'use client'

import { DevTool } from '@hookform/devtools'
import { Control } from 'react-hook-form'

export const HookFormDevTool = ({ control }: { control: Control<any> }) => {
  return <>{process.env.NODE_ENV === 'development' && <DevTool control={control} />}</>
}
