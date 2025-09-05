<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader type="authenticated" />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">{{ $t('transfer.title') }}</h1>
        <p class="text-gray-600 mt-2">{{ $t('transfer.subtitle') }}</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Transfer Form -->
        <div class="lg:col-span-2">
          <div class="bg-white rounded-lg shadow-md p-6">
            <h2 class="text-xl font-semibold text-gray-900 mb-6">{{ $t('transfer.newTransfer') }}</h2>
            
            <form @submit.prevent="handleTransfer" class="space-y-6">
              <!-- Recipient -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  {{ $t('transfer.recipient') }}
                </label>
                <div class="relative">
                  <div class="flex space-x-2">
                    <input
                      v-model="transferForm.recipient"
                      type="text"
                      class="form-input flex-1"
                      :placeholder="$t('transfer.enterEmailOrUserId')"
                      @input="searchRecipient"
                      required
                    >
                    <button
                      type="button"
                      @click="searchRecipient"
                      class="btn-secondary"
                      :disabled="!transferForm.recipient || transferForm.recipient.length < 2"
                    >
                      <i v-if="searchLoading" class="fas fa-spinner fa-spin mr-2"></i>
                      <i v-else class="fas fa-search mr-2"></i>
                      {{ searchLoading ? $t('transfer.searching') : $t('transfer.search') }}
                    </button>
                  </div>
                  
                  <!-- Search Results Dropdown -->
                  <div
                    v-if="searchResults.length > 0"
                    class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto"
                  >
                    <div
                      v-for="user in searchResults"
                      :key="user.id"
                      @click="selectRecipient(user)"
                      class="px-4 py-3 hover:bg-gray-100 cursor-pointer border-b border-gray-200 last:border-b-0"
                    >
                      <div class="flex items-center">
                        <div class="flex-shrink-0 h-8 w-8 bg-primary-100 rounded-full flex items-center justify-center">
                          <span class="text-primary-600 font-medium text-sm">
                            {{ user.username.charAt(0).toUpperCase() }}
                          </span>
                        </div>
                        <div class="ml-3">
                          <p class="text-sm font-medium text-gray-900">{{ user.username }}</p>
                          <p class="text-xs text-gray-500">{{ user.email }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
                <p class="text-xs text-gray-500 mt-1">
                  {{ $t('transfer.enterEmailOrUserIdHelp') }}
                </p>
              </div>

              <!-- Amount -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  {{ $t('common.amount') }}
                </label>
                <div class="relative">
                  <span class="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-500">
                    $
                  </span>
                  <input
                    v-model="transferForm.amount"
                    type="number"
                    step="0.01"
                    min="0.01"
                    :max="walletData.balance"
                    class="form-input pl-8 w-full"
                    placeholder="0.00"
                    required
                  >
                </div>
                <div class="flex justify-between text-xs text-gray-500 mt-1">
                  <span>{{ $t('transfer.available') }}: ${{ walletData.balance?.toFixed(2) || '0.00' }}</span>
                  <span v-if="transferForm.amount">{{ $t('transfer.fee') }}: ${{ calculateTransferFee(transferForm.amount) }}</span>
                </div>
                <div v-if="transferForm.amount" class="flex justify-between text-sm font-medium text-gray-700 mt-2 pt-2 border-t border-gray-200">
                  <span>{{ $t('transfer.totalAmount') }}:</span>
                  <span class="text-primary-600">${{ (parseFloat(transferForm.amount) + parseFloat(calculateTransferFee(transferForm.amount))).toFixed(2) }}</span>
                </div>
              </div>

              <!-- Description -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  {{ $t('transfer.description') }}
                </label>
                <input
                  v-model="transferForm.description"
                  type="text"
                  class="form-input w-full"
                  :placeholder="$t('transfer.descriptionPlaceholder')"
                  maxlength="100"
                >
                <p class="text-xs text-gray-500 mt-1">
                  {{ transferForm.description?.length || 0 }}/100 {{ $t('transfer.characters') }}
                </p>
              </div>

              <!-- Transfer Button -->
              <div class="pt-4">
                <button
                  type="submit"
                  class="btn-primary w-full"
                  :disabled="transferLoading || !isFormValid"
                >
                  <i v-if="transferLoading" class="fas fa-spinner fa-spin mr-2"></i>
                  <i v-else class="fas fa-paper-plane mr-2"></i>
                  {{ transferLoading ? $t('transfer.processingTransfer') : $t('transfer.sendMoney') }}
                </button>
              </div>
            </form>

            <!-- Error Message -->
            <div v-if="error" class="mt-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
              {{ error }}
            </div>

            <!-- Success Message -->
            <div v-if="success" class="mt-4 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
              {{ success }}
            </div>
          </div>
        </div>

        <!-- Transfer Info -->
        <div class="space-y-6">
          <!-- Balance Card -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ $t('transfer.yourBalance') }}</h3>
            <div class="text-center">
              <p class="text-3xl font-bold text-primary-600">
                ${{ walletData.balance?.toFixed(2) || '0.00' }}
              </p>
              <p class="text-sm text-gray-500 mt-1">
                {{ $t('transfer.availableForTransfer') }}
              </p>
            </div>
          </div>

          <!-- Transfer Limits -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ $t('transfer.transferLimits') }}</h3>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600">{{ $t('transfer.dailyLimit') }}:</span>
                <span class="font-medium">$10,000</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">{{ $t('transfer.monthlyLimit') }}:</span>
                <span class="font-medium">$50,000</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">{{ $t('transfer.minTransfer') }}:</span>
                <span class="font-medium">$1.00</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">{{ $t('transfer.transferFee') }}:</span>
                <span class="font-medium text-red-600">1% (Min $1, Max $50)</span>
              </div>
            </div>
          </div>

          <!-- Recent Transfers -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ $t('transfer.recentTransfers') }}</h3>
            <div v-if="recentTransfers.length === 0" class="text-center py-4">
              <i class="fas fa-inbox text-2xl text-gray-300"></i>
              <p class="text-gray-500 text-sm mt-2">{{ $t('transfer.noRecentTransfers') }}</p>
            </div>
            <div v-else class="space-y-3">
              <div
                v-for="transfer in recentTransfers.slice(0, 3)"
                :key="transfer.id"
                class="flex items-center justify-between p-3 border border-gray-200 rounded-lg hover:bg-gray-50"
              >
                <div>
                  <p class="text-sm font-medium text-gray-900">
                    <!-- SECURE: Description is now properly escaped -->
                    {{ transfer.description || 'Transfer' }}
                  </p>
                  <p class="text-xs text-gray-500">
                    {{ formatDate(transfer.created_at) }}
                  </p>
                </div>
                <div class="text-right">
                  <p class="text-sm font-semibold text-red-600">
                    -${{ parseFloat(transfer.amount).toFixed(2) }}
                  </p>
                  <p class="text-xs text-gray-500">
                    {{ transfer.status }}
                  </p>
                </div>
              </div>
            </div>
            <div class="mt-4 text-center">
              <router-link
                to="/transactions"
                class="text-primary-600 hover:text-primary-800 text-sm font-medium"
              >
                {{ $t('transfer.viewAll') }} →
              </router-link>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/AppHeader.vue'
import { walletService } from '@/services/wallet'
import { transactionService } from '@/services/transaction'
import { userService } from '@/services/user'

export default {
  name: 'Transfer',
  components: {
    AppHeader
  },
  setup() {
    const authStore = useAuthStore()
    const user = computed(() => authStore.user)
    
    const transferLoading = ref(false)
    const error = ref('')
    const success = ref('')
    
    const walletData = ref({
      balance: 0,
      transaction_count: 0,
      currency: 'USD'
    })
    
    const userWallet = ref(null)
    const recentTransfers = ref([])
    const searchResults = ref([])
    const searchLoading = ref(false)
    
    const transferForm = ref({
      recipient: '',
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

    const loadRecentTransfers = async () => {
      try {
        const response = await transactionService.getTransactions(20) // Get more transactions to filter
        
        // Filter for outgoing transfers only
        recentTransfers.value = response.filter(t => 
          t.type === 'transfer' && 
          t.wallet_id === userWallet.value?.id
        ).slice(0, 5) // Get only the first 5
        
      } catch (error) {
        console.error('Error loading recent transfers:', error)
        recentTransfers.value = []
      }
    }

    let searchTimeout = null

    const searchRecipient = async () => {
      if (!transferForm.value.recipient || transferForm.value.recipient.length < 2) {
        searchResults.value = []
        return
      }
      
      // Clear previous timeout
      if (searchTimeout) {
        clearTimeout(searchTimeout)
      }
      
      // Debounce search
      searchTimeout = setTimeout(async () => {
        searchLoading.value = true
        searchResults.value = []
        
        try {
          const results = await userService.searchUsers(transferForm.value.recipient)
          searchResults.value = results
        } catch (error) {
          console.error('Error searching users:', error)
          // Don't show error to user for search failures
        } finally {
          searchLoading.value = false
        }
      }, 300) // 300ms delay
    }

    const selectRecipient = (user) => {
      transferForm.value.recipient = user.email
      searchResults.value = []
    }

    const handleTransfer = async () => {
      transferLoading.value = true
      error.value = ''
      success.value = ''
      
      try {
        // Validate amount
        const amount = parseFloat(transferForm.value.amount)
        if (amount > walletData.value.balance) {
          throw new Error('Insufficient balance')
        }
        
        if (amount <= 0) {
          throw new Error('Amount must be greater than 0')
        }
        
        // Validate recipient
        if (!transferForm.value.recipient) {
          throw new Error('Recipient is required')
        }
        
        // Process transfer
        await walletService.transfer({
          recipient: transferForm.value.recipient,
          amount: amount,
          description: transferForm.value.description || 'Transfer'
        })
        
        // Success
        success.value = `Successfully transferred $${amount.toFixed(2)} to ${transferForm.value.recipient}`
        
        // Reset form
        transferForm.value = {
          recipient: '',
          amount: '',
          description: ''
        }
        
        // Reload data
        await loadWalletData()
        await loadRecentTransfers()
        
      } catch (err) {
        error.value = err.response?.data?.detail || err.message || 'Transfer failed. Please try again.'
      } finally {
        transferLoading.value = false
      }
    }

    const calculateTransferFee = (amount) => {
      if (!amount || parseFloat(amount) <= 0) return '0.00'
      
      const transferAmount = parseFloat(amount)
      const feePercentage = 0.01 // 1%
      let fee = transferAmount * feePercentage
      
      // Apply min/max limits
      if (fee < 1.0) fee = 1.0
      if (fee > 50.0) fee = 50.0
      
      return fee.toFixed(2)
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    const isFormValid = computed(() => {
      const amount = parseFloat(transferForm.value.amount)
      const fee = parseFloat(calculateTransferFee(transferForm.value.amount))
      const totalAmount = amount + fee
      
      return transferForm.value.recipient && 
             transferForm.value.amount && 
             amount > 0 &&
             totalAmount <= walletData.value.balance
    })

    onMounted(async () => {
      await loadWalletData()
      await loadRecentTransfers()
    })

    return {
      user,
      transferLoading,
      error,
      success,
      walletData,
      recentTransfers,
      transferForm,
      isFormValid,
      calculateTransferFee,
      searchRecipient,
      selectRecipient,
      handleTransfer,
      formatDate,
      searchResults,
      searchLoading
    }
  }
}
</script>
