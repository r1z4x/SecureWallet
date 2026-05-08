import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { walletService } from '@/services/wallet'
import { transactionService } from '@/services/transaction'

const MOCK_WALLET_DATA = {
  balance: 5420.75,
  currency: 'USD',
  transaction_count: 12,
  wallet_id: 'mock-wallet-001',
  created_at: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString(),
  updated_at: new Date().toISOString()
}

const MOCK_TRANSACTIONS = [
  {
    id: 'mock-tx-001',
    type: 'DEPOSIT',
    transaction_type: 'DEPOSIT',
    amount: 1000.00,
    description: 'Initial deposit',
    status: 'completed',
    created_at: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 'mock-tx-002',
    type: 'TRANSFER',
    transaction_type: 'TRANSFER',
    amount: 250.00,
    description: 'Transfer from alice@example.com',
    status: 'completed',
    created_at: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 'mock-tx-003',
    type: 'DEPOSIT',
    transaction_type: 'DEPOSIT',
    amount: 2500.00,
    description: 'Paycheck deposit',
    status: 'completed',
    created_at: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 'mock-tx-004',
    type: 'TRANSFER',
    transaction_type: 'TRANSFER',
    amount: 150.00,
    description: 'Transfer fee - sent to bob@example.com',
    status: 'completed',
    created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString()
  },
  {
    id: 'mock-tx-005',
    type: 'WITHDRAWAL',
    transaction_type: 'WITHDRAWAL',
    amount: 500.00,
    description: 'ATM withdrawal',
    status: 'completed',
    created_at: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString()
  }
]

export const useWalletStore = defineStore('wallet', () => {
  const walletData = ref({
    balance: 0,
    currency: 'USD',
    transaction_count: 0,
    wallet_id: null,
    created_at: null,
    updated_at: null
  })

  const transactions = ref([])
  const wallets = ref([])
  const loading = ref(false)
  const error = ref(null)
  const isUsingMockData = ref(false)

  const formattedBalance = computed(() => {
    return walletData.value.balance?.toFixed(2) || '0.00'
  })

  const hasWallet = computed(() => {
    return walletData.value.balance > 0 || walletData.value.wallet_id !== null
  })

  async function loadBalance() {
    loading.value = true
    error.value = null

    try {
      const response = await walletService.getBalance()
      walletData.value = {
        balance: response.balance ?? 0,
        currency: response.currency || 'USD',
        transaction_count: response.transaction_count ?? 0,
        wallet_id: response.wallet_id || null,
        created_at: response.created_at || null,
        updated_at: response.updated_at || null
      }
      isUsingMockData.value = false
    } catch (err) {
      console.warn('Backend wallet API unavailable, using mock data:', err.message)
      walletData.value = { ...MOCK_WALLET_DATA }
      isUsingMockData.value = true
    } finally {
      loading.value = false
    }
  }

  async function loadTransactions(limit = 50) {
    loading.value = true
    error.value = null

    try {
      const response = await transactionService.getTransactions(limit)
      transactions.value = Array.isArray(response) ? response : []
      if (!isUsingMockData.value) {
        isUsingMockData.value = false
      }
    } catch (err) {
      console.warn('Backend transactions API unavailable, using mock data:', err.message)
      transactions.value = MOCK_TRANSACTIONS
      isUsingMockData.value = true
    } finally {
      loading.value = false
    }
  }

  async function loadWallets() {
    try {
      const response = await walletService.getWallets()
      wallets.value = Array.isArray(response) ? response : []
    } catch (err) {
      console.warn('Backend wallets API unavailable')
      wallets.value = []
    }
  }

  async function refreshAll() {
    await Promise.all([loadBalance(), loadWallets()])
  }

  function reset() {
    walletData.value = {
      balance: 0,
      currency: 'USD',
      transaction_count: 0,
      wallet_id: null,
      created_at: null,
      updated_at: null
    }
    transactions.value = []
    wallets.value = []
    loading.value = false
    error.value = null
    isUsingMockData.value = false
  }

  return {
    walletData,
    transactions,
    wallets,
    loading,
    error,
    isUsingMockData,
    formattedBalance,
    hasWallet,
    loadBalance,
    loadTransactions,
    loadWallets,
    refreshAll,
    reset
  }
})
