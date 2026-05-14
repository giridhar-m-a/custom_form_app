import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { FormRender } from './FormRender'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useCreateResponse } from '@/hooks/queryHooks/useResponses'
import { FormField } from '@/types/form.types'

// Mock useCreateResponse
vi.mock('@/hooks/queryHooks/useResponses', () => ({
  useCreateResponse: vi.fn(),
}))

// Mock ScrollArea
vi.mock('../ui/scroll-area', () => ({
  ScrollArea: ({ children }: any) => <div data-testid="scroll-area">{children}</div>,
}))

// Mock next/navigation
vi.mock('next/navigation', () => ({
  useSearchParams: () => ({
    get: vi.fn().mockReturnValue('test-token'),
  }),
}))

const mockFields: FormField[] = [
  {
    fieldId: 'field-1',
    fieldLabel: 'Text Field',
    fieldType: 'text',
    isRequired: true,
    ordering: 1,
    formId: 'form-1',
    options: []
  }
]

describe('FormRender', () => {
  const mockCreateResponse = vi.fn()
  const mockOnSubmit = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useCreateResponse as any).mockReturnValue({
      mutate: mockCreateResponse,
      isPending: false,
    })
  })

  it('renders all provided fields', () => {
    render(
      <FormRender 
        fields={mockFields} 
        formId="form-1" 
        respondentId="resp-1" 
        token="token-1" 
      />
    )
    
    expect(screen.getByText('Text Field')).toBeDefined()
    expect(screen.getByRole('button', { name: /submit/i })).toBeDefined()
  })

  it('shows validation error for required field on submit', async () => {
    render(
      <FormRender 
        fields={mockFields} 
        formId="form-1" 
        respondentId="resp-1" 
        token="token-1" 
      />
    )
    
    const submitButton = screen.getByRole('button', { name: /submit/i })
    fireEvent.click(submitButton)

    // Validation is async in react-hook-form
    const errorMessage = await screen.findByText(/response is required/i)
    expect(errorMessage).toBeDefined()
    
    expect(mockCreateResponse).not.toHaveBeenCalled()
  })

  it('calls createResponse with correct data on valid submit', async () => {
    render(
      <FormRender 
        fields={mockFields} 
        formId="form-1" 
        respondentId="resp-1" 
        token="token-1" 
      />
    )
    
    const textInput = screen.getByLabelText(/text field/i)
    fireEvent.change(textInput, { target: { value: 'Valid Response' } })
    
    fireEvent.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() => {
      expect(mockCreateResponse).toHaveBeenCalledWith(expect.objectContaining({
        data: expect.objectContaining({
          responses: expect.arrayContaining([
            expect.objectContaining({
              responseText: 'Valid Response'
            })
          ])
        })
      }))
    })
  })

  it('calls custom onSubmit if provided', async () => {
    render(
      <FormRender 
        fields={mockFields} 
        formId="form-1" 
        respondentId="resp-1" 
        token="token-1" 
        onSubmit={mockOnSubmit}
      />
    )
    
    const textInput = screen.getByLabelText(/text field/i)
    fireEvent.change(textInput, { target: { value: 'Valid Response' } })
    
    fireEvent.click(screen.getByRole('button', { name: /submit/i }))

    await waitFor(() => {
      expect(mockOnSubmit).toHaveBeenCalled()
    })
  })
})
