const base = '/auth'

export const AUTH_ROUTES = {
  baseUrl: '/auth',
  login: {
    credentials: `${base}/login`,
    google: `${base}/google`
  },
  register: `${base}/register`
}
