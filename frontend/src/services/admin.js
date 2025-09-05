import { apiClient } from './auth'

export const adminService = {
  // Get all transactions for admin (system-wide)
  async getAdminTransactions(limit = 100) {
    const response = await apiClient.get(`/admin/transactions?limit=${limit}`)
    return response.data
  },

  // Get all users for admin
  async getAdminUsers() {
    const response = await apiClient.get('/admin/users')
    return response.data
  },

  // Get admin dashboard data
  async getAdminDashboard() {
    const response = await apiClient.get('/admin/dashboard')
    return response.data
  },

  // Save system settings
  async saveSystemSettings(settings) {
    const response = await apiClient.post('/admin/settings', settings)
    return response.data
  },

  // Get system settings
  async getSystemSettings() {
    const response = await apiClient.get('/admin/settings')
    return response.data
  },

  // Get all support tickets for admin
  async getAdminSupportTickets() {
    const response = await apiClient.get('/admin/support/tickets')
    return response.data
  },

  // Reply to a support ticket
  async replyToTicket(ticketId, message) {
    const response = await apiClient.post(`/admin/support/tickets/${ticketId}/reply`, {
      message: message
    })
    return response.data
  },

  // Resolve a support ticket
  async resolveTicket(ticketId) {
    const response = await apiClient.post(`/admin/support/tickets/${ticketId}/resolve`)
    return response.data
  }
}
