<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader type="authenticated" />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Profile Settings</h1>
        <p class="text-gray-600 mt-2">Manage your account information and preferences</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Profile Form -->
        <div class="lg:col-span-2">

          <div class="bg-white rounded-lg shadow-md p-6">
            <h2 class="text-xl font-semibold text-gray-900 mb-6">Personal Information</h2>
            
            <form @submit.prevent="handleUpdateProfile" class="space-y-6">
              <!-- Username -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Username
                </label>
                <input
                  v-model="profileForm.username"
                  type="text"
                  class="form-input w-full"
                  placeholder="Enter your username"
                  required
                >
              </div>

              <!-- Email -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Email Address
                </label>
                <input
                  v-model="profileForm.email"
                  type="email"
                  class="form-input w-full"
                  placeholder="Enter your email"
                  required
                >
              </div>

              <!-- Current Password -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Current Password
                </label>
                <input
                  v-model="profileForm.currentPassword"
                  type="password"
                  class="form-input w-full"
                  placeholder="Enter current password to confirm changes"
                >
                <p class="text-xs text-gray-500 mt-1">
                  Required when changing email or password
                </p>
              </div>

              <!-- New Password -->
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  New Password (Optional)
                </label>
                <input
                  v-model="profileForm.newPassword"
                  type="password"
                  class="form-input w-full"
                  placeholder="Enter new password"
                  minlength="8"
                >
                <p class="text-xs text-gray-500 mt-1">
                  Leave blank to keep current password
                </p>
              </div>

              <!-- Confirm New Password -->
              <div v-if="profileForm.newPassword">
                <label class="block text-sm font-medium text-gray-700 mb-2">
                  Confirm New Password
                </label>
                <input
                  v-model="profileForm.confirmPassword"
                  type="password"
                  class="form-input w-full"
                  placeholder="Confirm new password"
                  minlength="8"
                >
              </div>

              <!-- Update Button -->
              <div class="pt-4">
                <button
                  type="submit"
                  class="btn-primary w-full"
                  :disabled="updateLoading"
                >
                  <i v-if="updateLoading" class="fas fa-spinner fa-spin mr-2"></i>
                  <i v-else class="fas fa-save mr-2"></i>
                  {{ updateLoading ? 'Updating...' : 'Update Profile' }}
                </button>
              </div>
            </form>

            <!-- Error Message -->
            <div v-if="error" class="mt-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
              {{ error }}
            </div>

            <!-- Success Message -->
            <div v-if="success" class="mt-4 p-4 bg-green-100 border border-green-400 text-green-700 rounded">
              {{ success }}
            </div>
          </div>

          <!-- Security Settings -->
          <div class="bg-white rounded-lg shadow-md p-6 mt-6">
            <h2 class="text-xl font-semibold text-gray-900 mb-6">Security Settings</h2>
            
            <div class="space-y-4">
              <!-- Two-Factor Authentication -->
              <div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                <div>
                  <h3 class="text-lg font-medium text-gray-900">Two-Factor Authentication</h3>
                  <p class="text-sm text-gray-600">Add an extra layer of security to your account</p>
                </div>
                <div class="flex items-center">
                  <span class="text-sm mr-3" :class="userData.two_factor_enabled ? 'text-green-600' : 'text-gray-500'">
                    <i class="fas fa-circle text-xs mr-1" :class="userData.two_factor_enabled ? 'text-green-500' : 'text-gray-400'"></i>
                    {{ userData.two_factor_enabled ? 'Enabled' : 'Not enabled' }}
                  </span>
                  <button @click="open2FAModal" class="btn-secondary">
                    <i class="fas fa-shield-alt mr-2"></i>
                    {{ userData.two_factor_enabled ? 'Configure' : 'Enable' }}
                  </button>
                </div>
              </div>

              <!-- Login History -->
              <div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                <div>
                  <h3 class="text-lg font-medium text-gray-900">Login History</h3>
                  <p class="text-sm text-gray-600">View your recent login activity</p>
                </div>
                <button @click="showLoginHistoryModal = true" class="btn-secondary">
                  <i class="fas fa-history mr-2"></i>
                  View History
                </button>
              </div>

              <!-- Account Deletion -->
              <div class="flex items-center justify-between p-4 border border-red-200 rounded-lg bg-red-50">
                <div>
                  <h3 class="text-lg font-medium text-red-900">Delete Account</h3>
                  <p class="text-sm text-red-600">Permanently delete your account and all data</p>
                </div>
                <button 
                  @click="showDeleteAccount = true" 
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors duration-200"
                >
                  <i class="fas fa-trash mr-2"></i>
                  Delete Account
                </button>
              </div>

              <!-- Reset Database -->
              <div class="flex items-center justify-between p-4 border border-orange-200 rounded-lg bg-orange-50">
                <div>
                  <h3 class="text-lg font-medium text-orange-900">Reset Database</h3>
                  <p class="text-sm text-orange-600">Reset and recreate database with sample data</p>
                </div>
                <router-link 
                  to="/reset-database" 
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-orange-600 hover:bg-orange-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-orange-500 transition-colors duration-200"
                >
                  <i class="fas fa-database mr-2"></i>
                  Reset Database
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <!-- Profile Info -->
        <div class="space-y-6">
          <!-- Account Summary -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Account Summary</h3>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600">Member Since:</span>
                <span class="font-medium">{{ formatDate(userData.created_at) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Last Updated:</span>
                <span class="font-medium">{{ formatDate(userData.updated_at) }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Account Status:</span>
                <span class="text-green-600 font-medium">Active</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">User Type:</span>
                <span class="font-medium">{{ userData.is_admin ? 'Administrator' : 'User' }}</span>
              </div>
            </div>
          </div>

          <!-- Quick Stats -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Stats</h3>
            <div class="space-y-3">
              <div class="flex items-center">
                <i class="fas fa-wallet text-primary-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Current Balance</p>
                  <p class="text-lg font-bold text-primary-600">${{ walletData.balance?.toFixed(2) || '0.00' }}</p>
                </div>
              </div>
              <div class="flex items-center">
                <i class="fas fa-exchange-alt text-blue-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Total Transactions</p>
                  <p class="text-lg font-bold text-blue-600">{{ walletData.transaction_count || 0 }}</p>
                </div>
              </div>
              <div class="flex items-center">
                <i class="fas fa-calendar text-green-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Days Active</p>
                  <p class="text-lg font-bold text-green-600">{{ daysActive }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Help & Support -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Help & Support</h3>
            <div class="space-y-3">
              <router-link
                to="/support"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-headset text-primary-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Contact Support</p>
                  <p class="text-xs text-gray-500">Get help with your account</p>
                </div>
              </router-link>
              <router-link
                to="/help"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-book text-blue-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Help Center</p>
                  <p class="text-xs text-gray-500">Browse help articles</p>
                </div>
              </router-link>
              <router-link
                to="/faq"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-question-circle text-purple-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">FAQ</p>
                  <p class="text-xs text-gray-500">Frequently asked questions</p>
                </div>
              </router-link>
              <router-link
                to="/user-guide"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-user-graduate text-indigo-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">User Guide</p>
                  <p class="text-xs text-gray-500">Learn how to use SecureWallet</p>
                </div>
              </router-link>
              <router-link
                to="/security-tips"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-shield-alt text-red-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Security Tips</p>
                  <p class="text-xs text-gray-500">Keep your account secure</p>
                </div>
              </router-link>
              <router-link
                to="/terms"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-file-alt text-green-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Terms of Service</p>
                  <p class="text-xs text-gray-500">Read our terms</p>
                </div>
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- Two-Factor Authentication Modal -->
      <div v-if="show2FAModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="show2FAModal = false">
        <div class="bg-white rounded-lg p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold text-gray-900">Two-Factor Authentication</h2>
            <button @click="show2FAModal = false" class="text-gray-400 hover:text-gray-600">
              <i class="fas fa-times text-xl"></i>
            </button>
          </div>
          <TwoFactorAuth 
            :key="`2fa-${userData.two_factor_enabled}-${show2FAModal}`"
            :initial-status="userData.two_factor_enabled"
            @2fa-enabled="handle2FAEnabled" 
            @2fa-disabled="handle2FADisabled" 
          />
        </div>
      </div>

      <!-- Login History Modal -->
      <div v-if="showLoginHistoryModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="showLoginHistoryModal = false">
        <div class="bg-white rounded-lg p-6 max-w-4xl w-full mx-4 max-h-[90vh] overflow-y-auto" @click.stop>
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold text-gray-900">Login History</h2>
            <button @click="showLoginHistoryModal = false" class="text-gray-400 hover:text-gray-600">
              <i class="fas fa-times text-xl"></i>
            </button>
          </div>
          <LoginHistory />
        </div>
      </div>

      <!-- Delete Account Modal -->
      <div v-if="showDeleteAccount" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click="showDeleteAccount = false">
        <div class="bg-white rounded-lg p-6 max-w-md w-full mx-4" @click.stop>
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-semibold text-gray-900">Delete Account</h2>
            <button @click="showDeleteAccount = false" class="text-gray-400 hover:text-gray-600">
              <i class="fas fa-times text-xl"></i>
            </button>
          </div>
          <DeleteAccount @close="showDeleteAccount = false" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/AppHeader.vue'
import TwoFactorAuth from '@/components/TwoFactorAuth.vue'
import LoginHistory from '@/components/LoginHistory.vue'
import DeleteAccount from '@/components/DeleteAccount.vue'
import { userService } from '@/services/user'
import { walletService } from '@/services/wallet'

export default {
  name: 'Profile',
  components: {
    AppHeader,
    TwoFactorAuth,
    LoginHistory,
    DeleteAccount
  },
  setup() {
    const authStore = useAuthStore()
    const user = computed(() => authStore.user)
    
    const updateLoading = ref(false)
    const error = ref('')
    const success = ref('')
    
    const userData = ref({
      username: '',
      email: '',
      created_at: null,
      updated_at: null,
      is_admin: false,
      two_factor_enabled: false
    })
    
    const walletData = ref({
      balance: 0,
      transaction_count: 0
    })
    
    const profileForm = ref({
      username: '',
      email: '',
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    })

    const showDeleteAccount = ref(false)
    const show2FAModal = ref(false)
    const showLoginHistoryModal = ref(false)

    const loadUserData = async () => {
      try {
        const response = await userService.getCurrentUser()
        userData.value = response
        profileForm.value.username = response.username
        profileForm.value.email = response.email
      } catch (error) {
        console.error('Error loading user data:', error)
      }
    }

    const loadWalletData = async () => {
      try {
        const response = await walletService.getBalance()
        walletData.value = response
      } catch (error) {
        console.error('Error loading wallet data:', error)
      }
    }

    const handleUpdateProfile = async () => {
      updateLoading.value = true
      error.value = ''
      success.value = ''
      
      try {
        // Validate password confirmation
        if (profileForm.value.newPassword && profileForm.value.newPassword !== profileForm.value.confirmPassword) {
          throw new Error('New passwords do not match')
        }
        
        // Prepare update data
        const updateData = {
          username: profileForm.value.username,
          email: profileForm.value.email
        }
        
        if (profileForm.value.newPassword) {
          updateData.password = profileForm.value.newPassword
        }
        
        // Update profile
        await userService.updateCurrentUser(updateData)
        
        // Update auth store
        await authStore.getCurrentUser()
        
        success.value = 'Profile updated successfully!'
        
        // Clear form
        profileForm.value.currentPassword = ''
        profileForm.value.newPassword = ''
        profileForm.value.confirmPassword = ''
        
        // Reload data
        await loadUserData()
        
      } catch (err) {
        error.value = err.response?.data?.detail || err.message || 'Failed to update profile'
      } finally {
        updateLoading.value = false
      }
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    const scrollTo2FA = () => {
      document.querySelector('[data-section="2fa"]')?.scrollIntoView({ behavior: 'smooth' })
    }

    const scrollToLoginHistory = () => {
      document.querySelector('[data-section="login-history"]')?.scrollIntoView({ behavior: 'smooth' })
    }

    const handle2FAEnabled = () => {
      // Update the 2FA status in the sidebar
      // This will trigger a re-render of the sidebar
      show2FAModal.value = false
      // Reload user data to get updated 2FA status
      loadUserData()
      // Show success message
      success.value = 'Two-factor authentication has been enabled successfully.'
    }

    const handle2FADisabled = () => {
      // Update the 2FA status in the sidebar
      show2FAModal.value = false
      // Reload user data to get updated 2FA status
      loadUserData()
      // Show success message
      success.value = 'Two-factor authentication has been disabled successfully.'
    }

    const open2FAModal = () => {
      show2FAModal.value = true
      // Clear any previous success/error messages
      success.value = ''
      error.value = ''
    }

    const daysActive = computed(() => {
      if (!userData.value.created_at) return 0
      const created = new Date(userData.value.created_at)
      const now = new Date()
      const diffTime = Math.abs(now - created)
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
      return diffDays
    })

    const isModalOpen = computed(() => {
      return show2FAModal.value || showLoginHistoryModal.value || showDeleteAccount.value
    })

    // Watch for modal state changes to prevent body scroll
    watch(isModalOpen, (isOpen) => {
      if (isOpen) {
        document.body.style.overflow = 'hidden'
      } else {
        document.body.style.overflow = 'auto'
      }
    })

    onMounted(() => {
      loadUserData()
      loadWalletData()
      
      // Add keyboard event listener for ESC key
      const handleKeydown = (event) => {
        if (event.key === 'Escape') {
          show2FAModal.value = false
          showLoginHistoryModal.value = false
          showDeleteAccount.value = false
        }
      }
      
      document.addEventListener('keydown', handleKeydown)
      
      // Cleanup on unmount
      onUnmounted(() => {
        document.removeEventListener('keydown', handleKeydown)
      })
    })

    return {
      user,
      updateLoading,
      error,
      success,
      userData,
      walletData,
      profileForm,
      showDeleteAccount,
      show2FAModal,
      showLoginHistoryModal,
      daysActive,
      handleUpdateProfile,
      handle2FAEnabled,
      handle2FADisabled,
      open2FAModal,
      formatDate,
      scrollTo2FA,
      scrollToLoginHistory
    }
  }
}
</script>


