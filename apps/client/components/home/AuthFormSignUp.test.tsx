import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import AuthFormSignUp from './AuthFormSignUp'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useRegister } from '@/hooks/queryHooks/useAuth'

// Mock useRegister
vi.mock('@/hooks/queryHooks/useAuth', () => ({
  useRegister: vi.fn(),
}))

describe('AuthFormSignUp', () => {
  const mockRegister = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useRegister as any).mockReturnValue({
      mutate: mockRegister,
      isPending: false,
    })
  })

  it('renders sign up form correctly', () => {
    const { container } = render(<AuthFormSignUp />)
    expect(container.querySelector('#name')).toBeDefined()
    expect(container.querySelector('#email')).toBeDefined()
    expect(container.querySelector('#password')).toBeDefined()
    expect(container.querySelector('#confirmPassword')).toBeDefined()
    expect(screen.getByRole('button', { name: /sign up/i })).toBeDefined()
  })

  it('calls register with correct data on valid submit', async () => {
    const { container } = render(<AuthFormSignUp />)
    
    fireEvent.change(screen.getByLabelText(/full name/i), { target: { value: 'John Doe' } })
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'john@example.com' } })
    fireEvent.change(container.querySelector('#password')!, { target: { value: 'password123' } })
    fireEvent.change(container.querySelector('#confirmPassword')!, { target: { value: 'password123' } })
    
    fireEvent.click(screen.getByRole('button', { name: /sign up/i }))

    await waitFor(() => {
      expect(mockRegister).toHaveBeenCalledWith({
        name: 'John Doe',
        email: 'john@example.com',
        password: 'password123',
        confirmPassword: 'password123',
      })
    })
  })

  it('disables button when loading', () => {
    ;(useRegister as any).mockReturnValue({
      mutate: mockRegister,
      isPending: true,
    })
    
    render(<AuthFormSignUp />)
    expect(screen.getByRole('button', { name: /sign up/i })).toBeDisabled()
  })
})
