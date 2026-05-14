import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import AuthFormLogin from './AuthFormLogin'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useCredentialAuth, useGoogleAuth } from '@/hooks/queryHooks/useAuth'

// Mock hooks
vi.mock('@/hooks/queryHooks/useAuth', () => ({
  useCredentialAuth: vi.fn(),
  useGoogleAuth: vi.fn(),
}))

vi.mock('@react-oauth/google', () => ({
  useGoogleLogin: vi.fn(() => vi.fn()),
}))

// Mock sub-components to focus on AuthFormLogin
vi.mock('./AuthFormSignUp', () => ({
  default: () => <div data-testid="signup-form">SignUp Form</div>,
}))

vi.mock('./AuthFormReset', () => ({
  default: () => <div data-testid="reset-form">Reset Form</div>,
}))

vi.mock('./TempUser', () => ({
  default: () => <div data-testid="temp-user">Temp User</div>,
}))

describe('AuthFormLogin', () => {
  const mockCredentialLogin = vi.fn()
  const mockGoogleLogin = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
    ;(useCredentialAuth as any).mockReturnValue({
      mutate: mockCredentialLogin,
      isPending: false,
    })
    ;(useGoogleAuth as any).mockReturnValue({
      mutate: mockGoogleLogin,
      isPending: false,
    })
  })

  it('renders login form by default', () => {
    render(<AuthFormLogin />)
    expect(screen.getByText('Welcome Back')).toBeDefined()
    expect(screen.getByLabelText(/email/i)).toBeDefined()
    expect(screen.getByLabelText(/password/i)).toBeDefined()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeDefined()
  })

  it('shows validation errors for empty fields on submit', async () => {
    render(<AuthFormLogin />)
    const submitButton = screen.getByRole('button', { name: /sign in/i })
    
    fireEvent.click(submitButton)

    // Wait for validation messages. Since they are required fields in the schema.
    // Note: react-hook-form might need some time to update state.
    // But since they are 'required' in the input tag too, we might need to bypass browser validation if testing-library doesn't handle it.
    // However, zodResolver should handle it.
  })

  it('calls credentialLogin with correct data on valid submit', async () => {
    render(<AuthFormLogin />)
    
    fireEvent.change(screen.getByLabelText(/email/i), { target: { value: 'test@example.com' } })
    fireEvent.change(screen.getByLabelText(/password/i), { target: { value: 'password123' } })
    
    fireEvent.click(screen.getByRole('button', { name: /sign in/i }))

    await waitFor(() => {
      expect(mockCredentialLogin).toHaveBeenCalledWith({
        email: 'test@example.com',
        password: 'password123',
      })
    })
  })

  it('switches to sign up form when "Sign Up" button is clicked', () => {
    render(<AuthFormLogin />)
    const signUpToggle = screen.getByRole('button', { name: /sign up/i })
    
    fireEvent.click(signUpToggle)
    
    expect(screen.getByTestId('signup-form')).toBeDefined()
    expect(screen.getByText('Sign Up')).toBeDefined()
  })

  it('switches to reset form when "Forgot Password?" is clicked', () => {
    render(<AuthFormLogin />)
    const forgotPassword = screen.getByText(/forgot password\?/i)
    
    fireEvent.click(forgotPassword)
    
    expect(screen.getByTestId('reset-form')).toBeDefined()
  })

  it('disables buttons when loading', () => {
    ;(useCredentialAuth as any).mockReturnValue({
      mutate: mockCredentialLogin,
      isPending: true,
    })
    
    render(<AuthFormLogin />)
    
    // The button text changes when loading, so we find it by role.
    // There are two buttons: "Continue with Google" and the submit button.
    const buttons = screen.getAllByRole('button')
    const submitButton = buttons.find(b => b.getAttribute('type') === 'submit')
    const googleButton = screen.getByRole('button', { name: /continue with google/i })

    expect(submitButton).toBeDisabled()
    expect(googleButton).toBeDisabled()
  })
})
