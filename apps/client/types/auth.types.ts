import { User } from './user.types'

export interface AuthResponse {
  accessToken: string
  refreshToken: string
  user: User
}

export interface VerifyTokenResponse {
  userID: string
}

export interface RefreshTokenResponse {
  accessToken: string
  refreshToken: string
}
