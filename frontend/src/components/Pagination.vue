<template>
  <div v-if="totalItems > 0" class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
    <!-- Mobile pagination -->
    <div class="flex-1 flex justify-between sm:hidden">
      <button 
        @click="$emit('prev-page')" 
        :disabled="currentPage === 1"
        class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {{ $t('common.previous') }}
      </button>
      <button 
        @click="$emit('next-page')" 
        :disabled="currentPage === totalPages"
        class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
      >
        {{ $t('common.next') }}
      </button>
    </div>
    
    <!-- Desktop pagination -->
    <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
      <div>
        <p class="text-sm text-gray-700">
          {{ $t('common.showing') }} 
          <span class="font-medium">{{ startIndex + 1 }}</span>
          {{ $t('common.to') }}
          <span class="font-medium">{{ Math.min(endIndex, totalItems) }}</span>
          {{ $t('common.of') }}
          <span class="font-medium">{{ totalItems }}</span>
          {{ $t('common.results') }}
        </p>
      </div>
      <div>
        <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
          <button 
            @click="$emit('prev-page')" 
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <i class="fas fa-chevron-left"></i>
          </button>
          <template v-for="page in visiblePages" :key="page">
            <button 
              v-if="page !== '...'"
              @click="$emit('go-to-page', page)"
              :class="[
                'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                page === currentPage
                  ? 'z-10 bg-primary-50 border-primary-500 text-primary-600'
                  : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
              ]"
            >
              {{ page }}
            </button>
            <span 
              v-else
              class="relative inline-flex items-center px-4 py-2 border border-gray-300 bg-white text-sm font-medium text-gray-700"
            >
              ...
            </span>
          </template>
          <button 
            @click="$emit('next-page')" 
            :disabled="currentPage === totalPages"
            class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <i class="fas fa-chevron-right"></i>
          </button>
        </nav>
      </div>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

export default {
  name: 'Pagination',
  props: {
    currentPage: {
      type: Number,
      required: true
    },
    totalItems: {
      type: Number,
      required: true
    },
    itemsPerPage: {
      type: Number,
      default: 10
    }
  },
  emits: ['prev-page', 'next-page', 'go-to-page'],
  setup(props) {
    const { t } = useI18n()
    
    const totalPages = computed(() => Math.ceil(props.totalItems / props.itemsPerPage))
    const startIndex = computed(() => (props.currentPage - 1) * props.itemsPerPage)
    const endIndex = computed(() => startIndex.value + props.itemsPerPage)
    
    // Generate visible page numbers for pagination
    const visiblePages = computed(() => {
      const pages = []
      const total = totalPages.value
      const current = props.currentPage
      
      if (total <= 7) {
        // Show all pages if total is 7 or less
        for (let i = 1; i <= total; i++) {
          pages.push(i)
        }
      } else {
        // Always show first page
        pages.push(1)
        
        if (current > 3) {
          pages.push('...')
        }
        
        // Show pages around current page
        const start = Math.max(2, current - 1)
        const end = Math.min(total - 1, current + 1)
        
        for (let i = start; i <= end; i++) {
          if (i > 1 && i < total) {
            pages.push(i)
          }
        }
        
        if (current < total - 2) {
          pages.push('...')
        }
        
        // Always show last page
        if (total > 1) {
          pages.push(total)
        }
      }
      
      return pages
    })
    
    return {
      t,
      totalPages,
      startIndex,
      endIndex,
      visiblePages
    }
  }
}
</script>
