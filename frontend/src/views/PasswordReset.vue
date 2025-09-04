<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100">
    <!-- Main Content -->
    <div class="w-full">
      <div class="grid grid-cols-1 lg:grid-cols-2 items-stretch w-full min-h-screen gap-0">
        <!-- Left: Image + Copy -->
        <div class="hidden lg:block relative w-full h-full">
          <img 
            src="https://images.unsplash.com/photo-1553729459-efe14ef6055d?q=80&w=1600&auto=format&fit=crop" 
            alt="Password Reset"
            class="absolute inset-0 w-full h-full object-cover"
          >
          <div class="absolute inset-0 bg-white/50 backdrop-blur-[1px]"></div>
          <div class="relative z-10 h-full flex items-center">
            <div class="p-12 max-w-xl">
              <h2 class="text-4xl font-bold text-gray-900 mb-4">{{ $t('auth.resetPasswordTitle') }}</h2>
              <p class="text-gray-700 mb-8">{{ $t('auth.resetPasswordSubtitle') }}</p>
              <ul class="space-y-3 text-gray-800">
                <li class="flex items-start"><i class="fas fa-shield-alt text-primary-600 mt-1 mr-3"></i><span>{{ $t('security.tips.tip1') || 'Secure password reset process' }}</span></li>
                <li class="flex items-start"><i class="fas fa-envelope text-primary-600 mt-1 mr-3"></i><span>{{ $t('security.tips.tip2') || 'Reset link sent to your email' }}</span></li>
                <li class="flex items-start"><i class="fas fa-clock text-primary-600 mt-1 mr-3"></i><span>{{ $t('security.tips.tip3') || 'Quick and secure account recovery' }}</span></li>
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

            <form @submit.prevent="handlePasswordReset" class="bg-white p-8 rounded-lg shadow-lg">
              <div class="space-y-6">
                <div>
                  <label class="form-label">{{ $t('auth.email') }}</label>
                  <input 
                    v-model="form.email" 
                    type="email" 
                    class="form-input" 
                    :placeholder="$t('auth.emailPlaceholder')"
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
                    {{ loading ? $t('auth.sending') : $t('auth.sendResetLink') }}
                  </button>
                </div>
              </div>

              <div class="mt-6 text-center">
                <div class="text-gray-600 text-sm">
                  {{ $t('auth.rememberPassword') }} 
                  <router-link 
                    to="/auth/login" 
                    class="text-primary-600 hover:text-primary-800"
                  >
                    {{ $t('auth.signInLink') }}
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
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import LanguageSelector from '@/components/LanguageSelector.vue'

export default {
  name: 'PasswordReset',
  components: {
    LanguageSelector
  },
  setup() {
    const authStore = useAuthStore()
    const { t } = useI18n()
    
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
        success.value = t('auth.passwordResetSent')
        form.value.email = ''
      } catch (err) {
        error.value = err.response?.data?.error || err.message || t('auth.passwordResetFailed')
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
