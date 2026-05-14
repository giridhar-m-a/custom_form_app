import { authSlice, setTokens, clearTokens, getTokens } from './auth.slice'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import Cookies from 'js-cookie'

vi.mock('js-cookie', () => ({
  default: {
    get: vi.fn(),
    set: vi.fn(),
    remove: vi.fn(),
  },
}))

describe('authSlice', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('sets tokens and cookies correctly', () => {
    const initialState = { accessToken: null, refreshToken: null }
    const action = setTokens({ accessToken: 'access', refreshToken: 'refresh' })
    const state = authSlice.reducer(initialState, action)

    expect(state.accessToken).toBe('access')
    expect(state.refreshToken).toBe('refresh')
    expect(Cookies.set).toHaveBeenCalledWith('accessToken', 'access', expect.any(Object))
    expect(Cookies.set).toHaveBeenCalledWith('refreshToken', 'refresh', expect.any(Object))
  })

  it('clears tokens and cookies correctly', () => {
    const initialState = { accessToken: 'access', refreshToken: 'refresh' }
    const action = clearTokens()
    const state = authSlice.reducer(initialState, action)

    expect(state.accessToken).toBeNull()
    expect(state.refreshToken).toBeNull()
    expect(Cookies.remove).toHaveBeenCalledWith('accessToken')
    expect(Cookies.remove).toHaveBeenCalledWith('refreshToken')
  })

  it('selects tokens correctly', () => {
    const state = { accessToken: 'access', refreshToken: 'refresh' }
    const tokens = getTokens({ auth: state } as any)
    expect(tokens).toEqual({ accessToken: 'access', refreshToken: 'refresh' })
  })
})
