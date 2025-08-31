<template>
  <div class="min-h-screen bg-gray-50">
    <Navigation />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Admin Panel</h1>
        <p class="text-gray-600 mt-2">Manage users, transactions, and system settings</p>
      </div>

      <!-- Admin Stats -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-users text-blue-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Users</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.totalUsers }}</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-exchange-alt text-green-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Transactions</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.totalTransactions }}</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-dollar-sign text-yellow-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Volume</p>
              <p class="text-2xl font-bold text-gray-900">${{ stats.totalVolume?.toFixed(2) || '0.00' }}</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-ticket-alt text-purple-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Support Tickets</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.totalTickets }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Admin Tabs -->
      <div class="bg-white rounded-lg shadow-md">
        <div class="border-b border-gray-200">
          <nav class="-mb-px flex space-x-8 px-6">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              @click="activeTab = tab.id"
              :class="[
                'py-4 px-1 border-b-2 font-medium text-sm',
                activeTab === tab.id
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              <i :class="tab.icon" class="mr-2"></i>
              {{ tab.name }}
            </button>
          </nav>
        </div>

        <div class="p-6">
          <!-- Users Tab -->
          <div v-if="activeTab === 'users'" class="space-y-6">
            <div class="flex justify-between items-center">
              <h2 class="text-xl font-semibold text-gray-900">User Management</h2>
              <button class="btn-primary">
                <i class="fas fa-plus mr-2"></i>
                Add User
              </button>
            </div>
            
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Joined</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="user in users" :key="user.id">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <div class="flex-shrink-0 h-10 w-10">
                          <div class="h-10 w-10 rounded-full bg-primary-100 flex items-center justify-center">
                            <span class="text-primary-600 font-medium">{{ user.username.charAt(0).toUpperCase() }}</span>
                          </div>
                        </div>
                        <div class="ml-4">
                          <div class="text-sm font-medium text-gray-900">{{ user.username }}</div>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.email }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                        ]"
                      >
                        {{ user.is_active ? 'Active' : 'Inactive' }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {{ user.is_admin ? 'Admin' : 'User' }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ formatDate(user.created_at) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button class="text-primary-600 hover:text-primary-900 mr-3">Edit</button>
                      <button class="text-red-600 hover:text-red-900">Delete</button>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Transactions Tab -->
          <div v-if="activeTab === 'transactions'" class="space-y-6">
            <div class="flex justify-between items-center">
              <h2 class="text-xl font-semibold text-gray-900">Transaction History</h2>
              <div class="flex space-x-2">
                <select v-model="transactionFilter" class="form-input">
                  <option value="">All Types</option>
                  <option value="transfer">Transfer</option>
                  <option value="deposit">Deposit</option>
                  <option value="withdrawal">Withdrawal</option>
                </select>
                <button class="btn-secondary">
                  <i class="fas fa-download mr-2"></i>
                  Export
                </button>
              </div>
            </div>
            
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="transaction in filteredTransactions" :key="transaction.id">
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">#{{ transaction.id }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          getTransactionTypeColor(transaction.transaction_type)
                        ]"
                      >
                        {{ transaction.transaction_type }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      ${{ parseFloat(transaction.amount).toFixed(2) }}
                    </td>
                    <td class="px-6 py-4 text-sm text-gray-900">{{ transaction.description }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          transaction.status === 'completed' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                        ]"
                      >
                        {{ transaction.status }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ formatDate(transaction.created_at) }}
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- Support Tab -->
          <div v-if="activeTab === 'support'" class="space-y-6">
            <div class="flex justify-between items-center">
              <h2 class="text-xl font-semibold text-gray-900">Support Tickets</h2>
              <div class="flex space-x-2">
                <select v-model="ticketFilter" class="form-input">
                  <option value="">All Status</option>
                  <option value="open">Open</option>
                  <option value="in_progress">In Progress</option>
                  <option value="resolved">Resolved</option>
                  <option value="closed">Closed</option>
                </select>
              </div>
            </div>
            
            <div class="space-y-4">
              <div
                v-for="ticket in filteredTickets"
                :key="ticket.id"
                class="border border-gray-200 rounded-lg p-6 hover:bg-gray-50"
              >
                <div class="flex items-start justify-between">
                  <div class="flex-1">
                    <div class="flex items-center space-x-3 mb-2">
                      <h3 class="text-lg font-medium text-gray-900">#{{ ticket.id }} - {{ ticket.subject }}</h3>
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          getTicketStatusColor(ticket.status)
                        ]"
                      >
                        {{ ticket.status }}
                      </span>
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          getTicketPriorityColor(ticket.priority)
                        ]"
                      >
                        {{ ticket.priority }}
                      </span>
                    </div>
                    <p class="text-gray-600 mb-3">{{ ticket.message }}</p>
                    <div class="flex items-center space-x-4 text-sm text-gray-500">
                      <span>User ID: {{ ticket.user_id }}</span>
                      <span>Category: {{ ticket.category }}</span>
                      <span>{{ formatDate(ticket.created_at) }}</span>
                    </div>
                  </div>
                  <div class="ml-4">
                    <button class="btn-secondary mr-2">
                      <i class="fas fa-reply mr-2"></i>
                      Reply
                    </button>
                    <button class="btn-primary">
                      <i class="fas fa-check mr-2"></i>
                      Resolve
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Settings Tab -->
          <div v-if="activeTab === 'settings'" class="space-y-6">
            <h2 class="text-xl font-semibold text-gray-900">System Settings</h2>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="bg-gray-50 rounded-lg p-6">
                <h3 class="text-lg font-medium text-gray-900 mb-4">Security Settings</h3>
                <div class="space-y-4">
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">Two-Factor Authentication</span>
                    <button class="btn-secondary">Enable</button>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">Session Timeout</span>
                    <select class="form-input w-32">
                      <option>30 min</option>
                      <option>1 hour</option>
                      <option>2 hours</option>
                    </select>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">Password Policy</span>
                    <button class="btn-secondary">Configure</button>
                  </div>
                </div>
              </div>
              
              <div class="bg-gray-50 rounded-lg p-6">
                <h3 class="text-lg font-medium text-gray-900 mb-4">Transaction Limits</h3>
                <div class="space-y-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Daily Transfer Limit</label>
                    <input type="number" class="form-input w-full" value="10000" />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Monthly Transfer Limit</label>
                    <input type="number" class="form-input w-full" value="50000" />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">Minimum Transfer Amount</label>
                    <input type="number" class="form-input w-full" value="1" />
                  </div>
                </div>
              </div>
            </div>
            
            <div class="flex justify-end">
              <button class="btn-primary">
                <i class="fas fa-save mr-2"></i>
                Save Settings
              </button>
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
import Navigation from '@/components/Navigation.vue'
import { userService } from '@/services/user'
import { transactionService } from '@/services/transaction'
import { supportService } from '@/services/support'

export default {
  name: 'Admin',
  components: {
    Navigation
  },
  setup() {
    const authStore = useAuthStore()
    const user = computed(() => authStore.user)
    
    const activeTab = ref('users')
    const transactionFilter = ref('')
    const ticketFilter = ref('')
    
    const users = ref([])
    const transactions = ref([])
    const tickets = ref([])
    const stats = ref({
      totalUsers: 0,
      totalTransactions: 0,
      totalVolume: 0,
      totalTickets: 0
    })

    const tabs = [
      { id: 'users', name: 'Users', icon: 'fas fa-users' },
      { id: 'transactions', name: 'Transactions', icon: 'fas fa-exchange-alt' },
      { id: 'support', name: 'Support', icon: 'fas fa-headset' },
      { id: 'settings', name: 'Settings', icon: 'fas fa-cog' }
    ]

    const loadUsers = async () => {
      try {
        const response = await userService.getAllUsers()
        users.value = response
        stats.value.totalUsers = response.length
      } catch (error) {
        console.error('Error loading users:', error)
      }
    }

    const loadTransactions = async () => {
      try {
        const response = await transactionService.getTransactions(100)
        transactions.value = response
        stats.value.totalTransactions = response.length
        stats.value.totalVolume = response.reduce((sum, t) => sum + parseFloat(t.amount), 0)
      } catch (error) {
        console.error('Error loading transactions:', error)
      }
    }

    const loadTickets = async () => {
      try {
        const response = await supportService.getTickets()
        tickets.value = response
        stats.value.totalTickets = response.length
      } catch (error) {
        console.error('Error loading tickets:', error)
      }
    }

    const filteredTransactions = computed(() => {
      if (!transactionFilter.value) return transactions.value
      return transactions.value.filter(t => t.transaction_type === transactionFilter.value)
    })

    const filteredTickets = computed(() => {
      if (!ticketFilter.value) return tickets.value
      return tickets.value.filter(t => t.status === ticketFilter.value)
    })

    const getTransactionTypeColor = (type) => {
      switch (type) {
        case 'transfer':
          return 'bg-blue-100 text-blue-800'
        case 'deposit':
          return 'bg-green-100 text-green-800'
        case 'withdrawal':
          return 'bg-red-100 text-red-800'
        default:
          return 'bg-gray-100 text-gray-800'
      }
    }

    const getTicketStatusColor = (status) => {
      switch (status) {
        case 'open':
          return 'bg-orange-100 text-orange-800'
        case 'in_progress':
          return 'bg-blue-100 text-blue-800'
        case 'resolved':
          return 'bg-green-100 text-green-800'
        case 'closed':
          return 'bg-gray-100 text-gray-800'
        default:
          return 'bg-gray-100 text-gray-800'
      }
    }

    const getTicketPriorityColor = (priority) => {
      switch (priority) {
        case 'urgent':
          return 'bg-red-100 text-red-800'
        case 'high':
          return 'bg-orange-100 text-orange-800'
        case 'medium':
          return 'bg-yellow-100 text-yellow-800'
        case 'low':
          return 'bg-green-100 text-green-800'
        default:
          return 'bg-gray-100 text-gray-800'
      }
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

    onMounted(() => {
      loadUsers()
      loadTransactions()
      loadTickets()
    })

    return {
      user,
      activeTab,
      transactionFilter,
      ticketFilter,
      users,
      transactions,
      tickets,
      stats,
      tabs,
      filteredTransactions,
      filteredTickets,
      getTransactionTypeColor,
      getTicketStatusColor,
      getTicketPriorityColor,
      formatDate
    }
  }
}
</script>
