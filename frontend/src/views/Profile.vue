<template>
  <div class="min-h-screen bg-gray-50">
    <Navigation />
    
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
                  <span class="text-sm text-gray-500 mr-3">Not enabled</span>
                  <button class="btn-secondary">
                    <i class="fas fa-shield-alt mr-2"></i>
                    Enable
                  </button>
                </div>
              </div>

              <!-- Login History -->
              <div class="flex items-center justify-between p-4 border border-gray-200 rounded-lg">
                <div>
                  <h3 class="text-lg font-medium text-gray-900">Login History</h3>
                  <p class="text-sm text-gray-600">View your recent login activity</p>
                </div>
                <button class="btn-secondary">
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
                <button class="btn-danger">
                  <i class="fas fa-trash mr-2"></i>
                  Delete Account
                </button>
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
              <a
                href="#"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-book text-blue-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Help Center</p>
                  <p class="text-xs text-gray-500">Browse help articles</p>
                </div>
              </a>
              <a
                href="#"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-file-alt text-green-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Terms of Service</p>
                  <p class="text-xs text-gray-500">Read our terms</p>
                </div>
              </a>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Navigation from '@/components/Navigation.vue'
import { userService } from '@/services/user'
import { walletService } from '@/services/wallet'

export default {
  name: 'Profile',
  components: {
    Navigation
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
      is_admin: false
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

    const daysActive = computed(() => {
      if (!userData.value.created_at) return 0
      const created = new Date(userData.value.created_at)
      const now = new Date()
      const diffTime = Math.abs(now - created)
      const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
      return diffDays
    })

    onMounted(() => {
      loadUserData()
      loadWalletData()
    })

    return {
      user,
      updateLoading,
      error,
      success,
      userData,
      walletData,
      profileForm,
      daysActive,
      handleUpdateProfile,
      formatDate
    }
  }
}
</script>
