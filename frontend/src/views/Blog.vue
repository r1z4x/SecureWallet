<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100">
    <!-- Header -->
    <AppHeader type="blog" />

    <!-- Blog Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <!-- Page Header -->
      <div class="text-center mb-16">
        <h1 class="text-5xl font-bold text-gray-900 mb-6">{{ $t('blog.title') }}</h1>
        <p class="text-xl text-gray-600 max-w-3xl mx-auto">{{ $t('blog.subtitle') }}</p>
      </div>

      <!-- Search and Filters -->
      <div class="mb-12">
        <div class="flex flex-col md:flex-row gap-4 items-center justify-between">
          <div class="relative flex-1 max-w-md">
            <input
              v-model="searchQuery"
              type="text"
              :placeholder="$t('blog.searchPlaceholder')"
              @input="watchAndRefetch"
              class="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            >
            <i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
          </div>
          
          <div class="flex gap-2">
            <button
              v-for="category in categories"
              :key="category.value"
              @click="() => { selectedCategory = selectedCategory === category.value ? null : category.value; watchAndRefetch(); }"
              :class="[
                'px-4 py-2 rounded-lg font-medium transition-colors',
                selectedCategory === category.value
                  ? 'bg-primary-600 text-white'
                  : 'bg-white text-gray-700 hover:bg-gray-50 border border-gray-300'
              ]"
            >
              {{ category.label }}
            </button>
          </div>
        </div>
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
          @click="fetchBlogPosts"
          class="bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors"
        >
          {{ $t('blog.error.retry') }}
        </button>
      </div>

      <!-- Blog Grid -->
      <div v-if="!loading && !error" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        <article 
          v-for="post in filteredPosts" 
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
                  <div class="text-sm font-medium text-gray-900">{{ post.author.name }}</div>
                  <div class="text-xs text-gray-500">{{ post.author.title }}</div>
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

      <!-- Pagination -->
      <div v-if="!loading && !error && totalPages > 1" class="mt-16 flex justify-center">
        <div class="flex items-center space-x-2">
          <button
            @click="() => { currentPage = Math.max(1, currentPage - 1); fetchBlogPosts(); }"
            :disabled="currentPage === 1"
            class="px-4 py-2 rounded-lg border border-gray-300 text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <i class="fas fa-chevron-left mr-2"></i>
            {{ $t('blog.previous') }}
          </button>
          
          <div class="flex items-center space-x-1">
            <button
              v-for="page in visiblePages"
              :key="page"
              @click="() => { currentPage = page; fetchBlogPosts(); }"
              :class="[
                'px-4 py-2 rounded-lg font-medium transition-colors',
                currentPage === page
                  ? 'bg-primary-600 text-white'
                  : 'text-gray-700 hover:bg-gray-100'
              ]"
            >
              {{ page }}
            </button>
          </div>
          
          <button
            @click="() => { currentPage = Math.min(totalPages, currentPage + 1); fetchBlogPosts(); }"
            :disabled="currentPage === totalPages"
            class="px-4 py-2 rounded-lg border border-gray-300 text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ $t('blog.next') }}
            <i class="fas fa-chevron-right ml-2"></i>
          </button>
        </div>
      </div>

      <!-- No Results -->
      <div v-if="!loading && !error && filteredPosts.length === 0" class="text-center py-16">
        <i class="fas fa-search text-6xl text-gray-300 mb-6"></i>
        <h3 class="text-2xl font-bold text-gray-900 mb-4">{{ $t('blog.noResults.title') }}</h3>
        <p class="text-gray-600 mb-8">{{ $t('blog.noResults.description') }}</p>
        <button
          @click="clearFilters"
          class="bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors"
        >
          {{ $t('blog.noResults.clearFilters') }}
        </button>
      </div>
    </main>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import AppHeader from '@/components/AppHeader.vue'
import { blogService } from '@/services/blog'

export default {
  name: 'Blog',
  components: {
    AppHeader
  },
  setup() {
    const searchQuery = ref('')
    const selectedCategory = ref(null)
    const currentPage = ref(1)
    const postsPerPage = 9

    // Blog data from database
    const blogPosts = ref([])
    const loading = ref(false)
    const error = ref(null)

    const categories = ref([])

    // Computed properties
    const filteredPosts = computed(() => {
      let filtered = blogPosts.value

      // Filter by search query
      if (searchQuery.value) {
        filtered = filtered.filter(post =>
          post.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
          post.excerpt.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
          post.author.name.toLowerCase().includes(searchQuery.value.toLowerCase())
        )
      }

      // Filter by category
      if (selectedCategory.value) {
        filtered = filtered.filter(post => post.category === selectedCategory.value)
      }

      return filtered
    })

    const totalPages = computed(() => Math.ceil(filteredPosts.value.length / postsPerPage))

    const visiblePages = computed(() => {
      const pages = []
      const start = Math.max(1, currentPage.value - 2)
      const end = Math.min(totalPages.value, currentPage.value + 2)
      
      for (let i = start; i <= end; i++) {
        pages.push(i)
      }
      
      return pages
    })

    // Methods
    const getCategoryBadgeClass = (category) => {
      const classes = {
        security: 'bg-red-600',
        technology: 'bg-blue-600',
        finance: 'bg-green-600'
      }
      return classes[category] || 'bg-gray-600'
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    const clearFilters = () => {
      searchQuery.value = ''
      selectedCategory.value = null
      currentPage.value = 1
      fetchBlogPosts()
    }

    // Fetch blog posts from database
    const fetchBlogPosts = async () => {
      try {
        loading.value = true
        error.value = null
        
        const params = {
          page: currentPage.value,
          limit: postsPerPage
        }
        
        if (selectedCategory.value) {
          params.category = selectedCategory.value
        }
        
        if (searchQuery.value) {
          params.search = searchQuery.value
        }
        
        const response = await blogService.getPosts(params)
        blogPosts.value = response.posts || []
      } catch (err) {
        console.error('Error fetching blog posts:', err)
        error.value = 'Failed to load blog posts'
        blogPosts.value = []
      } finally {
        loading.value = false
      }
    }

    // Fetch categories from database
    const fetchCategories = async () => {
      try {
        const response = await blogService.getCategories()
        categories.value = response.map(cat => ({
          value: cat.slug,
          label: cat.name
        }))
      } catch (err) {
        console.error('Error fetching categories:', err)
        // Fallback to default categories
        categories.value = [
          { value: 'security', label: 'Security' },
          { value: 'technology', label: 'Technology' },
          { value: 'finance', label: 'Finance' }
        ]
      }
    }

    // Debounce function for search
    let searchTimeout = null
    const debouncedSearch = () => {
      if (searchTimeout) {
        clearTimeout(searchTimeout)
      }
      searchTimeout = setTimeout(() => {
        if (searchQuery.value || selectedCategory.value) {
          currentPage.value = 1
        }
        fetchBlogPosts()
      }, 500) // 500ms delay
    }

    // Watch for changes and refetch data
    const watchAndRefetch = () => {
      debouncedSearch()
    }

    // Initialize data on component mount
    onMounted(async () => {
      await fetchCategories()
      await fetchBlogPosts()
    })

    return {
      searchQuery,
      selectedCategory,
      currentPage,
      blogPosts,
      categories,
      filteredPosts,
      totalPages,
      visiblePages,
      loading,
      error,
      getCategoryBadgeClass,
      formatDate,
      clearFilters,
      fetchBlogPosts,
      fetchCategories,
      watchAndRefetch
    }
  }
}
</script>

<style scoped>
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
</style>
