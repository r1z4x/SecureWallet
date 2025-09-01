<template>
  <div>
    
    <div v-if="loading" class="text-center py-4">
      <i class="fas fa-spinner fa-spin text-2xl text-primary-600"></i>
      <p class="text-gray-600 mt-2">Loading...</p>
    </div>

    <div v-else-if="twoFactorEnabled" class="space-y-4">
      <div class="flex items-center">
        <i class="fas fa-shield-alt text-green-600 text-2xl mr-3"></i>
        <div>
          <h3 class="text-lg font-medium text-gray-900">2FA is Enabled</h3>
          <p class="text-gray-600">Your account is protected with two-factor authentication</p>
        </div>
      </div>

      <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
        <h4 class="text-sm font-medium text-yellow-800 mb-2">Disable 2FA</h4>
        <p class="text-sm text-yellow-700 mb-3">
          Enter your current 2FA code to disable two-factor authentication.
        </p>
        
        <div class="space-y-3">
          <input
            v-model="disableCode"
            type="text"
            placeholder="Enter 6-digit code"
            class="form-input w-full"
            maxlength="6"
            pattern="[0-9]{6}"
          >
          <button
            @click="disable2FA"
            :disabled="disableLoading || !disableCode"
            class="btn-danger w-full"
          >
            <i v-if="disableLoading" class="fas fa-spinner fa-spin mr-2"></i>
            <i v-else class="fas fa-times mr-2"></i>
            Disable 2FA
          </button>
        </div>
      </div>
    </div>

    <div v-else class="space-y-4">
      <div class="flex items-center">
        <i class="fas fa-shield-alt text-gray-400 text-2xl mr-3"></i>
        <div>
          <h3 class="text-lg font-medium text-gray-900">2FA is Disabled</h3>
          <p class="text-gray-600">Enable two-factor authentication for enhanced security</p>
        </div>
      </div>

      <div v-if="qrCodeUrl" class="bg-gray-50 border border-gray-200 rounded-lg p-4">
        <h4 class="text-sm font-medium text-gray-800 mb-2">Setup Instructions</h4>
        <ol class="text-sm text-gray-700 space-y-2 mb-4">
          <li>1. Download an authenticator app (Google Authenticator, Authy, etc.)</li>
          <li>2. Scan the QR code below with your authenticator app</li>
          <li>3. Enter the 6-digit code from your app to enable 2FA</li>
        </ol>

        <div class="text-center mb-4">
          <img v-if="qrCodeDataUrl" :src="qrCodeDataUrl" alt="QR Code" class="mx-auto border border-gray-300 rounded">
          <div v-else class="text-center py-8">
            <i class="fas fa-qrcode text-4xl text-gray-300"></i>
            <p class="text-gray-500 mt-2">Loading QR code...</p>
          </div>
        </div>

        <div class="space-y-3">
          <input
            v-model="enableCode"
            type="text"
            placeholder="Enter 6-digit code"
            class="form-input w-full"
            maxlength="6"
            pattern="[0-9]{6}"
          >
          <button
            @click="enable2FA"
            :disabled="enableLoading || !enableCode"
            class="btn-primary w-full"
          >
            <i v-if="enableLoading" class="fas fa-spinner fa-spin mr-2"></i>
            <i v-else class="fas fa-shield-alt mr-2"></i>
            Enable 2FA
          </button>
        </div>
      </div>

      <div v-else class="text-center py-4">
        <button
          @click="load2FAStatus"
          class="btn-primary"
        >
          <i class="fas fa-shield-alt mr-2"></i>
          Setup 2FA
        </button>
      </div>
    </div>

    <!-- Error/Success Messages -->
    <div v-if="error" class="mt-4 p-3 bg-red-50 border border-red-200 rounded-lg">
      <p class="text-sm text-red-700">{{ error }}</p>
    </div>

    <div v-if="success" class="mt-4 p-3 bg-green-50 border border-green-200 rounded-lg">
      <p class="text-sm text-green-700">{{ success }}</p>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { twoFactorService } from '@/services/twoFactor'
import QRCode from 'qrcode'

export default {
  name: 'TwoFactorAuth',
  props: {
    initialStatus: {
      type: Boolean,
      default: false
    }
  },
  emits: ['2fa-enabled', '2fa-disabled'],
  setup(props, { emit }) {
    const loading = ref(false)
    const twoFactorEnabled = ref(props.initialStatus)
    const qrCodeUrl = ref('')
    const qrCodeDataUrl = ref('')
    const secret = ref('')
    const enableCode = ref('')
    const disableCode = ref('')
    const enableLoading = ref(false)
    const disableLoading = ref(false)
    const error = ref('')
    const success = ref('')

    const load2FAStatus = async () => {
      loading.value = true
      error.value = ''
      try {
        const response = await twoFactorService.getStatus()
        twoFactorEnabled.value = response.two_factor_enabled
        
        // Only generate QR code if 2FA is not enabled
        if (!response.two_factor_enabled && response.qr_code_url) {
          qrCodeUrl.value = response.qr_code_url
          // SECURE: Don't store secret in frontend
          // secret.value = response.secret
          
          // Generate QR code from the URL
          try {
            const qrDataUrl = await QRCode.toDataURL(response.qr_code_url, {
              width: 200,
              margin: 2,
              color: {
                dark: '#000000',
                light: '#FFFFFF'
              }
            })
            qrCodeDataUrl.value = qrDataUrl
          } catch (qrError) {
            console.error('Failed to generate QR code:', qrError)
            error.value = 'Failed to generate QR code'
          }
        }
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to load 2FA status'
      } finally {
        loading.value = false
      }
    }

    const enable2FA = async () => {
      if (!enableCode.value || enableCode.value.length !== 6) {
        error.value = 'Please enter a valid 6-digit code'
        return
      }

      enableLoading.value = true
      error.value = ''
      success.value = ''

      try {
        await twoFactorService.enable(enableCode.value)
        success.value = '2FA enabled successfully!'
        twoFactorEnabled.value = true
        qrCodeUrl.value = ''
        qrCodeDataUrl.value = ''
        enableCode.value = ''
        
        // Emit event to parent component
        emit('2fa-enabled')
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to enable 2FA'
      } finally {
        enableLoading.value = false
      }
    }

    const disable2FA = async () => {
      if (!disableCode.value || disableCode.value.length !== 6) {
        error.value = 'Please enter a valid 6-digit code'
        return
      }

      disableLoading.value = true
      error.value = ''
      success.value = ''

      try {
        await twoFactorService.disable(disableCode.value)
        success.value = '2FA disabled successfully!'
        twoFactorEnabled.value = false
        disableCode.value = ''
        
        // Clear QR code and reset state
        qrCodeUrl.value = ''
        qrCodeDataUrl.value = ''
        secret.value = ''
        enableCode.value = ''
        
        // Emit event to parent component
        emit('2fa-disabled')
      } catch (err) {
        error.value = err.response?.data?.error || 'Failed to disable 2FA'
      } finally {
        disableLoading.value = false
      }
    }

    onMounted(() => {
      // Always load 2FA status to get the latest state
      load2FAStatus()
    })

    // Watch for changes in initialStatus prop
    watch(() => props.initialStatus, (newStatus) => {
      twoFactorEnabled.value = newStatus
      // If 2FA is disabled, reload status to show setup option
      if (!newStatus) {
        load2FAStatus()
      }
    })

    return {
      loading,
      twoFactorEnabled,
      qrCodeUrl,
      qrCodeDataUrl,
      secret,
      enableCode,
      disableCode,
      enableLoading,
      disableLoading,
      error,
      success,
      load2FAStatus,
      enable2FA,
      disable2FA
    }
  }
}
</script>
