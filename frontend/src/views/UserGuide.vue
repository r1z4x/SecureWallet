<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader type="authenticated" />
    
    <div class="container mx-auto px-4 py-8">
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-gray-900 mb-4">User Guide</h1>
        <p class="text-lg text-gray-600 max-w-3xl mx-auto">
          Learn how to use SecureWallet effectively. This comprehensive guide will walk you through all the features and help you get the most out of your digital wallet.
        </p>
      </div>

      <!-- Table of Contents -->
      <div class="max-w-4xl mx-auto mb-8">
        <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <h2 class="text-2xl font-semibold text-gray-900 mb-4">Table of Contents</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div
              v-for="(section, index) in sections"
              :key="section.id"
              class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors cursor-pointer"
              @click="scrollToSection(section.id)"
            >
              <div class="w-8 h-8 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center text-sm font-medium mr-3">
                {{ index + 1 }}
              </div>
              <div>
                <h3 class="font-medium text-gray-900">{{ section.title }}</h3>
                <p class="text-sm text-gray-500">{{ section.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Guide Sections -->
      <div class="max-w-4xl mx-auto space-y-8">
        <div
          v-for="(section, index) in sections"
          :key="section.id"
          :id="section.id"
          class="bg-white rounded-lg shadow-sm border border-gray-200"
        >
          <div class="p-6 border-b border-gray-200">
            <div class="flex items-center mb-4">
              <div class="w-10 h-10 bg-blue-100 text-blue-600 rounded-full flex items-center justify-center text-lg font-medium mr-4">
                {{ index + 1 }}
              </div>
              <h2 class="text-2xl font-semibold text-gray-900">{{ section.title }}</h2>
            </div>
            <p class="text-gray-600">{{ section.description }}</p>
          </div>
          
          <div class="p-6">
            <div v-for="(step, stepIndex) in section.steps" :key="stepIndex" class="mb-8 last:mb-0">
              <div class="flex items-start">
                <div class="w-8 h-8 bg-green-100 text-green-600 rounded-full flex items-center justify-center text-sm font-medium mr-4 mt-1 flex-shrink-0">
                  {{ stepIndex + 1 }}
                </div>
                <div class="flex-1">
                  <h3 class="text-lg font-medium text-gray-900 mb-2">{{ step.title }}</h3>
                  <p class="text-gray-600 mb-4">{{ step.description }}</p>
                  
                  <!-- Screenshot Placeholder -->
                  <div v-if="step.screenshot" class="mb-4">
                    <div class="bg-gray-100 border-2 border-dashed border-gray-300 rounded-lg p-8 text-center">
                      <i class="fas fa-image text-4xl text-gray-400 mb-2"></i>
                      <p class="text-gray-500">Screenshot: {{ step.screenshot }}</p>
                    </div>
                  </div>
                  
                  <!-- Tips -->
                  <div v-if="step.tips && step.tips.length > 0" class="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <h4 class="font-medium text-blue-900 mb-2 flex items-center">
                      <i class="fas fa-lightbulb mr-2"></i>
                      Pro Tips
                    </h4>
                    <ul class="text-sm text-blue-800 space-y-1">
                      <li v-for="(tip, tipIndex) in step.tips" :key="tipIndex" class="flex items-start">
                        <i class="fas fa-check-circle text-blue-600 mr-2 mt-0.5 flex-shrink-0"></i>
                        {{ tip }}
                      </li>
                    </ul>
                  </div>
                  
                  <!-- Warnings -->
                  <div v-if="step.warnings && step.warnings.length > 0" class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mt-4">
                    <h4 class="font-medium text-yellow-900 mb-2 flex items-center">
                      <i class="fas fa-exclamation-triangle mr-2"></i>
                      Important Notes
                    </h4>
                    <ul class="text-sm text-yellow-800 space-y-1">
                      <li v-for="(warning, warningIndex) in step.warnings" :key="warningIndex" class="flex items-start">
                        <i class="fas fa-info-circle text-yellow-600 mr-2 mt-0.5 flex-shrink-0"></i>
                        {{ warning }}
                      </li>
                    </ul>
                  </div>
                </div>
              </div>
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
              to="/security-tips"
              class="inline-flex items-center px-6 py-3 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
            >
              <i class="fas fa-shield-alt mr-2"></i>
              Security Tips
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue'
import AppHeader from '@/components/AppHeader.vue'

export default {
  name: 'UserGuide',
  components: {
    AppHeader
  },
  setup() {
    const sections = ref([
      {
        id: 'getting-started',
        title: 'Getting Started',
        description: 'Learn the basics of setting up and navigating your SecureWallet account.',
        steps: [
          {
            title: 'Create Your Account',
            description: 'Sign up for a SecureWallet account by providing your email address, username, and a strong password. You\'ll receive a verification email to activate your account.',
            screenshot: 'Account registration form',
            tips: [
              'Use a strong, unique password that you don\'t use elsewhere',
              'Keep your email address up to date for important notifications',
              'Choose a username that\'s easy to remember but hard to guess'
            ],
            warnings: [
              'Never share your login credentials with anyone',
              'Enable two-factor authentication for enhanced security'
            ]
          },
          {
            title: 'Complete Your Profile',
            description: 'After logging in, go to your Profile page to add additional information and customize your account settings.',
            screenshot: 'Profile settings page',
            tips: [
              'Add a profile picture to personalize your account',
              'Review and update your contact information regularly',
              'Set up your preferred currency and timezone'
            ]
          },
          {
            title: 'Enable Two-Factor Authentication',
            description: 'Add an extra layer of security to your account by enabling two-factor authentication (2FA).',
            screenshot: '2FA setup process',
            tips: [
              'Use an authenticator app like Google Authenticator or Authy',
              'Save your backup codes in a secure location',
              'Test your 2FA setup before logging out'
            ],
            warnings: [
              'Keep your backup codes safe - you\'ll need them if you lose your phone',
              'Don\'t share your 2FA codes with anyone'
            ]
          }
        ]
      },
      {
        id: 'wallet-management',
        title: 'Wallet Management',
        description: 'Learn how to manage your wallet, check balances, and understand your financial overview.',
        steps: [
          {
            title: 'Understanding Your Dashboard',
            description: 'Your dashboard provides a comprehensive overview of your account, including current balance, recent transactions, and quick action buttons.',
            screenshot: 'Dashboard overview',
            tips: [
              'Check your dashboard regularly to monitor your account activity',
              'Use the quick action buttons for common tasks',
              'Review recent transactions for any suspicious activity'
            ]
          },
          {
            title: 'Adding Money to Your Wallet',
            description: 'You can add funds to your wallet using various payment methods including credit cards, debit cards, or bank transfers.',
            screenshot: 'Deposit page',
            tips: [
              'Choose the payment method that works best for you',
              'Be aware of any processing fees for different payment methods',
              'Keep your payment information up to date'
            ],
            warnings: [
              'Only use payment methods that belong to you',
              'Be cautious of any unexpected charges'
            ]
          },
          {
            title: 'Withdrawing Funds',
            description: 'Withdraw money from your wallet to your bank account or other payment methods.',
            screenshot: 'Withdrawal page',
            tips: [
              'Verify your withdrawal details before confirming',
              'Note that withdrawals may take 1-3 business days',
              'Keep track of your withdrawal history'
            ],
            warnings: [
              'Ensure your bank account information is correct',
              'Be aware of any withdrawal limits or fees'
            ]
          }
        ]
      },
      {
        id: 'making-transfers',
        title: 'Making Transfers',
        description: 'Learn how to send money to other users, understand fees, and manage your transfers.',
        steps: [
          {
            title: 'Sending Money to Another User',
            description: 'Transfer money to other SecureWallet users by entering their username or email address.',
            screenshot: 'Transfer form',
            tips: [
              'Double-check the recipient\'s information before sending',
              'Add a note to your transfer for better record-keeping',
              'Use the search function to find users quickly'
            ],
            warnings: [
              'Transfers cannot be cancelled once completed',
              'Make sure you have sufficient funds including the transfer fee'
            ]
          },
          {
            title: 'Understanding Transfer Fees',
            description: 'Learn about the fees associated with transfers and how they are calculated.',
            screenshot: 'Fee calculation example',
            tips: [
              'Transfer fees are 1% of the transaction amount',
              'Minimum fee is $1, maximum fee is $50',
              'Fees are clearly displayed before you confirm the transfer'
            ]
          },
          {
            title: 'Managing Your Transfer History',
            description: 'View and manage your transfer history, including pending, completed, and failed transfers.',
            screenshot: 'Transaction history page',
            tips: [
              'Use filters to find specific transactions',
              'Export your transaction history for record-keeping',
              'Set up notifications for large transfers'
            ]
          }
        ]
      },
      {
        id: 'security-features',
        title: 'Security Features',
        description: 'Understand the security features available to protect your account and transactions.',
        steps: [
          {
            title: 'Login History Monitoring',
            description: 'Regularly review your login history to detect any unauthorized access to your account.',
            screenshot: 'Login history page',
            tips: [
              'Check your login history weekly',
              'Look for any logins from unfamiliar locations or devices',
              'Report any suspicious activity immediately'
            ],
            warnings: [
              'If you see an unfamiliar login, change your password immediately',
              'Contact support if you suspect unauthorized access'
            ]
          },
          {
            title: 'Password Security',
            description: 'Learn best practices for creating and maintaining secure passwords.',
            screenshot: 'Password strength indicator',
            tips: [
              'Use a mix of uppercase, lowercase, numbers, and symbols',
              'Avoid using personal information in your password',
              'Change your password regularly'
            ],
            warnings: [
              'Never reuse passwords from other accounts',
              'Don\'t share your password with anyone'
            ]
          },
          {
            title: 'Account Recovery',
            description: 'Set up account recovery options to regain access if you lose your credentials.',
            screenshot: 'Account recovery settings',
            tips: [
              'Keep your recovery email up to date',
              'Save your backup codes in a secure location',
              'Test your recovery process periodically'
            ]
          }
        ]
      },
      {
        id: 'troubleshooting',
        title: 'Troubleshooting',
        description: 'Common issues and how to resolve them.',
        steps: [
          {
            title: 'Login Problems',
            description: 'What to do if you can\'t log into your account.',
            screenshot: 'Login error messages',
            tips: [
              'Check that your email and password are correct',
              'Ensure your caps lock is off',
              'Try resetting your password if needed'
            ],
            warnings: [
              'Don\'t attempt to log in too many times - your account may be temporarily locked',
              'Contact support if you continue to have issues'
            ]
          },
          {
            title: 'Transfer Issues',
            description: 'Common problems with transfers and how to resolve them.',
            screenshot: 'Transfer error handling',
            tips: [
              'Check that you have sufficient funds',
              'Verify the recipient\'s information',
              'Ensure you\'re not exceeding transfer limits'
            ],
            warnings: [
              'Contact support immediately if money is deducted but not received',
              'Keep records of all transfer attempts'
            ]
          },
          {
            title: 'Getting Help',
            description: 'How to get help when you need it.',
            screenshot: 'Support contact options',
            tips: [
              'Check the FAQ section first for quick answers',
              'Use the support ticket system for detailed issues',
              'Include relevant details when contacting support'
            ]
          }
        ]
      }
    ])

    const scrollToSection = (sectionId) => {
      const element = document.getElementById(sectionId)
      if (element) {
        element.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }
    }

    return {
      sections,
      scrollToSection
    }
  }
}
</script>
