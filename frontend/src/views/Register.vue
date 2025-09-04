<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100">
    <!-- Main Content -->
    <div class="w-full">
      <div class="grid grid-cols-1 lg:grid-cols-2 items-stretch w-full min-h-screen gap-0">
        <!-- Left: Enhanced Design Panel -->
        <div class="hidden lg:block relative w-full h-full overflow-hidden">
          <!-- Animated Background with Working Animations -->
          <div class="absolute inset-0 bg-gradient-to-br from-blue-200 via-blue-300 to-indigo-300">
            <!-- Dynamic Gradient Flow -->
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-blue-400/20 to-transparent animate-gradient-flow"></div>
            
            <!-- Particle System - Data Security Flow -->
            <div class="particle-container">
              <div class="particle particle-1"></div>
              <div class="particle particle-2"></div>
              <div class="particle particle-3"></div>
              <div class="particle particle-4"></div>
              <div class="particle particle-5"></div>
              <div class="particle particle-6"></div>
              <div class="particle particle-7"></div>
              <div class="particle particle-8"></div>
            </div>
            
            <!-- Geometric Security Patterns -->
            <div class="geometric-container">
              <div class="geometric-shape shape-1"></div>
              <div class="geometric-shape shape-2"></div>
              <div class="geometric-shape shape-3"></div>
              <div class="geometric-shape shape-4"></div>
            </div>
            
            <!-- Data Flow Lines -->
            <div class="data-flow-container">
              <div class="data-line line-1"></div>
              <div class="data-line line-2"></div>
              <div class="data-line line-3"></div>
            </div>
          </div>
          
          <!-- Content Overlay -->
          <div class="relative z-10 h-full flex items-center">
            <div class="p-16 max-w-2xl">
              <!-- Enhanced Logo Section -->
              <div class="text-center lg:text-left mb-12">
                <div class="flex items-center justify-center lg:justify-start mb-6">
 
                  <div>
                    <img src="@/assets/logo.svg" alt="Secure Wallet Logo" class="w-full mr-4" style="width: 500px;">
                    <div class="flex items-center mt-2">
                      <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold bg-red-500/30 text-red-800 border border-red-400/40 backdrop-blur-sm">
                        <i class="fas fa-exclamation-triangle mr-1"></i>
                        {{ $t('common.vulnerable') }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Enhanced Features List -->
              <div class="space-y-6">
                <div class="group flex items-start p-4 rounded-xl bg-white/20 backdrop-blur-md border border-white/30 hover:bg-white/30 hover:border-white/50 transition-all duration-300 shadow-sm hover:shadow-md">
                  <div class="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl flex items-center justify-center mr-4 transition-all duration-300 shadow-lg group-hover:shadow-xl">
                    <i class="fas fa-user-shield text-white text-lg"></i>
                  </div>
                  <div>
                    <h3 class="text-gray-800 font-semibold text-lg mb-1 drop-shadow-sm">
                      {{ $t('security.tips.tip3') || 'Protect your account with strong credentials' }}
                    </h3>
                    <p class="text-gray-700 text-sm leading-relaxed drop-shadow-sm">
                      Description
                    </p>
                  </div>
                </div>

                <div class="group flex items-start p-4 rounded-xl bg-white/20 backdrop-blur-md border border-white/30 hover:bg-white/30 hover:border-white/50 transition-all duration-300 shadow-sm hover:shadow-md">
                  <div class="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-green-500 to-green-600 rounded-xl flex items-center justify-center mr-4 transition-all duration-300 shadow-lg group-hover:shadow-xl">
                    <i class="fas fa-check-circle text-white text-lg"></i>
                  </div>
                  <div>
                    <h3 class="text-gray-800 font-semibold text-lg mb-1 drop-shadow-sm">
                      {{ $t('landing.features.feature1.title') || 'Simple and secure onboarding' }}
                    </h3>
                    <p class="text-gray-700 text-sm leading-relaxed drop-shadow-sm">
                      Description
                    </p>
                  </div>
                </div>

                <div class="group flex items-start p-4 rounded-xl bg-white/20 backdrop-blur-md border border-white/30 hover:bg-white/30 hover:border-white/50 transition-all duration-300 shadow-sm hover:shadow-md">
                  <div class="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl flex items-center justify-center mr-4 transition-all duration-300 shadow-lg group-hover:shadow-xl">
                    <i class="fas fa-heart text-white text-lg"></i>
                  </div>
                  <div>
                    <h3 class="text-gray-800 font-semibold text-lg mb-1 drop-shadow-sm">
                      {{ $t('landing.features.feature2.title') || 'Privacy-first approach' }}
                    </h3>
                    <p class="text-gray-700 text-sm leading-relaxed drop-shadow-sm">
                      Description
                    </p>
                  </div>
                </div>
              </div>

              <!-- Trust Indicators -->
              <div class="mt-12 pt-8 border-t border-white/30">
                <div class="flex items-center justify-center lg:justify-start space-x-6 text-gray-700 text-sm drop-shadow-sm">
                  <div class="flex items-center">
                    <i class="fas fa-lock mr-2 text-green-600"></i>
                    <span>256-bit AES</span>
                  </div>
                  <div class="flex items-center">
                    <i class="fas fa-shield-alt mr-2 text-blue-600"></i>
                    <span>ISO 27001</span>
                  </div>
                  <div class="flex items-center">
                    <i class="fas fa-check-double mr-2 text-purple-600"></i>
                    <span>GDPR Compliant</span>
                  </div>
                </div>
              </div>
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

            <div class="text-center lg:text-left flex flex-col">
              <h2 class="text-4xl font-bold text-gray-900 mb-4">{{ $t('auth.registerTitle') }}</h2>
              <p class="text-gray-700" style="font-weight: 200;">{{ $t('auth.registerSubtitle') }}</p>
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

<style scoped>
/* Secure Wallet - Working Background Animations */

/* Dynamic Gradient Flow */
@keyframes gradient-flow {
  0% { transform: translateX(-100%); opacity: 0; }
  50% { opacity: 0.3; }
  100% { transform: translateX(100%); opacity: 0; }
}

.animate-gradient-flow {
  animation: gradient-flow 8s ease-in-out infinite;
}

/* Particle System - Data Security Flow */
.particle-container {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: linear-gradient(45deg, #60a5fa, #3b82f6, #1d4ed8);
  border-radius: 50%;
  box-shadow: 0 0 10px rgba(59, 130, 246, 0.6);
  animation: particle-float 12s linear infinite;
}

.particle-1 { top: 20%; left: 10%; animation-delay: 0s; }
.particle-2 { top: 35%; left: 15%; animation-delay: 1.5s; }
.particle-3 { top: 50%; left: 20%; animation-delay: 3s; }
.particle-4 { top: 65%; left: 25%; animation-delay: 4.5s; }
.particle-5 { top: 80%; left: 30%; animation-delay: 6s; }
.particle-6 { top: 25%; left: 35%; animation-delay: 7.5s; }
.particle-7 { top: 40%; left: 40%; animation-delay: 9s; }
.particle-8 { top: 55%; left: 45%; animation-delay: 10.5s; }

@keyframes particle-float {
  0% {
    transform: translateY(0) translateX(0) scale(0);
    opacity: 0;
  }
  10% {
    transform: translateY(-20px) translateX(10px) scale(1);
    opacity: 1;
  }
  90% {
    transform: translateY(-100px) translateX(50px) scale(1);
    opacity: 1;
  }
  100% {
    transform: translateY(-120px) translateX(60px) scale(0);
    opacity: 0;
  }
}

/* Geometric Security Patterns */
.geometric-container {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.geometric-shape {
  position: absolute;
  border: 2px solid rgba(59, 130, 246, 0.3);
  border-radius: 50%;
  animation: geometric-pulse 6s ease-in-out infinite;
}

.shape-1 {
  top: 15%;
  right: 20%;
  width: 60px;
  height: 60px;
  animation-delay: 0s;
}

.shape-2 {
  top: 45%;
  left: 15%;
  width: 40px;
  height: 40px;
  animation-delay: 2s;
}

.shape-3 {
  bottom: 25%;
  right: 15%;
  width: 50px;
  height: 50px;
  animation-delay: 4s;
}

.shape-4 {
  top: 70%;
  left: 25%;
  width: 30px;
  height: 30px;
  animation-delay: 1s;
}

@keyframes geometric-pulse {
  0%, 100% {
    transform: scale(1) rotate(0deg);
    opacity: 0.3;
    border-color: rgba(59, 130, 246, 0.3);
  }
  50% {
    transform: scale(1.2) rotate(180deg);
    opacity: 0.6;
    border-color: rgba(59, 130, 246, 0.6);
  }
}

/* Data Flow Lines */
.data-flow-container {
  position: absolute;
  inset: 0;
  pointer-events: none;
}

.data-line {
  position: absolute;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgba(59, 130, 246, 0.4), transparent);
  animation: data-flow 10s linear infinite;
}

.line-1 {
  top: 30%;
  left: 0;
  width: 100%;
  animation-delay: 0s;
}

.line-2 {
  top: 60%;
  left: 0;
  width: 100%;
  animation-delay: 3s;
}

.line-3 {
  top: 90%;
  left: 0;
  width: 100%;
  animation-delay: 6s;
}

@keyframes data-flow {
  0% {
    transform: translateX(-100%);
    opacity: 0;
  }
  50% {
    opacity: 1;
  }
  100% {
    transform: translateX(100%);
    opacity: 0;
  }
}

/* Hover Effects for Enhanced Interaction */
.particle-container:hover .particle {
  animation-duration: 6s;
  box-shadow: 0 0 20px rgba(59, 130, 246, 0.8);
}

.geometric-container:hover .geometric-shape {
  animation-duration: 3s;
  border-color: rgba(59, 130, 246, 0.8);
}

.data-flow-container:hover .data-line {
  animation-duration: 5s;
  background: linear-gradient(90deg, transparent, rgba(59, 130, 246, 0.8), transparent);
}
</style>
