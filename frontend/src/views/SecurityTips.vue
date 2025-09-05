<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader type="authenticated" />
    
    <div class="container mx-auto px-4 py-8">
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-gray-900 mb-4">Security Tips</h1>
        <p class="text-lg text-gray-600 max-w-3xl mx-auto">
          Protect your SecureWallet account with these essential security practices. Your security is our top priority.
        </p>
      </div>

      <!-- Security Status -->
      <div class="max-w-4xl mx-auto mb-8">
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h2 class="text-2xl font-semibold text-gray-900 mb-4 flex items-center">
            <i class="fas fa-shield-alt text-green-600 mr-3"></i>
            Your Security Status
          </h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div
              v-for="(check, index) in securityChecks"
              :key="index"
              class="flex items-center p-4 border border-gray-200 rounded-lg"
              :class="check.completed ? 'bg-green-50 border-green-200' : 'bg-yellow-50 border-yellow-200'"
            >
              <div class="flex-shrink-0 mr-4">
                <i
                  :class="[
                    'fas text-2xl',
                    check.completed ? 'fa-check-circle text-green-600' : 'fa-exclamation-circle text-yellow-600'
                  ]"
                ></i>
              </div>
              <div class="flex-1">
                <h3 class="font-medium text-gray-900">{{ check.title }}</h3>
                <p class="text-sm text-gray-600">{{ check.description }}</p>
              </div>
              <div class="flex-shrink-0">
                <span
                  class="px-3 py-1 text-xs font-medium rounded-full"
                  :class="check.completed ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'"
                >
                  {{ check.completed ? 'Completed' : 'Pending' }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Security Categories -->
      <div class="max-w-4xl mx-auto space-y-8">
        <div
          v-for="(category, index) in securityCategories"
          :key="category.id"
          class="bg-white rounded-lg shadow-sm border border-gray-200"
        >
          <div class="p-6 border-b border-gray-200">
            <div class="flex items-center mb-4">
              <div class="w-12 h-12 rounded-full flex items-center justify-center text-white text-xl font-medium mr-4"
                   :class="category.colorClass">
                <i :class="category.icon"></i>
              </div>
              <div>
                <h2 class="text-2xl font-semibold text-gray-900">{{ category.title }}</h2>
                <p class="text-gray-600">{{ category.description }}</p>
              </div>
            </div>
          </div>
          
          <div class="p-6">
            <div v-for="(tip, tipIndex) in category.tips" :key="tipIndex" class="mb-8 last:mb-0">
              <div class="flex items-start">
                <div class="w-8 h-8 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center text-sm font-medium mr-4 mt-1 flex-shrink-0">
                  {{ tipIndex + 1 }}
                </div>
                <div class="flex-1">
                  <h3 class="text-lg font-medium text-gray-900 mb-2">{{ tip.title }}</h3>
                  <p class="text-gray-600 mb-4">{{ tip.description }}</p>
                  
                  <!-- Action Items -->
                  <div v-if="tip.actions && tip.actions.length > 0" class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
                    <h4 class="font-medium text-blue-900 mb-3 flex items-center">
                      <i class="fas fa-tasks mr-2"></i>
                      Action Items
                    </h4>
                    <ul class="space-y-2">
                      <li
                        v-for="(action, actionIndex) in tip.actions"
                        :key="actionIndex"
                        class="flex items-start"
                      >
                        <i class="fas fa-arrow-right text-blue-600 mr-2 mt-1 flex-shrink-0"></i>
                        <span class="text-sm text-blue-800">{{ action }}</span>
                      </li>
                    </ul>
                  </div>
                  
                  <!-- Examples -->
                  <div v-if="tip.examples && tip.examples.length > 0" class="bg-gray-50 border border-gray-200 rounded-lg p-4">
                    <h4 class="font-medium text-gray-900 mb-3 flex items-center">
                      <i class="fas fa-lightbulb mr-2"></i>
                      Examples
                    </h4>
                    <div class="space-y-2">
                      <div
                        v-for="(example, exampleIndex) in tip.examples"
                        :key="exampleIndex"
                        class="flex items-start"
                      >
                        <div class="w-2 h-2 bg-green-500 rounded-full mr-3 mt-2 flex-shrink-0"></div>
                        <div>
                          <p class="text-sm font-medium text-gray-900">{{ example.title }}</p>
                          <p class="text-sm text-gray-600">{{ example.description }}</p>
                        </div>
                      </div>
                    </div>
                  </div>
                  
                  <!-- Warnings -->
                  <div v-if="tip.warnings && tip.warnings.length > 0" class="bg-red-50 border border-red-200 rounded-lg p-4 mt-4">
                    <h4 class="font-medium text-red-900 mb-3 flex items-center">
                      <i class="fas fa-exclamation-triangle mr-2"></i>
                      Security Warnings
                    </h4>
                    <ul class="space-y-2">
                      <li
                        v-for="(warning, warningIndex) in tip.warnings"
                        :key="warningIndex"
                        class="flex items-start"
                      >
                        <i class="fas fa-times-circle text-red-600 mr-2 mt-0.5 flex-shrink-0"></i>
                        <span class="text-sm text-red-800">{{ warning }}</span>
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Emergency Contacts -->
      <div class="max-w-4xl mx-auto mt-12">
        <div class="bg-red-50 border border-red-200 rounded-lg p-6">
          <h2 class="text-2xl font-semibold text-red-900 mb-4 flex items-center">
            <i class="fas fa-exclamation-triangle mr-3"></i>
            Emergency Security Contacts
          </h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="bg-white rounded-lg p-4 border border-red-200">
              <h3 class="font-medium text-red-900 mb-2">Suspicious Activity</h3>
              <p class="text-sm text-red-700 mb-3">If you notice any suspicious activity on your account:</p>
              <ul class="text-sm text-red-700 space-y-1">
                <li>• Change your password immediately</li>
                <li>• Enable two-factor authentication</li>
                <li>• Contact support right away</li>
                <li>• Review your recent transactions</li>
              </ul>
            </div>
            <div class="bg-white rounded-lg p-4 border border-red-200">
              <h3 class="font-medium text-red-900 mb-2">Account Compromise</h3>
              <p class="text-sm text-red-700 mb-3">If you suspect your account has been compromised:</p>
              <ul class="text-sm text-red-700 space-y-1">
                <li>• Contact support immediately</li>
                <li>• Freeze your account if possible</li>
                <li>• Document all suspicious activity</li>
                <li>• Monitor your other accounts</li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <!-- AppHeader -->
      <div class="max-w-4xl mx-auto mt-12">
        <div class="flex flex-col sm:flex-row justify-between items-center gap-4">
          <router-link
            to="/help"
            class="inline-flex items-center px-6 py-3 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors"
          >
            <i class="fas fa-arrow-left mr-2"></i>
            Back to Help Center
          </router-link>
          
          <div class="flex gap-4">
            <router-link
              to="/faq"
              class="inline-flex items-center px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              <i class="fas fa-question-circle mr-2"></i>
              View FAQ
            </router-link>
            <router-link
              to="/user-guide"
              class="inline-flex items-center px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
            >
              <i class="fas fa-user-graduate mr-2"></i>
              User Guide
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/AppHeader.vue'

export default {
  name: 'SecurityTips',
  components: {
    AppHeader
  },
  setup() {
    const authStore = useAuthStore()
    const user = computed(() => authStore.user)

    const securityChecks = ref([
      {
        title: 'Two-Factor Authentication',
        description: 'Enable 2FA for enhanced account security',
        completed: computed(() => user.value?.two_factor_enabled || false)
      },
      {
        title: 'Strong Password',
        description: 'Use a strong, unique password',
        completed: true // This would need to be checked against password strength
      },
      {
        title: 'Email Verification',
        description: 'Verify your email address',
        completed: true // Assuming email is verified if user is logged in
      },
      {
        title: 'Login History Review',
        description: 'Regularly check your login history',
        completed: false // This is a manual check
      }
    ])

    const securityCategories = ref([
      {
        id: 'passwords',
        title: 'Password Security',
        description: 'Create and maintain strong passwords to protect your account.',
        icon: 'fas fa-key',
        colorClass: 'bg-blue-600',
        tips: [
          {
            title: 'Create Strong Passwords',
            description: 'Use passwords that are difficult for others to guess but easy for you to remember.',
            actions: [
              'Use at least 12 characters',
              'Include uppercase and lowercase letters',
              'Add numbers and special characters',
              'Avoid personal information',
              'Use different passwords for each account'
            ],
            examples: [
              {
                title: 'Good Password',
                description: 'MyF@voriteC@t2024! - combines personal meaning with complexity'
              },
              {
                title: 'Better Password',
                description: 'Tr3@sur3Hunt3r#2024 - uses a phrase with substitutions'
              }
            ],
            warnings: [
              'Never use "password" or "123456"',
              'Don\'t use your name, birthday, or common words',
              'Avoid using the same password across multiple accounts',
              'Don\'t share your password with anyone'
            ]
          },
          {
            title: 'Password Management',
            description: 'Keep your passwords secure and up to date.',
            actions: [
              'Change your password every 3-6 months',
              'Use a password manager for secure storage',
              'Enable password recovery options',
              'Never write passwords on paper',
              'Log out from shared devices'
            ],
            warnings: [
              'Don\'t store passwords in plain text files',
              'Avoid using browser password autofill on public computers',
              'Don\'t share passwords via email or text messages'
            ]
          }
        ]
      },
      {
        id: 'authentication',
        title: 'Two-Factor Authentication',
        description: 'Add an extra layer of security to your account.',
        icon: 'fas fa-shield-alt',
        colorClass: 'bg-green-600',
        tips: [
          {
            title: 'Enable 2FA',
            description: 'Two-factor authentication requires both your password and a second verification method.',
            actions: [
              'Go to your Profile settings',
              'Click on "Security" section',
              'Enable Two-Factor Authentication',
              'Scan the QR code with your authenticator app',
              'Enter the verification code',
              'Save your backup codes securely'
            ],
            examples: [
              {
                title: 'Authenticator Apps',
                description: 'Google Authenticator, Authy, Microsoft Authenticator'
              },
              {
                title: 'Backup Codes',
                description: 'Keep these codes safe for account recovery'
              }
            ],
            warnings: [
              'Don\'t share your 2FA codes with anyone',
              'Keep your backup codes in a secure location',
              'Don\'t use SMS 2FA if possible (use authenticator apps)',
              'Test your 2FA setup before logging out'
            ]
          },
          {
            title: '2FA Best Practices',
            description: 'Use two-factor authentication effectively.',
            actions: [
              'Use authenticator apps instead of SMS',
              'Keep backup codes in a secure location',
              'Set up 2FA on multiple devices if possible',
              'Regularly review your 2FA settings'
            ],
            warnings: [
              'Don\'t disable 2FA unless absolutely necessary',
              'Be careful when changing phones or devices',
              'Contact support if you lose access to your 2FA'
            ]
          }
        ]
      },
      {
        id: 'account-monitoring',
        title: 'Account Monitoring',
        description: 'Regularly monitor your account for suspicious activity.',
        icon: 'fas fa-eye',
        colorClass: 'bg-purple-600',
        tips: [
          {
            title: 'Review Login History',
            description: 'Check your login history regularly to identify unauthorized access.',
            actions: [
              'Visit your Profile page weekly',
              'Click on "Login History"',
              'Look for unfamiliar locations or devices',
              'Check login times for unusual patterns',
              'Report suspicious activity immediately'
            ],
            examples: [
              {
                title: 'Suspicious Activity',
                description: 'Logins from unknown locations, unusual times, or multiple failed attempts'
              },
              {
                title: 'Normal Activity',
                description: 'Logins from your usual devices and locations'
              }
            ],
            warnings: [
              'Don\'t ignore login notifications',
              'Change your password if you see suspicious activity',
              'Contact support immediately for unauthorized access'
            ]
          },
          {
            title: 'Transaction Monitoring',
            description: 'Monitor your transactions for any unauthorized activity.',
            actions: [
              'Review your transaction history regularly',
              'Set up transaction notifications',
              'Verify all transfers before confirming',
              'Report unauthorized transactions immediately'
            ],
            warnings: [
              'Don\'t ignore unexpected transaction notifications',
              'Verify recipient information before sending money',
              'Contact support for any suspicious transactions'
            ]
          }
        ]
      },
      {
        id: 'device-security',
        title: 'Device Security',
        description: 'Keep your devices secure to protect your account.',
        icon: 'fas fa-laptop',
        colorClass: 'bg-orange-600',
        tips: [
          {
            title: 'Secure Your Devices',
            description: 'Protect the devices you use to access your account.',
            actions: [
              'Keep your operating system updated',
              'Install and update antivirus software',
              'Use a firewall',
              'Enable device encryption',
              'Lock your devices with passwords or biometrics'
            ],
            examples: [
              {
                title: 'Device Lock',
                description: 'Use fingerprint, face recognition, or strong PIN'
              },
              {
                title: 'Software Updates',
                description: 'Enable automatic updates for security patches'
              }
            ],
            warnings: [
              'Don\'t access your account on public computers',
              'Avoid using public Wi-Fi for financial transactions',
              'Don\'t share your devices with others',
              'Be careful with app permissions'
            ]
          },
          {
            title: 'Browser Security',
            description: 'Use secure browsers and keep them updated.',
            actions: [
              'Use updated browsers (Chrome, Firefox, Safari, Edge)',
              'Enable automatic updates',
              'Clear browser data regularly',
              'Use incognito/private mode on shared devices',
              'Disable saved passwords on public computers'
            ],
            warnings: [
              'Don\'t use outdated browsers',
              'Avoid browser extensions from unknown sources',
              'Don\'t save passwords on public computers',
              'Be careful with browser autofill'
            ]
          }
        ]
      },
      {
        id: 'phishing-protection',
        title: 'Phishing Protection',
        description: 'Learn to identify and avoid phishing attempts.',
        icon: 'fas fa-fish',
        colorClass: 'bg-red-600',
        tips: [
          {
            title: 'Identify Phishing Attempts',
            description: 'Learn to recognize common phishing tactics.',
            actions: [
              'Check the sender\'s email address carefully',
              'Look for spelling and grammar errors',
              'Verify URLs before clicking',
              'Don\'t trust urgent or threatening messages',
              'Contact the company directly to verify'
            ],
            examples: [
              {
                title: 'Suspicious Email',
                description: 'Urgent message claiming your account is locked'
              },
              {
                title: 'Fake URL',
                description: 'securewallet-secure.com instead of securewallet.com'
              }
            ],
            warnings: [
              'Never click links in suspicious emails',
              'Don\'t enter credentials on unfamiliar websites',
              'Don\'t trust emails asking for personal information',
              'Be especially careful with financial emails'
            ]
          },
          {
            title: 'Safe Browsing Practices',
            description: 'Browse the internet safely to avoid security threats.',
            actions: [
              'Only visit trusted websites',
              'Check for HTTPS in the URL',
              'Look for security indicators in your browser',
              'Use bookmarks for important sites',
              'Be careful with downloads'
            ],
            warnings: [
              'Don\'t download files from unknown sources',
              'Avoid clicking on pop-up ads',
              'Don\'t enter personal information on suspicious sites',
              'Be careful with social media links'
            ]
          }
        ]
      }
    ])

    return {
      securityChecks,
      securityCategories
    }
  }
}
</script>
