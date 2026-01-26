import { User } from '@/types/user.types'
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

const initialState: User = {
  fullName: '',
  email: '',
  profilePic: '',
  profilePicId: '',
  createdAt: '',
  updatedAt: '',
  id: ''
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
      state.profilePicId = action.payload.profilePicId
    }
  },
  selectors: {
    getMe: state => state
  }
})

export const { setMe } = meSlice.actions
export const { getMe } = meSlice.selectors
export default meSlice.reducer
