'use client'

import { CustomLoader } from '@/components/common/CustomLoader'
import { Button } from '@/components/ui/button'
import { useDownloadResponseFiles } from '@/hooks/queryHooks/useResponses'
import { DownloadIcon } from 'lucide-react'

export const DownloadButton = ({ filePath, fileName }: { filePath: string; fileName: string }) => {
  const { mutate, isPending } = useDownloadResponseFiles()
  const handleDownload = () => {
    mutate({ filePath, fileName })
  }
  return (
    <Button disabled={isPending} onClick={handleDownload} variant={'outline'}>
      {!isPending && <DownloadIcon />}
      {isPending && <CustomLoader />}
    </Button>
  )
}
