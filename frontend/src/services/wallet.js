import { apiClient } from './auth'

export const walletService = {
  // Get wallet balance and transaction count
  async getBalance() {
    const response = await apiClient.get('/wallets/balance')
    return response.data
  },

  // Get all user wallets
  async getWallets() {
    const response = await apiClient.get('/wallets/')
    return response.data
  },

  // Get specific wallet by ID
  async getWallet(walletId) {
    const response = await apiClient.get(`/wallets/${walletId}`)
    return response.data
  },

  // Create new wallet
  async createWallet(walletData) {
    const response = await apiClient.post('/wallets/', walletData)
    return response.data
  },

  // Update wallet
  async updateWallet(walletId, walletData) {
    const response = await apiClient.put(`/wallets/${walletId}`, walletData)
    return response.data
  },

  // Delete wallet
  async deleteWallet(walletId) {
    const response = await apiClient.delete(`/wallets/${walletId}`)
    return response.data
  },

  // Transfer funds
  async transfer(transferData) {
    const response = await apiClient.post('/wallets/transfer', transferData)
    return response.data
  },

  // Deposit funds
  async deposit(depositData) {
    const response = await apiClient.post('/wallets/deposit', depositData)
    return response.data
  },

  // Withdraw funds
  async withdraw(withdrawalData) {
    const response = await apiClient.post('/wallets/withdraw', withdrawalData)
    return response.data
  },

  // Get wallet transactions
  async getWalletTransactions(walletId, limit = 50) {
    const response = await apiClient.get(`/wallets/${walletId}/transactions?limit=${limit}`)
    return response.data
  }
}
