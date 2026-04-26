import { User } from '@/types/user.types'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

const initialState: User = {
  fullName: '',
  email: '',
  profilePic: '',
  createdAt: '',
  updatedAt: '',
  id: '',
  isTemp: false
}

export const meSlice = createSlice({
  name: 'me',
  initialState,
  reducers: {
    setMe: (state, action: PayloadAction<User>) => {
      console.log('setMe: ', action.payload)
      state.createdAt = action.payload.createdAt
      state.updatedAt = action.payload.updatedAt
      state.id = action.payload.id
      state.fullName = action.payload.fullName
      state.email = action.payload.email
      state.profilePic = action.payload.profilePic
      state.isTemp = action.payload.isTemp
    },
    setMyName: (state, action: PayloadAction<string>) => {
      state.fullName = action.payload
    }
  },
  selectors: {
    getMe: state => state,
    getIsTemp: state => state.isTemp
  }
})

export const { setMe, setMyName } = meSlice.actions
export const { getMe, getIsTemp } = meSlice.selectors
export default meSlice.reducer
