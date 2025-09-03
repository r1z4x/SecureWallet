import { defineStore } from 'pinia'
import { ref, computed, watch, nextTick } from 'vue'
import { authService } from '@/services/auth'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(null)
  const loading = ref(false)

  // Initialize token from localStorage
  const initToken = () => {
    const storedToken = localStorage.getItem('token')
    if (storedToken) {
      token.value = storedToken
    }
  }

  // Call initToken immediately
  initToken()

  // Watch for token changes and clear user data when token is cleared
  watch(token, (newToken) => {
    if (!newToken) {
      user.value = null
    }
  })

  // Watch for user changes
  watch(user, (newUser, oldUser) => {
    // User change detected
  })

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.is_admin || false)
  const isUserLoaded = computed(() => !!user.value)

  // Actions
  async function login(credentials) {
    loading.value = true
    try {
      // Clear any existing user data before login
      user.value = null
      
      const response = await authService.login(credentials)
      
      // Check if 2FA is required
      if (response && response.requires_2fa) {
        return response
      }
      
      // Set token and user data immediately
      token.value = response.access_token
      localStorage.setItem('token', response.access_token)
      
      // If user data is not in response, fetch it using getCurrentUser
      if (response.user) {
        user.value = response.user
      } else {
        // Fetch user data using the new token
        try {
          // Small delay to ensure axios interceptor picks up the new token
          await new Promise(resolve => setTimeout(resolve, 10))
          
          // Create a new axios instance with the token to ensure it's used
          const axios = (await import('axios')).default
          const tempApi = axios.create({
            baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
            timeout: 10000,
            headers: {
              'Authorization': `Bearer ${response.access_token}`
            }
          })
          
          const userResponse = await tempApi.get('/auth/me')
          const userData = userResponse.data
          user.value = userData
        } catch (error) {
          console.error('Failed to fetch user data after login:', error)
          // Don't throw error, just log it
        }
      }
      
      // Ensure reactive update by triggering a change
      await nextTick()
      
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  async function login2FA(userId, code) {
    loading.value = true
    try {
      // Clear any existing user data before 2FA login
      user.value = null
      
      const response = await authService.login2FA(userId, code)
      
      // Set token and user data immediately
      token.value = response.access_token
      localStorage.setItem('token', response.access_token)
      
      // If user data is not in response, fetch it using getCurrentUser
      if (response.user) {
        user.value = response.user
      } else {
        // Fetch user data using the new token
        try {
          // Small delay to ensure axios interceptor picks up the new token
          await new Promise(resolve => setTimeout(resolve, 10))
          
          // Create a new axios instance with the token to ensure it's used
          const axios = (await import('axios')).default
          const tempApi = axios.create({
            baseURL: import.meta.env.VITE_API_BASE_URL || '/api',
            timeout: 10000,
            headers: {
              'Authorization': `Bearer ${response.access_token}`
            }
          })
          
          const userResponse = await tempApi.get('/auth/me')
          const userData = userResponse.data
          user.value = userData
        } catch (error) {
          console.error('Failed to fetch user data after login2FA:', error)
          // Don't throw error, just log it
        }
      }
      
      // Ensure reactive update by triggering a change
      await nextTick()
      
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  async function register(userData) {
    loading.value = true
    try {
      const response = await authService.register(userData)
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    loading.value = true
    try {
      if (token.value) {
        await authService.logout()
      }
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      // Clear all auth data
      token.value = null
      user.value = null
      localStorage.removeItem('token')
      loading.value = false
    }
  }

  async function getCurrentUser() {
    if (!token.value) return null
    
    loading.value = true
    try {
      const response = await authService.getCurrentUser()
      user.value = response
      
      // Ensure reactive update
      await nextTick()
      
      return response
    } catch (error) {
      console.error('Get current user error:', error)
      // Clear user data on error
      user.value = null
      logout()
      throw error
    } finally {
      loading.value = false
    }
  }

  async function refreshToken() {
    if (!token.value) return false
    
    try {
      const response = await authService.refreshToken()
      token.value = response.access_token
      localStorage.setItem('token', response.access_token)
      return true
    } catch (error) {
      console.error('Token refresh error:', error)
      logout()
      return false
    }
  }

  async function requestPasswordReset(email) {
    loading.value = true
    try {
      const response = await authService.requestPasswordReset(email)
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // Initialize auth state
  async function init() {
    // Clear any existing user data first
    user.value = null
    
    // Ensure token is loaded from localStorage
    initToken()
    
    if (token.value) {
      try {
        await getCurrentUser()
        
        // Ensure reactive update
        await nextTick()
      } catch (error) {
        console.error('Failed to get current user during init:', error)
        // If getCurrentUser fails, clear the token
        token.value = null
        localStorage.removeItem('token')
      }
    }
  }

  return {
    // State
    user,
    token,
    loading,
    
    // Getters
    isAuthenticated,
    isAdmin,
    isUserLoaded,
    
    // Actions
    login,
    login2FA,
    register,
    logout,
    getCurrentUser,
    refreshToken,
    requestPasswordReset,
    init,
    initToken
  }
})
