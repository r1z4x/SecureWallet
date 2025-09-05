<template>
  <nav class="bg-white shadow-lg">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex">
          <div class="flex-shrink-0 flex items-center">
            <img src="@/assets/logo.svg" alt="Logo" class="h-12 w-full max-w-64 mr-3">
            <sup class="text-red-500 text-xs">{{ $t('common.vulnerable') }}</sup>
          </div>
          
          <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
            <router-link
              to="/dashboard"
              class="nav-link"
              :class="{ 'active': $route.path === '/dashboard' }"
            >
              <i class="fas fa-tachometer-alt mr-2"></i>
              {{ $t('nav.dashboard') }}
            </router-link>
            
            <router-link
              to="/wallet"
              class="nav-link"
              :class="{ 'active': $route.path === '/wallet' }"
            >
              <i class="fas fa-wallet mr-2"></i>
              {{ $t('nav.wallet') }}
            </router-link>
            
            <router-link
              to="/transactions"
              class="nav-link"
              :class="{ 'active': $route.path === '/transactions' }"
            >
              <i class="fas fa-exchange-alt mr-2"></i>
              {{ $t('nav.transactions') }}
            </router-link>
            
            <router-link
              to="/transfer"
              class="nav-link"
              :class="{ 'active': $route.path === '/transfer' }"
            >
              <i class="fas fa-paper-plane mr-2"></i>
              {{ $t('nav.transfer') }}
            </router-link>
            

            
            <router-link
              v-if="isAdmin"
              to="/admin"
              class="nav-link"
              :class="{ 'active': $route.path === '/admin' }"
            >
              <i class="fas fa-shield-alt mr-2"></i>
              {{ $t('nav.admin') }}
            </router-link>
            

          </div>
        </div>
        
        <div class="flex items-center space-x-4">
          <!-- Language Selector -->
          <div class="relative">
            <button
              @click="toggleLanguageDropdown"
              class="flex items-center justify-center w-8 h-8 rounded-full hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors"
              :title="getLanguageName(currentLocale)"
            >
              <img
                :src="getFlagIcon(currentLocale)"
                :alt="getLanguageName(currentLocale)"
                class="w-6 h-6 rounded-full object-cover shadow-sm"
              />
            </button>

            <!-- Language Dropdown -->
            <div
              v-if="languageDropdownOpen"
              class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-50 border border-gray-200"
            >
              <div class="py-1">
                <button
                  v-for="locale in availableLocales"
                  :key="locale.code"
                  @click="changeLanguage(locale.code)"
                  class="flex items-center space-x-3 w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900"
                  :class="{ 'bg-indigo-50 text-indigo-700': currentLocale === locale.code }"
                >
                  <span class="w-5 h-5 rounded-full overflow-hidden">
                    <img
                      :src="locale.flag"
                      :alt="locale.name"
                      class="w-full h-full object-cover"
                    />
                  </span>
                  <span>{{ locale.name }}</span>
                  <span v-if="currentLocale === locale.code" class="text-indigo-600">
                    <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                      <path
                        fill-rule="evenodd"
                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                        clip-rule="evenodd"
                      />
                    </svg>
                  </span>
                </button>
              </div>
            </div>
          </div>

          <!-- Desktop User Dropdown -->
          <div class="hidden sm:flex sm:items-center">
            <div class="flex-shrink-0 relative">
              <div class="relative">
                <button
                  @click="toggleUserDropdown"
                  class="flex items-center text-sm text-gray-600 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 rounded-md p-2"
                >
                  <span class="mr-2">{{ $t('common.welcome') }}, {{ user?.username || user?.email || 'User' }}</span>
                  <i class="fas fa-chevron-down text-xs"></i>
                </button>
                
                <!-- Dropdown Menu -->
                <div
                  v-if="userDropdownOpen"
                  class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg py-1 z-50 border border-gray-200"
                >
                  <router-link
                    to="/profile"
                    class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    @click="userDropdownOpen = false"
                  >
                    <i class="fas fa-user mr-3"></i>
                    {{ $t('nav.profile') }}
                  </router-link>
                  
                  <router-link
                    to="/support"
                    class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                    @click="userDropdownOpen = false"
                  >
                    <i class="fas fa-headset mr-3"></i>
                    {{ $t('nav.support') }}
                  </router-link>
                  
                  <div class="border-t border-gray-200 my-1"></div>
                  
                  <button
                    @click="handleLogout"
                    class="flex items-center w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                  >
                    <i class="fas fa-sign-out-alt mr-3"></i>
                    {{ $t('nav.logout') }}
                  </button>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Mobile menu button -->
          <div class="sm:hidden">
            <button
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500"
            >
              <i class="fas fa-bars" v-if="!mobileMenuOpen"></i>
              <i class="fas fa-times" v-else></i>
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Mobile menu -->
    <div class="sm:hidden" v-show="mobileMenuOpen">
      <div class="pt-2 pb-3 space-y-1">
        <!-- Mobile User Dropdown -->
        <div class="px-4 py-2 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-gray-900">
              {{ $t('common.welcome') }}, {{ isUserLoaded ? (user?.username || user?.email || 'User') : 'Loading...' }}
            </span>
            <button
              @click="toggleUserDropdown"
              class="text-gray-400 hover:text-gray-600"
            >
              <i class="fas fa-chevron-down"></i>
            </button>
          </div>
          
          <div v-if="userDropdownOpen" class="mt-2 space-y-1">
            <router-link
              to="/profile"
              class="block px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded-md"
              @click="mobileMenuOpen = false; userDropdownOpen = false"
            >
              <i class="fas fa-user mr-2"></i>
              {{ $t('nav.profile') }}
            </router-link>
            
            <router-link
              to="/support"
              class="block px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded-md"
              @click="mobileMenuOpen = false; userDropdownOpen = false"
            >
              <i class="fas fa-headset mr-2"></i>
              {{ $t('nav.support') }}
            </router-link>
            
            <button
              @click="handleLogout"
              class="block w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded-md"
            >
              <i class="fas fa-sign-out-alt mr-2"></i>
              {{ $t('nav.logout') }}
            </button>
          </div>
        </div>
        <router-link
          to="/dashboard"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/dashboard' }"
        >
          <i class="fas fa-tachometer-alt mr-2"></i>
          {{ $t('nav.dashboard') }}
        </router-link>
        
        <router-link
          to="/wallet"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/wallet' }"
        >
          <i class="fas fa-wallet mr-2"></i>
          {{ $t('nav.wallet') }}
        </router-link>
        
        <router-link
          to="/transactions"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/transactions' }"
        >
          <i class="fas fa-exchange-alt mr-2"></i>
          {{ $t('nav.transactions') }}
        </router-link>
        
        <router-link
          to="/transfer"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/transfer' }"
        >
          <i class="fas fa-paper-plane mr-2"></i>
          {{ $t('nav.transfer') }}
        </router-link>
        

        
        <router-link
          v-if="isAdmin"
          to="/admin"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/admin' }"
        >
          <i class="fas fa-shield-alt mr-2"></i>
          {{ $t('nav.admin') }}
        </router-link>
        

      </div>
    </div>
  </nav>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useI18n } from 'vue-i18n'

export default {
  name: 'Navigation',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const { locale } = useI18n()
    const mobileMenuOpen = ref(false)
    const userDropdownOpen = ref(false)
    const languageDropdownOpen = ref(false)
    
    const user = computed(() => authStore.user)
    const isAdmin = computed(() => authStore.isAdmin)
    const currentLocale = computed(() => locale.value)
    const isUserLoaded = computed(() => authStore.isUserLoaded)
    
    const availableLocales = [
      {
        code: 'en',
        name: 'English',
        flag: 'https://flagcdn.com/w40/gb.png'
      },
      {
        code: 'tr',
        name: 'Türkçe',
        flag: 'https://flagcdn.com/w40/tr.png'
      },
      {
        code: 'es',
        name: 'Español',
        flag: 'https://flagcdn.com/w40/es.png'
      }
    ]
    
    const changeLanguage = (newLocale) => {
      locale.value = newLocale
      localStorage.setItem('locale', newLocale)
      languageDropdownOpen.value = false
    }
    
    const toggleLanguageDropdown = () => {
      languageDropdownOpen.value = !languageDropdownOpen.value
      userDropdownOpen.value = false
    }
    
    const toggleUserDropdown = () => {
      userDropdownOpen.value = !userDropdownOpen.value
      languageDropdownOpen.value = false
    }
    
    const getFlagIcon = (code) => {
      const locale = availableLocales.find(l => l.code === code)
      return locale ? locale.flag : availableLocales[0].flag
    }

    const getLanguageName = (code) => {
      const locale = availableLocales.find(l => l.code === code)
      return locale ? locale.name : availableLocales[0].name
    }
    
    const handleLogout = async () => {
      await authStore.logout()
      router.push('/auth/login')
    }
    
    // Close dropdown when clicking outside
    const closeDropdown = (event) => {
      if (!event.target.closest('.relative')) {
        userDropdownOpen.value = false
        languageDropdownOpen.value = false
      }
    }
    
    // Add click outside listener
    onMounted(() => {
      document.addEventListener('click', closeDropdown)
    })
    
    onUnmounted(() => {
      document.removeEventListener('click', closeDropdown)
    })
    
    return {
      user,
      isAdmin,
      mobileMenuOpen,
      userDropdownOpen,
      languageDropdownOpen,
      currentLocale,
      availableLocales,
      changeLanguage,
      toggleLanguageDropdown,
      getFlagIcon,
      getLanguageName,
      toggleUserDropdown,
      handleLogout,
      isUserLoaded
    }
  }
}
</script>

<style scoped>
.nav-link {
  @apply inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium text-gray-500 hover:text-gray-700 hover:border-gray-300;
}

.nav-link.active {
  @apply border-primary-500 text-primary-600;
}

.mobile-nav-link {
  @apply block pl-3 pr-4 py-2 border-l-4 text-base font-medium text-gray-500 hover:text-gray-700 hover:bg-gray-50 hover:border-gray-300;
}

.mobile-nav-link.active {
  @apply border-primary-500 text-primary-600 bg-primary-50;
}

.btn-secondary {
  @apply inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500;
}
</style>
