<template>
  <div class="min-h-screen bg-gray-50">
    <Navigation />
    
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="text-center mb-12">
        <h1 class="text-4xl font-bold text-gray-900 mb-4">Terms of Service</h1>
        <p class="text-xl text-gray-600">Last updated: {{ lastUpdated }}</p>
      </div>

      <!-- Table of Contents -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 class="text-2xl font-bold text-gray-900 mb-4">Table of Contents</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-2">
            <a href="#acceptance" class="block text-primary-600 hover:text-primary-800">1. Acceptance of Terms</a>
            <a href="#description" class="block text-primary-600 hover:text-primary-800">2. Description of Service</a>
            <a href="#eligibility" class="block text-primary-600 hover:text-primary-800">3. Eligibility</a>
            <a href="#registration" class="block text-primary-600 hover:text-primary-800">4. Registration and Account</a>
            <a href="#security" class="block text-primary-600 hover:text-primary-800">5. Security</a>
            <a href="#fees" class="block text-primary-600 hover:text-primary-800">6. Fees and Charges</a>
          </div>
          <div class="space-y-2">
            <a href="#transactions" class="block text-primary-600 hover:text-primary-800">7. Transactions</a>
            <a href="#prohibited" class="block text-primary-600 hover:text-primary-800">8. Prohibited Activities</a>
            <a href="#privacy" class="block text-primary-600 hover:text-primary-800">9. Privacy Policy</a>
            <a href="#liability" class="block text-primary-600 hover:text-primary-800">10. Limitation of Liability</a>
            <a href="#termination" class="block text-primary-600 hover:text-primary-800">11. Termination</a>
            <a href="#changes" class="block text-primary-600 hover:text-primary-800">12. Changes to Terms</a>
          </div>
        </div>
      </div>

      <!-- Terms Content -->
      <div class="space-y-8">
        <!-- Acceptance of Terms -->
        <div id="acceptance" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">1. Acceptance of Terms</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              By accessing and using SecureWallet ("the Service"), you accept and agree to be bound by the terms and provision of this agreement. If you do not agree to abide by the above, please do not use this service.
            </p>
            <p class="text-gray-700">
              These Terms of Service ("Terms") govern your use of our website and services. By using our Service, you agree to these Terms in full.
            </p>
          </div>
        </div>

        <!-- Description of Service -->
        <div id="description" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">2. Description of Service</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              SecureWallet is a digital wallet service that allows users to:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li>Create and manage digital wallets</li>
              <li>Send and receive money between users</li>
              <li>Deposit funds into their wallet</li>
              <li>View transaction history</li>
              <li>Enable two-factor authentication for enhanced security</li>
            </ul>
            <p class="text-gray-700">
              The Service is provided "as is" and we reserve the right to modify, suspend, or discontinue the Service at any time.
            </p>
          </div>
        </div>

        <!-- Eligibility -->
        <div id="eligibility" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">3. Eligibility</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              To use our Service, you must:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li>Be at least 18 years old</li>
              <li>Have the legal capacity to enter into binding agreements</li>
              <li>Provide accurate and complete information during registration</li>
              <li>Comply with all applicable laws and regulations</li>
            </ul>
            <p class="text-gray-700">
              We reserve the right to refuse service to anyone for any reason at our discretion.
            </p>
          </div>
        </div>

        <!-- Registration and Account -->
        <div id="registration" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">4. Registration and Account</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              When you create an account, you must provide accurate, current, and complete information. You are responsible for:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li>Maintaining the security of your account credentials</li>
              <li>All activities that occur under your account</li>
              <li>Notifying us immediately of any unauthorized use</li>
              <li>Keeping your contact information up to date</li>
            </ul>
            <p class="text-gray-700">
              You may not share your account credentials with anyone else or allow others to access your account.
            </p>
          </div>
        </div>

        <!-- Security -->
        <div id="security" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">5. Security</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              We implement industry-standard security measures to protect your account and transactions. However, you are also responsible for:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li>Using a strong, unique password</li>
              <li>Enabling two-factor authentication when available</li>
              <li>Keeping your device and software updated</li>
              <li>Not accessing your account from public or unsecured networks</li>
            </ul>
            <p class="text-gray-700">
              We are not liable for any losses resulting from your failure to maintain adequate security measures.
            </p>
          </div>
        </div>

        <!-- Fees and Charges -->
        <div id="fees" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">6. Fees and Charges</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              The following fees apply to our Service:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li><strong>Transfer Fee:</strong> 1% of transfer amount (minimum $1.00, maximum $50.00)</li>
              <li><strong>Deposit Fee:</strong> No fee for deposits</li>
              <li><strong>Account Maintenance:</strong> No monthly or annual fees</li>
            </ul>
            <p class="text-gray-700 mb-4">
              All fees are clearly displayed before you complete any transaction. Fees are deducted from your wallet balance.
            </p>
            <p class="text-gray-700">
              We reserve the right to modify our fee structure with 30 days' notice to users.
            </p>
          </div>
        </div>

        <!-- Transactions -->
        <div id="transactions" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">7. Transactions</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              All transactions are subject to the following terms:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li>Transactions are processed immediately but may be subject to review</li>
              <li>We may hold or reverse transactions for security or compliance reasons</li>
              <li>You are responsible for ensuring accurate recipient information</li>
              <li>Completed transactions cannot be cancelled or reversed</li>
              <li>Daily transfer limit: $10,000</li>
              <li>Monthly transfer limit: $50,000</li>
            </ul>
            <p class="text-gray-700">
              We are not responsible for losses due to incorrect recipient information provided by you.
            </p>
          </div>
        </div>

        <!-- Prohibited Activities -->
        <div id="prohibited" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">8. Prohibited Activities</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              You may not use our Service for any illegal or unauthorized purpose. Prohibited activities include:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li>Money laundering or terrorist financing</li>
              <li>Fraudulent transactions or scams</li>
              <li>Violation of any applicable laws or regulations</li>
              <li>Attempting to gain unauthorized access to our systems</li>
              <li>Interfering with the Service or other users</li>
              <li>Using the Service for gambling or illegal activities</li>
            </ul>
            <p class="text-gray-700">
              Violation of these terms may result in account suspension or termination and legal action.
            </p>
          </div>
        </div>

        <!-- Privacy Policy -->
        <div id="privacy" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">9. Privacy Policy</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              Your privacy is important to us. Our Privacy Policy explains how we collect, use, and protect your information.
            </p>
            <p class="text-gray-700 mb-4">
              We collect information you provide during registration and use of our Service, including:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li>Personal identification information (name, email, etc.)</li>
              <li>Transaction history and wallet activity</li>
              <li>Device and usage information</li>
              <li>Communication records with our support team</li>
            </ul>
            <p class="text-gray-700">
              We use this information to provide our Service, comply with legal obligations, and improve our platform.
            </p>
          </div>
        </div>

        <!-- Limitation of Liability -->
        <div id="liability" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">10. Limitation of Liability</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              To the maximum extent permitted by law, SecureWallet shall not be liable for any indirect, incidental, special, consequential, or punitive damages, including but not limited to:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li>Loss of profits, data, or use</li>
              <li>Business interruption</li>
              <li>Damages resulting from unauthorized access or use</li>
              <li>Any other damages arising from your use of the Service</li>
            </ul>
            <p class="text-gray-700">
              Our total liability to you for any claims shall not exceed the amount you paid us in fees during the 12 months preceding the claim.
            </p>
          </div>
        </div>

        <!-- Termination -->
        <div id="termination" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">11. Termination</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              You may terminate your account at any time by contacting our support team. We may terminate or suspend your account:
            </p>
            <ul class="list-disc list-inside text-gray-700 space-y-2 mb-4">
              <li>For violation of these Terms</li>
              <li>For suspicious or fraudulent activity</li>
              <li>To comply with legal requirements</li>
              <li>At our discretion for any reason</li>
            </ul>
            <p class="text-gray-700">
              Upon termination, you will lose access to your account and any remaining funds may be held for a period as required by law.
            </p>
          </div>
        </div>

        <!-- Changes to Terms -->
        <div id="changes" class="bg-white rounded-lg shadow-md p-8">
          <h2 class="text-2xl font-bold text-gray-900 mb-4">12. Changes to Terms</h2>
          <div class="prose prose-gray max-w-none">
            <p class="text-gray-700 mb-4">
              We reserve the right to modify these Terms at any time. Changes will be effective immediately upon posting on our website.
            </p>
            <p class="text-gray-700 mb-4">
              We will notify you of significant changes via email or through our Service. Your continued use of the Service after changes constitutes acceptance of the new Terms.
            </p>
            <p class="text-gray-700">
              If you do not agree to the modified Terms, you must stop using our Service and contact us to close your account.
            </p>
          </div>
        </div>
      </div>

      <!-- Contact Information -->
      <div class="mt-12 bg-gradient-to-r from-primary-600 to-primary-700 rounded-lg shadow-lg p-8 text-center text-white">
        <h3 class="text-2xl font-bold mb-4">Questions About These Terms?</h3>
        <p class="text-primary-100 mb-6">If you have any questions about these Terms of Service, please contact our support team</p>
        <router-link
          to="/support"
          class="inline-flex items-center px-6 py-3 bg-white text-primary-600 font-medium rounded-lg hover:bg-gray-50 transition-colors"
        >
          <i class="fas fa-envelope mr-2"></i>
          Contact Support
        </router-link>
      </div>
    </div>
  </div>
</template>

<script>
import Navigation from '@/components/Navigation.vue'

export default {
  name: 'TermsOfService',
  components: {
    Navigation
  },
  setup() {
    const lastUpdated = 'December 2024'

    return {
      lastUpdated
    }
  }
}
</script>
