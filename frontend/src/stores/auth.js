import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authService } from '@/services/auth'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const token = ref(localStorage.getItem('token'))
  const loading = ref(false)

  // Getters
  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.is_admin || false)

  // Actions
  async function login(credentials) {
    loading.value = true
    try {
      const response = await authService.login(credentials)
      
      // Check if 2FA is required
      if (response.requires_2fa) {
        return response
      }
      
      token.value = response.access_token
      user.value = response.user
      localStorage.setItem('token', response.access_token)
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
      const response = await authService.login2FA(userId, code)
      token.value = response.access_token
      user.value = response.user
      localStorage.setItem('token', response.access_token)
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
      return response
    } catch (error) {
      console.error('Get current user error:', error)
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
    if (token.value) {
      await getCurrentUser()
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
    
    // Actions
    login,
    login2FA,
    register,
    logout,
    getCurrentUser,
    refreshToken,
    requestPasswordReset,
    init
  }
})
