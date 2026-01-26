import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import Cookies from 'js-cookie'

interface AuthState {
  accessToken: string | null
  refreshToken: string | null
}

const initialState: AuthState = {
  accessToken: Cookies.get('accessToken') || null,
  refreshToken: Cookies.get('refreshToken') || null
}

export const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setTokens: (state, action: PayloadAction<{ accessToken: string; refreshToken: string }>) => {
      state.accessToken = action.payload.accessToken
      state.refreshToken = action.payload.refreshToken
      try {
        Cookies.set('accessToken', action.payload.accessToken, {
          expires: new Date(Date.now() + 60 * 60 * 24 * 7 * 1000)
        })
        Cookies.set('refreshToken', action.payload.refreshToken, {
          expires: new Date(Date.now() + 60 * 60 * 24 * 7 * 1000)
        })
      } catch (e) {
        console.log('error setting cookies', e)
      }
    },
    clearTokens: state => {
      state.accessToken = null
      state.refreshToken = null
      Cookies.remove('accessToken')
      Cookies.remove('refreshToken')
    }
  },
  selectors: {
    getTokens: state => ({
      accessToken: state.accessToken,
      refreshToken: state.refreshToken
    })
  }
})

export const { setTokens, clearTokens } = authSlice.actions

export const { getTokens } = authSlice.selectors

export default authSlice.reducer
