export interface FormField {
  fieldId: string
  fieldLabel: string
  fieldType: string
  isRequired: boolean
  ordering: number
  formId: string
  options: FieldOption[]
}

export interface FieldOption {
  optionId: string
  optionLabel: string
  ordering: number
  fieldId: string
}

export interface FormFieldOptionRecord {
  [key: string] : FieldOption
}
