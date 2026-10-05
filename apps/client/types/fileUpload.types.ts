interface FileUploadResponse {
  fileInfo: FileInfo
  signedUrl: string
}

interface FileInfo {
  fileName: string
  filePath: string
  fileSize: number
  fileType: string
}
