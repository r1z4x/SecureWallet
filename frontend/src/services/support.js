import { apiClient } from './auth'

export const supportService = {
  // Create a new support ticket
  async createTicket(ticketData) {
    const response = await apiClient.post('/support/tickets', ticketData)
    return response.data
  },

  // Get all tickets for current user
  async getTickets() {
    const response = await apiClient.get('/support/tickets')
    return response.data
  },

  // Get specific ticket by ID
  async getTicket(ticketId) {
    const response = await apiClient.get(`/support/tickets/${ticketId}`)
    return response.data
  },

  // Update ticket
  async updateTicket(ticketId, ticketData) {
    const response = await apiClient.put(`/support/tickets/${ticketId}`, ticketData)
    return response.data
  },

  // Close ticket
  async closeTicket(ticketId) {
    const response = await apiClient.post(`/support/tickets/${ticketId}/close`)
    return response.data
  },

  // Add reply to ticket
  async addReply(ticketId, replyData) {
    const response = await apiClient.post(`/support/tickets/${ticketId}/replies`, replyData)
    return response.data
  },

  // Get ticket replies
  async getTicketReplies(ticketId) {
    const response = await apiClient.get(`/support/tickets/${ticketId}/replies`)
    return response.data
  }
}
