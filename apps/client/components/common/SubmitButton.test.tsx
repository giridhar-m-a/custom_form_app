import { render, screen } from '@testing-library/react'
import { SubmitButton } from './SubmitButton'
import { describe, it, expect, vi } from 'vitest'

// Mock the Button component if needed, or just use it
vi.mock('../ui/button', () => ({
  Button: ({ children, disabled, type, ...props }: any) => (
    <button type={type} disabled={disabled} {...props}>
      {children}
    </button>
  ),
}))

vi.mock('./CustomLoader', () => ({
  CustomLoader: () => <span data-testid="loader">Loading...</span>,
}))

describe('SubmitButton', () => {
  it('renders children correctly', () => {
    render(<SubmitButton>Submit</SubmitButton>)
    expect(screen.getByText('Submit')).toBeDefined()
  })

  it('shows loader when isLoading is true', () => {
    render(<SubmitButton isLoading={true}>Submit</SubmitButton>)
    expect(screen.getByTestId('loader')).toBeDefined()
  })

  it('is disabled when isLoading is true', () => {
    render(<SubmitButton isLoading={true}>Submit</SubmitButton>)
    const button = screen.getByRole('button')
    expect(button).toHaveProperty('disabled', true)
  })

  it('is disabled when disabled prop is true', () => {
    render(<SubmitButton disabled={true}>Submit</SubmitButton>)
    const button = screen.getByRole('button')
    expect(button).toHaveProperty('disabled', true)
  })
})
