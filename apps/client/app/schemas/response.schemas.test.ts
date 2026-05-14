import { ResponseItemSchema, FormSubmissionSchema, BaseResponseItemSchema } from './response.schemas'
import { describe, it, expect } from 'vitest'

describe('ResponseItemSchema', () => {
  describe('Text / TextArea', () => {
    it('validates required text fields', () => {
      const schema = ResponseItemSchema({ type: 'text', isRequired: true, isMultiple: false })
      const result = schema.safeParse({ formFieldId: '1', responseText: '', responseOptions: [], responseFiles: [] })
      expect(result.success).toBe(false)
      if (!result.success) {
        expect(result.error.format().responseText?._errors).toContain('Response is required')
      }
    })

    it('allows empty text when not required', () => {
      const schema = ResponseItemSchema({ type: 'text', isRequired: false, isMultiple: false })
      const result = schema.safeParse({ formFieldId: '1', responseText: '', responseOptions: [], responseFiles: [] })
      expect(result.success).toBe(true)
    })
  })

  describe('Email', () => {
    it('validates email format', () => {
      const schema = ResponseItemSchema({ type: 'email', isRequired: true, isMultiple: false })
      const result = schema.safeParse({ formFieldId: '1', responseText: 'invalid-email', responseOptions: [], responseFiles: [] })
      expect(result.success).toBe(false)
      if (!result.success) {
        expect(result.error.format().responseText?._errors).toContain('Invalid email')
      }
    })

    it('accepts valid email', () => {
      const schema = ResponseItemSchema({ type: 'email', isRequired: true, isMultiple: false })
      const result = schema.safeParse({ formFieldId: '1', responseText: 'test@example.com', responseOptions: [], responseFiles: [] })
      expect(result.success).toBe(true)
    })
  })

  describe('Number / Slider', () => {
    it('validates numeric values', () => {
      const schema = ResponseItemSchema({ type: 'number', isRequired: true, isMultiple: false })
      const result = schema.safeParse({ formFieldId: '1', responseText: 'abc', responseOptions: [], responseFiles: [] })
      expect(result.success).toBe(false)
      if (!result.success) {
        expect(result.error.format().responseText?._errors).toContain('Must be a valid number')
      }
    })

    it('accepts valid numeric string', () => {
      const schema = ResponseItemSchema({ type: 'number', isRequired: true, isMultiple: false })
      const result = schema.safeParse({ formFieldId: '1', responseText: '123', responseOptions: [], responseFiles: [] })
      expect(result.success).toBe(true)
    })
  })

  describe('Options (Radio / Dropdown)', () => {
    it('requires exactly one option for single-select', () => {
      const schema = ResponseItemSchema({ type: 'radio', isRequired: true, isMultiple: false })
      
      // Empty
      let result = schema.safeParse({ formFieldId: '1', responseText: '', responseOptions: [], responseFiles: [] })
      expect(result.success).toBe(false)

      // Multiple
      result = schema.safeParse({ 
        formFieldId: '1', 
        responseText: '', 
        responseOptions: [{ optionId: 'opt1' }, { optionId: 'opt2' }], 
        responseFiles: [] 
      })
      expect(result.success).toBe(false)

      // One
      result = schema.safeParse({ 
        formFieldId: '1', 
        responseText: '', 
        responseOptions: [{ optionId: 'opt1' }], 
        responseFiles: [] 
      })
      expect(result.success).toBe(true)
    })
  })

  describe('Files / Media', () => {
    const validFile = { fileName: 'test.jpg', filePath: '/path/to/test.jpg', fileSize: 1024, fileType: 'image/jpeg' }

    it('requires at least one file when isRequired is true', () => {
      const schema = ResponseItemSchema({ type: 'image', isRequired: true, isMultiple: true })
      const result = schema.safeParse({ formFieldId: '1', responseText: '', responseOptions: [], responseFiles: [] })
      expect(result.success).toBe(false)
    })

    it('enforces max 1 file when isMultiple is false', () => {
      const schema = ResponseItemSchema({ type: 'image', isRequired: true, isMultiple: false })
      const result = schema.safeParse({ 
        formFieldId: '1', 
        responseText: '', 
        responseOptions: [], 
        responseFiles: [validFile, validFile] 
      })
      expect(result.success).toBe(false)
    })
  })
})

describe('FormSubmissionSchema', () => {
  it('validates structural integrity of form submission', () => {
    const validPayload = {
      formId: 'form-123',
      respondentId: 'user-456',
      responses: [
        { formFieldId: 'field-1', responseText: 'hello' }
      ]
    }
    const result = FormSubmissionSchema.safeParse(validPayload)
    expect(result.success).toBe(true)
  })

  it('fails if responses array is empty', () => {
    const invalidPayload = {
      formId: 'form-123',
      respondentId: 'user-456',
      responses: []
    }
    const result = FormSubmissionSchema.safeParse(invalidPayload)
    expect(result.success).toBe(false)
  })
})
