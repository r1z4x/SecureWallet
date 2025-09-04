<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100">
    <!-- Main Content -->
    <div class="w-full">
      <div class="grid grid-cols-1 lg:grid-cols-2 items-stretch w-full min-h-screen gap-0">
        <!-- Left: Image + Copy -->
        <div class="hidden lg:block relative w-full h-full">
          <img 
            src="https://images.unsplash.com/photo-1517245386807-bb43f82c33c4?q=80&w=1600&auto=format&fit=crop" 
            alt="Create Account"
            class="absolute inset-0 w-full h-full object-cover"
          >
          <div class="absolute inset-0 bg-white/50 backdrop-blur-[1px]"></div>
          <div class="relative z-10 h-full flex items-center">
            <div class="p-12 max-w-xl">
              <h2 class="text-4xl font-bold text-gray-900 mb-4">{{ $t('auth.registerTitle') }}</h2>
              <p class="text-gray-700 mb-8">{{ $t('auth.registerSubtitle') }}</p>
              <ul class="space-y-3 text-gray-800">
                <li class="flex items-start"><i class="fas fa-user-shield text-primary-600 mt-1 mr-3"></i><span>{{ $t('security.tips.tip3') || 'Protect your account with strong credentials' }}</span></li>
                <li class="flex items-start"><i class="fas fa-check-circle text-primary-600 mt-1 mr-3"></i><span>{{ $t('landing.features.feature1.description') || 'Simple and secure onboarding experience' }}</span></li>
                <li class="flex items-start"><i class="fas fa-heart text-primary-600 mt-1 mr-3"></i><span>{{ $t('landing.features.feature2.description') || 'Privacy‑first approach to your data' }}</span></li>
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

            <form @submit.prevent="handleRegister" class="bg-white p-8 rounded-lg shadow-lg">
              <div class="space-y-6">
                <div>
                  <label class="form-label">{{ $t('auth.username') }}</label>
                  <input 
                    v-model="form.username" 
                    type="text" 
                    class="form-input" 
                    :placeholder="$t('auth.usernamePlaceholder')"
                    required
                  >
                </div>
                
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
                  <label class="form-label">{{ $t('auth.password') }}</label>
                  <input 
                    v-model="form.password" 
                    type="password" 
                    class="form-input" 
                    :placeholder="$t('auth.passwordPlaceholder')"
                    required
                  >
                </div>
                
                <div>
                  <label class="form-label">{{ $t('auth.confirmPassword') }}</label>
                  <input 
                    v-model="form.confirm_password" 
                    type="password" 
                    class="form-input" 
                    :placeholder="$t('auth.confirmPasswordPlaceholder')"
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
                    {{ loading ? $t('auth.creatingAccount') : $t('auth.createAccountButton') }}
                  </button>
                </div>
              </div>

              <div class="mt-6 text-center">
                <div class="text-gray-600 text-sm">
                  {{ $t('auth.alreadyHaveAccount') }} 
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
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'
import LanguageSelector from '@/components/LanguageSelector.vue'

export default {
  name: 'Register',
  components: {
    LanguageSelector
  },
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const { t } = useI18n()
    
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
        error.value = t('auth.passwordsDoNotMatch')
        return
      }
      
      loading.value = true
      error.value = ''
      
      try {
        await authStore.register(form.value)
        router.push('/auth/login')
      } catch (err) {
        error.value = err.response?.data?.detail || t('auth.registrationFailed')
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
