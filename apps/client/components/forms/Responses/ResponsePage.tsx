'use client'

import { Badge } from '@/components/ui/badge'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { useGetFormById, useGetFormFields } from '@/hooks/queryHooks/useFormApp'
import { useSingleResponse } from '@/hooks/queryHooks/useResponses'
import { FormField } from '@/types/form.types'
import { format, isValid, parseISO } from 'date-fns'
import {
  CalendarClock,
  CheckCircle2,
  Download,
  ExternalLink,
  FileText,
  LayoutList,
  Loader,
  Star,
  XCircle
} from 'lucide-react'
import { DownloadButton } from './DownloadButton'

export const ResponsePage = ({ params }: { params: { id: string; submissionId: string } }) => {
  const { id, submissionId } = params
  const { data: form, isLoading: isLoadingForm } = useGetFormById(id)
  const { data: fields, isLoading: isLoadingFields } = useGetFormFields(id)
  const { data: response, isLoading: isLoadingResponse } = useSingleResponse(submissionId)
  const isLoading = isLoadingForm || isLoadingFields || isLoadingResponse

  const fieldMap = new Map<string, FormField>()
  fields?.data?.forEach(field => {
    fieldMap.set(field.fieldId, field)
  })

  if (isLoading) {
    return (
      <div className="w-full min-h-[50vh] flex flex-col justify-center items-center gap-4 text-muted-foreground">
        <Loader className="w-8 h-8 animate-spin text-primary" />
        <p className="text-sm font-medium animate-pulse">Loading response details...</p>
      </div>
    )
  }

  return (
    <div className="mx-auto py-8 px-4 sm:px-6 lg:px-8 animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="mb-8 space-y-2">
        <h1 className="text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
          {form?.data?.title || 'Form Response'}
        </h1>
        <p className="text-muted-foreground text-lg flex items-center gap-2 flex-wrap">
          Submission ID: <span className="font-mono text-sm px-2 py-1 bg-muted rounded-md">{submissionId}</span>
        </p>
      </div>
      <div className="space-y-6">
          {response?.data?.responses?.map(responseItem => {
            const field = fieldMap.get(responseItem.formFieldId)
            if (!field) return null

            const isQuiz = field.options?.some(opt => opt.isAnswer)

            return (
              <Card
                key={responseItem.formFieldId}
                className={`overflow-hidden border-border/50 bg-card/50 shadow-sm hover:shadow-md transition-all duration-200`}>
                <CardHeader className="bg-muted/30">
                  <CardTitle className="text-lg flex flex-col sm:flex-row sm:items-center justify-between gap-2">
                    <div className="flex items-center gap-2">
                      <FileText className="w-4 h-4 text-primary shrink-0" />
                      <span className="leading-tight">{field.fieldLabel}</span>
                    </div>
                    {field.isRequired && (
                      <Badge variant="secondary" className="w-fit text-xs font-normal">
                        Required
                      </Badge>
                    )}
                  </CardTitle>
                </CardHeader>
                <CardContent className="pt-6">
                  {/* Text Response mapping to basic styles or specialized visualizations */}
                  {responseItem.responseText && (
                    <div className="bg-background rounded-lg border border-border/50 p-4 text-foreground/90 whitespace-pre-wrap leading-relaxed shadow-sm">
                      {field.fieldType === 'color' ? (
                        <div className="flex items-center gap-3">
                          <div
                            className="w-8 h-8 rounded-full border shadow-sm"
                            style={{ backgroundColor: responseItem.responseText }}
                          />
                          <span className="font-mono uppercase">{responseItem.responseText}</span>
                        </div>
                      ) : field.fieldType === 'url' ? (
                        <a
                          href={
                            responseItem.responseText.startsWith('http')
                              ? responseItem.responseText
                              : `https://${responseItem.responseText}`
                          }
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-primary hover:underline inline-flex items-center gap-1">
                          {responseItem.responseText}
                          <ExternalLink className="w-3 h-3" />
                        </a>
                      ) : field.fieldType === 'rating' ? (
                        <div className="flex items-center gap-2">
                          <span className="text-2xl font-bold">{responseItem.responseText}</span>
                          <Star className="w-6 h-6 fill-yellow-400 text-yellow-500" />
                        </div>
                      ) : ['date', 'time', 'datetime'].includes(field.fieldType) ? (
                        <div className="flex items-center gap-2">
                          <CalendarClock className="w-5 h-5 text-muted-foreground mr-1" />
                          <span className="font-medium text-foreground">
                            {(() => {
                              try {
                                let d = parseISO(responseItem.responseText)
                                if (!isValid(d)) d = new Date(responseItem.responseText)
                                if (
                                  !isValid(d) &&
                                  field.fieldType === 'time' &&
                                  /^\d{2}:\d{2}/.test(responseItem.responseText)
                                ) {
                                  d = new Date(`1970-01-01T${responseItem.responseText}`)
                                }
                                if (isValid(d)) {
                                  if (field.fieldType === 'date') return format(d, 'PPP')
                                  if (field.fieldType === 'datetime') return format(d, 'PPp')
                                  if (field.fieldType === 'time') return format(d, 'p')
                                }
                                return responseItem.responseText
                              } catch (e) {
                                return responseItem.responseText
                              }
                            })()}
                          </span>
                        </div>
                      ) : (
                        responseItem.responseText
                      )}
                    </div>
                  )}

                  {/* File/Media Response */}
                  {responseItem.responseFiles && responseItem.responseFiles.length > 0 && (
                    <div className="flex flex-col gap-3">
                      {responseItem.responseFiles.map((file, idx) => (
                        <div
                          key={idx}
                          className="flex items-center justify-between p-3 rounded-lg border border-border bg-background group hover:border-primary/50 transition-colors shadow-sm">
                          <div className="flex items-center gap-3 truncate">
                            <div className="p-2 rounded-md bg-primary/10 text-primary">
                              <Download className="w-4 h-4" />
                            </div>
                            <span className="text-sm font-medium truncate" title={file.filePath}>
                              {file.filePath.split('/').pop() || 'Download File'}
                            </span>
                          </div>
                          <DownloadButton filePath={file.filePath} fileName={file.fileName} />
                        </div>
                      ))}
                    </div>
                  )}

                  {/* Options Response (Multi-choice, Checks, Dropdowns) */}
                  {responseItem.responseOptions && responseItem.responseOptions.length > 0 && (
                    <div className="flex flex-wrap gap-2 mt-1">
                      {responseItem.responseOptions.map(option => {
                        const optionData = field?.options?.find(opt => opt.optionId === option.formOptionId)
                        const isCorrect = optionData?.isAnswer

                        return (
                          <Badge
                            key={option.formOptionId}
                            variant={isQuiz ? (isCorrect ? 'default' : 'destructive') : 'default'}
                            className={`px-3 py-1.5 text-sm font-medium shadow-sm flex items-center gap-1.5 ${
                              isQuiz && isCorrect ? 'bg-green-600 hover:bg-green-700 text-white' : ''
                            }`}>
                            {optionData?.optionLabel || 'Selected Option'}
                            {isQuiz && isCorrect && <CheckCircle2 className="w-3.5 h-3.5" />}
                            {isQuiz && !isCorrect && <XCircle className="w-3.5 h-3.5" />}
                          </Badge>
                        )
                      })}
                    </div>
                  )}

                  {/* Empty State */}
                  {!responseItem.responseText &&
                    (!responseItem.responseFiles || responseItem.responseFiles.length === 0) &&
                    (!responseItem.responseOptions || responseItem.responseOptions.length === 0) && (
                      <p className="text-muted-foreground text-sm italic py-2">No response provided</p>
                    )}
                </CardContent>
              </Card>
            )
          })}

          {(!response?.data?.responses || response.data.responses.length === 0) && (
            <div className="text-center py-12 px-4 border border-dashed rounded-xl bg-muted/20">
              <LayoutList className="w-12 h-12 text-muted-foreground/50 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-foreground">No responses found</h3>
              <p className="text-muted-foreground mt-1">This submission appears to be empty.</p>
            </div>
          )}
        </div>
    </div>
  )
}
