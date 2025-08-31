import { apiClient } from './auth'

export const dataManagementService = {
  // Get database statistics
  async getStats() {
    const response = await apiClient.get('/data/stats')
    return response.data
  },

  // Reset database completely (clear, init, create sample data)
  async resetDatabase() {
    const response = await apiClient.post('/data/reset-database')
    return response.data
  }
}
