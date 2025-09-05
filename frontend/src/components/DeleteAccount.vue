<template>
  <div>
    
    <div class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex items-center">
        <i class="fas fa-exclamation-triangle text-red-600 text-xl mr-3"></i>
        <div>
          <h3 class="text-lg font-medium text-red-800">Warning</h3>
          <p class="text-red-700 mt-1">
            This action cannot be undone. All your data, including wallets, transactions, and personal information will be permanently deleted.
          </p>
        </div>
      </div>
    </div>

    <form @submit.prevent="handleDeleteAccount" class="space-y-6">
      <!-- Password -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Current Password
        </label>
        <input
          v-model="form.password"
          type="password"
          class="form-input w-full"
          placeholder="Enter your current password"
          required
        >
      </div>

      <!-- Confirmation -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Confirmation
        </label>
        <input
          v-model="form.confirm"
          type="text"
          class="form-input w-full"
          placeholder="Type 'DELETE MY ACCOUNT' to confirm"
          required
        >
        <p class="text-xs text-gray-500 mt-1">
          Type exactly: <span class="font-mono text-red-600">DELETE MY ACCOUNT</span>
        </p>
      </div>

      <!-- Final Warning -->
      <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <div class="flex items-start">
          <i class="fas fa-info-circle text-yellow-600 text-lg mr-3 mt-0.5"></i>
          <div>
            <h4 class="text-sm font-medium text-yellow-800 mb-1">Before you delete your account:</h4>
            <ul class="text-sm text-yellow-700 space-y-1">
              <li>• Transfer or withdraw all funds from your wallets</li>
              <li>• Download any important transaction history</li>
              <li>• Cancel any pending transactions</li>
              <li>• This action is irreversible</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Delete Button -->
      <div class="pt-4">
        <button
          type="submit"
          class="inline-flex items-center justify-center w-full px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
          :disabled="loading || !isFormValid"
        >
          <i v-if="loading" class="fas fa-spinner fa-spin mr-2"></i>
          <i v-else class="fas fa-trash-alt mr-2"></i>
          {{ loading ? 'Deleting Account...' : 'Delete My Account' }}
        </button>
      </div>
    </form>

    <!-- Error/Success Messages -->
    <div v-if="error" class="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg">
      <p class="text-sm text-red-700">{{ error }}</p>
    </div>

    <div v-if="success" class="mt-4 p-3 bg-green-50 border border-green-200 rounded-lg">
      <p class="text-sm text-green-700">{{ success }}</p>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { userService } from '@/services/user'

export default {
  name: 'DeleteAccount',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    const form = ref({
      password: '',
      confirm: ''
    })
    
    const loading = ref(false)
    const error = ref('')
    const success = ref('')

    const isFormValid = computed(() => {
      return form.value.password && form.value.confirm === 'DELETE MY ACCOUNT'
    })

    const handleDeleteAccount = async () => {
      if (!isFormValid.value) {
        error.value = 'Please fill in all fields correctly'
        return
      }

      loading.value = true
      error.value = ''
      success.value = ''

      try {
        await userService.deleteAccount(form.value.password, form.value.confirm)
        success.value = 'Account deleted successfully. You will be redirected to the login page.'
        
        // Clear form
        form.value.password = ''
        form.value.confirm = ''
        
        // Logout and redirect after a delay
        setTimeout(async () => {
          await authStore.logout()
          router.push('/auth/login')
        }, 2000)
        
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to delete account'
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      loading,
      error,
      success,
      isFormValid,
      handleDeleteAccount
    }
  }
}
</script>
