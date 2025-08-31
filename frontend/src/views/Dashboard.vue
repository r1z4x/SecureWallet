<template>
  <div class="min-h-screen bg-gray-50">
    <Navigation />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Welcome Section -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">
          Welcome back, {{ user?.username || user?.email || 'User' }}!
        </h1>
        <p class="text-gray-600 mt-2">Here's your financial overview</p>
      </div>

      <!-- Balance Card -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-lg font-semibold text-gray-900">Current Balance</h2>
            <p class="text-3xl font-bold text-primary-600">
              ${{ walletData.balance?.toFixed(2) || '0.00' }}
            </p>
            <p class="text-sm text-gray-500 mt-1">
              {{ walletData.transaction_count || 0 }} transactions total
            </p>
          </div>
          <div class="text-right">
            <div class="text-sm text-gray-500">Currency</div>
            <div class="text-lg font-semibold text-gray-900">
              {{ walletData.currency || 'USD' }}
            </div>
          </div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-paper-plane text-primary-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">Send Money</h3>
              <p class="text-gray-600 text-sm">Transfer funds to another user</p>
            </div>
          </div>
          <div class="mt-4">
            <router-link
              to="/transfer"
              class="btn-primary w-full"
            >
              <i class="fas fa-arrow-right mr-2"></i>
              Transfer Now
            </router-link>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-plus-circle text-green-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">Deposit</h3>
              <p class="text-gray-600 text-sm">Add money to your wallet</p>
            </div>
          </div>
          <div class="mt-4">
            <button
              @click="showDepositModal = true"
              class="btn-secondary w-full"
            >
              <i class="fas fa-plus mr-2"></i>
              Deposit
            </button>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-exchange-alt text-blue-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <h3 class="text-lg font-semibold text-gray-900">Transactions</h3>
              <p class="text-gray-600 text-sm">View your transaction history</p>
            </div>
          </div>
          <div class="mt-4">
            <router-link
              to="/transactions"
              class="btn-secondary w-full"
            >
              <i class="fas fa-list mr-2"></i>
              View All
            </router-link>
          </div>
        </div>
      </div>

      <!-- Recent Transactions -->
      <div class="bg-white rounded-lg shadow-md">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">Recent Transactions</h3>
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
              v-for="transaction in transactions.slice(0, 5)"
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
          
          <div v-if="transactions.length > 5" class="mt-6 text-center">
            <router-link
              to="/transactions"
              class="text-primary-600 hover:text-primary-800 text-sm font-medium"
            >
              View all transactions →
            </router-link>
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
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Navigation from '@/components/Navigation.vue'
import { walletService } from '@/services/wallet'
import { transactionService } from '@/services/transaction'

export default {
  name: 'Dashboard',
  components: {
    Navigation
  },
  setup() {
    const authStore = useAuthStore()
    const user = computed(() => authStore.user)
    
    const loading = ref(false)
    const depositLoading = ref(false)
    const showDepositModal = ref(false)
    
    const walletData = ref({
      balance: 0,
      transaction_count: 0,
      currency: 'USD'
    })
    
    const userWallet = ref(null)
    const transactions = ref([])
    
    const depositForm = ref({
      amount: '',
      description: ''
    })

    const loadWalletData = async () => {
      try {
        const response = await walletService.getBalance()
        walletData.value = response
        
        // Also get user's wallet details
        const wallets = await walletService.getWallets()
        userWallet.value = wallets[0] // Get first wallet
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
      // Determine if this is incoming or outgoing for the current user
      if (!userWallet.value) return 'text-gray-600'
      
      const isIncoming = transaction.to_wallet_id === userWallet.value.id
      const isOutgoing = transaction.from_wallet_id === userWallet.value.id
      
      if (isIncoming) return 'text-green-600'
      if (isOutgoing) return 'text-red-600'
      return 'text-gray-600'
    }

    const formatAmount = (transaction) => {
      const amount = parseFloat(transaction.amount)
      
      if (!userWallet.value) return `$${amount.toFixed(2)}`
      
      const isIncoming = transaction.to_wallet_id === userWallet.value.id
      const isOutgoing = transaction.from_wallet_id === userWallet.value.id
      
      if (isIncoming) return `+$${amount.toFixed(2)}`
      if (isOutgoing) return `-$${amount.toFixed(2)}`
      return `$${amount.toFixed(2)}`
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    onMounted(async () => {
      await loadWalletData()
      await loadTransactions()
    })

    return {
      user,
      loading,
      depositLoading,
      showDepositModal,
      walletData,
      transactions,
      depositForm,
      handleDeposit,
      getTransactionIcon,
      getTransactionColor,
      getAmountColor,
      formatAmount,
      formatDate
    }
  }
}
</script>
