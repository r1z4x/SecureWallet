<template>
  <div class="min-h-screen bg-gray-50">
    <Navigation />
    
    <div class="container mx-auto px-4 py-8">
      <!-- Header -->
      <div class="text-center mb-8">
        <h1 class="text-4xl font-bold text-gray-900 mb-4">Frequently Asked Questions</h1>
        <p class="text-lg text-gray-600 max-w-2xl mx-auto">
          Find answers to the most common questions about SecureWallet. Can't find what you're looking for? 
          <router-link to="/support" class="text-blue-600 hover:text-blue-800 font-medium">Contact our support team</router-link>.
        </p>
      </div>

      <!-- Search Bar -->
      <div class="max-w-2xl mx-auto mb-8">
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search FAQs..."
            class="w-full pl-12 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <i class="fas fa-search text-gray-400"></i>
          </div>
        </div>
      </div>

      <!-- Quick Navigation -->
      <div class="flex flex-wrap justify-center gap-4 mb-8">
        <button
          v-for="category in categories"
          :key="category.id"
          @click="scrollToCategory(category.id)"
          class="px-4 py-2 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
        >
          {{ category.name }}
        </button>
      </div>

      <!-- FAQ Categories -->
      <div class="max-w-4xl mx-auto space-y-8">
        <div
          v-for="category in filteredCategories"
          :key="category.id"
          :id="category.id"
          class="bg-white rounded-lg shadow-sm border border-gray-200"
        >
          <div class="p-6 border-b border-gray-200">
            <h2 class="text-2xl font-semibold text-gray-900 flex items-center">
              <i :class="category.icon" class="mr-3 text-blue-600"></i>
              {{ category.name }}
            </h2>
            <p class="text-gray-600 mt-2">{{ category.description }}</p>
          </div>
          
          <div class="divide-y divide-gray-200">
            <div
              v-for="(faq, index) in category.faqs"
              :key="index"
              class="p-6"
            >
              <button
                @click="toggleFAQ(category.id, index)"
                class="w-full text-left flex justify-between items-start"
              >
                <h3 class="text-lg font-medium text-gray-900 pr-4">
                  {{ faq.question }}
                </h3>
                <i
                  :class="[
                    'fas fa-chevron-down text-gray-400 transition-transform duration-200',
                    { 'rotate-180': openFAQs[`${category.id}-${index}`] }
                  ]"
                ></i>
              </button>
              <div
                v-show="openFAQs[`${category.id}-${index}`]"
                class="mt-4 text-gray-600 leading-relaxed"
                v-html="faq.answer"
              ></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Still Need Help -->
      <div class="max-w-2xl mx-auto mt-12 text-center">
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-6">
          <h3 class="text-xl font-semibold text-blue-900 mb-2">Still Need Help?</h3>
          <p class="text-blue-700 mb-4">
            Can't find the answer you're looking for? Our support team is here to help.
          </p>
          <div class="flex flex-col sm:flex-row gap-4 justify-center">
            <router-link
              to="/support"
              class="inline-flex items-center px-6 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              <i class="fas fa-headset mr-2"></i>
              Contact Support
            </router-link>
            <router-link
              to="/help"
              class="inline-flex items-center px-6 py-3 bg-white border border-blue-600 text-blue-600 rounded-lg hover:bg-blue-50 transition-colors"
            >
              <i class="fas fa-arrow-left mr-2"></i>
              Back to Help Center
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import Navigation from '@/components/Navigation.vue'

export default {
  name: 'FAQ',
  components: {
    Navigation
  },
  setup() {
    const searchQuery = ref('')
    const openFAQs = ref({})

    const categories = ref([
      {
        id: 'account',
        name: 'Account & Security',
        icon: 'fas fa-shield-alt',
        description: 'Questions about account management, security features, and authentication.',
        faqs: [
          {
            question: 'How do I enable two-factor authentication?',
            answer: 'Go to your Profile settings and click on "Security". You can enable 2FA using an authenticator app like Google Authenticator or Authy. Scan the QR code with your app and enter the verification code to complete setup.'
          },
          {
            question: 'What should I do if I forget my password?',
            answer: 'Click on "Forgot Password" on the login page. Enter your email address and we\'ll send you a password reset link. Make sure to check your spam folder if you don\'t receive the email immediately.'
          },
          {
            question: 'How can I change my email address?',
            answer: 'Go to Profile settings and click on "Edit Profile". You can update your email address there. You\'ll need to verify the new email address before the change takes effect.'
          },
          {
            question: 'Is my personal information secure?',
            answer: 'Yes, we use industry-standard encryption and security measures to protect your personal information. All data is encrypted in transit and at rest. We never share your personal information with third parties without your consent.'
          },
          {
            question: 'How do I delete my account?',
            answer: 'Go to Profile settings and scroll down to "Account Management". Click on "Delete Account" and follow the confirmation process. Please note that this action is irreversible and will permanently delete all your data.'
          }
        ]
      },
      {
        id: 'transactions',
        name: 'Transactions & Payments',
        icon: 'fas fa-exchange-alt',
        description: 'Information about making transfers, transaction fees, and payment processing.',
        faqs: [
          {
            question: 'How do I send money to another user?',
            answer: 'Go to the Transfer page and enter the recipient\'s username or email. Enter the amount and review the transfer fee. Click "Send Money" to complete the transaction. The recipient will receive the money instantly.'
          },
          {
            question: 'What are the transfer fees?',
            answer: 'Transfer fees are 1% of the transaction amount, with a minimum of $1 and maximum of $50. For example, a $100 transfer will have a $1 fee, while a $1000 transfer will have a $10 fee.'
          },
          {
            question: 'How long do transfers take?',
            answer: 'Transfers between SecureWallet users are processed instantly. The recipient will see the funds in their account immediately after the transfer is completed.'
          },
          {
            question: 'Can I cancel a transfer?',
            answer: 'Unfortunately, transfers cannot be cancelled once they are completed. Please double-check all details before confirming your transfer. If you sent money to the wrong person, contact our support team immediately.'
          },
          {
            question: 'What currencies do you support?',
            answer: 'We currently support USD, EUR, GBP, CAD, AUD, CHF, and CNY. You can view your balance in different currencies on your Wallet page.'
          },
          {
            question: 'How do I view my transaction history?',
            answer: 'Go to the Transactions page to view your complete transaction history. You can filter transactions by type, status, and date range. You can also search for specific transactions.'
          }
        ]
      },
      {
        id: 'wallet',
        name: 'Wallet & Balance',
        icon: 'fas fa-wallet',
        description: 'Questions about wallet management, balance updates, and currency conversion.',
        faqs: [
          {
            question: 'How do I add money to my wallet?',
            answer: 'Go to the Deposit page and select your preferred payment method. You can use a credit card, debit card, or bank transfer. Follow the instructions to complete your deposit.'
          },
          {
            question: 'How do I withdraw money from my wallet?',
            answer: 'Go to your Wallet page and click on "Withdraw". Enter the amount and select your withdrawal method. Withdrawals to bank accounts typically take 1-3 business days.'
          },
          {
            question: 'What are the minimum and maximum transaction limits?',
            answer: 'Minimum transfer amount is $1. Maximum transfer amount is $10,000 per transaction. Daily transfer limit is $50,000. These limits may vary based on your account verification level.'
          },
          {
            question: 'How do I check my account balance?',
            answer: 'Your current balance is displayed on the Dashboard and Wallet pages. You can also view your balance in different currencies and see recent transactions.'
          },
          {
            question: 'What happens if I have insufficient funds?',
            answer: 'If you don\'t have enough funds for a transfer, the transaction will be declined. Make sure to check your balance before making transfers and consider the transfer fee in your calculations.'
          }
        ]
      },
      {
        id: 'technical',
        name: 'Technical Issues',
        icon: 'fas fa-tools',
        description: 'Help with technical problems, app issues, and troubleshooting.',
        faqs: [
          {
            question: 'The app is not loading properly. What should I do?',
            answer: 'Try refreshing the page or clearing your browser cache. If the problem persists, try using a different browser or device. Contact support if the issue continues.'
          },
          {
            question: 'I\'m having trouble logging in. What can I do?',
            answer: 'Make sure you\'re using the correct email and password. If you have 2FA enabled, ensure your authenticator app is working properly. Try resetting your password if needed.'
          },
          {
            question: 'The page is showing an error. What does this mean?',
            answer: 'Error messages usually indicate a temporary issue with our servers or your internet connection. Try refreshing the page or waiting a few minutes before trying again.'
          },
          {
            question: 'How do I update my browser?',
            answer: 'For Chrome: Click the three dots menu → Help → About Google Chrome. For Firefox: Click the menu button → Help → About Firefox. For Safari: Click Safari → About Safari.'
          },
          {
            question: 'Is the app compatible with mobile devices?',
            answer: 'Yes, our web app is fully responsive and works on all mobile devices and tablets. For the best experience, we recommend using the latest version of your mobile browser.'
          }
        ]
      },
      {
        id: 'privacy',
        name: 'Privacy & Legal',
        icon: 'fas fa-user-shield',
        description: 'Information about privacy policies, terms of service, and legal matters.',
        faqs: [
          {
            question: 'How do you protect my privacy?',
            answer: 'We follow strict privacy policies and never share your personal information with third parties without your consent. All data is encrypted and stored securely. You can review our full Privacy Policy for more details.'
          },
          {
            question: 'What information do you collect?',
            answer: 'We collect information necessary to provide our services, including your name, email, phone number, and transaction data. We do not collect unnecessary personal information.'
          },
          {
            question: 'Can I request a copy of my data?',
            answer: 'Yes, you can request a copy of your personal data by contacting our support team. We\'ll provide you with all the information we have about your account within 30 days.'
          },
          {
            question: 'How do I report a security concern?',
            answer: 'If you notice any suspicious activity or have security concerns, contact our support team immediately. You can also report issues through the Security section in your Profile settings.'
          },
          {
            question: 'What are your terms of service?',
            answer: 'Our Terms of Service outline the rules and guidelines for using SecureWallet. You can read the full terms on our Terms of Service page. By using our service, you agree to these terms.'
          }
        ]
      }
    ])

    const filteredCategories = computed(() => {
      if (!searchQuery.value) return categories.value

      return categories.value.map(category => ({
        ...category,
        faqs: category.faqs.filter(faq =>
          faq.question.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
          faq.answer.toLowerCase().includes(searchQuery.value.toLowerCase())
        )
      })).filter(category => category.faqs.length > 0)
    })

    const toggleFAQ = (categoryId, index) => {
      const key = `${categoryId}-${index}`
      openFAQs.value[key] = !openFAQs.value[key]
    }

    const scrollToCategory = (categoryId) => {
      const element = document.getElementById(categoryId)
      if (element) {
        element.scrollIntoView({ behavior: 'smooth', block: 'start' })
      }
    }

    onMounted(() => {
      // Open first FAQ in each category by default
      categories.value.forEach(category => {
        if (category.faqs.length > 0) {
          openFAQs.value[`${category.id}-0`] = true
        }
      })
    })

    return {
      searchQuery,
      categories,
      filteredCategories,
      openFAQs,
      toggleFAQ,
      scrollToCategory
    }
  }
}
</script>
