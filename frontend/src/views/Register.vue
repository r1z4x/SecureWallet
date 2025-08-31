<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100 flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8">
      <div class="text-center">
        <i class="fas fa-user-plus text-primary-600 text-5xl mb-4"></i>
        <h2 class="text-3xl font-bold text-gray-900">Create Account</h2>
        <p class="text-gray-600 mt-2">Join SecureWallet today</p>
      </div>

      <form @submit.prevent="handleRegister" class="bg-white p-8 rounded-lg shadow-lg">
        <div class="space-y-6">
          <div>
            <label class="form-label">Username</label>
            <input 
              v-model="form.username" 
              type="text" 
              class="form-input" 
              placeholder="Choose a username"
              required
            >
          </div>
          
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
            <label class="form-label">Password</label>
            <input 
              v-model="form.password" 
              type="password" 
              class="form-input" 
              placeholder="Create a password"
              required
            >
          </div>
          
          <div>
            <label class="form-label">Confirm Password</label>
            <input 
              v-model="form.confirm_password" 
              type="password" 
              class="form-input" 
              placeholder="Confirm your password"
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
              <i v-else class="fas fa-user-plus mr-2"></i>
              {{ loading ? 'Creating Account...' : 'Create Account' }}
            </button>
          </div>
        </div>

        <div class="mt-6 text-center">
          <div class="text-gray-600 text-sm">
            Already have an account? 
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
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'Register',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    const form = ref({
      username: '',
      email: '',
      password: '',
      confirm_password: ''
    })
    
    const loading = ref(false)
    const error = ref('')

    const handleRegister = async () => {
      if (form.value.password !== form.value.confirm_password) {
        error.value = 'Passwords do not match.'
        return
      }
      
      loading.value = true
      error.value = ''
      
      try {
        await authStore.register(form.value)
        router.push('/auth/login')
      } catch (err) {
        error.value = err.response?.data?.detail || 'Registration failed. Please try again.'
      } finally {
        loading.value = false
      }
    }

    return {
      form,
      loading,
      error,
      handleRegister
    }
  }
}
</script>
