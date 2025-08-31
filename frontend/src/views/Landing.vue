<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100">
    <!-- Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <div class="flex items-center">
            <i class="fas fa-wallet text-primary-600 text-3xl mr-3"></i>
            <h1 class="text-2xl font-bold text-gray-900">SecureWallet</h1>
          </div>
          <div class="flex space-x-4">
            <router-link to="/auth/login" class="btn-primary">
              <i class="fas fa-sign-in-alt mr-2"></i>Sign In
            </router-link>
            <router-link to="/auth/register" class="btn-secondary">
              <i class="fas fa-user-plus mr-2"></i>Create Account
            </router-link>
          </div>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <div class="text-center">
        <div class="mb-12">
          <i class="fas fa-wallet text-primary-600 text-8xl mb-8"></i>
          <h1 class="text-5xl font-bold text-gray-900 mb-6">
            Welcome to SecureWallet
          </h1>
          <p class="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
            Your secure digital wallet for fast and reliable transactions. 
            Experience the future of banking with cutting-edge security and seamless user experience.
          </p>
        </div>

        <div class="flex justify-center space-x-6 mb-16">
          <router-link to="/auth/login" class="btn-primary text-lg px-8 py-4">
            <i class="fas fa-sign-in-alt mr-2"></i>Get Started
          </router-link>
          <router-link to="/auth/register" class="btn-secondary text-lg px-8 py-4">
            <i class="fas fa-user-plus mr-2"></i>Create Account
          </router-link>
        </div>

        <!-- Features -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mt-20">
          <div class="card text-center">
            <i class="fas fa-shield-alt text-primary-600 text-4xl mb-4"></i>
            <h3 class="text-xl font-bold text-gray-900 mb-2">Secure Banking</h3>
            <p class="text-gray-600">
              Bank-grade security with end-to-end encryption to protect your financial data.
            </p>
          </div>
          
          <div class="card text-center">
            <i class="fas fa-bolt text-primary-600 text-4xl mb-4"></i>
            <h3 class="text-xl font-bold text-gray-900 mb-2">Fast Transactions</h3>
            <p class="text-gray-600">
              Instant money transfers and real-time transaction processing.
            </p>
          </div>
          
          <div class="card text-center">
            <i class="fas fa-mobile-alt text-primary-600 text-4xl mb-4"></i>
            <h3 class="text-xl font-bold text-gray-900 mb-2">Mobile First</h3>
            <p class="text-gray-600">
              Optimized for mobile devices with responsive design and intuitive interface.
            </p>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
export default {
  name: 'Landing'
}
</script>
