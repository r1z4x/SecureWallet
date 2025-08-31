import { apiClient } from './auth'

export const loginHistoryService = {
  // Get login history
  async getLoginHistory(limit = 50) {
    const response = await apiClient.get(`/login-history/?limit=${limit}`)
    return response.data
  },

  // Get recent login history
  async getRecentLoginHistory(limit = 10) {
    const response = await apiClient.get(`/login-history/recent/?limit=${limit}`)
    return response.data
  }
}
