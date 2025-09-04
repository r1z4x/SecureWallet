<template>
  <div class="min-h-screen bg-gradient-to-br from-primary-50 to-blue-100">

    <!-- Header -->
    <header class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-6">
          <div class="flex items-center">
            <img src="@/assets/logo.svg" alt="Logo" class="h-10 max-w-64 w-auto mr-3">
            <h1 class="text-2xl font-bold text-gray-900">
              <sup class="text-red-500 text-xs">{{ $t('common.vulnerable') }}</sup>
            </h1>
          </div>
          <div class="flex items-center space-x-4">
            <LanguageSelector />
            <router-link to="/" class="btn-secondary">
              <i class="fas fa-home mr-2"></i>{{ $t('nav.home') }}
            </router-link>
          </div>
        </div>
      </div>
    </header>

    <!-- Blog Post Content -->
    <main class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
      <!-- Back to Blog Link -->
      <div class="mb-8">
        <router-link 
          to="/blog" 
          class="inline-flex items-center text-primary-600 hover:text-primary-800 transition-colors font-medium"
        >
          <i class="fas fa-arrow-left mr-2"></i>
          <span>{{ $t('blog.backToBlog') }}</span>
        </router-link>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="text-center py-16">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
        <p class="mt-4 text-gray-600">{{ $t('blog.loading') }}</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="text-center py-16">
        <i class="fas fa-exclamation-triangle text-6xl text-red-300 mb-6"></i>
        <h3 class="text-2xl font-bold text-gray-900 mb-4">{{ $t('blog.error.title') }}</h3>
        <p class="text-gray-600 mb-8">{{ error }}</p>
        <button
          @click="fetchPost"
          class="bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors"
        >
          {{ $t('blog.error.retry') }}
        </button>
      </div>

      <!-- Article Content -->
      <article v-else-if="post" class="bg-white rounded-2xl shadow-lg overflow-hidden">
        <!-- Featured Image -->
        <div class="relative">
          <img :src="post.image" :alt="post.title" class="w-full h-96 object-cover">
          <div class="absolute top-6 left-6">
            <span :class="getCategoryBadgeClass(post.category)" class="text-white text-sm px-4 py-2 rounded-full font-semibold">
              {{ $t(`blog.categories.${post.category}`) }}
            </span>
          </div>
        </div>

        <!-- Article Content -->
        <div class="p-8">
          <!-- Meta Information -->
          <div class="flex items-center text-sm text-gray-500 mb-6">
            <div class="flex items-center mr-6">
              <i class="fas fa-calendar mr-2"></i>
              <span>{{ formatDate(post.date) }}</span>
            </div>
            <div class="flex items-center mr-6">
              <i class="fas fa-clock mr-2"></i>
              <span>{{ post.readTime }} {{ $t('blog.readTime') }}</span>
            </div>
            <div class="flex items-center">
              <i class="fas fa-user mr-2"></i>
              <span>{{ post.author.name }}</span>
            </div>
          </div>

          <!-- Title -->
          <h1 class="text-4xl font-bold text-gray-900 mb-6 leading-tight">{{ post.title }}</h1>

          <!-- Author Info -->
          <div class="flex items-center mb-8 p-6 bg-gray-50 rounded-xl">
            <img :src="post.author.avatar" :alt="post.author.name" class="w-16 h-16 rounded-full mr-4">
            <div>
              <div class="text-lg font-semibold text-gray-900">{{ post.author.name }}</div>
              <div class="text-gray-600">{{ post.author.title }}</div>
              <div class="text-sm text-gray-500">{{ post.author.bio }}</div>
            </div>
          </div>

          <!-- Article Body -->
          <div class="prose prose-lg max-w-none">
            <div v-html="post.content"></div>
          </div>

          <!-- Tags -->
          <div class="mt-8 pt-8 border-t border-gray-200">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ $t('blog.tags') }}</h3>
            <div class="flex flex-wrap gap-2">
              <span 
                v-for="tag in post.tags" 
                :key="tag"
                class="px-3 py-1 bg-primary-100 text-primary-700 rounded-full text-sm font-medium"
              >
                {{ tag }}
              </span>
            </div>
          </div>

          <!-- Share Buttons -->
          <div class="mt-8 pt-8 border-t border-gray-200">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ $t('blog.share') }}</h3>
            <div class="flex space-x-4">
              <button class="flex items-center px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors">
                <i class="fab fa-twitter mr-2"></i>
                {{ $t('blog.shareTwitter') }}
              </button>
              <button class="flex items-center px-4 py-2 bg-blue-800 text-white rounded-lg hover:bg-blue-900 transition-colors">
                <i class="fab fa-facebook mr-2"></i>
                {{ $t('blog.shareFacebook') }}
              </button>
              <button class="flex items-center px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors">
                <i class="fab fa-linkedin mr-2"></i>
                {{ $t('blog.shareLinkedIn') }}
              </button>
            </div>
          </div>
        </div>
      </article>

      <!-- Related Posts -->
      <div v-if="relatedPosts.length > 0" class="mt-16">
        <h2 class="text-3xl font-bold text-gray-900 mb-8 text-center">{{ $t('blog.relatedPosts') }}</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
          <article 
            v-for="relatedPost in relatedPosts" 
            :key="relatedPost.id"
            class="group bg-white rounded-xl shadow-lg overflow-hidden transform hover:scale-105 transition-all duration-300 hover:shadow-xl"
          >
            <div class="relative">
              <img :src="relatedPost.image" :alt="relatedPost.title" class="w-full h-48 object-cover group-hover:scale-110 transition-transform duration-300">
              <div class="absolute top-4 left-4" :class="getCategoryBadgeClass(relatedPost.category)">
                <span class="text-white text-xs px-3 py-1 rounded-full font-semibold">
                  {{ $t(`blog.categories.${relatedPost.category}`) }}
                </span>
              </div>
            </div>
            
            <div class="p-6">
              <h3 class="text-lg font-bold text-gray-900 mb-3 group-hover:text-primary-600 transition-colors line-clamp-2">
                {{ relatedPost.title }}
              </h3>
              
              <p class="text-gray-600 mb-4 line-clamp-3">
                {{ relatedPost.excerpt }}
              </p>
              
              <router-link 
                :to="`/blog/${relatedPost.slug}`" 
                class="inline-flex items-center text-primary-600 font-semibold hover:text-primary-800 transition-colors"
              >
                {{ $t('blog.readMore') }}
                <i class="fas fa-arrow-right ml-2 group-hover:translate-x-1 transition-transform"></i>
              </router-link>
            </div>
          </article>
        </div>
      </div>

      <!-- Comments Section -->
      <div class="mt-16 bg-white rounded-2xl shadow-lg p-8">
        <h2 class="text-3xl font-bold text-gray-900 mb-8">{{ $t('blog.comments.title') }}</h2>
        
        <!-- Comment Form -->
        <div class="mb-8 p-6 bg-gray-50 rounded-xl">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ $t('blog.comments.addComment') }}</h3>
          <form @submit.prevent="addComment">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <input
                v-model="commentForm.name"
                type="text"
                :placeholder="$t('blog.comments.namePlaceholder')"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                required
              >
              <input
                v-model="commentForm.email"
                type="email"
                :placeholder="$t('blog.comments.emailPlaceholder')"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                required
              >
            </div>
            <textarea
              v-model="commentForm.content"
              :placeholder="$t('blog.comments.contentPlaceholder')"
              rows="4"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-transparent mb-4"
              required
            ></textarea>
            <button
              type="submit"
              :disabled="submittingComment"
              class="bg-primary-600 text-white px-6 py-3 rounded-lg font-semibold hover:bg-primary-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="submittingComment">
                <i class="fas fa-spinner fa-spin mr-2"></i>
                {{ $t('blog.comments.submitting') }}
              </span>
              <span v-else>
                {{ $t('blog.comments.submit') }}
              </span>
            </button>
          </form>
        </div>

        <!-- Comments List -->
        <div v-if="comments.length > 0" class="space-y-6">
          <div 
            v-for="comment in comments" 
            :key="comment.id"
            class="p-6 bg-gray-50 rounded-xl"
          >
            <div class="flex items-start space-x-4">
              <div class="w-12 h-12 bg-primary-600 rounded-full flex items-center justify-center text-white font-semibold">
                {{ comment.name.charAt(0).toUpperCase() }}
              </div>
              <div class="flex-1">
                <div class="flex items-center space-x-2 mb-2">
                  <span class="font-semibold text-gray-900">{{ comment.name }}</span>
                  <span class="text-sm text-gray-500">{{ formatDate(comment.created_at) }}</span>
                </div>
                <p class="text-gray-700">{{ comment.content }}</p>
              </div>
            </div>
          </div>

          <!-- Comments Pagination -->
          <div v-if="commentPagination.pages > 1" class="mt-8 pt-6 border-t border-gray-200">
            <div class="flex items-center justify-between">
              <div class="text-sm text-gray-500">
                {{ $t('common.showing') }} {{ (commentPagination.page - 1) * commentPagination.limit + 1 }} 
                {{ $t('common.to') }} {{ Math.min(commentPagination.page * commentPagination.limit, commentPagination.total) }} 
                {{ $t('common.of') }} {{ commentPagination.total }} {{ $t('common.results') }}
              </div>
              
              <div class="flex items-center space-x-2">
                <!-- Previous Button -->
                <button
                  @click="changeCommentPage(commentPagination.page - 1)"
                  :disabled="commentPagination.page <= 1"
                  class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <i class="fas fa-chevron-left mr-1"></i>
                  {{ $t('blog.previous') }}
                </button>

                <!-- Page Numbers -->
                <div class="flex items-center space-x-1">
                  <button
                    v-for="pageNum in getVisibleCommentPages()"
                    :key="pageNum"
                    @click="changeCommentPage(pageNum)"
                    :class="[
                      'px-3 py-2 text-sm font-medium rounded-md',
                      pageNum === commentPagination.page
                        ? 'bg-primary-600 text-white'
                        : 'text-gray-500 bg-white border border-gray-300 hover:bg-gray-50'
                    ]"
                  >
                    {{ pageNum }}
                  </button>
                </div>

                <!-- Next Button -->
                <button
                  @click="changeCommentPage(commentPagination.page + 1)"
                  :disabled="commentPagination.page >= commentPagination.pages"
                  class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {{ $t('blog.next') }}
                  <i class="fas fa-chevron-right ml-1"></i>
                </button>
              </div>
            </div>
          </div>
                </div>

        <!-- No Comments -->
        <div v-else-if="!loadingComments" class="text-center py-8 text-gray-500">
          <i class="fas fa-comments text-4xl mb-4 opacity-50"></i>
          <p>{{ $t('blog.comments.noComments') }}</p>
        </div>

        <!-- Comments Loading -->
        <div v-else class="text-center py-8">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          <p class="mt-2 text-gray-600">{{ $t('blog.comments.loading') }}</p>
        </div>
      </div>
    </main>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import LanguageSelector from '@/components/LanguageSelector.vue'
import { blogService } from '@/services/blog.js'

export default {
  name: 'BlogPost',
  components: {
    LanguageSelector
  },
  setup() {
    const route = useRoute()
    
    // Reactive state
    const post = ref(null)
    const relatedPosts = ref([])
    const comments = ref([])
    const loading = ref(true)
    const loadingComments = ref(true)
    const error = ref(null)
    const submittingComment = ref(false)
    
    // Comment pagination state
    const commentPagination = ref({
      page: 1,
      limit: 10,
      total: 0,
      pages: 0
    })
    
    const commentForm = ref({
      name: '',
      email: '',
      content: ''
    })

    // Fetch blog post data
    const fetchPost = async () => {
      try {
        loading.value = true
        error.value = null
        
        const response = await blogService.getPost(route.params.slug)
        post.value = response
        relatedPosts.value = response.relatedPosts || []
        
        // Fetch comments for this post
        await fetchComments()
      } catch (err) {
        error.value = err.message || 'Failed to fetch blog post'
        console.error('Error fetching blog post:', err)
      } finally {
        loading.value = false
      }
    }

    // Fetch comments for the current post
    const fetchComments = async (page = 1) => {
      try {
        loadingComments.value = true
        const response = await blogService.getComments(route.params.slug, { page, limit: commentPagination.value.limit })
        comments.value = response.comments || []
        
        // Update pagination state
        commentPagination.value = {
          page: response.page || page,
          limit: response.limit || commentPagination.value.limit,
          total: response.total || 0,
          pages: response.pages || 0
        }
      } catch (err) {
        console.error('Error fetching comments:', err)
        comments.value = []
      } finally {
        loadingComments.value = false
      }
    }

    // Add a new comment
    const addComment = async () => {
      try {
        submittingComment.value = true
        
        const newComment = {
          name: commentForm.value.name,
          email: commentForm.value.email,
          content: commentForm.value.content
        }
        
        await blogService.addComment(route.params.slug, newComment)
        
        // Reset form
        commentForm.value = {
          name: '',
          email: '',
          content: ''
        }
        
        // Refresh comments and go to first page
        await fetchComments(1)
        
        // Show success message (you can add a toast notification here)
        console.log('Comment added successfully')
      } catch (err) {
        console.error('Error adding comment:', err)
        // Show error message (you can add a toast notification here)
        alert('Failed to add comment. Please try again.')
      } finally {
        submittingComment.value = false
      }
    }

    // Methods
    const getCategoryBadgeClass = (category) => {
      const classes = {
        security: 'bg-red-600',
        technology: 'bg-blue-600',
        finance: 'bg-green-600'
      }
      return classes[category] || 'bg-gray-600'
    }

    // Change comment page
    const changeCommentPage = async (page) => {
      if (page < 1 || page > commentPagination.value.pages) return
      await fetchComments(page)
    }

    // Get visible page numbers for pagination
    const getVisibleCommentPages = () => {
      const current = commentPagination.value.page
      const total = commentPagination.value.pages
      const delta = 2
      
      const range = []
      const rangeWithDots = []
      
      for (let i = Math.max(2, current - delta); i <= Math.min(total - 1, current + delta); i++) {
        range.push(i)
      }
      
      if (current - delta > 2) {
        rangeWithDots.push(1, '...')
      } else {
        rangeWithDots.push(1)
      }
      
      rangeWithDots.push(...range)
      
      if (current + delta < total - 1) {
        rangeWithDots.push('...', total)
      } else if (total > 1) {
        rangeWithDots.push(total)
      }
      
      return rangeWithDots
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    // Fetch data on component mount
    onMounted(() => {
      fetchPost()
    })

    return {
      post,
      relatedPosts,
      comments,
      loading,
      loadingComments,
      error,
      commentForm,
      submittingComment,
      commentPagination,
      getCategoryBadgeClass,
      formatDate,
      addComment,
      fetchPost,
      changeCommentPage,
      getVisibleCommentPages
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

.prose {
  color: #374151;
}

.prose h2 {
  color: #111827;
  font-size: 1.5rem;
  font-weight: 700;
  margin-top: 2rem;
  margin-bottom: 1rem;
}

.prose p {
  margin-bottom: 1rem;
  line-height: 1.75;
}
</style>
