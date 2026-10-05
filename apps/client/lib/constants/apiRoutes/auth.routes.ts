const base = '/auth'

export const AUTH_ROUTES = {
  baseUrl: '/auth',
  login: {
    credentials: `${base}/login`,
    google: `${base}/google`
  },
  register: `${base}/register`,
  verify: `${base}/verify`,
  refreshToken: `${base}/refresh-token`,
  resetRequest: `${base}/request-password-reset`,
  resetPassword: `${base}/reset-password`,
  tempUser: `${base}/temp-user-auth`
}
