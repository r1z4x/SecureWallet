<template>
  <!-- Public Pages Header -->
  <header v-if="type === 'landing' || type === 'blog' || type === 'blog-post'" class="bg-white shadow-sm">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-6">
        <!-- Logo Section -->
        <div class="flex items-center">
          <img src="@/assets/logo.svg" alt="Logo" class="h-12 w-full max-w-64 mr-3">
          <sup class="text-red-500 text-xs">{{ $t('common.vulnerable') }}</sup>
        </div>
        
        <!-- Actions Section -->
        <div class="flex items-center space-x-4">
          <LanguageSelector />

          <template v-if="type === 'blog'">
            <router-link to="/" class="btn-secondary">
              <i class="fas fa-home mr-2"></i>{{ $t('nav.home') }}
            </router-link>
          </template>
          
          <!-- Dynamic Actions based on type -->
          <template v-if="type === 'landing' || type === 'blog'">
            <router-link to="/auth/login" class="btn-primary">
              <i class="fas fa-sign-in-alt mr-2"></i>{{ $t('auth.signIn') }}
            </router-link>
            <router-link to="/auth/register" class="btn-secondary">
              <i class="fas fa-user-plus mr-2"></i>{{ $t('auth.register') }}
            </router-link>
          </template>

          <!-- Custom actions slot -->
          <slot name="actions"></slot>
        </div>
      </div>
    </div>
  </header>

  <!-- Authenticated Pages Navigation -->
  <Navigation v-else-if="type === 'authenticated'" />
</template>

<script>
import LanguageSelector from '@/components/LanguageSelector.vue'
import Navigation from '@/components/Navigation.vue'

export default {
  name: 'AppHeader',
  components: {
    LanguageSelector,
    Navigation
  },
  props: {
    type: {
      type: String,
      default: 'landing',
      validator: (value) => {
        return ['landing', 'blog', 'blog-post', 'authenticated', 'custom'].includes(value)
      }
    }
  }
}
</script>
