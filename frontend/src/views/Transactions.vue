<template>
  <div class="min-h-screen bg-gray-50">
    <Navigation />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Transaction History</h1>
        <p class="text-gray-600 mt-2">View and manage your financial transactions</p>
      </div>

      <!-- Filters and Search -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Transaction Type</label>
            <select
              v-model="filters.type"
              class="form-input w-full"
              @change="loadTransactions"
            >
              <option value="">All Types</option>
              <option value="TRANSFER">Transfer</option>
              <option value="DEPOSIT">Deposit</option>
              <option value="WITHDRAWAL">Withdrawal</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
            <select
              v-model="filters.status"
              class="form-input w-full"
              @change="loadTransactions"
            >
              <option value="">All Status</option>
              <option value="COMPLETED">Completed</option>
              <option value="PENDING">Pending</option>
              <option value="FAILED">Failed</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Date Range</label>
            <select
              v-model="filters.dateRange"
              class="form-input w-full"
              @change="loadTransactions"
            >
              <option value="7">Last 7 days</option>
              <option value="30">Last 30 days</option>
              <option value="90">Last 90 days</option>
              <option value="365">Last year</option>
              <option value="all">All time</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Search</label>
            <input
              v-model="filters.search"
              type="text"
              class="form-input w-full"
              placeholder="Search transactions..."
              @input="debounceSearch"
            >
          </div>
        </div>
      </div>

      <!-- Transaction Summary -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-exchange-alt text-blue-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Transactions</p>
              <p class="text-2xl font-bold text-gray-900">{{ summary.total }}</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-plus-circle text-green-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Incoming</p>
              <p class="text-2xl font-bold text-green-600">${{ summary.incoming?.toFixed(2) || '0.00' }}</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-minus-circle text-red-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Outgoing</p>
              <p class="text-2xl font-bold text-red-600">${{ summary.outgoing?.toFixed(2) || '0.00' }}</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-chart-line text-purple-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Net Flow</p>
              <p class="text-2xl font-bold" :class="summary.net >= 0 ? 'text-green-600' : 'text-red-600'">
                ${{ summary.net?.toFixed(2) || '0.00' }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Transactions List -->
      <div class="bg-white rounded-lg shadow-md">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-semibold text-gray-900">
              Transactions ({{ transactions.length }})
            </h3>
            <div class="flex space-x-2">
              <button
                @click="exportTransactions"
                class="btn-secondary"
                :disabled="transactions.length === 0"
              >
                <i class="fas fa-download mr-2"></i>
                Export
              </button>
            </div>
          </div>
        </div>
        
        <div class="p-6">
          <div v-if="loading" class="text-center py-8">
            <i class="fas fa-spinner fa-spin text-2xl text-gray-400"></i>
            <p class="text-gray-500 mt-2">Loading transactions...</p>
          </div>
          
          <div v-else-if="transactions.length === 0" class="text-center py-8">
            <i class="fas fa-inbox text-4xl text-gray-300"></i>
            <p class="text-gray-500 mt-2">No transactions found</p>
            <p class="text-gray-400 text-sm">Try adjusting your filters or make your first transaction</p>
          </div>
          
          <div v-else class="space-y-4">
            <div
              v-for="transaction in transactions"
              :key="transaction.id"
              class="border border-gray-200 rounded-lg p-6 hover:bg-gray-50 transition-colors"
            >
              <div class="flex items-center justify-between">
                <div class="flex items-center space-x-4">
                  <div class="flex-shrink-0">
                    <div class="w-12 h-12 rounded-full flex items-center justify-center"
                         :class="getTransactionBgColor(transaction)">
                      <i
                        :class="getTransactionIcon(transaction)"
                        class="text-white text-lg"
                      ></i>
                    </div>
                  </div>
                  
                  <div>
                    <h4 class="text-lg font-medium text-gray-900">
                      {{ transaction.description || 'Transaction' }}
                    </h4>
                    <div class="flex items-center space-x-4 text-sm text-gray-500 mt-1">
                      <span>{{ formatDate(transaction.created_at) }}</span>
                      <span class="flex items-center">
                        <i class="fas fa-circle text-xs mr-1"></i>
                        {{ transaction.transaction_type }}
                      </span>
                      <span class="flex items-center">
                        <i class="fas fa-circle text-xs mr-1"></i>
                        {{ transaction.status }}
                      </span>
                    </div>
                    
                    <!-- Transaction Details -->
                    <div class="mt-2 text-sm text-gray-600">
                      <div v-if="transaction.from_wallet && transaction.to_wallet">
                        <span v-if="transaction.from_wallet.user_id === user?.id">
                          To: {{ transaction.to_wallet.user?.email || 'Unknown' }}
                        </span>
                        <span v-else>
                          From: {{ transaction.from_wallet.user?.email || 'Unknown' }}
                        </span>
                      </div>
                      <div v-else-if="transaction.transaction_type === 'DEPOSIT'">
                        External deposit
                      </div>
                      <div v-else-if="transaction.transaction_type === 'WITHDRAWAL'">
                        External withdrawal
                      </div>
                    </div>
                  </div>
                </div>
                
                <div class="text-right">
                  <p
                    class="text-xl font-bold"
                    :class="getAmountColor(transaction)"
                  >
                    {{ formatAmount(transaction) }}
                  </p>
                  <p class="text-sm text-gray-500 mt-1">
                    Transaction ID: {{ transaction.id }}
                  </p>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Pagination -->
          <div v-if="totalPages > 1" class="mt-8 flex justify-center">
            <nav class="flex space-x-2">
              <button
                @click="changePage(currentPage - 1)"
                :disabled="currentPage === 1"
                class="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Previous
              </button>
              
              <button
                v-for="page in visiblePages"
                :key="page"
                @click="changePage(page)"
                :class="[
                  'px-3 py-2 border rounded-md text-sm font-medium',
                  page === currentPage
                    ? 'bg-primary-600 text-white border-primary-600'
                    : 'border-gray-300 text-gray-700 hover:bg-gray-50'
                ]"
              >
                {{ page }}
              </button>
              
              <button
                @click="changePage(currentPage + 1)"
                :disabled="currentPage === totalPages"
                class="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Next
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Navigation from '@/components/Navigation.vue'
import { transactionService } from '@/services/transaction'

export default {
  name: 'Transactions',
  components: {
    Navigation
  },
  setup() {
    const authStore = useAuthStore()
    const user = computed(() => authStore.user)
    
    const loading = ref(false)
    const transactions = ref([])
    const currentPage = ref(1)
    const totalPages = ref(1)
    const itemsPerPage = 20
    
    const filters = ref({
      type: '',
      status: '',
      dateRange: '30',
      search: ''
    })
    
    const summary = ref({
      total: 0,
      incoming: 0,
      outgoing: 0,
      net: 0
    })

    const loadTransactions = async () => {
      loading.value = true
      try {
        const response = await transactionService.getTransactions(100) // Get more for filtering
        transactions.value = response
        
        // Apply filters
        let filtered = response
        
        if (filters.value.type) {
          filtered = filtered.filter(t => t.transaction_type === filters.value.type)
        }
        
        if (filters.value.status) {
          filtered = filtered.filter(t => t.status === filters.value.status)
        }
        
        if (filters.value.search) {
          const search = filters.value.search.toLowerCase()
          filtered = filtered.filter(t => 
            t.description?.toLowerCase().includes(search) ||
            t.id.toString().includes(search)
          )
        }
        
        // Calculate summary
        calculateSummary(filtered)
        
        // Apply pagination
        const start = (currentPage.value - 1) * itemsPerPage
        const end = start + itemsPerPage
        transactions.value = filtered.slice(start, end)
        totalPages.value = Math.ceil(filtered.length / itemsPerPage)
        
      } catch (error) {
        console.error('Error loading transactions:', error)
      } finally {
        loading.value = false
      }
    }

    const calculateSummary = (transactions) => {
      let incoming = 0
      let outgoing = 0
      
      transactions.forEach(transaction => {
        const amount = parseFloat(transaction.amount)
        const isIncoming = transaction.to_wallet && transaction.to_wallet.user_id === user.value?.id
        
        if (isIncoming) {
          incoming += amount
        } else {
          outgoing += amount
        }
      })
      
      summary.value = {
        total: transactions.length,
        incoming,
        outgoing,
        net: incoming - outgoing
      }
    }

    const changePage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
        loadTransactions()
      }
    }

    const debounceSearch = () => {
      clearTimeout(searchTimeout.value)
      searchTimeout.value = setTimeout(() => {
        currentPage.value = 1
        loadTransactions()
      }, 300)
    }

    const searchTimeout = ref(null)

    const exportTransactions = () => {
      // Simple CSV export
      const csvContent = [
        ['Date', 'Type', 'Description', 'Amount', 'Status', 'Transaction ID'],
        ...transactions.value.map(t => [
          formatDate(t.created_at),
          t.transaction_type,
          t.description || '',
          formatAmount(t),
          t.status,
          t.id
        ])
      ].map(row => row.join(',')).join('\n')
      
      const blob = new Blob([csvContent], { type: 'text/csv' })
      const url = window.URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `transactions-${new Date().toISOString().split('T')[0]}.csv`
      a.click()
      window.URL.revokeObjectURL(url)
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

    const getTransactionBgColor = (transaction) => {
      switch (transaction.transaction_type) {
        case 'TRANSFER':
          return 'bg-blue-600'
        case 'DEPOSIT':
          return 'bg-green-600'
        case 'WITHDRAWAL':
          return 'bg-red-600'
        default:
          return 'bg-gray-600'
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
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    const visiblePages = computed(() => {
      const pages = []
      const start = Math.max(1, currentPage.value - 2)
      const end = Math.min(totalPages.value, currentPage.value + 2)
      
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
      
      return pages
    })

    onMounted(() => {
      loadTransactions()
    })

    return {
      user,
      loading,
      transactions,
      currentPage,
      totalPages,
      filters,
      summary,
      visiblePages,
      loadTransactions,
      changePage,
      debounceSearch,
      exportTransactions,
      getTransactionIcon,
      getTransactionBgColor,
      getAmountColor,
      formatAmount,
      formatDate
    }
  }
}
</script>
