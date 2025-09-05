<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader type="authenticated" />
    
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
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
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

    <!-- Confirmation Overlay -->
    <div v-if="isPageBlocked" class="fixed inset-0 bg-black bg-opacity-50 z-[10000] flex items-center justify-center">
      <transition-group name="confirmation" tag="div" class="space-y-4">
        <div
          v-for="confirmation in pendingConfirmations"
          :key="confirmation.id"
          class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6"
        >
          <div class="flex items-center mb-4">
            <div class="flex-shrink-0">
              <i 
                :class="[
                  'h-6 w-6',
                  {
                    'fas fa-exclamation-triangle text-yellow-500': confirmation.type === 'warning',
                    'fas fa-exclamation-circle text-red-500': confirmation.type === 'error',
                    'fas fa-question-circle text-blue-500': confirmation.type === 'info',
                    'fas fa-check-circle text-green-500': confirmation.type === 'success'
                  }
                ]"
              ></i>
            </div>
            <div class="ml-3">
              <h3 class="text-lg font-medium text-gray-900">
                {{ confirmation.title }}
              </h3>
            </div>
          </div>
          
          <div class="mb-6">
            <p class="text-sm text-gray-600 break-words">
              {{ confirmation.message }}
            </p>
          </div>
          
          <div class="flex justify-end space-x-3">
            <button
              @click="handleConfirmation(confirmation.id, false)"
              :class="[
                'px-4 py-2 text-sm font-medium rounded-md text-white focus:outline-none focus:ring-2 focus:ring-offset-2',
                confirmation.cancelClass
              ]"
            >
              {{ confirmation.cancelText }}
            </button>
            <button
              @click="handleConfirmation(confirmation.id, true)"
              :class="[
                'px-4 py-2 text-sm font-medium rounded-md text-white focus:outline-none focus:ring-2 focus:ring-offset-2',
                confirmation.confirmClass
              ]"
            >
              {{ confirmation.confirmText }}
            </button>
          </div>
        </div>
      </transition-group>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import AppHeader from '@/components/AppHeader.vue'
import { dataManagementService } from '@/services/dataManagement'

export default {
  name: 'ResetDatabase',
  components: {
    AppHeader
  },
  setup() {
    const { t } = useI18n()
    const loading = ref(false)
    const stats = ref({})
    const success = ref('')
    const error = ref('')
    
    // Confirmation states
    const pendingConfirmations = ref([])
    const confirmationId = ref(0)
    const isPageBlocked = ref(false)

    // Confirmation System
    const showConfirmation = (message, options = {}) => {
      return new Promise((resolve) => {
        const id = ++confirmationId.value
        const confirmation = {
          id,
          message: String(message),
          type: options.type || 'warning',
          title: options.title || t('common.confirmation'),
          confirmText: options.confirmText || t('common.yes'),
          cancelText: options.cancelText || t('common.no'),
          confirmClass: options.confirmClass || 'bg-red-500 hover:bg-red-600',
          cancelClass: options.cancelClass || 'bg-gray-500 hover:bg-gray-600',
          timestamp: Date.now(),
          resolve
        }
        
        pendingConfirmations.value.push(confirmation)
        isPageBlocked.value = true
        
        return confirmation
      })
    }

    const handleConfirmation = (id, confirmed) => {
      const index = pendingConfirmations.value.findIndex(c => c.id === id)
      if (index > -1) {
        const confirmation = pendingConfirmations.value[index]
        confirmation.resolve(confirmed)
        pendingConfirmations.value.splice(index, 1)
        
        // Unblock page if no more confirmations
        if (pendingConfirmations.value.length === 0) {
          isPageBlocked.value = false
        }
      }
    }

    const confirmAction = (message, options = {}) => {
      return showConfirmation(message, options)
    }

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
      try {
        const confirmed = await confirmAction(t('resetDatabase.confirmReset'), {
          type: 'error',
          title: t('resetDatabase.resetDatabase'),
          confirmText: t('resetDatabase.reset'),
          cancelText: t('common.cancel'),
          confirmClass: 'bg-red-500 hover:bg-red-600',
          cancelClass: 'bg-gray-500 hover:bg-gray-600'
        })
        
        if (!confirmed) {
          return
        }
        
        loading.value = true
        error.value = ''
        success.value = ''
        
        try {
          // Use the new single endpoint that handles everything in the backend
          const response = await dataManagementService.resetDatabase()
          
          success.value = response.message || t('resetDatabase.resetSuccess')
          stats.value = response.stats || {}
        } catch (err) {
          console.error('Reset database error:', err)
          if (err.response?.data?.error) {
            error.value = `${t('resetDatabase.resetError')}: ${err.response.data.error}`
          } else if (err.response?.data?.details) {
            error.value = `${t('resetDatabase.resetError')}: ${err.response.data.details}`
          } else if (err.message) {
            error.value = `${t('resetDatabase.resetError')}: ${err.message}`
          } else {
            error.value = `${t('resetDatabase.resetError')}: ${t('common.unknownError')}`
          }
        } finally {
          loading.value = false
        }
      } catch (error) {
        console.error('Confirmation error:', error)
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
      resetDatabase,
      // Confirmations
      pendingConfirmations,
      isPageBlocked,
      confirmAction,
      handleConfirmation
    }
  }
}
</script>

<style scoped>
/* Confirmation animations */
.confirmation-enter-active,
.confirmation-leave-active {
  transition: all 0.3s ease;
}

.confirmation-enter-from {
  opacity: 0;
  transform: scale(0.9) translateY(-20px);
}

.confirmation-leave-to {
  opacity: 0;
  transform: scale(0.9) translateY(-20px);
}

.confirmation-move {
  transition: transform 0.3s ease;
}

/* Page blocking styles */
.page-blocked {
  pointer-events: none;
  user-select: none;
}
</style>
