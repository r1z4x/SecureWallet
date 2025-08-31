import { apiClient } from './auth'

export const transactionService = {
  // Get all transactions for current user
  async getTransactions(limit = 50) {
    const response = await apiClient.get(`/transactions/?limit=${limit}`)
    return response.data
  },

  // Get specific transaction by ID
  async getTransaction(transactionId) {
    const response = await apiClient.get(`/transactions/${transactionId}`)
    return response.data
  },

  // Create new transaction
  async createTransaction(transactionData) {
    const response = await apiClient.post('/transactions/', transactionData)
    return response.data
  },

  // Update transaction
  async updateTransaction(transactionId, transactionData) {
    const response = await apiClient.put(`/transactions/${transactionId}`, transactionData)
    return response.data
  },

  // Delete transaction
  async deleteTransaction(transactionId) {
    const response = await apiClient.delete(`/transactions/${transactionId}`)
    return response.data
  }
}
