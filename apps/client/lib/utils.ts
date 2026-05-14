import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

const minioServer = process.env.NEXT_PUBLIC_MINIO_URL

export const getFileUrl = (path?: string) => {
  // const minioServer = process.env.MINIO_SERVER || 'minio.custom-form-app.home'
  return path ? `${minioServer}/custom-form-app/${path}` : ''
}

export const validateFile = ({
  file,
  size,
  type,
  setFile
}: {
  file?: File
  size?: number
  type: Set<string>
  setFile: (file: File) => void
}): { error?: string; valid: boolean } => {
  if (!file) return { error: 'No file selected', valid: false }
  if (size && file.size > size) return { error: 'File size exceeds limit', valid: false }
  if (!type.has(file.type)) return { error: 'Invalid file type', valid: false }

  setFile(file)

  return { valid: true }
}
