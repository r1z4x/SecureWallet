import { apiClient } from './auth'

export const userService = {
  // Get current user profile
  async getCurrentUser() {
    const response = await apiClient.get('/auth/me')
    return response.data
  },

  // Update current user profile
  async updateCurrentUser(userData) {
    const response = await apiClient.put('/users/me', userData)
    return response.data
  },

  // Get user by ID (admin only)
  async getUser(userId) {
    const response = await apiClient.get(`/users/${userId}`)
    return response.data
  },

  // Update user by ID (admin only)
  async updateUser(userId, userData) {
    const response = await apiClient.put(`/users/${userId}`, userData)
    return response.data
  },

  // Delete user by ID (admin only)
  async deleteUser(userId) {
    const response = await apiClient.delete(`/users/${userId}`)
    return response.data
  },

  // Get all users (admin only)
  async getAllUsers() {
    const response = await apiClient.get('/users/')
    return response.data
  },

  // Change password
  async changePassword(passwordData) {
    const response = await apiClient.post('/auth/change-password', passwordData)
    return response.data
  },

  // Request password reset
  async requestPasswordReset(email) {
    const response = await apiClient.post('/auth/password-reset', { email })
    return response.data
  },

  // Reset password with token
  async resetPassword(token, newPassword) {
    const response = await apiClient.post('/auth/password-reset/confirm', {
      token,
      new_password: newPassword
    })
    return response.data
  },

  // Search users
  async searchUsers(query) {
    const response = await apiClient.get(`/users/search?q=${encodeURIComponent(query)}`)
    return response.data
  }
}
