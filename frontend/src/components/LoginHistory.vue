<template>
  <div>
    
    <div v-if="loading" class="text-center py-4">
      <i class="fas fa-spinner fa-spin text-2xl text-primary-600"></i>
      <p class="text-gray-600 mt-2">Loading...</p>
    </div>

    <div v-else-if="error" class="text-center py-4">
      <i class="fas fa-exclamation-triangle text-2xl text-red-600"></i>
      <p class="text-red-600 mt-2">{{ error }}</p>
    </div>

    <div v-else-if="loginHistory.length === 0" class="text-center py-8">
      <i class="fas fa-history text-4xl text-gray-300"></i>
      <p class="text-gray-500 mt-2">No login history found</p>
    </div>

    <div v-else class="space-y-4">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-medium text-gray-900">Recent Login Attempts</h3>
        <button
          @click="loadFullHistory"
          class="text-primary-600 hover:text-primary-800 text-sm font-medium"
        >
          View All
        </button>
      </div>

      <div class="space-y-3">
        <div
          v-for="login in loginHistory"
          :key="login.id"
          class="border border-gray-200 rounded-lg p-4 hover:bg-gray-50"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <div
                :class="{
                  'text-green-600': login.status === 'success',
                  'text-red-600': login.status === 'failed',
                  'text-yellow-600': login.status === 'blocked'
                }"
              >
                <i
                  :class="{
                    'fas fa-check-circle': login.status === 'success',
                    'fas fa-times-circle': login.status === 'failed',
                    'fas fa-exclamation-triangle': login.status === 'blocked'
                  }"
                  class="text-xl"
                ></i>
              </div>
              <div>
                <p class="font-medium text-gray-900">
                  {{ login.status === 'success' ? 'Successful Login' : 'Failed Login' }}
                </p>
                <p class="text-sm text-gray-500">
                  {{ formatDate(login.created_at) }}
                </p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-sm font-medium text-gray-900">{{ login.ip_address }}</p>
              <p class="text-xs text-gray-500">{{ login.location }}</p>
            </div>
          </div>
          
          <div class="mt-3 pt-3 border-t border-gray-100">
            <p class="text-xs text-gray-500 truncate">
              {{ login.user_agent }}
            </p>
          </div>
        </div>
      </div>

      <div v-if="showLoadMore" class="text-center mt-6">
        <button
          @click="loadMore"
          :disabled="loadingMore"
          class="btn-secondary"
        >
          <i v-if="loadingMore" class="fas fa-spinner fa-spin mr-2"></i>
          <i v-else class="fas fa-plus mr-2"></i>
          {{ loadingMore ? 'Loading...' : 'Load More' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { loginHistoryService } from '@/services/loginHistory'

export default {
  name: 'LoginHistory',
  setup() {
    const loading = ref(false)
    const loadingMore = ref(false)
    const loginHistory = ref([])
    const error = ref('')
    const showLoadMore = ref(false)
    const currentLimit = ref(10)

    const loadLoginHistory = async (limit = 10) => {
      loading.value = true
      error.value = ''
      try {
        const response = await loginHistoryService.getRecentLoginHistory(limit)
        loginHistory.value = response.recent_login_history || []
        showLoadMore.value = response.total > limit
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load login history'
      } finally {
        loading.value = false
      }
    }

    const loadFullHistory = async () => {
      loading.value = true
      error.value = ''
      try {
        const response = await loginHistoryService.getLoginHistory(50)
        loginHistory.value = response.login_history || []
        showLoadMore.value = false
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load full login history'
      } finally {
        loading.value = false
      }
    }

    const loadMore = async () => {
      loadingMore.value = true
      try {
        currentLimit.value += 10
        const response = await loginHistoryService.getRecentLoginHistory(currentLimit.value)
        loginHistory.value = response.recent_login_history || []
        showLoadMore.value = response.total > currentLimit.value
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load more login history'
      } finally {
        loadingMore.value = false
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
      loadLoginHistory()
    })

    return {
      loading,
      loadingMore,
      loginHistory,
      error,
      showLoadMore,
      loadLoginHistory,
      loadFullHistory,
      loadMore,
      formatDate
    }
  }
}
</script>
