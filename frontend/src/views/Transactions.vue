<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader type="authenticated" />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">{{ $t('transactions.title') }}</h1>
                  <p class="text-gray-600 mt-2">{{ $t('transactions.subtitle') }}</p>
      </div>

      <!-- Filters and Search -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
        <div class="mb-4">
          <h3 class="text-lg font-semibold text-gray-900 mb-2">{{ $t('transactions.filters') }}</h3>
          <p class="text-sm text-gray-600">{{ $t('transactions.filtersDescription') }}</p>
        </div>
        
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <!-- Transaction Type Filter -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">
              <i class="fas fa-filter mr-2 text-gray-400"></i>
              {{ $t('transactions.transactionType') }}
            </label>
            <div class="relative">
              <select
                v-model="filters.type"
                class="form-select"
                @change="loadTransactions"
              >
                <option value="">{{ $t('transactions.allTypes') }}</option>
                <option value="TRANSFER">{{ $t('transactions.transfer') }}</option>
                <option value="DEPOSIT">{{ $t('transactions.deposit') }}</option>
                <option value="WITHDRAWAL">{{ $t('transactions.withdrawal') }}</option>
              </select>
              <div class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                <i class="fas fa-chevron-down text-gray-400 text-xs"></i>
              </div>
            </div>
          </div>
          
          <!-- Status Filter -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">
              <i class="fas fa-circle mr-2 text-gray-400"></i>
              {{ $t('transactions.status') }}
            </label>
            <div class="relative">
              <select
                v-model="filters.status"
                class="form-select"
                @change="loadTransactions"
              >
                <option value="">{{ $t('transactions.allStatus') }}</option>
                <option value="COMPLETED">{{ $t('transactions.completed') }}</option>
                <option value="PENDING">{{ $t('transactions.pending') }}</option>
                <option value="FAILED">{{ $t('transactions.failed') }}</option>
              </select>
              <div class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                <i class="fas fa-chevron-down text-gray-400 text-xs"></i>
              </div>
            </div>
          </div>
          
          <!-- Date Range Filter -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">
              <i class="fas fa-calendar mr-2 text-gray-400"></i>
              {{ $t('transactions.dateRange') }}
            </label>
            <div class="relative">
              <select
                v-model="filters.dateRange"
                class="form-select"
                @change="loadTransactions"
              >
                <option value="7">{{ $t('transactions.last7Days') }}</option>
                <option value="30">{{ $t('transactions.last30Days') }}</option>
                <option value="90">{{ $t('transactions.last90Days') }}</option>
                <option value="365">{{ $t('transactions.lastYear') }}</option>
                <option value="all">{{ $t('transactions.allTime') }}</option>
              </select>
              <div class="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                <i class="fas fa-chevron-down text-gray-400 text-xs"></i>
              </div>
            </div>
          </div>
          
          <!-- Search Filter -->
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-700">
              <i class="fas fa-search mr-2 text-gray-400"></i>
              {{ $t('transactions.search') }}
            </label>
            <div class="relative">
              <input
                v-model="filters.search"
                type="text"
                class="form-input pl-10"
                :placeholder="$t('transactions.searchPlaceholder')"
                @input="debounceSearch"
              >
              <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
                <i class="fas fa-search text-gray-400 text-sm"></i>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Clear Filters Button -->
        <div class="mt-4 flex justify-end">
          <button
            @click="clearFilters"
            class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors"
          >
            <i class="fas fa-times mr-2"></i>
            {{ $t('transactions.clearFilters') }}
          </button>
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
              <p class="text-sm font-medium text-gray-600">{{ $t('transactions.totalTransactions') }}</p>
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
              <p class="text-sm font-medium text-gray-600">{{ $t('transactions.totalIncoming') }}</p>
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
              <p class="text-sm font-medium text-gray-600">{{ $t('transactions.totalOutgoing') }}</p>
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
              <p class="text-sm font-medium text-gray-600">{{ $t('transactions.netFlow') }}</p>
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
              {{ $t('transactions.transactions') }} ({{ transactions.length }})
            </h3>
            <div class="flex space-x-2">
              <button
                @click="exportTransactions"
                class="btn-secondary"
                :disabled="transactions.length === 0"
              >
                <i class="fas fa-download mr-2"></i>
                {{ $t('transactions.export') }}
              </button>
            </div>
          </div>
        </div>
        
        <div class="p-6">
          <div v-if="loading" class="text-center py-8">
            <i class="fas fa-spinner fa-spin text-2xl text-gray-400"></i>
            <p class="text-gray-500 mt-2">{{ $t('transactions.loadingTransactions') }}</p>
          </div>
          
          <div v-else-if="transactions.length === 0" class="text-center py-8">
            <i class="fas fa-inbox text-4xl text-gray-300"></i>
            <p class="text-gray-500 mt-2">{{ $t('transactions.noTransactionsFound') }}</p>
            <p class="text-gray-400 text-sm">{{ $t('transactions.tryAdjustingFilters') }}</p>
          </div>
          
          <div v-else class="space-y-4">
            <div
              v-for="transaction in transactions"
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
                    <!-- SECURE: Description is now properly escaped -->
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
          
          <!-- Pagination -->
          <div v-if="totalPages > 1" class="mt-8 flex justify-center">
            <nav class="flex space-x-2">
              <button
                @click="changePage(currentPage - 1)"
                :disabled="currentPage === 1"
                class="px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {{ $t('transactions.previous') }}
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
                {{ $t('transactions.next') }}
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
import AppHeader from '@/components/AppHeader.vue'
import { transactionService } from '@/services/transaction'

export default {
  name: 'Transactions',
  components: {
    AppHeader
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

    const clearFilters = () => {
      filters.value = {
        type: '',
        status: '',
        dateRange: '30',
        search: ''
      }
      currentPage.value = 1
      loadTransactions()
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
      // Determine if this is incoming or outgoing for the current user
      if (!user.value) return 'text-gray-600'
      
      // For TRANSFER transactions, check if it's outgoing (negative) or incoming (positive)
      if (transaction.transaction_type === 'TRANSFER') {
        // If description contains "fee", it's likely an outgoing transfer
        if (transaction.description && transaction.description.includes('fee')) {
          return 'text-red-600'
        }
        // If description contains "from", it's likely an incoming transfer
        if (transaction.description && transaction.description.includes('from')) {
          return 'text-green-600'
        }
      }
      
      // For other transaction types, they're positive
      return 'text-green-600'
    }

    const formatAmount = (transaction) => {
      const amount = parseFloat(transaction.amount)
      
      // For TRANSFER transactions, show negative for outgoing
      if (transaction.transaction_type === 'TRANSFER') {
        if (transaction.description && transaction.description.includes('fee')) {
          return `-$${amount.toFixed(2)}`
        }
        if (transaction.description && transaction.description.includes('from')) {
          return `+$${amount.toFixed(2)}`
        }
      }
      
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
      clearFilters,
      exportTransactions,
      getTransactionIcon,
      getTransactionColor,
      getTransactionBgColor,
      getAmountColor,
      formatAmount,
      formatDate
    }
  }
}
</script>
