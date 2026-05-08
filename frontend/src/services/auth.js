import { apiClient } from '@/lib/api'

export { apiClient }

export const authService = {
  async login(credentials) {
    const response = await apiClient.post('/auth/login', credentials)
    return response.data
  },

  async login2FA(userId, code) {
    const response = await apiClient.post('/auth/login/2fa', { user_id: userId, code })
    return response.data
  },

  async register(userData) {
    const response = await apiClient.post('/auth/register', userData)
    return response.data
  },

  async logout() {
    const response = await apiClient.post('/auth/logout')
    return response.data
  },

  async getCurrentUser() {
    const response = await apiClient.get('/auth/me')
    return response.data
  },

  async refreshToken() {
    const refreshToken = localStorage.getItem('refresh_token')
    const response = await apiClient.post('/auth/refresh', {
      refresh_token: refreshToken
    })
    return response.data
  },

  async requestPasswordReset(email) {
    const response = await apiClient.post('/auth/password-reset', { email })
    return response.data
  },

  async resetPassword(token, newPassword) {
    const response = await apiClient.post('/auth/password-reset/confirm', {
      token,
      new_password: newPassword
    })
    return response.data
  },

  async getOAuthProviders() {
    const response = await apiClient.get('/auth/oauth/providers')
    return response.data
  },

  async getLoginMethods(userId) {
    const response = await apiClient.get(`/auth/login-methods/${userId}`)
    return response.data
  },

  async getMyLoginMethods() {
    const response = await apiClient.get('/auth/login-methods')
    return response.data
  },

  async setOAuthPassword(data) {
    const response = await apiClient.post('/auth/oauth/set-password', data)
    return response.data
  },

  getOAuthAuthorizeURL(provider) {
    const apiBase = import.meta.env.VITE_API_BASE_URL || window.location.origin
    return `${apiBase}/api/auth/oauth/${provider}/authorize`
  }
}
