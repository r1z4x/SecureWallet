<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <i class="fas fa-key text-primary-600 text-5xl mb-4"></i>
        <h2 class="text-3xl font-bold text-gray-900">Reset Password</h2>
        <p class="text-gray-600 mt-2">Enter your email to receive reset instructions</p>
      </div>

      <form @submit.prevent="handlePasswordReset" class="bg-white p-8 rounded-lg shadow-lg">
        <div class="space-y-6">
          <div>
            <label class="form-label">Email</label>
            <input 
              v-model="form.email" 
              type="email" 
              class="form-input" 
              placeholder="Enter your email"
              required
            >
          </div>

          <div>
            <button 
              type="submit" 
              class="btn-primary w-full"
              :disabled="loading"
            >
              <i v-if="loading" class="fas fa-spinner fa-spin mr-2"></i>
              <i v-else class="fas fa-paper-plane mr-2"></i>
              {{ loading ? 'Sending...' : 'Send Reset Link' }}
            </button>
          </div>
        </div>

        <div class="mt-6 text-center">
          <div class="text-gray-600 text-sm">
            Remember your password? 
            <router-link 
              to="/auth/login" 
              class="text-primary-600 hover:text-primary-800"
            >
              Sign in
            </router-link>
          </div>
        </div>
      </form>

      <!-- Error Alert -->
      <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        {{ error }}
      </div>

      <!-- Success Alert -->
      <div v-if="success" class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">
        {{ success }}
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'PasswordReset',
  setup() {
    const authStore = useAuthStore()
    
    const form = ref({
      email: ''
    })
    
    const loading = ref(false)
    const error = ref('')
    const success = ref('')

    const handlePasswordReset = async () => {
      loading.value = true
      error.value = ''
      success.value = ''
      
      try {
        await authStore.requestPasswordReset(form.value.email)
        success.value = 'Password reset email sent! Please check your inbox.'
        form.value.email = ''
      } catch (err) {
        error.value = err.response?.data?.detail || 'Failed to send reset email. Please try again.'
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      loading,
      error,
      success,
      handlePasswordReset
    }
  }
}
</script>
