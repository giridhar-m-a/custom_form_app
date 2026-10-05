'use client'
import { useState, useRef, useCallback, useEffect } from 'react'
import { useMutation } from '@tanstack/react-query'
import axios, { AxiosError } from 'axios'
import { Upload, X, Copy, Check, FileText, Image, Film, Music, Archive, Table, Code, LucideIcon } from 'lucide-react'
import Cookies from 'js-cookie'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { Badge } from '@/components/ui/badge'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { Separator } from '@/components/ui/separator'
import { ApiResponse } from '@/types/api.types'

// ─── Types ────────────────────────────────────────────────────────────────────
interface UploadParams {
  file: File
  onProgress: (pct: number) => void
  path: string
  token: string
}

type UploadResponse = ApiResponse<FileUploadResponse>

interface MutationParams {
  file: File
  path: string
}

interface ApiErrorResponse {
  message?: string
}

// ─── API ──────────────────────────────────────────────────────────────────────
const uploadFile = async ({ file, onProgress, path, token }: UploadParams): Promise<UploadResponse> => {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('path', path)

  const { data } = await axios.post<UploadResponse>(
    `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/v1/files/file-upload`,
    formData,
    {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: `Bearer ${token}`
      },
      onUploadProgress: e => {
        const total = e.total ?? e.loaded
        onProgress(Math.round((e.loaded * 100) / total))
      }
    }
  )

  return data
}

// ─── Helpers ──────────────────────────────────────────────────────────────────
const formatSize = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 ** 2) return `${(bytes / 1024).toFixed(1)} KB`
  return `${(bytes / 1024 ** 2).toFixed(1)} MB`
}

type ExtMap = Record<string, LucideIcon>

const EXT_MAP: ExtMap = {
  pdf: FileText,
  txt: FileText,
  md: FileText,
  png: Image,
  jpg: Image,
  jpeg: Image,
  gif: Image,
  webp: Image,
  svg: Image,
  mp4: Film,
  mov: Film,
  avi: Film,
  mp3: Music,
  wav: Music,
  ogg: Music,
  zip: Archive,
  rar: Archive,
  '7z': Archive,
  csv: Table,
  xlsx: Table,
  xls: Table,
  js: Code,
  ts: Code,
  jsx: Code,
  tsx: Code,
  json: Code
}

const getFileIcon = (name: string): LucideIcon => {
  const ext = name.split('.').pop()?.toLowerCase() ?? ''
  return EXT_MAP[ext] ?? FileText
}

// ─── Props ────────────────────────────────────────────────────────────────────
interface FileUploadProps {
  /** Server-side destination path passed along with the upload */
  uploadPath: string
  /** Callback called with the returned file path on success */
  handleResponse?: (val: FileInfo) => void
  /**
   * When true, the upload starts automatically as soon as a file is selected.
   * The manual "Upload File" button is hidden.
   * Defaults to false.
   */
  autoUpload?: boolean
  token: string
  /** Accepted file types (e.g. "image/*") */
  accept?: string
}

// ─── Component ────────────────────────────────────────────────────────────────
export default function FileUpload({
  uploadPath,
  token,
  handleResponse,
  autoUpload = false,
  accept
}: FileUploadProps) {
  const [file, setFile] = useState<File | null>(null)
  const [progress, setProgress] = useState(0)
  const [dragOver, setDragOver] = useState(false)
  const [copied, setCopied] = useState(false)
  const [localError, setLocalError] = useState<string | null>(null)
  const inputRef = useRef<HTMLInputElement>(null)

  const { mutate, data, isPending, isSuccess, isError, error, reset } = useMutation<
    UploadResponse,
    AxiosError<ApiErrorResponse>,
    MutationParams
  >({
    mutationFn: ({ file, path }) => uploadFile({ file, onProgress: setProgress, path, token }),
    onSuccess: data => {
      setProgress(100)
      handleResponse?.(data.data?.fileInfo as FileInfo)
    },
    onError: () => setProgress(0)
  })

  // ── Auto-upload: fires whenever `file` changes and autoUpload is enabled ──
  useEffect(() => {
    if (autoUpload && file) {
      mutate({ file, path: uploadPath })
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [file])

  const validateFile = (f: File): boolean => {
    if (!accept) return true
    
    const acceptedTypes = accept.split(',').map(t => t.trim())
    const fileType = f.type
    const fileName = f.name.toLowerCase()

    return acceptedTypes.some(type => {
      if (type.endsWith('/*')) {
        const baseType = type.replace('/*', '')
        return fileType.startsWith(baseType)
      }
      if (type.startsWith('.')) {
        return fileName.endsWith(type.toLowerCase())
      }
      return fileType === type
    })
  }

  const handleFile = useCallback(
    (f: File | null | undefined) => {
      if (!f) return
      
      setLocalError(null)
      if (!validateFile(f)) {
        setLocalError(`Invalid file type. Expected: ${accept}`)
        return
      }

      setFile(f)
      setProgress(0)
      reset()
    },
    [reset, accept]
  )

  const handleDrop = useCallback(
    (e: React.DragEvent<HTMLDivElement>) => {
      e.preventDefault()
      setDragOver(false)
      const f = e.dataTransfer.files[0]
      if (f) handleFile(f)
    },
    [handleFile]
  )

  const handleRemove = () => {
    setFile(null)
    setProgress(0)
    setLocalError(null)
    reset()
    if (inputRef.current) inputRef.current.value = ''
  }

  const handleCopy = () => {
    if (!data?.data?.fileInfo?.filePath) return
    navigator.clipboard.writeText(data.data.fileInfo.filePath)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  const Icon: LucideIcon = file ? getFileIcon(file.name) : FileText

  return (
    <div className="w-full bg-zinc-950 flex items-center justify-center p-6">
      <Card className="w-full bg-zinc-900 border-zinc-800 shadow-2xl">
        <CardHeader className="pb-4">
          <div className="flex items-center gap-2">
            <div className="p-2 rounded-lg bg-zinc-800">
              <Upload className="w-4 h-4 text-zinc-300" />
            </div>
            <div>
              <CardTitle className="text-zinc-100 text-lg font-semibold tracking-tight">Upload File</CardTitle>
              <CardDescription className="text-zinc-500 text-xs mt-0.5">
                {autoUpload
                  ? 'Select a file — upload starts automatically'
                  : 'Single file upload with progress tracking'}
              </CardDescription>
            </div>
          </div>
        </CardHeader>

        <Separator className="bg-zinc-800" />

        <CardContent className="pt-5 space-y-4">
          {/* Drop Zone */}
          <div
            onDragOver={e => {
              e.preventDefault()
              setDragOver(true)
            }}
            onDragLeave={() => setDragOver(false)}
            onDrop={handleDrop}
            onClick={() => {
              if (!file) inputRef.current?.click()
            }}
            className={[
              'relative rounded-xl border-2 border-dashed px-6 py-10 text-center transition-all duration-200 cursor-pointer select-none',
              dragOver
                ? 'border-zinc-400 bg-zinc-800/60'
                : file
                  ? 'border-zinc-700 bg-zinc-800/30 cursor-default'
                  : 'border-zinc-700 bg-zinc-800/20 hover:border-zinc-600 hover:bg-zinc-800/40'
            ].join(' ')}>
            <input
              ref={inputRef}
              type="file"
              accept={accept}
              className="hidden"
              onChange={e => handleFile(e.target.files?.[0])}
            />

            {file ? (
              <div className="flex items-center gap-3 text-left">
                <div className="shrink-0 p-2.5 rounded-lg bg-zinc-700/60">
                  <Icon className="w-5 h-5 text-zinc-300" />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="text-sm font-medium text-zinc-200 truncate">{file.name}</p>
                  <p className="text-xs text-zinc-500 mt-0.5">{formatSize(file.size)}</p>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  className="shrink-0 h-7 w-7 text-zinc-500 hover:text-red-400 hover:bg-red-400/10"
                  onClick={e => {
                    e.stopPropagation()
                    handleRemove()
                  }}>
                  <X className="w-3.5 h-3.5" />
                </Button>
              </div>
            ) : (
              <div className="space-y-2">
                <div className="mx-auto w-10 h-10 rounded-full bg-zinc-800 flex items-center justify-center">
                  <Upload className="w-4 h-4 text-zinc-500" />
                </div>
                <div>
                  <p className="text-sm text-zinc-400">
                    Drop your file here, or{' '}
                    <span className="text-zinc-200 underline underline-offset-2 font-medium">browse</span>
                  </p>
                  <p className="text-xs text-zinc-600 mt-1">
                    {accept ? `Accepted: ${accept}` : autoUpload ? 'Upload begins immediately on selection' : 'Any file format accepted'}
                  </p>
                </div>
              </div>
            )}
          </div>

          {/* Progress */}
          {(isPending || (progress > 0 && !isError && !isSuccess)) && (
            <div className="space-y-1.5">
              <div className="flex justify-between items-center">
                <span className="text-xs text-zinc-500">Uploading…</span>
                <span className="text-xs font-medium text-zinc-300">{progress}%</span>
              </div>
              <Progress value={progress} className="h-1.5 bg-zinc-800 [&>div]:bg-zinc-300" />
            </div>
          )}

          {/* Manual Upload Button — hidden when autoUpload is true */}
          {!autoUpload && (
            <Button
              className="w-full bg-zinc-100 text-zinc-900 hover:bg-white font-semibold tracking-tight transition-all"
              disabled={!file || isPending}
              onClick={() => {
                if (file) mutate({ file, path: uploadPath })
              }}>
              {isPending ? (
                <span className="flex items-center gap-2">
                  <span className="w-3.5 h-3.5 border-2 border-zinc-400 border-t-zinc-900 rounded-full animate-spin" />
                  Uploading
                </span>
              ) : (
                <span className="flex items-center gap-2">
                  <Upload className="w-3.5 h-3.5" />
                  Upload File
                </span>
              )}
            </Button>
          )}

          {/* Success */}
          {isSuccess && data?.data?.fileInfo.filePath && (
            <div className="rounded-lg bg-emerald-950/40 border border-emerald-800/50 p-3.5 space-y-2">
              <div className="flex items-center gap-1.5">
                <Badge
                  variant="outline"
                  className="text-emerald-400 border-emerald-700/60 bg-emerald-900/30 text-[10px] px-1.5 py-0">
                  ✓ uploaded
                </Badge>
                <span className="text-xs text-zinc-500">File path returned</span>
              </div>
              <div className="flex items-center gap-2">
                <code className="flex-1 text-xs text-emerald-300 bg-zinc-900/60 rounded px-2.5 py-1.5 truncate font-mono">
                  {data.data?.fileInfo.filePath}
                </code>
                <Button
                  variant="ghost"
                  size="icon"
                  className="shrink-0 h-7 w-7 text-zinc-500 hover:text-zinc-200 hover:bg-zinc-700"
                  onClick={handleCopy}>
                  {copied ? <Check className="w-3.5 h-3.5 text-emerald-400" /> : <Copy className="w-3.5 h-3.5" />}
                </Button>
              </div>
            </div>
          )}

          {/* Error */}
          {(isError || localError) && (
            <Alert className="bg-red-950/30 border-red-800/50 text-red-300 py-3">
              <AlertDescription className="text-xs">
                {localError || error?.response?.data?.message || error?.message || 'Upload failed. Please try again.'}
              </AlertDescription>
            </Alert>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
