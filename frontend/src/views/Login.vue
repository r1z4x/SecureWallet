<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100">

    <!-- Main Content -->
    <div class="w-full">
      <div class="grid grid-cols-1 lg:grid-cols-2 items-stretch w-full min-h-screen gap-0">
        <!-- Left: Image + Copy -->
        <div class="hidden lg:block relative w-full h-full">
          <img 
            src="https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?q=80&w=1600&auto=format&fit=crop" 
            alt="Secure Wallet"
            class="absolute inset-0 w-full h-full object-cover blur-sm"
          >
          <div class="absolute inset-0 bg-white/50 backdrop-blur-[1px]"></div>
          <div class="relative z-10 h-full flex items-center">
            <div class="p-12 max-w-xl">
              <h2 class="text-5xl font-bold text-gray-900 mb-6 leading-tight">Welcome Back</h2>
              <p class="text-xl text-gray-700 mb-10 leading-relaxed">Sign in to your account</p>
              <ul class="space-y-4 text-gray-800">
                <li class="flex items-start">
                  <i class="fas fa-shield-alt text-primary-600 mt-1 mr-4 text-lg"></i>
                  <span class="text-lg font-medium">{{ $t('security.tips.tip1') || 'Two‑factor authentication supported' }}</span>
                </li>
                <li class="flex items-start">
                  <i class="fas fa-lock text-primary-600 mt-1 mr-4 text-lg"></i>
                  <span class="text-lg font-medium">{{ $t('security.tips.tip2') || 'Your data is encrypted at rest and in transit' }}</span>
                </li>
                <li class="flex items-start">
                  <i class="fas fa-bolt text-primary-600 mt-1 mr-4 text-lg"></i>
                  <span class="text-lg font-medium">{{ $t('landing.features.feature3.description') || 'Fast and reliable access to your wallet' }}</span>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Right: Form -->
        <div class="flex items-center justify-center w-full h-full p-4 lg:p-8 bg-transparent">
          <div class="max-w-md w-full space-y-8">

            <div class="mb-6">
              <div class="flex justify-between items-center">
                <router-link 
                   to="/" 
                   class="flex items-center text-primary-600 hover:text-primary-800 transition-colors"
                >
                  <i class="fas fa-arrow-left mr-2"></i>
                  <span>{{ $t('auth.backToHome') }}</span>
                </router-link>
                <LanguageSelector />
              </div>
            </div>

            <div class="text-center lg:text-left flex">
              <img src="@/assets/logo.svg" alt="Logo" class="max-w-64 w-full mx-auto lg:mx-0 mr-3">
              <sup class="text-red-500 text-xs ml-2 font-bold">{{ $t('common.vulnerable') }}</sup>
            </div>

            <form @submit.prevent="handleLogin" class="bg-white p-8 rounded-lg shadow-lg">
              <div class="space-y-6">
                <div>
                  <label class="form-label">{{ $t('auth.usernameOrEmail') }}</label>
                  <input 
                    v-model="form.username" 
                    type="text" 
                    class="form-input" 
                    :placeholder="$t('auth.usernameOrEmailPlaceholder')"
                    required
                    :disabled="requires2FA"
                  >
                </div>
                
                <div>
                  <label class="form-label">{{ $t('auth.password') }}</label>
                  <input 
                    v-model="form.password" 
                    type="password" 
                    class="form-input" 
                    :placeholder="$t('auth.passwordPlaceholder')"
                    required
                    :disabled="requires2FA"
                  >
                </div>

                <!-- 2FA Code Input -->
                <div v-if="requires2FA">
                  <label class="form-label">{{ $t('auth.twoFactorCode') }}</label>
                  <input 
                    v-model="form.code2FA" 
                    type="text" 
                    class="form-input" 
                    :placeholder="$t('auth.twoFactorCodePlaceholder')"
                    maxlength="6"
                    pattern="[0-9]{6}"
                    required
                  >
                  <p class="text-xs text-gray-500 mt-1">
                    {{ $t('auth.twoFactorCodeHelp') }}
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
                    {{ loading ? $t('auth.signingIn') : requires2FA ? $t('auth.verify2FA') : $t('auth.signInButton') }}
                  </button>
                </div>
              </div>

              <div class="mt-6 text-center space-y-2">
                <router-link 
                  to="/auth/password-reset" 
                  class="text-primary-600 hover:text-primary-800 text-sm block"
                >
                  {{ $t('auth.forgotPassword') }}
                </router-link>
                <div class="text-gray-600 text-sm">
                  {{ $t('auth.dontHaveAccount') }} 
                  <router-link 
                    to="/auth/register" 
                    class="text-primary-600 hover:text-primary-800"
                  >
                    {{ $t('auth.signUpLink') }}
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
      </div>
    </div>
  </div>
</template>

<script>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import LanguageSelector from '@/components/LanguageSelector.vue'

export default {
  name: 'Login',
  components: {
    LanguageSelector
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const { t } = useI18n()
    
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
          // Clear form data after successful login
          form.value = { username: '', password: '', code2FA: '' }
          requires2FA.value = false
          userId2FA.value = null
          
          // Wait for user data to be fully loaded in store
          await new Promise(resolve => setTimeout(resolve, 200))
          
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
          
          // Clear form data after successful login
          form.value = { username: '', password: '', code2FA: '' }
          
          // Wait for user data to be fully loaded in store
          await new Promise(resolve => setTimeout(resolve, 200))
          
          router.push('/dashboard')
        }
      } catch (err) {
        error.value = err.response?.data?.error || err.message || t('auth.invalidCredentials')
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
