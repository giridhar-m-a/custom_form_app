import React, { useState } from 'react'
import { Upload, FileText, CheckCircle2, AlertCircle, X, Send } from 'lucide-react'
import { useBulkCreateInvitation } from '@/hooks/queryHooks/useInvitations'
import { SubmitButton } from '@/components/common/SubmitButton'

interface BulkInvitationFormProps {
  formId: string
  setInviteOpen: (open: boolean) => void
}

const BulkInvitationForm = ({ formId, setInviteOpen }: BulkInvitationFormProps) => {
  const [file, setFile] = useState<File | null>(null)
  const [error, setError] = useState<string | null>(null)
  const { mutate: createBulkInvitation, isPending } = useBulkCreateInvitation()

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0]
    if (!selectedFile) return

    if (selectedFile.type !== 'text/csv' && !selectedFile.name.endsWith('.csv')) {
      setError('Please upload a valid CSV file.')
      return
    }

    setFile(selectedFile)
    setError(null)
  }

  const removeFile = () => {
    setFile(null)
    setError(null)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!file) return

    const formData = new FormData()
    formData.append('file', file)
    createBulkInvitation(
      { data: formData, formId },
      {
        onSuccess: () => {
          setInviteOpen?.(false)
        }
      }
    )
  }

  return (
    <div className="flex items-center justify-center text-foreground">
      <div className="w-full max-w-2xl overflow-hidden">
        <form className="p-8 space-y-6" onSubmit={handleSubmit}>
          {/* Upload Area */}
          {!file ? (
            <div className="relative group">
              <input
                type="file"
                accept=".csv"
                onChange={handleFileChange}
                className="absolute inset-0 w-full h-full opacity-0 cursor-pointer z-10"
              />
              <div
                className={`border-2 border-dashed rounded-xl p-10 flex flex-col items-center justify-center transition-all ${
                  error
                    ? 'border-destructive/50 bg-destructive/10'
                    : 'border-border bg-muted/30 group-hover:border-primary/50 group-hover:bg-primary/5'
                }`}>
                <div
                  className={`p-4 rounded-full mb-4 ${
                    error ? 'bg-destructive/20 text-destructive' : 'bg-primary/10 text-primary'
                  }`}>
                  <Upload size={32} />
                </div>
                <p className="font-medium text-foreground">Click to upload or drag and drop</p>
                <p className="text-xs text-muted-foreground mt-2 uppercase tracking-wider font-semibold">
                  Only CSV files accepted
                </p>
              </div>
            </div>
          ) : (
            <div className="flex items-center justify-between p-4 bg-primary/5 border border-primary/20 rounded-lg">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-primary text-primary-foreground rounded-md">
                  <FileText size={20} />
                </div>
                <div>
                  <p className="text-sm font-semibold text-foreground">{file.name}</p>
                  <p className="text-xs text-muted-foreground">{(file.size / 1024).toFixed(1)} KB</p>
                </div>
              </div>
              <button
                type="button"
                onClick={removeFile}
                className="p-1 hover:bg-primary/10 rounded-full text-primary transition-colors">
                <X size={20} />
              </button>
            </div>
          )}

          {/* Error Message */}
          {error && (
            <div className="flex items-center gap-2 text-destructive text-sm bg-destructive/10 p-3 rounded-lg border border-destructive/20">
              <AlertCircle size={16} />
              {error}
            </div>
          )}

          {/* Submit Button */}
          <div className="pt-2">
            <SubmitButton
              type="submit"
              disabled={!file || isPending}
              isLoading={isPending}
              className={`w-full py-3 px-4 rounded-xl font-bold flex items-center justify-center gap-2 transition-all shadow-lg ${
                !file
                  ? 'bg-muted text-muted-foreground cursor-not-allowed shadow-none'
                  : 'bg-primary text-primary-foreground hover:bg-primary/90 active:scale-[0.98]'
              }`}>
              <>
                <Send size={18} /> Submit CSV Data
              </>
            </SubmitButton>
          </div>
        </form>
      </div>
    </div>
  )
}

export default BulkInvitationForm
