export interface ResponseFilter {
  page: number
  limit: number
  search?: string
}

export interface ResponseList {
  invitedEmail: string
  invitedName: string
  respondentId: string
  submissionId: string
  submittedAt: string
}

export interface ResponseFiles {
  fileId: string
  fileName: string
  filePath: string
  fileSize: number
  fileType: string
  fileUploadedAt: string
}

export interface ResponseOptions {
  formOptionId: string
  id: string
}

export interface ResponseFields {
  formFieldId: string
  responseText: string
  responseOptions: ResponseOptions[]
  responseFiles: ResponseFiles[]
}

export interface ResponseDetails {
  submissionId: string
  formId: string
  respondentId: string
  submittedAt: string
  responses: ResponseFields[]
}
