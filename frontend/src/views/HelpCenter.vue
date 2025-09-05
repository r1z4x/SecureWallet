<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader type="authenticated" />
    
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <h1 class="text-4xl font-bold text-gray-900 mb-4">Help Center</h1>
        <p class="text-xl text-gray-600">Find answers to your questions and get support</p>
      </div>

      <!-- Search -->
      <div class="mb-8">
        <div class="relative max-w-2xl mx-auto">
          <input
            v-model="searchQuery"
            type="text"
            class="w-full pl-12 pr-4 py-3 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            placeholder="Search for help topics..."
            @input="filterFAQs"
          >
          <div class="absolute inset-y-0 left-0 flex items-center pl-4">
            <i class="fas fa-search text-gray-400"></i>
          </div>
        </div>
      </div>

      <!-- Quick Links -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        <div class="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
          <i class="fas fa-question-circle text-3xl text-blue-600 mb-4"></i>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">FAQ</h3>
          <p class="text-gray-600 text-sm mb-4">Find answers to common questions</p>
          <router-link to="/faq" class="text-blue-600 hover:text-blue-800 font-medium">Browse FAQ →</router-link>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
          <i class="fas fa-book text-3xl text-green-600 mb-4"></i>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">User Guide</h3>
          <p class="text-gray-600 text-sm mb-4">Step-by-step instructions</p>
          <router-link to="/user-guide" class="text-green-600 hover:text-green-800 font-medium">Read Guide →</router-link>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
          <i class="fas fa-shield-alt text-3xl text-purple-600 mb-4"></i>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">Security Tips</h3>
          <p class="text-gray-600 text-sm mb-4">Keep your account secure</p>
          <router-link to="/security-tips" class="text-purple-600 hover:text-purple-800 font-medium">Learn More →</router-link>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6 text-center hover:shadow-lg transition-shadow">
          <i class="fas fa-headset text-3xl text-orange-600 mb-4"></i>
          <h3 class="text-lg font-semibold text-gray-900 mb-2">Contact Support</h3>
          <p class="text-gray-600 text-sm mb-4">Can't find what you're looking for? Contact our support team</p>
          <router-link to="/support" class="text-orange-600 hover:text-orange-800 font-medium">
            Contact Us →
          </router-link>
        </div>
      </div>

      <!-- FAQ Sections -->
      <div class="space-y-8">
        <!-- Getting Started -->
        <div id="getting-started" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
            <i class="fas fa-wallet text-primary-600 mr-3"></i>
            Getting Started
          </h2>
          
          <div class="space-y-4">
            <div
              v-for="faq in filteredFAQs.gettingStarted"
              :key="faq.id"
              class="border border-gray-200 rounded-lg"
            >
              <button
                @click="toggleFAQ(faq.id)"
                class="w-full px-6 py-4 text-left flex justify-between items-center hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 rounded-lg"
              >
                <span class="font-medium text-gray-900">{{ faq.question }}</span>
                <i :class="['fas', expandedFAQs.includes(faq.id) ? 'fa-chevron-up' : 'fa-chevron-down', 'text-gray-400']"></i>
              </button>
              <div v-if="expandedFAQs.includes(faq.id)" class="px-6 pb-4">
                <p class="text-gray-600">{{ faq.answer }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Security -->
        <div id="security" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
            <i class="fas fa-shield-alt text-green-600 mr-3"></i>
            Security & Privacy
          </h2>
          
          <div class="space-y-4">
            <div
              v-for="faq in filteredFAQs.security"
              :key="faq.id"
              class="border border-gray-200 rounded-lg"
            >
              <button
                @click="toggleFAQ(faq.id)"
                class="w-full px-6 py-4 text-left flex justify-between items-center hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 rounded-lg"
              >
                <span class="font-medium text-gray-900">{{ faq.question }}</span>
                <i :class="['fas', expandedFAQs.includes(faq.id) ? 'fa-chevron-up' : 'fa-chevron-down', 'text-gray-400']"></i>
              </button>
              <div v-if="expandedFAQs.includes(faq.id)" class="px-6 pb-4">
                <p class="text-gray-600">{{ faq.answer }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Transactions -->
        <div id="transactions" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
            <i class="fas fa-exchange-alt text-blue-600 mr-3"></i>
            Transactions & Transfers
          </h2>
          
          <div class="space-y-4">
            <div
              v-for="faq in filteredFAQs.transactions"
              :key="faq.id"
              class="border border-gray-200 rounded-lg"
            >
              <button
                @click="toggleFAQ(faq.id)"
                class="w-full px-6 py-4 text-left flex justify-between items-center hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 rounded-lg"
              >
                <span class="font-medium text-gray-900">{{ faq.question }}</span>
                <i :class="['fas', expandedFAQs.includes(faq.id) ? 'fa-chevron-up' : 'fa-chevron-down', 'text-gray-400']"></i>
              </button>
              <div v-if="expandedFAQs.includes(faq.id)" class="px-6 pb-4">
                <p class="text-gray-600">{{ faq.answer }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Account Management -->
        <div id="account" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
            <i class="fas fa-user-cog text-purple-600 mr-3"></i>
            Account Management
          </h2>
          
          <div class="space-y-4">
            <div
              v-for="faq in filteredFAQs.account"
              :key="faq.id"
              class="border border-gray-200 rounded-lg"
            >
              <button
                @click="toggleFAQ(faq.id)"
                class="w-full px-6 py-4 text-left flex justify-between items-center hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 rounded-lg"
              >
                <span class="font-medium text-gray-900">{{ faq.question }}</span>
                <i :class="['fas', expandedFAQs.includes(faq.id) ? 'fa-chevron-up' : 'fa-chevron-down', 'text-gray-400']"></i>
              </button>
              <div v-if="expandedFAQs.includes(faq.id)" class="px-6 pb-4">
                <p class="text-gray-600">{{ faq.answer }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Contact Support -->
      <div class="mt-12 bg-gradient-to-r from-primary-600 to-primary-700 rounded-lg shadow-lg p-8 text-center text-white">
        <h3 class="text-2xl font-bold mb-4">Still Need Help?</h3>
        <p class="text-primary-100 mb-6">Our support team is here to help you with any questions or issues</p>
        <router-link
          to="/support"
          class="inline-flex items-center px-6 py-3 bg-white text-primary-600 font-medium rounded-lg hover:bg-gray-50 transition-colors"
        >
          <i class="fas fa-headset mr-2"></i>
          Contact Support
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import AppHeader from '@/components/AppHeader.vue'

export default {
  name: 'HelpCenter',
  components: {
    AppHeader
  },
  setup() {
    const searchQuery = ref('')
    const expandedFAQs = ref([])

    // FAQ Data
    const faqs = {
      gettingStarted: [
        {
          id: 'gs1',
          question: 'How do I create an account?',
          answer: 'To create an account, click on the "Sign Up" button on the homepage. You\'ll need to provide your email address, choose a username, and create a secure password. After registration, you\'ll receive a confirmation email to verify your account.'
        },
        {
          id: 'gs2',
          question: 'How do I add money to my wallet?',
          answer: 'You can add money to your wallet by going to the Dashboard and clicking on the "Deposit" button. Enter the amount you want to deposit and add an optional description. The funds will be added to your wallet immediately.'
        },
        {
          id: 'gs3',
          question: 'What is the minimum deposit amount?',
          answer: 'The minimum deposit amount is $1.00. There are no maximum limits for deposits, but large amounts may require additional verification for security purposes.'
        },
        {
          id: 'gs4',
          question: 'How do I enable two-factor authentication?',
          answer: 'Go to your Profile page and click on the "Two-Factor Authentication" section. Follow the setup instructions to scan the QR code with your authenticator app and enter the verification code to enable 2FA.'
        }
      ],
      security: [
        {
          id: 'sec1',
          question: 'How secure is my account?',
          answer: 'Your account is protected by industry-standard security measures including encryption, secure authentication, and regular security audits. We also recommend enabling two-factor authentication for additional protection.'
        },
        {
          id: 'sec2',
          question: 'What should I do if I suspect unauthorized activity?',
          answer: 'If you notice any suspicious activity, immediately change your password, enable two-factor authentication if not already enabled, and contact our support team. We\'ll investigate the issue and take appropriate action.'
        },
        {
          id: 'sec3',
          question: 'How do I change my password?',
          answer: 'Go to your Profile page and click on "Change Password". Enter your current password and your new password. Make sure your new password is strong and unique.'
        },
        {
          id: 'sec4',
          question: 'Is my personal information safe?',
          answer: 'Yes, we take your privacy seriously. Your personal information is encrypted and stored securely. We never share your information with third parties without your explicit consent.'
        }
      ],
      transactions: [
        {
          id: 'tx1',
          question: 'How do I send money to another user?',
          answer: 'Go to the Transfer page and enter the recipient\'s email address or user ID. Enter the amount you want to send and add an optional description. Review the transfer fee and total amount before confirming the transfer.'
        },
        {
          id: 'tx2',
          question: 'What are the transfer fees?',
          answer: 'Transfer fees are 1% of the transfer amount, with a minimum fee of $1.00 and a maximum fee of $50.00. For example, a $100 transfer will have a $1 fee, while a $1000 transfer will have a $10 fee.'
        },
        {
          id: 'tx3',
          question: 'How long do transfers take?',
          answer: 'Transfers between users are processed instantly. You\'ll receive a confirmation email and the recipient will see the funds in their wallet immediately.'
        },
        {
          id: 'tx4',
          question: 'What are the transfer limits?',
          answer: 'Daily transfer limit is $10,000 and monthly limit is $50,000. Minimum transfer amount is $1.00. These limits help ensure security and prevent fraud.'
        },
        {
          id: 'tx5',
          question: 'Can I cancel a transfer?',
          answer: 'Once a transfer is completed, it cannot be cancelled. However, if you notice an error, contact our support team immediately and we\'ll investigate the issue.'
        }
      ],
      account: [
        {
          id: 'acc1',
          question: 'How do I update my profile information?',
          answer: 'Go to your Profile page and click on "Edit Profile". You can update your username, email address, and other personal information. Some changes may require verification.'
        },
        {
          id: 'acc2',
          question: 'Can I have multiple wallets?',
          answer: 'Currently, each user can have one main wallet. We\'re working on adding support for multiple wallets in future updates.'
        },
        {
          id: 'acc3',
          question: 'How do I delete my account?',
          answer: 'To delete your account, contact our support team. Please note that account deletion is permanent and all data will be permanently removed. Make sure to withdraw any remaining funds before deletion.'
        },
        {
          id: 'acc4',
          question: 'What currencies are supported?',
          answer: 'Currently, we support USD (US Dollar) as the primary currency. We\'re working on adding support for additional currencies in the future.'
        }
      ]
    }

    const filteredFAQs = computed(() => {
      if (!searchQuery.value) return faqs

      const query = searchQuery.value.toLowerCase()
      const filtered = {}

      Object.keys(faqs).forEach(category => {
        filtered[category] = faqs[category].filter(faq =>
          faq.question.toLowerCase().includes(query) ||
          faq.answer.toLowerCase().includes(query)
        )
      })

      return filtered
    })

    const toggleFAQ = (faqId) => {
      const index = expandedFAQs.value.indexOf(faqId)
      if (index > -1) {
        expandedFAQs.value.splice(index, 1)
      } else {
        expandedFAQs.value.push(faqId)
      }
    }

    const scrollToSection = (sectionId) => {
      const element = document.getElementById(sectionId)
      if (element) {
        element.scrollIntoView({ behavior: 'smooth' })
      }
    }

    const filterFAQs = () => {
      // This function is called on input, but the filtering is handled by computed property
    }

    onMounted(() => {
      // Expand first FAQ in each section by default
      Object.keys(faqs).forEach(category => {
        if (faqs[category].length > 0) {
          expandedFAQs.value.push(faqs[category][0].id)
        }
      })
    })

    return {
      searchQuery,
      expandedFAQs,
      filteredFAQs,
      toggleFAQ,
      scrollToSection,
      filterFAQs
    }
  }
}
</script>
