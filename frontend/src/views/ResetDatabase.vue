<template>
  <div class="min-h-screen bg-gray-50">
    <Navigation />
    
    <!-- Header -->
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Reset Database</h1>
            <p class="text-sm text-gray-600 mt-1">Reset and recreate database with sample data</p>
          </div>
          <div class="flex items-center space-x-4">
            <button
              @click="resetDatabase"
              :disabled="loading"
              class="btn-danger flex items-center"
            >
              <i v-if="loading" class="fas fa-spinner fa-spin mr-2"></i>
              <i v-else class="fas fa-database mr-2"></i>
              {{ loading ? 'Resetting...' : 'Reset Database' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Database Statistics -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
          <h2 class="text-xl font-semibold text-gray-900 mb-6">Database Statistics</h2>
          
          <div v-if="loading" class="text-center py-8">
            <i class="fas fa-spinner fa-spin text-2xl text-primary-600"></i>
            <p class="text-gray-600 mt-2">Loading statistics...</p>
          </div>
          
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div class="flex items-center justify-between p-4 bg-blue-50 rounded-lg">
              <div class="flex items-center">
                <i class="fas fa-users text-blue-600 text-xl mr-3"></i>
                <div>
                  <p class="text-lg font-bold text-blue-600">{{ stats.users || 0 }}</p>
                  <p class="text-sm font-medium text-blue-900">Total Users</p>
                </div>
              </div>
            </div>
            
            <div class="flex items-center justify-between p-4 bg-green-50 rounded-lg">
              <div class="flex items-center">
                <i class="fas fa-wallet text-green-600 text-xl mr-3"></i>
                <div>
                  <p class="text-lg font-bold text-green-600">{{ stats.wallets || 0 }}</p>
                  <p class="text-sm font-medium text-green-900">Total Wallets</p>
                </div>
              </div>
            </div>
            
            <div class="flex items-center justify-between p-4 bg-purple-50 rounded-lg">
              <div class="flex items-center">
                <i class="fas fa-exchange-alt text-purple-600 text-xl mr-3"></i>
                <div>
                  <p class="text-lg font-bold text-purple-600">{{ stats.transactions || 0 }}</p>
                  <p class="text-sm font-medium text-purple-900">Total Transactions</p>
                </div>
              </div>
            </div>
            
            <div class="flex items-center justify-between p-4 bg-orange-50 rounded-lg">
              <div class="flex items-center">
                <i class="fas fa-history text-orange-600 text-xl mr-3"></i>
                <div>
                  <p class="text-lg font-bold text-orange-600">{{ stats.login_history || 0 }}</p>
                  <p class="text-sm font-medium text-orange-900">Login History</p>
                </div>
              </div>
            </div>
          </div>
        </div>

      <!-- Data Requirements -->
      <div class="mt-8 bg-white rounded-lg shadow-md p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Data Requirements</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="space-y-3">
            <div class="flex items-center">
              <i class="fas fa-check-circle text-green-600 mr-3"></i>
              <span class="text-gray-700">Minimum 500 records in database</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-check-circle text-green-600 mr-3"></i>
              <span class="text-gray-700">Consistent and realistic data</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-check-circle text-green-600 mr-3"></i>
              <span class="text-gray-700">Multiple user accounts</span>
            </div>
          </div>
          <div class="space-y-3">
            <div class="flex items-center">
              <i class="fas fa-check-circle text-green-600 mr-3"></i>
              <span class="text-gray-700">Transaction history</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-check-circle text-green-600 mr-3"></i>
              <span class="text-gray-700">Login history tracking</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-check-circle text-green-600 mr-3"></i>
              <span class="text-gray-700">Wallet balances</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Success/Error Messages -->
      <div v-if="success" class="mt-6 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
        {{ success }}
      </div>
      
      <div v-if="error" class="mt-6 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import Navigation from '@/components/Navigation.vue'
import { dataManagementService } from '@/services/dataManagement'

export default {
  name: 'ResetDatabase',
  components: {
    Navigation
  },
  setup() {
    const loading = ref(false)
    const stats = ref({})
    const success = ref('')
    const error = ref('')

    const loadStats = async () => {
      try {
        const response = await dataManagementService.getStats()
        stats.value = response.stats || {}
      } catch (err) {
        console.error('Error loading stats:', err)
        if (err.response?.data?.error) {
          error.value = `Failed to load database statistics: ${err.response.data.error}`
        } else if (err.message) {
          error.value = `Failed to load database statistics: ${err.message}`
        } else {
          error.value = 'Failed to load database statistics: Unknown error occurred'
        }
      }
    }

    const resetDatabase = async () => {
      if (!confirm('Are you sure you want to reset the database? This will clear all data and recreate sample data. This action cannot be undone.')) {
        return
      }
      
      loading.value = true
      error.value = ''
      success.value = ''
      
      try {
        // Use the new single endpoint that handles everything in the backend
        const response = await dataManagementService.resetDatabase()
        
        success.value = response.message || 'Database reset successfully! All data has been cleared and recreated.'
        stats.value = response.stats || {}
      } catch (err) {
        console.error('Reset database error:', err)
        if (err.response?.data?.error) {
          error.value = `Failed to reset database: ${err.response.data.error}`
        } else if (err.response?.data?.details) {
          error.value = `Failed to reset database: ${err.response.data.details}`
        } else if (err.message) {
          error.value = `Failed to reset database: ${err.message}`
        } else {
          error.value = 'Failed to reset database: Unknown error occurred'
        }
      } finally {
        loading.value = false
      }
    }

    onMounted(() => {
      loadStats()
    })

    return {
      loading,
      stats,
      success,
      error,
      resetDatabase
    }
  }
}
</script>
