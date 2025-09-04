<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100">
    <!-- Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <div class="text-center lg:text-left flex">
              <img src="@/assets/logo.svg" alt="Logo" class="max-w-64 w-full mx-auto lg:mx-0">
              <sup class="text-red-500 text-xs ml-2 font-bold">{{ $t('common.vulnerable') }}</sup>
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
        <div class="mb-16">
          
          <!-- Enhanced Title -->
          <h1 class="text-6xl md:text-7xl font-bold text-gray-900 mb-8 leading-tight">
            <span class="bg-gradient-to-r from-primary-600 to-blue-600 bg-clip-text text-transparent">
              {{ $t('landing.welcomeTitle') }}
            </span>
            <br>
            <span class="inline-block items-center justify-center">
              <img src="@/assets/logo.svg" alt="Logo" class="max-w-96 w-full mx-auto lg:mx-0">
            </span>
          </h1>
          
          <!-- Enhanced Subtitle -->
          <p class="text-xl md:text-2xl text-gray-600 mb-12 max-w-4xl mx-auto leading-relaxed">
            {{ $t('landing.welcomeSubtitle') }}
          </p>
        </div>

        <!-- Enhanced CTA Buttons -->
        <div class="flex flex-col sm:flex-row justify-center items-center space-y-4 sm:space-y-0 sm:space-x-6 mb-20">
          <router-link 
            to="/auth/login" 
            class="group relative overflow-hidden bg-gradient-to-r from-primary-600 to-blue-600 text-white text-lg px-10 py-4 rounded-full font-semibold shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-300"
          >
            <span class="relative z-10 flex items-center">
              <i class="fas fa-sign-in-alt mr-3 group-hover:rotate-12 transition-transform duration-300"></i>
              {{ $t('landing.getStarted') }}
            </span>
            <div class="absolute inset-0 bg-gradient-to-r from-primary-700 to-blue-700 transform scale-x-0 group-hover:scale-x-100 transition-transform duration-300 origin-left"></div>
          </router-link>
          
          <router-link 
            to="/auth/register" 
            class="group relative overflow-hidden bg-white text-primary-600 text-lg px-10 py-4 rounded-full font-semibold border-2 border-primary-600 shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-300 hover:bg-primary-50"
          >
            <span class="relative z-10 flex items-center">
              <i class="fas fa-user-plus mr-3 group-hover:rotate-12 transition-transform duration-300"></i>
              {{ $t('auth.register') }}
            </span>
          </router-link>
        </div>

        <!-- Enhanced Features Section -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mt-24">
          <div class="group card text-center p-8 transform hover:scale-105 transition-all duration-300 hover:shadow-2xl">
            <div class="relative mb-6">
              <div class="absolute inset-0 bg-primary-100 rounded-full blur-xl opacity-50 group-hover:opacity-75 transition-opacity duration-300"></div>
              <i class="fas fa-shield-alt text-primary-600 text-5xl relative z-10 group-hover:scale-110 transition-transform duration-300"></i>
            </div>
            <h3 class="text-2xl font-bold text-gray-900 mb-4">{{ $t('landing.features.secureBanking.title') }}</h3>
            <p class="text-gray-600 leading-relaxed">
              {{ $t('landing.features.secureBanking.description') }}
            </p>
          </div>
          
          <div class="group card text-center p-8 transform hover:scale-105 transition-all duration-300 hover:shadow-2xl">
            <div class="relative mb-6">
              <div class="absolute inset-0 bg-blue-100 rounded-full blur-xl opacity-50 group-hover:opacity-75 transition-opacity duration-300"></div>
              <i class="fas fa-bolt text-blue-600 text-5xl relative z-10 group-hover:scale-110 transition-transform duration-300"></i>
            </div>
            <h3 class="text-2xl font-bold text-gray-900 mb-4">{{ $t('landing.features.fastTransactions.title') }}</h3>
            <p class="text-gray-600 leading-relaxed">
              {{ $t('landing.features.fastTransactions.description') }}
            </p>
          </div>
          
          <div class="group card text-center p-8 transform hover:scale-105 transition-all duration-300 hover:shadow-2xl">
            <div class="relative mb-6">
              <div class="absolute inset-0 bg-green-100 rounded-full blur-xl opacity-50 group-hover:opacity-75 transition-opacity duration-300"></div>
              <i class="fas fa-mobile-alt text-green-600 text-5xl relative z-10 group-hover:scale-110 transition-transform duration-300"></i>
            </div>
            <h3 class="text-2xl font-bold text-gray-900 mb-4">{{ $t('landing.features.mobileFirst.title') }}</h3>
            <p class="text-gray-600 leading-relaxed">
              {{ $t('landing.features.mobileFirst.description') }}
            </p>
          </div>
        </div>

        <!-- New Stats Section -->
        <div class="mt-24 bg-gradient-to-r from-primary-50 to-blue-50 rounded-3xl p-12">
          <h2 class="text-3xl font-bold text-gray-900 mb-12 text-center">{{ $t('landing.stats.title') }}</h2>
          <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div class="text-center">
              <div class="text-4xl font-bold text-primary-600 mb-2">10K+</div>
              <div class="text-gray-600">{{ $t('landing.stats.users') }}</div>
            </div>
            <div class="text-center">
              <div class="text-4xl font-bold text-blue-600 mb-2">50K+</div>
              <div class="text-gray-600">{{ $t('landing.stats.transactions') }}</div>
            </div>
            <div class="text-center">
              <div class="text-4xl font-bold text-green-600 mb-2">99.9%</div>
              <div class="text-gray-600">{{ $t('landing.stats.uptime') }}</div>
            </div>
            <div class="text-center">
              <div class="text-4xl font-bold text-purple-600 mb-2">24/7</div>
              <div class="text-gray-600">{{ $t('landing.stats.support') }}</div>
            </div>
          </div>
        </div>

        <!-- Latest Blog Posts Section -->
        <div class="mt-24">
          <div class="text-center mb-16">
            <h2 class="text-4xl font-bold text-gray-900 mb-4">{{ $t('blog.title') }}</h2>
            <p class="text-xl text-gray-600 max-w-2xl mx-auto">{{ $t('blog.subtitle') }}</p>
          </div>
          
          <!-- Loading State -->
          <div v-if="loading" class="text-center py-16">
            <div class="inline-flex items-center px-4 py-2 font-semibold leading-6 text-primary-600">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-primary-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ $t('blog.loading') }}
            </div>
          </div>

          <!-- Error State -->
          <div v-else-if="error" class="text-center py-16">
            <i class="fas fa-exclamation-triangle text-6xl text-red-500 mb-6"></i>
            <h3 class="text-2xl font-bold text-gray-900 mb-4">{{ $t('blog.error.title') }}</h3>
            <p class="text-gray-600 mb-8">{{ error }}</p>
            <button
              @click="fetchLatestPosts"
              class="bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors"
            >
              {{ $t('blog.error.retry') }}
            </button>
          </div>

          <!-- Blog Posts Grid -->
          <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-8">
            <article 
              v-for="post in latestPosts" 
              :key="post.id"
              class="group bg-white rounded-2xl shadow-lg overflow-hidden transform hover:scale-105 transition-all duration-300 hover:shadow-2xl"
            >
              <div class="relative">
                <img :src="post.image" :alt="post.title" class="w-full h-48 object-cover group-hover:scale-110 transition-transform duration-300">
                <div class="absolute top-4 left-4" :class="getCategoryBadgeClass(post.category)">
                  <span class="text-white text-xs px-3 py-1 rounded-full font-semibold">
                    {{ $t(`blog.categories.${post.category}`) }}
                  </span>
                </div>
              </div>
              <div class="p-6">
                <div class="flex items-center text-sm text-gray-500 mb-3">
                  <i class="fas fa-calendar mr-2"></i>
                  <span>{{ formatDate(post.date) }}</span>
                  <span class="mx-2">•</span>
                  <i class="fas fa-clock mr-2"></i>
                  <span>{{ post.readTime }} {{ $t('blog.readTime') }}</span>
                </div>
                <h3 class="text-xl font-bold text-gray-900 mb-3 group-hover:text-primary-600 transition-colors line-clamp-2">
                  {{ post.title }}
                </h3>
                <p class="text-gray-600 mb-4 line-clamp-3">
                  {{ post.excerpt }}
                </p>
                <div class="flex items-center justify-between">
                  <div class="flex items-center">
                    <img :src="post.author.avatar" :alt="post.author.name" class="w-8 h-8 rounded-full mr-3">
                    <div>
                      <div class="text-sm font-medium text-gray-900 text-left">{{ post.author.name }}</div>
                      <div class="text-xs text-gray-500 text-left">{{ post.author.title }}</div>
                    </div>
                  </div>
                  <router-link 
                    :to="`/blog/${post.slug}`" 
                    class="inline-flex items-center text-primary-600 font-semibold hover:text-primary-800 transition-colors"
                  >
                    {{ $t('blog.readMore') }}
                    <i class="fas fa-arrow-right ml-2 group-hover:translate-x-1 transition-transform"></i>
                  </router-link>
                </div>
              </div>
            </article>
          </div>

          <div class="text-center mt-12">
            <router-link to="/blog" class="inline-flex items-center bg-primary-600 text-white px-8 py-3 rounded-full font-semibold hover:bg-primary-700 transition-colors transform hover:scale-105">
              {{ $t('blog.readMore') }}
              <i class="fas fa-arrow-right ml-2"></i>
            </router-link>
          </div>
        </div>

        <!-- Testimonials Section -->
        <div class="mt-24 bg-gradient-to-br from-gray-50 to-blue-50 rounded-3xl p-12">
          <div class="text-center mb-16">
            <h2 class="text-4xl font-bold text-gray-900 mb-4">{{ $t('landing.testimonials.title') }}</h2>
            <p class="text-xl text-gray-600 max-w-2xl mx-auto">{{ $t('landing.testimonials.subtitle') }}</p>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div class="bg-white p-8 rounded-2xl shadow-lg text-center">
              <div class="flex justify-center mb-6">
                <div class="w-16 h-16 bg-gradient-to-br from-primary-400 to-blue-500 rounded-full flex items-center justify-center">
                  <i class="fas fa-user text-white text-2xl"></i>
                </div>
              </div>
              <p class="text-gray-600 mb-6 italic">"{{ $t('landing.testimonials.testimonial1.text') }}"</p>
              <div>
                <div class="font-semibold text-gray-900">{{ $t('landing.testimonials.testimonial1.name') }}</div>
                <div class="text-sm text-gray-500">{{ $t('landing.testimonials.testimonial1.position') }}</div>
              </div>
            </div>

            <div class="bg-white p-8 rounded-2xl shadow-lg text-center">
              <div class="flex justify-center mb-6">
                <div class="w-16 h-16 bg-gradient-to-br from-green-400 to-teal-500 rounded-full flex items-center justify-center">
                  <i class="fas fa-user text-white text-2xl"></i>
                </div>
              </div>
              <p class="text-gray-600 mb-6 italic">"{{ $t('landing.testimonials.testimonial2.text') }}"</p>
              <div>
                <div class="font-semibold text-gray-900">{{ $t('landing.testimonials.testimonial2.name') }}</div>
                <div class="text-sm text-gray-500">{{ $t('landing.testimonials.testimonial2.position') }}</div>
              </div>
            </div>

            <div class="bg-white p-8 rounded-2xl shadow-lg text-center">
              <div class="flex justify-center mb-6">
                <div class="w-16 h-16 bg-gradient-to-br from-purple-400 to-pink-500 rounded-full flex items-center justify-center">
                  <i class="fas fa-user text-white text-2xl"></i>
                </div>
              </div>
              <p class="text-gray-600 mb-6 italic">"{{ $t('landing.testimonials.testimonial3.text') }}"</p>
              <div>
                <div class="font-semibold text-gray-900">{{ $t('landing.testimonials.testimonial3.name') }}</div>
                <div class="text-sm text-gray-500">{{ $t('landing.testimonials.testimonial3.position') }}</div>
              </div>
            </div>
          </div>
        </div>



        <!-- CTA Section -->
        <div class="mt-24 bg-gradient-to-r from-primary-600 to-blue-600 rounded-3xl p-16 text-center text-white">
          <h2 class="text-4xl font-bold mb-6">{{ $t('landing.cta.title') }}</h2>
          <p class="text-xl mb-8 opacity-90 max-w-2xl mx-auto">{{ $t('landing.cta.subtitle') }}</p>
          <div class="flex flex-col sm:flex-row justify-center items-center space-y-4 sm:space-y-0 sm:space-x-6">
            <router-link 
              to="/auth/register" 
              class="bg-white text-primary-600 px-8 py-4 rounded-full font-semibold hover:bg-gray-100 transition-colors transform hover:scale-105"
            >
              {{ $t('landing.cta.primaryButton') }}
            </router-link>
            <router-link 
              to="/support" 
              class="border-2 border-white text-white px-8 py-4 rounded-full font-semibold hover:bg-white hover:text-primary-600 transition-colors transform hover:scale-105"
            >
              {{ $t('landing.cta.secondaryButton') }}
            </router-link>
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
              <img src="@/assets/logo.svg" alt="Logo" class="max-w-64 w-full mr-3">
              <sup class="text-red-400 text-xs">{{ $t('common.vulnerable') }}</sup>
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
import { ref, onMounted } from 'vue'
import LanguageSelector from '@/components/LanguageSelector.vue'
import { blogService } from '@/services/blog'

export default {
  name: 'Landing',
  components: {
    LanguageSelector
  },
  setup() {
    const latestPosts = ref([])
    const loading = ref(false)
    const error = ref(null)

    // Fetch latest blog posts
    const fetchLatestPosts = async () => {
      try {
        loading.value = true
        error.value = null
        
        const response = await blogService.getPosts({ limit: 3 })
        latestPosts.value = response.posts || []
      } catch (err) {
        console.error('Error fetching latest posts:', err)
        error.value = 'Failed to load latest posts'
        latestPosts.value = []
      } finally {
        loading.value = false
      }
    }

    // Get category badge class
    const getCategoryBadgeClass = (category) => {
      const classes = {
        security: 'bg-red-600',
        technology: 'bg-blue-600',
        finance: 'bg-green-600'
      }
      return classes[category] || 'bg-gray-600'
    }

    // Format date
    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    // Initialize data on component mount
    onMounted(async () => {
      await fetchLatestPosts()
    })

    return {
      latestPosts,
      loading,
      error,
      fetchLatestPosts,
      getCategoryBadgeClass,
      formatDate
    }
  }
}
</script>

<style scoped>
/* Ensure sup text displays properly above the main text */
sup {
  vertical-align: super;
  font-size: 0.6em;
  line-height: 1;
  position: relative;
  top: -0.5em;
}

/* Hero section specific sup styling */
h1 sup {
  font-size: 0.256em;
  top: -1em;
}

/* Header specific sup styling */
header h1 sup {
  font-size: 0.4em;
  top: -0.8em;
}

/* Footer specific sup styling */
footer sup {
  font-size: 0.4em;
  top: -0.8em;
}

/* Enhanced card styling */
.card {
  @apply bg-white rounded-2xl shadow-lg border border-gray-100;
  backdrop-filter: blur(10px);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.card:hover {
  @apply shadow-2xl;
  transform: translateY(-8px);
}

/* Gradient text animation */
.bg-gradient-to-r {
  background-size: 200% 200%;
  animation: gradient-shift 3s ease infinite;
}

@keyframes gradient-shift {
  0%, 100% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
}

/* Enhanced button animations */
.group:hover .group-hover\:rotate-12 {
  transform: rotate(12deg);
}

/* Stats section enhancement */
.stats-card {
  background: linear-gradient(135deg, rgba(59, 130, 246, 0.05) 0%, rgba(147, 51, 234, 0.05) 100%);
  backdrop-filter: blur(20px);
}

/* Responsive improvements */
@media (max-width: 768px) {
  .text-6xl {
    font-size: 3rem;
  }
  
  .text-7xl {
    font-size: 3.5rem;
  }
  
  .py-20 {
    padding-top: 5rem;
    padding-bottom: 5rem;
  }
}

/* Smooth scrolling */
html {
  scroll-behavior: smooth;
}

/* Enhanced focus states */
button:focus,
a:focus {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
}

/* Loading animation for icons */
.fas {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Enhanced hover effects */
.transform {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Glass morphism effect for cards */
.card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .hero h1 {
    font-size: 2.5rem;
  }
  
  .hero p {
    font-size: 1.1rem;
  }
  
  .features-grid {
    grid-template-columns: 1fr;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .hero h1 {
    font-size: 2rem;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
}
</style>
