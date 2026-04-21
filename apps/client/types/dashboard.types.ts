export interface SubmissionsByMonth {
  month: string
  totalSubmissions: number
}

export interface DashboardData {
  totalForms: number
  totalSubmissions: number
  totalActiveForms: number
  totalClosedForms: number
  totalInvitations: number
  submissionsByMonth: SubmissionsByMonth[]
}
