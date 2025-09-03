import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api'

// Create axios instance
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000
})

// Export the api instance for use in other services
export const apiClient = api

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor removed - router guard handles auth redirects

export const authService = {
  async login(credentials) {
    const response = await api.post('/auth/login', credentials)
    return response.data
  },

  async login2FA(userId, code) {
    const response = await api.post('/auth/login/2fa', { user_id: userId, code })
    return response.data
  },

  async register(userData) {
    const response = await api.post('/auth/register', userData)
    return response.data
  },

  async logout() {
    const response = await api.post('/auth/logout')
    return response.data
  },

  async getCurrentUser() {
    const response = await api.get('/auth/me')
    return response.data
  },

  async refreshToken() {
    const response = await api.post('/auth/refresh')
    return response.data
  },

  async requestPasswordReset(email) {
    const response = await api.post('/auth/password-reset', { email })
    return response.data
  },

  async resetPassword(token, newPassword) {
    const response = await api.post('/auth/password-reset/confirm', {
      token,
      new_password: newPassword
    })
    return response.data
  }
}
