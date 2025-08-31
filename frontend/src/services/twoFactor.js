import { apiClient } from './auth'

export const twoFactorService = {
  // Get 2FA status
  async getStatus() {
    const response = await apiClient.get('/2fa/status')
    return response.data
  },

  // Enable 2FA
  async enable(code) {
    const response = await apiClient.post('/2fa/enable', { code })
    return response.data
  },

  // Disable 2FA
  async disable(code) {
    const response = await apiClient.post('/2fa/disable', { code })
    return response.data
  },

  // Verify 2FA code
  async verify(code) {
    const response = await apiClient.post('/2fa/verify', { code })
    return response.data
  },

  // Login with 2FA
  async login2FA(userId, code) {
    const response = await apiClient.post('/auth/login/2fa', { user_id: userId, code })
    return response.data
  }
}
