<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100">
    <!-- Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <div class="flex items-center">
            <i class="fas fa-wallet text-primary-600 text-3xl mr-3"></i>
            <h1 class="text-2xl font-bold text-gray-900">{{ $t('common.appName') }} <sup class="text-red-500 text-xs">{{ $t('common.vulnerable') }}</sup></h1>
          </div>
          <div class="flex items-center space-x-4">
            <LanguageSelector />
            <router-link to="/auth/login" class="btn-primary">
              <i class="fas fa-sign-in-alt mr-2"></i>{{ $t('auth.signIn') }}
            </router-link>
            <router-link to="/auth/register" class="btn-secondary">
              <i class="fas fa-user-plus mr-2"></i>{{ $t('auth.register') }}
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
            {{ $t('landing.welcomeTitle') }} {{ $t('common.appName') }} <sup class="text-red-500 text-xs">{{ $t('common.vulnerable') }}</sup>
          </h1>
          <p class="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
            {{ $t('landing.welcomeSubtitle') }}
          </p>
        </div>

        <div class="flex justify-center space-x-6 mb-16">
          <router-link to="/auth/login" class="btn-primary text-lg px-8 py-4">
            <i class="fas fa-sign-in-alt mr-2"></i>{{ $t('landing.getStarted') }}
          </router-link>
          <router-link to="/auth/register" class="btn-secondary text-lg px-8 py-4">
            <i class="fas fa-user-plus mr-2"></i>{{ $t('auth.register') }}
          </router-link>
        </div>

        <!-- Features -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mt-20">
          <div class="card text-center">
            <i class="fas fa-shield-alt text-primary-600 text-4xl mb-4"></i>
            <h3 class="text-xl font-bold text-gray-900 mb-2">{{ $t('landing.features.secureBanking.title') }}</h3>
            <p class="text-gray-600">
              {{ $t('landing.features.secureBanking.description') }}
            </p>
          </div>
          
          <div class="card text-center">
            <i class="fas fa-bolt text-primary-600 text-4xl mb-4"></i>
            <h3 class="text-xl font-bold text-gray-900 mb-2">{{ $t('landing.features.fastTransactions.title') }}</h3>
            <p class="text-gray-600">
              {{ $t('landing.features.fastTransactions.description') }}
            </p>
          </div>
          
          <div class="card text-center">
            <i class="fas fa-mobile-alt text-primary-600 text-4xl mb-4"></i>
            <h3 class="text-xl font-bold text-gray-900 mb-2">{{ $t('landing.features.mobileFirst.title') }}</h3>
            <p class="text-gray-600">
              {{ $t('landing.features.mobileFirst.description') }}
            </p>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-gray-800 text-white py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div class="col-span-1 md:col-span-2">
            <div class="flex items-center mb-4">
              <i class="fas fa-wallet text-primary-400 text-2xl mr-2"></i>
              <span class="text-xl font-bold">{{ $t('common.appName') }} <sup class="text-red-400 text-xs">{{ $t('common.vulnerable') }}</sup></span>
            </div>
            <p class="text-gray-300 mb-4">
              {{ $t('landing.footerDescription') }}
            </p>
            <div class="flex space-x-4">
              <a href="#" class="text-gray-400 hover:text-white">
                <i class="fab fa-twitter text-xl"></i>
              </a>
              <a href="#" class="text-gray-400 hover:text-white">
                <i class="fab fa-facebook text-xl"></i>
              </a>
              <a href="#" class="text-gray-400 hover:text-white">
                <i class="fab fa-linkedin text-xl"></i>
              </a>
            </div>
          </div>
          
          <div>
            <h3 class="text-lg font-semibold mb-4">{{ $t('landing.quickLinks') }}</h3>
            <ul class="space-y-2">
              <li>
                <router-link to="/help" class="text-gray-300 hover:text-white">
                  {{ $t('landing.helpCenter') }}
                </router-link>
              </li>
              <li>
                <router-link to="/support" class="text-gray-300 hover:text-white">
                  {{ $t('landing.contactSupport') }}
                </router-link>
              </li>
              <li>
                <a href="#" class="text-gray-300 hover:text-white">
                  {{ $t('landing.aboutUs') }}
                </a>
              </li>
            </ul>
          </div>
          
          <div>
            <h3 class="text-lg font-semibold mb-4">{{ $t('landing.legal') }}</h3>
            <ul class="space-y-2">
              <li>
                <router-link to="/terms" class="text-gray-300 hover:text-white">
                  {{ $t('landing.termsOfService') }}
                </router-link>
              </li>
              <li>
                <a href="#" class="text-gray-300 hover:text-white">
                  {{ $t('landing.privacyPolicy') }}
                </a>
              </li>
              <li>
                <a href="#" class="text-gray-300 hover:text-white">
                  {{ $t('landing.cookiePolicy') }}
                </a>
              </li>
            </ul>
          </div>
        </div>
        
        <div class="border-t border-gray-700 mt-8 pt-8 text-center">
          <p class="text-gray-400">
            {{ $t('landing.copyright') }}
          </p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import LanguageSelector from '@/components/LanguageSelector.vue'

export default {
  name: 'Landing',
  components: {
    LanguageSelector
  }
}
</script>
