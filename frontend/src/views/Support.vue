<template>
  <div class="min-h-screen bg-gray-50">
    <Navigation />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">Customer Support</h1>
        <p class="text-gray-600 mt-2">Get help with your account and transactions</p>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Main Content -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Help Center Link -->
          <div class="bg-gradient-to-r from-primary-600 to-primary-700 rounded-lg shadow-md p-6 text-white">
            <div class="flex items-center justify-between">
              <div>
                <h2 class="text-xl font-semibold">Looking for Quick Answers?</h2>
                <p class="text-primary-100 mt-1">Check our Help Center for frequently asked questions</p>
              </div>
              <router-link
                to="/help"
                class="inline-flex items-center px-4 py-2 bg-white text-primary-600 font-medium rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-question-circle mr-2"></i>
                Visit Help Center
              </router-link>
            </div>
          </div>

          <!-- Create Ticket Button -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <div class="flex items-center justify-between">
              <div>
                <h2 class="text-xl font-semibold text-gray-900">Need Help?</h2>
                <p class="text-gray-600 mt-1">Create a support ticket and we'll get back to you</p>
              </div>
              <button
                @click="showCreateModal = true"
                class="btn-primary"
              >
                <i class="fas fa-plus mr-2"></i>
                Create New Ticket
              </button>
            </div>
          </div>

          <!-- My Tickets -->
          <div class="bg-white rounded-lg shadow-md">
            <div class="px-6 py-4 border-b border-gray-200">
              <h2 class="text-xl font-semibold text-gray-900">My Support Tickets</h2>
            </div>
            
            <div class="p-6">
              <div v-if="loading" class="text-center py-8">
                <i class="fas fa-spinner fa-spin text-2xl text-gray-400"></i>
                <p class="text-gray-500 mt-2">Loading tickets...</p>
              </div>
              
              <div v-else-if="tickets.length === 0" class="text-center py-8">
                <i class="fas fa-inbox text-4xl text-gray-300"></i>
                <p class="text-gray-500 mt-2">No support tickets yet</p>
                <p class="text-gray-400 text-sm">Create your first ticket above</p>
              </div>
              
              <div v-else class="space-y-4">
                <div
                  v-for="ticket in tickets"
                  :key="ticket.id"
                  class="border border-gray-200 rounded-lg p-6 hover:bg-gray-50 transition-colors"
                >
                  <div class="flex items-start justify-between">
                    <div class="flex-1">
                      <div class="flex items-center space-x-3 mb-2">
                        <h3 class="text-lg font-medium text-gray-900">
                          {{ ticket.subject }}
                        </h3>
                        <span
                          class="px-2 py-1 text-xs font-medium rounded-full"
                          :class="getStatusColor(ticket.status)"
                        >
                          {{ ticket.status }}
                        </span>
                        <span
                          class="px-2 py-1 text-xs font-medium rounded-full"
                          :class="getPriorityColor(ticket.priority)"
                        >
                          {{ ticket.priority }}
                        </span>
                      </div>
                      
                      <p class="text-gray-600 mb-3">
                        {{ ticket.description }}
                      </p>
                      
                      <div class="flex items-center space-x-4 text-sm text-gray-500">
                        <span>
                          <i class="fas fa-calendar mr-1"></i>
                          {{ formatDate(ticket.created_at) }}
                        </span>
                        <span>
                          <i class="fas fa-hashtag mr-1"></i>
                          #{{ ticket.id }}
                        </span>
                      </div>
                    </div>
                    
                    <div class="ml-4">
                      <button
                        @click="viewTicket(ticket)"
                        class="btn-secondary"
                      >
                        <i class="fas fa-eye mr-2"></i>
                        View
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Sidebar -->
        <div class="space-y-6">
          <!-- Quick Help -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Help</h3>
            <div class="space-y-3">
              <router-link
                to="/faq"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-question-circle text-blue-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">FAQ</p>
                  <p class="text-xs text-gray-500">Frequently asked questions</p>
                </div>
              </router-link>
              <router-link
                to="/user-guide"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-book text-green-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">User Guide</p>
                  <p class="text-xs text-gray-500">How to use our platform</p>
                </div>
              </router-link>
              <router-link
                to="/security-tips"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-shield-alt text-purple-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Security Tips</p>
                  <p class="text-xs text-gray-500">Keep your account safe</p>
                </div>
              </router-link>
              <router-link
                to="/help"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-life-ring text-orange-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Help Center</p>
                  <p class="text-xs text-gray-500">Browse all help articles</p>
                </div>
              </router-link>
              <router-link
                to="/terms"
                class="flex items-center p-3 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
              >
                <i class="fas fa-file-alt text-indigo-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Terms of Service</p>
                  <p class="text-xs text-gray-500">Read our terms</p>
                </div>
              </router-link>
            </div>
          </div>

          <!-- Contact Information -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Contact Information</h3>
            <div class="space-y-3">
              <div class="flex items-center">
                <i class="fas fa-envelope text-primary-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Email Support</p>
                  <p class="text-xs text-gray-500">support@securewallet.com</p>
                </div>
              </div>
              <div class="flex items-center">
                <i class="fas fa-phone text-primary-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Phone Support</p>
                  <p class="text-xs text-gray-500">+1 (555) 123-4567</p>
                </div>
              </div>
              <div class="flex items-center">
                <i class="fas fa-clock text-primary-600 mr-3"></i>
                <div>
                  <p class="text-sm font-medium text-gray-900">Support Hours</p>
                  <p class="text-xs text-gray-500">24/7 Available</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Support Stats -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Support Stats</h3>
            <div class="space-y-3">
              <div class="flex justify-between">
                <span class="text-gray-600">Total Tickets:</span>
                <span class="font-medium">{{ tickets.length }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Open Tickets:</span>
                <span class="font-medium text-orange-600">{{ openTicketsCount }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Resolved:</span>
                <span class="font-medium text-green-600">{{ resolvedTicketsCount }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Response Time:</span>
                <span class="font-medium">~2 hours</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Ticket Modal -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showCreateModal = false"
    >
      <div class="relative top-10 mx-auto p-5 border w-11/12 max-w-2xl shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold text-gray-900">
              Create Support Ticket
            </h3>
            <button
              @click="showCreateModal = false"
              class="text-gray-400 hover:text-gray-600"
            >
              <i class="fas fa-times text-xl"></i>
            </button>
          </div>
          
          <form @submit.prevent="handleCreateTicket" class="space-y-6">
            <!-- Subject -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Subject
              </label>
              <input
                v-model="ticketForm.subject"
                type="text"
                class="form-input w-full"
                placeholder="Brief description of your issue"
                required
              >
            </div>



            <!-- Priority -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Priority
              </label>
              <select
                v-model="ticketForm.priority"
                class="form-input w-full"
                required
              >
                <option value="low">Low</option>
                <option value="medium">Medium</option>
                <option value="high">High</option>
                <option value="urgent">Urgent</option>
              </select>
            </div>

            <!-- Message -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Message
              </label>
              <textarea
                v-model="ticketForm.message"
                rows="6"
                class="form-input w-full"
                placeholder="Please describe your issue in detail..."
                required
              ></textarea>
              <p class="text-xs text-gray-500 mt-1">
                {{ ticketForm.message?.length || 0 }}/1000 characters
              </p>
            </div>

            <!-- Error Message -->
            <div v-if="error" class="p-4 bg-red-100 border border-red-400 text-red-700 rounded">
              {{ error }}
            </div>

            <!-- Success Message -->
            <div v-if="success" class="p-4 bg-green-100 border border-green-400 text-green-700 rounded">
              {{ success }}
            </div>

            <!-- Submit Button -->
            <div class="flex justify-end space-x-3 pt-4">
              <button
                type="button"
                @click="showCreateModal = false"
                class="btn-secondary"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="btn-primary"
                :disabled="createLoading"
              >
                <i v-if="createLoading" class="fas fa-spinner fa-spin mr-2"></i>
                <i v-else class="fas fa-paper-plane mr-2"></i>
                {{ createLoading ? 'Creating Ticket...' : 'Submit Ticket' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Ticket Detail Modal -->
    <div
      v-if="showTicketModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showTicketModal = false"
    >
      <div class="relative top-10 mx-auto p-5 border w-11/12 max-w-2xl shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-semibold text-gray-900">
              Ticket #{{ selectedTicket?.id }}
            </h3>
            <button
              @click="showTicketModal = false"
              class="text-gray-400 hover:text-gray-600"
            >
              <i class="fas fa-times text-xl"></i>
            </button>
          </div>
          
          <div v-if="selectedTicket" class="space-y-4">
            <div>
              <h4 class="font-medium text-gray-900">{{ selectedTicket.subject }}</h4>
              <div class="flex items-center space-x-3 mt-2">
                <span
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="getStatusColor(selectedTicket.status)"
                >
                  {{ selectedTicket.status }}
                </span>
                <span
                  class="px-2 py-1 text-xs font-medium rounded-full"
                  :class="getPriorityColor(selectedTicket.priority)"
                >
                  {{ selectedTicket.priority }}
                </span>
              </div>
            </div>
            
            <div>
              <h5 class="font-medium text-gray-900 mb-2">Message:</h5>
              <p class="text-gray-600 bg-gray-50 p-3 rounded">
                {{ selectedTicket.description }}
              </p>
            </div>
            
            <div class="text-sm text-gray-500">
              <p>Created: {{ formatDate(selectedTicket.created_at) }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Navigation from '@/components/Navigation.vue'
import { supportService } from '@/services/support'

export default {
  name: 'Support',
  components: {
    Navigation
  },
  setup() {
    const authStore = useAuthStore()
    const user = computed(() => authStore.user)
    
    const loading = ref(false)
    const createLoading = ref(false)
    const error = ref('')
    const success = ref('')
    const showCreateModal = ref(false)
    const showTicketModal = ref(false)
    const selectedTicket = ref(null)
    
    const tickets = ref([])
    
    const ticketForm = ref({
      subject: '',
      priority: 'medium',
      message: ''
    })

    const loadTickets = async () => {
      loading.value = true
      try {
        const response = await supportService.getTickets()
        console.log('Support tickets response:', response) // Debug log
        
        // Ensure tickets is always an array
        if (Array.isArray(response)) {
          tickets.value = response
        } else if (response && Array.isArray(response.tickets)) {
          tickets.value = response.tickets
        } else if (response && Array.isArray(response.data)) {
          tickets.value = response.data
        } else {
          console.warn('Unexpected response format:', response)
          tickets.value = []
        }
      } catch (error) {
        console.error('Error loading tickets:', error)
        tickets.value = []
      } finally {
        loading.value = false
      }
    }

    const handleCreateTicket = async () => {
      createLoading.value = true
      error.value = ''
      success.value = ''
      
      try {
        await supportService.createTicket(ticketForm.value)
        
        success.value = 'Support ticket created successfully!'
        
        // Reset form
        ticketForm.value = {
          subject: '',
          priority: 'medium',
          message: ''
        }
        
        // Reload tickets
        await loadTickets()
        
        // Close modal after successful creation
        setTimeout(() => {
          showCreateModal.value = false
          success.value = ''
        }, 2000)
        
      } catch (err) {
        error.value = err.response?.data?.detail || err.message || 'Failed to create ticket'
      } finally {
        createLoading.value = false
      }
    }

    const viewTicket = (ticket) => {
      selectedTicket.value = ticket
      showTicketModal.value = true
    }

    const getStatusColor = (status) => {
      switch (status) {
        case 'open':
          return 'bg-orange-100 text-orange-800'
        case 'in_progress':
          return 'bg-blue-100 text-blue-800'
        case 'resolved':
          return 'bg-green-100 text-green-800'
        case 'closed':
          return 'bg-gray-100 text-gray-800'
        default:
          return 'bg-gray-100 text-gray-800'
      }
    }

    const getPriorityColor = (priority) => {
      switch (priority) {
        case 'urgent':
          return 'bg-red-100 text-red-800'
        case 'high':
          return 'bg-orange-100 text-orange-800'
        case 'medium':
          return 'bg-yellow-100 text-yellow-800'
        case 'low':
          return 'bg-green-100 text-green-800'
        default:
          return 'bg-gray-100 text-gray-800'
      }
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    const openTicketsCount = computed(() => {
      if (!Array.isArray(tickets.value)) return 0
      return tickets.value.filter(t => t.status === 'open').length
    })

    const resolvedTicketsCount = computed(() => {
      if (!Array.isArray(tickets.value)) return 0
      return tickets.value.filter(t => t.status === 'resolved').length
    })

    onMounted(() => {
      loadTickets()
    })

    return {
      user,
      loading,
      createLoading,
      error,
      success,
      showCreateModal,
      showTicketModal,
      selectedTicket,
      tickets,
      ticketForm,
      openTicketsCount,
      resolvedTicketsCount,
      handleCreateTicket,
      viewTicket,
      getStatusColor,
      getPriorityColor,
      formatDate
    }
  }
}
</script>
