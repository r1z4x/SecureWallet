<template>
  <div class="min-h-screen bg-gray-50">
    <Navigation />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">My Wallet</h1>
        <p class="text-gray-600 mt-2">Manage your digital wallet and transactions</p>
      </div>

      <!-- Wallet Overview -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8 mb-8">
        <!-- Balance Card -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">Current Balance</h3>
            <i class="fas fa-wallet text-primary-600 text-2xl"></i>
          </div>
          <div class="text-center">
            <p class="text-4xl font-bold text-primary-600">
              ${{ walletData.balance?.toFixed(2) || '0.00' }}
            </p>
            <p class="text-sm text-gray-500 mt-2">
              {{ walletData.currency || 'USD' }}
            </p>
          </div>
          <div class="mt-6 space-y-2">
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Total Transactions:</span>
              <span class="font-medium">{{ walletData.transaction_count || 0 }}</span>
            </div>
            <div class="flex justify-between text-sm">
              <span class="text-gray-600">Wallet Status:</span>
              <span class="text-green-600 font-medium">Active</span>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
          <div class="space-y-3">
            <button
              @click="showDepositModal = true"
              class="w-full btn-primary"
            >
              <i class="fas fa-plus mr-2"></i>
              Deposit Funds
            </button>
            <button
              @click="showWithdrawModal = true"
              class="w-full btn-secondary"
            >
              <i class="fas fa-minus mr-2"></i>
              Withdraw Funds
            </button>
            <router-link
              to="/transfer"
              class="w-full btn-secondary block text-center"
            >
              <i class="fas fa-paper-plane mr-2"></i>
              Send Money
            </router-link>
          </div>
        </div>

        <!-- Wallet Info -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Wallet Information</h3>
          <div class="space-y-3">
            <div class="flex justify-between">
              <span class="text-gray-600">Wallet ID:</span>
              <span class="font-mono text-sm">{{ walletData.wallet_id || 'N/A' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">Created:</span>
              <span class="text-sm">{{ formatDate(walletData.created_at) || 'N/A' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-gray-600">Last Updated:</span>
              <span class="text-sm">{{ formatDate(walletData.updated_at) || 'N/A' }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Transactions -->
      <div class="bg-white rounded-lg shadow-md">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-semibold text-gray-900">Recent Transactions</h3>
            <router-link
              to="/transactions"
              class="text-primary-600 hover:text-primary-800 text-sm font-medium"
            >
              View All →
            </router-link>
          </div>
        </div>
        <div class="p-6">
          <div v-if="loading" class="text-center py-8">
            <i class="fas fa-spinner fa-spin text-2xl text-gray-400"></i>
            <p class="text-gray-500 mt-2">Loading transactions...</p>
          </div>
          
          <div v-else-if="transactions.length === 0" class="text-center py-8">
            <i class="fas fa-inbox text-4xl text-gray-300"></i>
            <p class="text-gray-500 mt-2">No transactions yet</p>
            <p class="text-gray-400 text-sm">Your transaction history will appear here</p>
          </div>
          
          <div v-else class="space-y-4">
            <div
              v-for="transaction in transactions.slice(0, 10)"
              :key="transaction.id"
              class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50"
            >
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <i
                    :class="[getTransactionIcon(transaction), getTransactionColor(transaction)]"
                    class="text-lg"
                  ></i>
                </div>
                <div class="ml-4">
                  <p class="text-sm font-medium text-gray-900">
                    {{ transaction.description || 'Transaction' }}
                  </p>
                  <p class="text-xs text-gray-500">
                    {{ formatDate(transaction.created_at) }}
                  </p>
                </div>
              </div>
              <div class="text-right">
                <p
                  class="text-sm font-semibold"
                  :class="getAmountColor(transaction)"
                >
                  {{ formatAmount(transaction) }}
                </p>
                <p class="text-xs text-gray-500">
                  {{ transaction.status }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Deposit Modal -->
    <div
      v-if="showDepositModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showDepositModal = false"
    >
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Deposit Funds</h3>
          <form @submit.prevent="handleDeposit">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Amount
              </label>
              <input
                v-model="depositForm.amount"
                type="number"
                step="0.01"
                min="0.01"
                class="form-input w-full"
                placeholder="Enter amount"
                required
              >
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Description (Optional)
              </label>
              <input
                v-model="depositForm.description"
                type="text"
                class="form-input w-full"
                placeholder="e.g., Bank transfer"
              >
            </div>
            <div class="flex justify-end space-x-3">
              <button
                type="button"
                @click="showDepositModal = false"
                class="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="btn-primary"
                :disabled="depositLoading"
              >
                <i v-if="depositLoading" class="fas fa-spinner fa-spin mr-2"></i>
                {{ depositLoading ? 'Processing...' : 'Deposit' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Withdraw Modal -->
    <div
      v-if="showWithdrawModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showWithdrawModal = false"
    >
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Withdraw Funds</h3>
          <form @submit.prevent="handleWithdraw">
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Amount
              </label>
              <input
                v-model="withdrawForm.amount"
                type="number"
                step="0.01"
                min="0.01"
                :max="walletData.balance"
                class="form-input w-full"
                placeholder="Enter amount"
                required
              >
              <p class="text-xs text-gray-500 mt-1">
                Available balance: ${{ walletData.balance?.toFixed(2) || '0.00' }}
              </p>
            </div>
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Description (Optional)
              </label>
              <input
                v-model="withdrawForm.description"
                type="text"
                class="form-input w-full"
                placeholder="e.g., Bank withdrawal"
              >
            </div>
            <div class="flex justify-end space-x-3">
              <button
                type="button"
                @click="showWithdrawModal = false"
                class="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="btn-primary"
                :disabled="withdrawLoading"
              >
                <i v-if="withdrawLoading" class="fas fa-spinner fa-spin mr-2"></i>
                {{ withdrawLoading ? 'Processing...' : 'Withdraw' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Navigation from '@/components/Navigation.vue'
import { walletService } from '@/services/wallet'
import { transactionService } from '@/services/transaction'

export default {
  name: 'Wallet',
  components: {
    Navigation
  },
  setup() {
    const authStore = useAuthStore()
    const user = computed(() => authStore.user)
    
    const loading = ref(false)
    const depositLoading = ref(false)
    const withdrawLoading = ref(false)
    const showDepositModal = ref(false)
    const showWithdrawModal = ref(false)
    
    const walletData = ref({
      balance: 0,
      transaction_count: 0,
      currency: 'USD',
      wallet_id: null,
      created_at: null,
      updated_at: null
    })
    
    const transactions = ref([])
    
    const depositForm = ref({
      amount: '',
      description: ''
    })

    const withdrawForm = ref({
      amount: '',
      description: ''
    })

    const loadWalletData = async () => {
      try {
        const response = await walletService.getBalance()
        walletData.value = response
      } catch (error) {
        console.error('Error loading wallet data:', error)
      }
    }

    const loadTransactions = async () => {
      loading.value = true
      try {
        const response = await transactionService.getTransactions(10)
        transactions.value = response
      } catch (error) {
        console.error('Error loading transactions:', error)
      } finally {
        loading.value = false
      }
    }

    const handleDeposit = async () => {
      depositLoading.value = true
      try {
        await walletService.deposit(depositForm.value)
        await loadWalletData()
        await loadTransactions()
        showDepositModal.value = false
        depositForm.value = { amount: '', description: '' }
      } catch (error) {
        console.error('Error processing deposit:', error)
      } finally {
        depositLoading.value = false
      }
    }

    const handleWithdraw = async () => {
      withdrawLoading.value = true
      try {
        await walletService.withdraw(withdrawForm.value)
        await loadWalletData()
        await loadTransactions()
        showWithdrawModal.value = false
        withdrawForm.value = { amount: '', description: '' }
      } catch (error) {
        console.error('Error processing withdrawal:', error)
      } finally {
        withdrawLoading.value = false
      }
    }

    const getTransactionIcon = (transaction) => {
      switch (transaction.transaction_type) {
        case 'TRANSFER':
          return 'fas fa-exchange-alt'
        case 'DEPOSIT':
          return 'fas fa-plus-circle'
        case 'WITHDRAWAL':
          return 'fas fa-minus-circle'
        default:
          return 'fas fa-circle'
      }
    }

    const getTransactionColor = (transaction) => {
      switch (transaction.transaction_type) {
        case 'TRANSFER':
          return 'text-blue-600'
        case 'DEPOSIT':
          return 'text-green-600'
        case 'WITHDRAWAL':
          return 'text-red-600'
        default:
          return 'text-gray-600'
      }
    }

    const getAmountColor = (transaction) => {
      const isIncoming = transaction.to_wallet && transaction.to_wallet.user_id === user.value?.id
      return isIncoming ? 'text-green-600' : 'text-red-600'
    }

    const formatAmount = (transaction) => {
      const amount = parseFloat(transaction.amount)
      const isIncoming = transaction.to_wallet && transaction.to_wallet.user_id === user.value?.id
      const prefix = isIncoming ? '+' : '-'
      return `${prefix}$${amount.toFixed(2)}`
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    onMounted(() => {
      loadWalletData()
      loadTransactions()
    })

    return {
      user,
      loading,
      depositLoading,
      withdrawLoading,
      showDepositModal,
      showWithdrawModal,
      walletData,
      transactions,
      depositForm,
      withdrawForm,
      handleDeposit,
      handleWithdraw,
      getTransactionIcon,
      getTransactionColor,
      getAmountColor,
      formatAmount,
      formatDate
    }
  }
}
</script>
