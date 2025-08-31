<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <i class="fas fa-wallet text-primary-600 text-5xl mb-4"></i>
        <h2 class="text-3xl font-bold text-gray-900">Welcome Back</h2>
        <p class="text-gray-600 mt-2">Sign in to your account</p>
      </div>

      <form @submit.prevent="handleLogin" class="bg-white p-8 rounded-lg shadow-lg">
        <div class="space-y-6">
          <div>
            <label class="form-label">Username or Email</label>
            <input 
              v-model="form.username" 
              type="text" 
              class="form-input" 
              placeholder="Enter your username or email"
              required
              :disabled="requires2FA"
            >
          </div>
          
          <div>
            <label class="form-label">Password</label>
            <input 
              v-model="form.password" 
              type="password" 
              class="form-input" 
              placeholder="Enter your password"
              required
              :disabled="requires2FA"
            >
          </div>

          <!-- 2FA Code Input -->
          <div v-if="requires2FA">
            <label class="form-label">2FA Code</label>
            <input 
              v-model="form.code2FA" 
              type="text" 
              class="form-input" 
              placeholder="Enter 6-digit code"
              maxlength="6"
              pattern="[0-9]{6}"
              required
            >
            <p class="text-xs text-gray-500 mt-1">
              Enter the 6-digit code from your authenticator app
            </p>
          </div>

          <div>
            <button 
              type="submit" 
              class="btn-primary w-full"
              :disabled="loading"
            >
              <i v-if="loading" class="fas fa-spinner fa-spin mr-2"></i>
              <i v-else-if="requires2FA" class="fas fa-shield-alt mr-2"></i>
              <i v-else class="fas fa-sign-in-alt mr-2"></i>
              {{ loading ? 'Signing In...' : requires2FA ? 'Verify 2FA' : 'Sign In' }}
            </button>
          </div>
        </div>

        <div class="mt-6 text-center space-y-2">
          <router-link 
            to="/auth/password-reset" 
            class="text-primary-600 hover:text-primary-800 text-sm block"
          >
            Forgot your password?
          </router-link>
          <div class="text-gray-600 text-sm">
            Don't have an account? 
            <router-link 
              to="/auth/register" 
              class="text-primary-600 hover:text-primary-800"
            >
              Sign up
            </router-link>
          </div>
        </div>
      </form>

      <!-- Error Alert -->
      <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">
        {{ error }}
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'Login',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    const form = ref({
      username: '',
      password: '',
      code2FA: ''
    })
    
    const loading = ref(false)
    const error = ref('')
    const requires2FA = ref(false)
    const userId2FA = ref(null)

    const handleLogin = async () => {
      loading.value = true
      error.value = ''
      
      try {
        if (requires2FA.value) {
          // 2FA verification
          await authStore.login2FA(userId2FA.value, form.value.code2FA)
          router.push('/dashboard')
        } else {
          // Initial login
          const response = await authStore.login(form.value)
          
          // Check if 2FA is required
          if (response && response.requires_2fa) {
            requires2FA.value = true
            userId2FA.value = response.user_id
            form.value.code2FA = ''
            return
          }
          
          router.push('/dashboard')
        }
      } catch (err) {
        error.value = err.response?.data?.error || err.message || 'Login failed'
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      loading,
      error,
      requires2FA,
      handleLogin
    }
  }
}
</script>
