import { meSlice, setMe, setMyName, getMe, getIsTemp } from './me.slice'
import { describe, it, expect } from 'vitest'
import { User } from '@/types/user.types'

const mockUser: User = {
  id: 'user-1',
  fullName: 'John Doe',
  email: 'john@example.com',
  profilePic: 'pic.jpg',
  createdAt: '2023-01-01',
  updatedAt: '2023-01-02',
  isTemp: false,
}

describe('meSlice', () => {
  it('sets user data correctly', () => {
    const initialState = { id: '', fullName: '', email: '', profilePic: '', createdAt: '', updatedAt: '', isTemp: false }
    const action = setMe(mockUser)
    const state = meSlice.reducer(initialState, action)

    expect(state).toEqual(mockUser)
  })

  it('sets my name correctly', () => {
    const initialState = { ...mockUser }
    const action = setMyName('Jane Doe')
    const state = meSlice.reducer(initialState, action)

    expect(state.fullName).toBe('Jane Doe')
    expect(state.email).toBe('john@example.com') // Other fields should remain unchanged
  })

  it('selects me correctly', () => {
    const state = { ...mockUser }
    const me = getMe({ me: state } as any)
    expect(me).toEqual(mockUser)
  })

  it('selects isTemp correctly', () => {
    const state = { ...mockUser, isTemp: true }
    const isTemp = getIsTemp({ me: state } as any)
    expect(isTemp).toBe(true)
  })
})
