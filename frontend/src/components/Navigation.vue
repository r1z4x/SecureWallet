<template>
  <nav class="bg-white shadow-lg">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between h-16">
        <div class="flex">
          <div class="flex-shrink-0 flex items-center">
            <i class="fas fa-wallet text-primary-600 text-2xl mr-2"></i>
            <span class="text-xl font-bold text-gray-900">SecureWallet</span>
          </div>
          
          <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
            <router-link
              to="/dashboard"
              class="nav-link"
              :class="{ 'active': $route.path === '/dashboard' }"
            >
              <i class="fas fa-tachometer-alt mr-2"></i>
              Dashboard
            </router-link>
            
            <router-link
              to="/wallet"
              class="nav-link"
              :class="{ 'active': $route.path === '/wallet' }"
            >
              <i class="fas fa-wallet mr-2"></i>
              Wallet
            </router-link>
            
            <router-link
              to="/transactions"
              class="nav-link"
              :class="{ 'active': $route.path === '/transactions' }"
            >
              <i class="fas fa-exchange-alt mr-2"></i>
              Transactions
            </router-link>
            
            <router-link
              to="/transfer"
              class="nav-link"
              :class="{ 'active': $route.path === '/transfer' }"
            >
              <i class="fas fa-paper-plane mr-2"></i>
              Transfer
            </router-link>
            
            <router-link
              to="/profile"
              class="nav-link"
              :class="{ 'active': $route.path === '/profile' }"
            >
              <i class="fas fa-user mr-2"></i>
              Profile
            </router-link>
            
            <router-link
              to="/support"
              class="nav-link"
              :class="{ 'active': $route.path === '/support' }"
            >
              <i class="fas fa-headset mr-2"></i>
              Support
            </router-link>
            
            <router-link
              v-if="isAdmin"
              to="/admin"
              class="nav-link"
              :class="{ 'active': $route.path === '/admin' }"
            >
              <i class="fas fa-shield-alt mr-2"></i>
              Admin
            </router-link>
          </div>
        </div>
        
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <span class="text-sm text-gray-600 mr-4">
              Welcome, {{ user?.username || user?.email || 'User' }}
            </span>
            <button
              @click="handleLogout"
              class="btn-secondary"
            >
              <i class="fas fa-sign-out-alt mr-2"></i>
              Logout
            </button>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Mobile menu -->
    <div class="sm:hidden" v-show="mobileMenuOpen">
      <div class="pt-2 pb-3 space-y-1">
        <router-link
          to="/dashboard"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/dashboard' }"
        >
          <i class="fas fa-tachometer-alt mr-2"></i>
          Dashboard
        </router-link>
        
        <router-link
          to="/wallet"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/wallet' }"
        >
          <i class="fas fa-wallet mr-2"></i>
          Wallet
        </router-link>
        
        <router-link
          to="/transactions"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/transactions' }"
        >
          <i class="fas fa-exchange-alt mr-2"></i>
          Transactions
        </router-link>
        
        <router-link
          to="/transfer"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/transfer' }"
        >
          <i class="fas fa-paper-plane mr-2"></i>
          Transfer
        </router-link>
        
        <router-link
          to="/profile"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/profile' }"
        >
          <i class="fas fa-user mr-2"></i>
          Profile
        </router-link>
        
        <router-link
          to="/support"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/support' }"
        >
          <i class="fas fa-headset mr-2"></i>
          Support
        </router-link>
        
        <router-link
          v-if="isAdmin"
          to="/admin"
          class="mobile-nav-link"
          :class="{ 'active': $route.path === '/admin' }"
        >
          <i class="fas fa-shield-alt mr-2"></i>
          Admin
        </router-link>
      </div>
    </div>
  </nav>
</template>

<script>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

export default {
  name: 'Navigation',
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    const mobileMenuOpen = ref(false)
    
    const user = computed(() => authStore.user)
    const isAdmin = computed(() => authStore.isAdmin)
    
    const handleLogout = async () => {
      await authStore.logout()
      router.push('/auth/login')
    }
    
    return {
      user,
      isAdmin,
      mobileMenuOpen,
      handleLogout
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
