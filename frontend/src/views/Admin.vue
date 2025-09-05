<template>
  <div class="min-h-screen bg-gray-50">
    <AppHeader type="authenticated" />
    
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Admin Access Check -->
      <div v-if="!user?.is_admin" class="mb-8">
        <div class="bg-red-50 border border-red-200 rounded-md p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <i class="fas fa-exclamation-triangle text-red-400"></i>
            </div>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">
                {{ $t('admin.accessDenied') }}
              </h3>
              <p class="text-sm text-red-700 mt-1">
                {{ $t('admin.adminPrivilegesRequired') }}
              </p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Header -->
      <div v-if="user?.is_admin" class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">{{ $t('admin.title') }}</h1>
        <p class="text-gray-600 mt-2">{{ $t('admin.subtitle') }}</p>
      </div>

      <!-- Admin Stats -->
      <div v-if="user?.is_admin" class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-users text-blue-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">{{ $t('admin.totalUsers') }}</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.totalUsers }}</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-exchange-alt text-green-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">{{ $t('admin.totalTransactions') }}</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.totalTransactions }}</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-dollar-sign text-yellow-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">{{ $t('admin.totalVolume') }}</p>
              <p class="text-2xl font-bold text-gray-900">${{ stats.totalVolume?.toFixed(2) || '0.00' }}</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <i class="fas fa-ticket-alt text-purple-600 text-2xl"></i>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">{{ $t('admin.supportTickets') }}</p>
              <p class="text-2xl font-bold text-gray-900">{{ stats.totalTickets }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Admin Tabs -->
      <div v-if="user?.is_admin" class="bg-white rounded-lg shadow-md">
        <div class="border-b border-gray-200">
          <nav class="-mb-px flex space-x-8 px-6">
            <button
              v-for="tab in tabs"
              :key="tab.id"
              @click="activeTab = tab.id"
              :class="[
                'py-4 px-1 border-b-2 font-medium text-sm',
                activeTab === tab.id
                  ? 'border-primary-500 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              <i :class="tab.icon" class="mr-2"></i>
              {{ tab.name }}
            </button>
          </nav>
        </div>

        <div class="p-6">
          <!-- Users Tab -->
          <div v-if="activeTab === 'users'" class="space-y-6">
            <div class="flex justify-between items-center">
              <h2 class="text-xl font-semibold text-gray-900">{{ $t('admin.userManagement') }}</h2>
              <div class="flex items-center space-x-4">
                <div class="flex items-center space-x-2">
                  <label class="text-sm text-gray-700">{{ $t('common.showing') }}:</label>
                  <select 
                    v-model="itemsPerPage" 
                    @change="changeItemsPerPage(itemsPerPage)"
                    class="form-input w-20 text-sm"
                  >
                    <option :value="5">5</option>
                    <option :value="10">10</option>
                    <option :value="25">25</option>
                    <option :value="50">50</option>
                  </select>
                  <span class="text-sm text-gray-700">{{ $t('common.results') }}</span>
                </div>
                <button @click="addUser" class="btn-primary">
                  <i class="fas fa-plus mr-2"></i>
                  {{ $t('admin.addUser') }}
                </button>
              </div>
            </div>
            
            <div class="overflow-x-auto">
              <!-- Debug Info -->
              <div v-if="users.length === 0" class="text-center py-8">
                <i class="fas fa-users text-4xl text-gray-300 mb-4"></i>
                <p class="text-gray-500">{{ $t('admin.noUsersFound') }}</p>
                <p class="text-gray-400 text-sm mt-2">{{ $t('admin.usersWillAppearHere') }}</p>
              </div>
              
              <table v-else class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('common.user') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('auth.email') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('common.status') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('common.role') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('admin.joined') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('admin.actions') }}</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="user in paginatedUsers" :key="user.id">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <div class="flex-shrink-0 h-10 w-10">
                          <div class="h-10 w-10 rounded-full bg-primary-100 flex items-center justify-center">
                            <span class="text-primary-600 font-medium">{{ user.username ? user.username.charAt(0).toUpperCase() : '?' }}</span>
                          </div>
                        </div>
                        <div class="ml-4">
                          <div class="text-sm font-medium text-gray-900">{{ user.username || 'Unknown User' }}</div>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ user.email || 'No email' }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          user.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                        ]"
                      >
                        {{ (user.is_active === true) ? $t('admin.active') : $t('admin.inactive') }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {{ (user.is_admin === true) ? $t('admin.admin') : $t('admin.user') }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ user.created_at ? formatDate(user.created_at) : 'Unknown' }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                      <button @click="editUser(user.id)" class="text-primary-600 hover:text-primary-900 mr-3">{{ $t('common.edit') }}</button>
                      <button @click="deleteUser(user.id)" class="text-red-600 hover:text-red-900">{{ $t('common.delete') }}</button>
                    </td>
                  </tr>
                </tbody>
              </table>
              
              <!-- Pagination -->
              <Pagination 
                :current-page="currentPage"
                :total-items="users.length"
                :items-per-page="itemsPerPage"
                @prev-page="prevPage"
                @next-page="nextPage"
                @go-to-page="goToPage"
              />
            </div>
          </div>

          <!-- Transactions Tab -->
          <div v-if="activeTab === 'transactions'" class="space-y-6">
            <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center space-y-4 sm:space-y-0">
              <h2 class="text-xl font-semibold text-gray-900">{{ $t('admin.transactionHistory') }}</h2>
              <div class="flex flex-col sm:flex-row space-y-3 sm:space-y-0 sm:space-x-3 w-full sm:w-auto">
                <div class="flex flex-col space-y-1">
                  <select 
                    v-model="transactionFilter" 
                    class="form-input w-full sm:w-48 border-gray-300 focus:border-primary-500 focus:ring-primary-500"
                  >
                    <option value="">{{ $t('admin.allTypes') }}</option>
                    <option value="transfer">{{ $t('admin.transfer') }}</option>
                    <option value="deposit">{{ $t('admin.deposit') }}</option>
                    <option value="withdrawal">{{ $t('admin.withdrawal') }}</option>
                  </select>
                </div>
                <div class="flex flex-col space-y-1">
                  <button 
                    @click="exportTransactions" 
                    class="btn-secondary w-full sm:w-auto bg-blue-600 hover:bg-blue-700 border-blue-600 hover:border-blue-700 text-white"
                    :disabled="filteredTransactions.length === 0"
                  >
                    <i class="fas fa-download mr-2"></i>
                    {{ $t('transactions.export') }}
                  </button>
                </div>
              </div>
            </div>
            
            <div class="overflow-x-auto">
              <!-- Debug Info -->
              <div v-if="transactions.length === 0" class="text-center py-8">
                <i class="fas fa-exchange-alt text-4xl text-gray-300 mb-4"></i>
                <p class="text-gray-500">{{ $t('admin.noTransactionsFound') }}</p>
                <p class="text-gray-400 text-sm mt-2">{{ $t('admin.transactionsWillAppearHere') }}</p>
              </div>
              
              <table v-else class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('admin.id') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('admin.user') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('admin.type') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('admin.amount') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('admin.description') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('admin.status') }}</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">{{ $t('admin.date') }}</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="transaction in paginatedTransactions" :key="transaction.id">
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">#{{ transaction.id }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {{ transaction.wallet?.user?.email || 'Unknown User' }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          getTransactionTypeColor(transaction.transaction_type)
                        ]"
                      >
                        {{ transaction.transaction_type }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      ${{ parseFloat(transaction.amount).toFixed(2) }}
                    </td>
                    <td class="px-6 py-4 text-sm text-gray-900">
                      <!-- VULNERABLE: Stored XSS - description rendered without sanitization for admin exploitation -->
                      <span v-html="transaction.description || 'No description'"></span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          transaction.status === 'completed' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                        ]"
                      >
                        {{ transaction.status === 'completed' ? $t('admin.completed') : $t('admin.pending') }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ formatDate(transaction.created_at) }}
                    </td>
                  </tr>
                </tbody>
              </table>
              
              <!-- Pagination -->
              <Pagination 
                :current-page="currentPage"
                :total-items="filteredTransactions.length"
                :items-per-page="itemsPerPage"
                @prev-page="prevPage"
                @next-page="nextPage"
                @go-to-page="goToPage"
              />
            </div>
          </div>

          <!-- Support Tab -->
          <div v-if="activeTab === 'support'" class="space-y-6">
            <div class="flex justify-between items-center">
              <h2 class="text-xl font-semibold text-gray-900">{{ $t('admin.supportTickets') }}</h2>
              <div class="flex space-x-2">
                <select v-model="ticketFilter" class="form-input">
                  <option value="">{{ $t('admin.allStatus') }}</option>
                  <option value="open">{{ $t('admin.open') }}</option>
                  <option value="in_progress">{{ $t('admin.inProgress') }}</option>
                  <option value="resolved">{{ $t('admin.resolved') }}</option>
                  <option value="closed">{{ $t('admin.closed') }}</option>
                </select>
              </div>
            </div>
            
            <div class="space-y-4">
              <!-- Debug Info -->
              <div v-if="tickets.length === 0" class="text-center py-8">
                <i class="fas fa-ticket-alt text-4xl text-gray-300 mb-4"></i>
                <p class="text-gray-500">{{ $t('admin.noTicketsFound') }}</p>
                <p class="text-gray-400 text-sm mt-2">{{ $t('admin.ticketsWillAppearHere') }}</p>
              </div>
              
              <div
                v-for="ticket in paginatedTickets"
                :key="ticket.id"
                class="border border-gray-200 rounded-lg p-6 hover:bg-gray-50"
              >
                <div class="flex items-start justify-between">
                  <div class="flex-1">
                    <div class="flex items-center space-x-3 mb-2">
                      <h3 class="text-lg font-medium text-gray-900">{{ ticket.subject }} <span class="text-sm text-gray-500 font-normal">(#{{ ticket.id.toLowerCase() }})</span></h3>
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          getTicketStatusColor(ticket.status)
                        ]"
                      >
                        {{ ticket.status }}
                      </span>
                      <span
                        :class="[
                          'px-2 py-1 text-xs font-medium rounded-full',
                          getTicketPriorityColor(ticket.priority)
                        ]"
                      >
                        {{ ticket.priority }}
                      </span>
                    </div>
                    <p class="text-gray-600 mb-3">{{ ticket.message }}</p>
                    <div class="flex items-center space-x-4 text-sm text-gray-500">
                      <span>{{ $t('admin.userId') }}: {{ ticket.user?.username || ticket.user_id }} ({{ ticket.user?.email || 'Unknown' }})</span>
                      <span>{{ $t('admin.category') }}: {{ ticket.category }}</span>
                      <span>{{ formatDate(ticket.created_at) }}</span>
                    </div>
                  </div>
                  <div class="ml-4">
                    <button @click="replyToTicket(ticket.id)" class="btn-secondary mr-2">
                      <i class="fas fa-reply mr-2"></i>
                      {{ $t('admin.reply') }}
                    </button>
                    <button @click="resolveTicket(ticket.id)" class="btn-primary">
                      <i class="fas fa-check mr-2"></i>
                      {{ $t('admin.resolve') }}
                    </button>
                  </div>
                </div>
              </div>
              
              <!-- Pagination -->
              <Pagination 
                :current-page="currentPage"
                :total-items="filteredTickets.length"
                :items-per-page="itemsPerPage"
                @prev-page="prevPage"
                @next-page="nextPage"
                @go-to-page="goToPage"
              />
            </div>
          </div>

          <!-- Settings Tab -->
          <div v-if="activeTab === 'settings'" class="space-y-6">
            <h2 class="text-xl font-semibold text-gray-900">{{ $t('admin.systemSettings') }}</h2>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div class="bg-gray-50 rounded-lg p-6">
                <h3 class="text-lg font-medium text-gray-900 mb-4">{{ $t('admin.securitySettings') }}</h3>
                <div class="space-y-4">
                  <div class="flex items-center justify-between">
                    <div class="flex flex-col">
                      <span class="text-sm font-medium text-gray-700">{{ $t('admin.twoFactorAuth') }}</span>
                      <span class="text-xs text-gray-500">{{ $t('admin.twoFactorAuthDescription') }}</span>
                    </div>
                    <div class="flex items-center">
                      <span class="text-sm text-gray-600 mr-3">
                        {{ systemSettings.security.twoFactorEnabled ? $t('admin.enabled') : $t('admin.disabled') }}
                      </span>
                      <button
                        @click="toggleTwoFactor"
                        :class="[
                          'relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2',
                          systemSettings.security.twoFactorEnabled ? 'bg-indigo-600' : 'bg-gray-200'
                        ]"
                      >
                        <span
                          :class="[
                            'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                            systemSettings.security.twoFactorEnabled ? 'translate-x-6' : 'translate-x-1'
                          ]"
                        />
                      </button>
                    </div>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">{{ $t('admin.sessionTimeout') }}</span>
                    <select 
                      v-model="systemSettings.security.sessionTimeout"
                      class="form-input w-40"
                    >
                      <option :value="30">{{ $t('admin.sessionTimeout30min') }}</option>
                      <option :value="60">{{ $t('admin.sessionTimeout1hour') }}</option>
                      <option :value="120">{{ $t('admin.sessionTimeout2hours') }}</option>
                      <option :value="1440">{{ $t('admin.sessionTimeout24hours') }}</option>
                      <option :value="10080">{{ $t('admin.sessionTimeout1week') }}</option>
                    </select>
                  </div>
                  <div class="flex items-center justify-between">
                    <span class="text-sm text-gray-600">{{ $t('admin.passwordPolicy') }}</span>
                    <button @click="showPasswordPolicyModal = true" class="btn-secondary">
                      {{ $t('admin.configure') }}
                    </button>
                  </div>
                </div>
              </div>
              
              <div class="bg-gray-50 rounded-lg p-6">
                <h3 class="text-lg font-medium text-gray-900 mb-4">{{ $t('admin.transactionLimits') }}</h3>
                <div class="space-y-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('admin.dailyTransferLimit') }}</label>
                    <input 
                      v-model="systemSettings.transactionLimits.dailyTransferLimit"
                      type="number" 
                      class="form-input w-full" 
                      min="100"
                      step="100"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('admin.monthlyTransferLimit') }}</label>
                    <input 
                      v-model="systemSettings.transactionLimits.monthlyTransferLimit"
                      type="number" 
                      class="form-input w-full" 
                      min="1000"
                      step="1000"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('admin.minTransferAmount') }}</label>
                    <input 
                      v-model="systemSettings.transactionLimits.minTransferAmount"
                      type="number" 
                      class="form-input w-full" 
                      min="0.01"
                      step="0.01"
                    />
                  </div>
                </div>
              </div>
            </div>
            
            <div class="flex justify-end items-center">
              <div class="text-sm text-gray-500 italic">
                <i class="fas fa-info-circle mr-1"></i>
                {{ $t('admin.autoSaveInfo') }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add User Modal -->
    <div
      v-if="showAddUserModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showAddUserModal = false"
    >
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ $t('admin.addUser') }}</h3>
          <form @submit.prevent="submitAddUser" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('auth.username') }}</label>
              <input 
                v-model="newUser.username"
                type="text" 
                class="form-input w-full" 
                required
                minlength="3"
                maxlength="50"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('auth.email') }}</label>
              <input 
                v-model="newUser.email"
                type="email" 
                class="form-input w-full" 
                required
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('auth.password') }}</label>
              <input 
                v-model="newUser.password"
                type="password" 
                class="form-input w-full" 
                required
                minlength="8"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('common.role') }}</label>
              <select v-model="newUser.is_admin" class="form-input w-full">
                <option :value="false">{{ $t('admin.user') }}</option>
                <option :value="true">{{ $t('admin.admin') }}</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('common.status') }}</label>
              <select v-model="newUser.is_active" class="form-input w-full">
                <option :value="true">{{ $t('admin.active') }}</option>
                <option :value="false">{{ $t('admin.inactive') }}</option>
              </select>
            </div>
            <div class="flex justify-end space-x-3 mt-6">
              <button 
                type="button"
                @click="showAddUserModal = false"
                class="btn-secondary"
              >
                {{ $t('common.cancel') }}
              </button>
              <button 
                type="submit"
                class="btn-primary"
                :disabled="addUserLoading"
              >
                <span v-if="addUserLoading">{{ $t('common.processing') }}</span>
                <span v-else>{{ $t('admin.addUser') }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Edit User Modal -->
    <div
      v-if="showEditUserModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showEditUserModal = false"
    >
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ $t('admin.editUser') }}</h3>
          <form @submit.prevent="submitEditUser" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('auth.username') }}</label>
              <input 
                v-model="editUserData.username"
                type="text" 
                class="form-input w-full" 
                required
                minlength="3"
                maxlength="50"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('auth.email') }}</label>
              <input 
                v-model="editUserData.email"
                type="email" 
                class="form-input w-full" 
                required
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('common.role') }}</label>
              <select v-model="editUserData.is_admin" class="form-input w-full">
                <option :value="false">{{ $t('admin.user') }}</option>
                <option :value="true">{{ $t('admin.admin') }}</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('common.status') }}</label>
              <select v-model="editUserData.is_active" class="form-input w-full">
                <option :value="true">{{ $t('admin.active') }}</option>
                <option :value="false">{{ $t('admin.inactive') }}</option>
              </select>
            </div>
            <div class="flex justify-end space-x-3 mt-6">
              <button 
                type="button"
                @click="showEditUserModal = false"
                class="btn-secondary"
              >
                {{ $t('common.cancel') }}
              </button>
              <button 
                type="submit"
                class="btn-primary"
                :disabled="editUserLoading"
              >
                <span v-if="editUserLoading">{{ $t('common.processing') }}</span>
                <span v-else>{{ $t('admin.updateUser') }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Password Policy Modal -->
    <div
      v-if="showPasswordPolicyModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showPasswordPolicyModal = false"
    >
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">{{ $t('admin.passwordPolicy') }}</h3>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">{{ $t('admin.minPasswordLength') }}</label>
              <input 
                v-model="systemSettings.security.passwordPolicy.minLength"
                type="number" 
                class="form-input w-full" 
                min="6"
                max="32"
              />
            </div>
            <div class="space-y-2">
              <label class="block text-sm font-medium text-gray-700">{{ $t('admin.passwordRequirements') }}</label>
              <div class="space-y-2">
                <label class="flex items-center">
                  <input 
                    v-model="systemSettings.security.passwordPolicy.requireUppercase"
                    type="checkbox" 
                    class="form-checkbox mr-2"
                  />
                  <span class="text-sm text-gray-700">{{ $t('admin.requireUppercase') }}</span>
                </label>
                <label class="flex items-center">
                  <input 
                    v-model="systemSettings.security.passwordPolicy.requireLowercase"
                    type="checkbox" 
                    class="form-checkbox mr-2"
                  />
                  <span class="text-sm text-gray-700">{{ $t('admin.requireLowercase') }}</span>
                </label>
                <label class="flex items-center">
                  <label class="flex items-center">
                    <input 
                      v-model="systemSettings.security.passwordPolicy.requireNumbers"
                      type="checkbox" 
                      class="form-checkbox mr-2"
                    />
                    <span class="text-sm text-gray-700">{{ $t('admin.requireNumbers') }}</span>
                  </label>
                </label>
                <label class="flex items-center">
                  <input 
                    v-model="systemSettings.security.passwordPolicy.requireSpecialChars"
                    type="checkbox" 
                    class="form-checkbox mr-2"
                  />
                  <span class="text-sm text-gray-700">{{ $t('admin.requireSpecialChars') }}</span>
                </label>
              </div>
            </div>
          </div>
          <div class="flex justify-end space-x-3 mt-6">
            <button 
              @click="showPasswordPolicyModal = false"
              class="btn-secondary"
            >
              {{ $t('common.cancel') }}
            </button>
            <button 
              @click="savePasswordPolicy"
              class="btn-primary"
            >
              {{ $t('common.save') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Reply to Ticket Modal -->
    <div
      v-if="showReplyModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click="showReplyModal = false"
    >
      <div class="relative top-20 mx-auto p-5 border w-full max-w-2xl shadow-lg rounded-md bg-white" @click.stop>
        <div class="mt-3">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">
              <i class="fas fa-reply mr-2 text-blue-600"></i>
              {{ $t('admin.reply') }} - {{ selectedTicket?.subject }} <span class="text-sm text-gray-500 font-normal">(#{{ selectedTicket?.id?.toLowerCase() }})</span>
            </h3>
            <button 
              @click="showReplyModal = false"
              class="text-gray-400 hover:text-gray-600"
            >
              <i class="fas fa-times text-xl"></i>
            </button>
          </div>
          
          <!-- Ticket Info -->
          <div class="bg-gray-50 rounded-lg p-4 mb-4">
            <div class="flex items-center space-x-4 text-sm text-gray-600 mb-2">
              <span><strong>{{ $t('admin.userId') }}:</strong> {{ selectedTicket?.user?.username || selectedTicket?.user_id }} ({{ selectedTicket?.user?.email || 'Unknown' }})</span>
              <span><strong>{{ $t('admin.priority') }}:</strong> 
                <span :class="getTicketPriorityColor(selectedTicket?.priority)" class="px-2 py-1 text-xs font-medium rounded-full ml-1">
                  {{ selectedTicket?.priority }}
                </span>
              </span>
              <span><strong>{{ $t('admin.status') }}:</strong> 
                <span :class="getTicketStatusColor(selectedTicket?.status)" class="px-2 py-1 text-xs font-medium rounded-full ml-1">
                  {{ selectedTicket?.status }}
                </span>
              </span>
            </div>
            <div class="text-sm text-gray-700">
              <strong>{{ $t('admin.originalMessage') }}:</strong>
              <p class="mt-1 p-2 bg-white rounded border">{{ selectedTicket?.message }}</p>
            </div>
          </div>

          <!-- Reply Form -->
          <form @submit.prevent="submitReply" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                {{ $t('admin.replyMessage') }}
              </label>
              <textarea 
                v-model="replyMessage"
                class="form-input w-full h-32 resize-none" 
                :placeholder="$t('admin.replyPlaceholder')"
                required
                maxlength="1000"
              ></textarea>
              <div class="text-right text-xs text-gray-500 mt-1">
                {{ replyMessage.length }}/1000
              </div>
            </div>
            
            <div class="flex justify-end space-x-3 mt-6">
              <button 
                type="button"
                @click="showReplyModal = false"
                class="btn-secondary"
              >
                {{ $t('common.cancel') }}
              </button>
              <button 
                type="submit"
                class="btn-primary"
                :disabled="replyLoading || !replyMessage.trim()"
              >
                <span v-if="replyLoading">
                  <i class="fas fa-spinner fa-spin mr-2"></i>
                  {{ $t('common.sending') }}
                </span>
                <span v-else>
                  <i class="fas fa-paper-plane mr-2"></i>
                  {{ $t('admin.sendReply') }}
                </span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Enhanced Notification Container -->
    <div class="notification-container">
      <!-- Clear All Button (only show if there are notifications) -->
      <div v-if="notifications.length > 1" class="mb-2">
        <button
          @click="clearAllNotifications"
          class="text-xs text-gray-500 hover:text-gray-700 bg-white px-2 py-1 rounded shadow-sm border"
          :title="$t('admin.clearAllNotifications')"
        >
          <i class="fas fa-times mr-1"></i>
          {{ $t('admin.clearAll') }}
        </button>
      </div>
      
      <transition-group name="notification" tag="div">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          :class="[
            'notification-item w-full bg-white shadow-lg rounded-lg ring-1 ring-black ring-opacity-5 overflow-hidden',
            {
              'border-l-4 border-green-400': notification.type === 'success',
              'border-l-4 border-red-400': notification.type === 'error',
              'border-l-4 border-yellow-400': notification.type === 'warning',
              'border-l-4 border-blue-400': notification.type === 'info'
            }
          ]"
        >
          <div class="p-4">
            <div class="flex items-start">
              <div class="flex-shrink-0">
                <i 
                  :class="[
                    'h-5 w-5',
                    {
                      'fas fa-check-circle text-green-400': notification.type === 'success',
                      'fas fa-exclamation-circle text-red-400': notification.type === 'error',
                      'fas fa-exclamation-triangle text-yellow-400': notification.type === 'warning',
                      'fas fa-info-circle text-blue-400': notification.type === 'info'
                    }
                  ]"
                ></i>
              </div>
              <div class="ml-3 w-0 flex-1 pt-0.5">
                <p class="text-sm font-medium text-gray-900 break-words">
                  {{ notification.message }}
                </p>
                <!-- Show timestamp for debugging -->
                <p v-if="notification.persistent" class="text-xs text-gray-500 mt-1">
                  {{ $t('admin.persistentNotification') }}
                </p>
              </div>
              <div class="ml-4 flex-shrink-0 flex items-start space-x-1">
                <!-- Action buttons if any -->
                <div v-if="notification.actions && notification.actions.length > 0" class="flex space-x-1">
                  <button
                    v-for="action in notification.actions"
                    :key="action.label"
                    @click="action.handler"
                    class="text-xs px-2 py-1 rounded text-white"
                    :class="action.class || 'bg-blue-500 hover:bg-blue-600'"
                  >
                    {{ action.label }}
                  </button>
                </div>
                <!-- Close button -->
                <button
                  @click="removeNotification(notification.id)"
                  class="bg-white rounded-md inline-flex text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                  :title="$t('admin.closeNotification')"
                >
                  <span class="sr-only">{{ $t('admin.close') }}</span>
                  <i class="fas fa-times h-5 w-5"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </transition-group>
      
      <!-- Queue indicator -->
      <div v-if="notificationQueue.length > 0" class="mt-2 text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded">
        <i class="fas fa-clock mr-1"></i>
        {{ $t('admin.notificationsInQueue', { count: notificationQueue.length }) }}
      </div>
    </div>

    <!-- Confirmation Overlay -->
    <div v-if="isPageBlocked" class="fixed inset-0 bg-black bg-opacity-50 z-[10000] flex items-center justify-center">
      <transition-group name="confirmation" tag="div" class="space-y-4">
        <div
          v-for="confirmation in pendingConfirmations"
          :key="confirmation.id"
          class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6"
        >
          <div class="flex items-center mb-4">
            <div class="flex-shrink-0">
              <i 
                :class="[
                  'h-6 w-6',
                  {
                    'fas fa-exclamation-triangle text-yellow-500': confirmation.type === 'warning',
                    'fas fa-exclamation-circle text-red-500': confirmation.type === 'error',
                    'fas fa-question-circle text-blue-500': confirmation.type === 'info',
                    'fas fa-check-circle text-green-500': confirmation.type === 'success'
                  }
                ]"
              ></i>
            </div>
            <div class="ml-3">
              <h3 class="text-lg font-medium text-gray-900">
                {{ confirmation.title }}
              </h3>
            </div>
          </div>
          
          <div class="mb-6">
            <p class="text-sm text-gray-600 break-words">
              {{ confirmation.message }}
            </p>
          </div>
          
          <div class="flex justify-end space-x-3">
            <button
              @click="handleConfirmation(confirmation.id, false)"
              :class="[
                'px-4 py-2 text-sm font-medium rounded-md text-white focus:outline-none focus:ring-2 focus:ring-offset-2',
                confirmation.cancelClass
              ]"
            >
              {{ confirmation.cancelText }}
            </button>
            <button
              @click="handleConfirmation(confirmation.id, true)"
              :class="[
                'px-4 py-2 text-sm font-medium rounded-md text-white focus:outline-none focus:ring-2 focus:ring-offset-2',
                confirmation.confirmClass
              ]"
            >
              {{ confirmation.confirmText }}
            </button>
          </div>
        </div>
      </transition-group>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/AppHeader.vue'
import Pagination from '@/components/Pagination.vue'
import { userService } from '@/services/user'
import { transactionService } from '@/services/transaction'
import { supportService } from '@/services/support'
import { adminService } from '@/services/admin'

export default {
  name: 'Admin',
  components: {
    AppHeader,
    Pagination
  },
  setup() {
    const { t } = useI18n()
    const authStore = useAuthStore()
    const user = computed(() => authStore.user)
    
    // Check if current user is admin
    const isAdminUser = computed(() => user.value?.is_admin === true)
    
    // State variables
    const activeTab = ref('users')
    const transactionFilter = ref('')
    const ticketFilter = ref('')
    
    const users = ref([])
    const transactions = ref([])
    const tickets = ref([])
    const stats = ref({
      totalUsers: 0,
      totalTransactions: 0,
      totalVolume: 0,
      totalTickets: 0
    })

    // Settings state
    const systemSettings = ref({
      security: {
        twoFactorEnabled: false,
        sessionTimeout: 30,
        passwordPolicy: {
          minLength: 8,
          requireUppercase: true,
          requireLowercase: true,
          requireNumbers: true,
          requireSpecialChars: true
        }
      },
      transactionLimits: {
        dailyTransferLimit: 10000,
        monthlyTransferLimit: 50000,
        minTransferAmount: 1
      }
    })

    // Load settings from localStorage on component mount
    const loadSettingsFromStorage = () => {
      try {
        const savedSettings = localStorage.getItem('adminSystemSettings')
        if (savedSettings) {
          const parsedSettings = JSON.parse(savedSettings)
          systemSettings.value = { ...systemSettings.value, ...parsedSettings }
        }
      } catch (error) {
        console.warn('Error loading settings from localStorage:', error)
      }
    }

    // Save settings to localStorage
    const saveSettingsToStorage = (settings) => {
      try {
        localStorage.setItem('adminSystemSettings', JSON.stringify(settings))
      } catch (error) {
        console.warn('Error saving settings to localStorage:', error)
      }
    }

    // Modal states
    const showPasswordPolicyModal = ref(false)
    const showAddUserModal = ref(false)
    const showEditUserModal = ref(false)
    const showReplyModal = ref(false)
    const selectedTicket = ref(null)
    const replyMessage = ref('')
    const replyLoading = ref(false)
    
    // Notification states
    const notifications = ref([])
    const notificationId = ref(0)
    const notificationQueue = ref([])
    const maxNotifications = 5
    const notificationTimeouts = ref(new Map())
    
    // Confirmation states
    const pendingConfirmations = ref([])
    const confirmationId = ref(0)
    const isPageBlocked = ref(false)
    
    // Auto-save timeout
    const autoSaveTimeout = ref(null)

    // New user form state
    const newUser = ref({
      username: '',
      email: '',
      password: '',
      is_admin: false,
      is_active: true
    })
    const addUserLoading = ref(false)
    const editUserLoading = ref(false)
    const selectedUser = ref(null)
    const editUserData = ref({
      id: '',
      username: '',
      email: '',
      is_admin: false,
      is_active: true
    })

    // Pagination states
    const currentPage = ref(1)
    const itemsPerPage = ref(10)
    
    // Pagination computed properties
    const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage.value)
    const endIndex = computed(() => startIndex.value + itemsPerPage.value)
    
    const paginatedUsers = computed(() => {
      return users.value.slice(startIndex.value, endIndex.value)
    })
    
    const paginatedTransactions = computed(() => {
      return filteredTransactions.value.slice(startIndex.value, endIndex.value)
    })
    
    const paginatedTickets = computed(() => {
      return filteredTickets.value.slice(startIndex.value, endIndex.value)
    })


    // Watch for filter changes to reset pagination
    watch(transactionFilter, () => {
      currentPage.value = 1
    })

    watch(ticketFilter, () => {
      currentPage.value = 1
    })

    // Watch for tab changes to reset pagination
    watch(activeTab, () => {
      currentPage.value = 1
    })

    const tabs = [
      { id: 'users', name: t('admin.users'), icon: 'fas fa-users' },
      { id: 'transactions', name: t('admin.transactions'), icon: 'fas fa-exchange-alt' },
      { id: 'support', name: t('admin.support'), icon: 'fas fa-headset' },
      { id: 'settings', name: t('admin.settings'), icon: 'fas fa-cog' }
    ]

    const loadUsers = async () => {
      try {

        
        // Check if current user is admin
        if (!user.value?.is_admin) {
          console.error('User is not admin, cannot load all users')
          throw new Error('Admin privileges required')
        }
        
        const response = await userService.getAllUsers()
        users.value = response || []
                  stats.value.totalUsers = users.value.length
        
        // Reset to first page when loading new data
        currentPage.value = 1
        
        // If no users loaded, show message
        if (users.value.length === 0) {
          // No users found
        }
      } catch (error) {
        console.error('Error loading users:', error)
        
        // If it's an admin privilege error, show appropriate message
        if (error.message === 'Admin privileges required') {
          users.value = []
          stats.value.totalUsers = 0
          return
        }
        
        // On error, set empty array
        users.value = []
        stats.value.totalUsers = 0
      }
    }

    const loadTransactions = async () => {
      try {

        // Use admin service to get ALL transactions from the system
        const response = await adminService.getAdminTransactions(100)
        transactions.value = response || []
        stats.value.totalTransactions = transactions.value.length
        stats.value.totalVolume = transactions.value.reduce((sum, t) => sum + parseFloat(t.amount), 0)
        
        // Reset to first page when loading new data
        currentPage.value = 1
        
        // If no transactions loaded, add mock data for testing
        if (transactions.value.length === 0) {
          transactions.value = [
            {
              id: 1,
              transaction_type: 'transfer',
              amount: '100.00',
              description: 'Test transfer',
              status: 'completed',
              created_at: new Date().toISOString()
            },
            {
              id: 2,
              transaction_type: 'deposit',
              amount: '500.00',
              description: 'Test deposit',
              status: 'completed',
              created_at: new Date(Date.now() - 86400000).toISOString()
            },
            {
              id: 3,
              transaction_type: 'withdrawal',
              amount: '50.00',
              description: 'Test withdrawal',
              status: 'pending',
              created_at: new Date(Date.now() - 172800000).toISOString()
            }
          ]
          stats.value.totalTransactions = transactions.value.length
          stats.value.totalVolume = transactions.value.reduce((sum, t) => sum + parseFloat(t.amount), 0)
        }
      } catch (error) {
        console.error('Error loading admin transactions:', error)
        // Add mock data on error
        transactions.value = [
          {
            id: 1,
            transaction_type: 'transfer',
            amount: '100.00',
            description: 'Test transfer',
            status: 'completed',
            created_at: new Date().toISOString()
          }
        ]
        stats.value.totalTransactions = transactions.value.length
        stats.value.totalVolume = transactions.value.reduce((sum, t) => sum + parseFloat(t.amount), 0)
      }
    }

    const loadTickets = async () => {
      try {
        // Use admin service to get ALL support tickets from the system
        const response = await adminService.getAdminSupportTickets()
        tickets.value = response || []
        stats.value.totalTickets = tickets.value.length
        
        // Reset to first page when loading new data
        currentPage.value = 1
        
        // If no tickets loaded, add mock data for testing
        if (tickets.value.length === 0) {
          tickets.value = [
            {
              id: '550e8400-e29b-41d4-a716-446655440001',
              subject: 'Login issue',
              message: 'I cannot login to my account',
              status: 'open',
              priority: 'high',
              user_id: '550e8400-e29b-41d4-a716-446655440002',
              category: 'Authentication',
              created_at: new Date().toISOString(),
              user: {
                id: '550e8400-e29b-41d4-a716-446655440002',
                username: 'testuser',
                email: 'test@example.com'
              }
            },
            {
              id: '550e8400-e29b-41d4-a716-446655440003',
              subject: 'Transfer problem',
              message: 'My transfer is stuck',
              status: 'in_progress',
              priority: 'medium',
              user_id: '550e8400-e29b-41d4-a716-446655440004',
              category: 'Transactions',
              created_at: new Date(Date.now() - 86400000).toISOString(),
              user: {
                id: '550e8400-e29b-41d4-a716-446655440004',
                username: 'user2',
                email: 'user2@example.com'
              }
            },
            {
              id: '550e8400-e29b-41d4-a716-446655440005',
              subject: 'Password reset',
              message: 'Need help with password reset',
              status: 'resolved',
              priority: 'low',
              user_id: '550e8400-e29b-41d4-a716-446655440006',
              category: 'Account',
              created_at: new Date(Date.now() - 172800000).toISOString(),
              user: {
                id: '550e8400-e29b-41d4-a716-446655440006',
                username: 'admin',
                email: 'admin@example.com'
              }
            }
          ]
          stats.value.totalTickets = tickets.value.length
        }
      } catch (error) {
        console.error('Error loading tickets:', error)
        // Add mock data on error
        tickets.value = [
          {
            id: '550e8400-e29b-41d4-a716-446655440007',
            subject: 'Test ticket',
            message: 'This is a test support ticket',
            status: 'open',
            priority: 'medium',
            user_id: '550e8400-e29b-41d4-a716-446655440008',
            category: 'General',
            created_at: new Date().toISOString(),
            user: {
              id: '550e8400-e29b-41d4-a716-446655440008',
              username: 'testuser',
              email: 'test@example.com'
            }
          }
        ]
        stats.value.totalTickets = tickets.value.length
      }
    }

    const loadData = async () => {
      try {
        await Promise.all([
          loadUsers(),
          loadTransactions(),
          loadTickets()
        ])
        
        // Load system settings from backend first
        try {
          const settingsResponse = await adminService.getSystemSettings()
          if (settingsResponse) {
            systemSettings.value = {
              ...systemSettings.value,
              ...settingsResponse
            }
          }
        } catch (settingsError) {
          console.warn('Could not load system settings:', settingsError)
        }
        
        // Then load from localStorage to override backend values with user's local changes
        loadSettingsFromStorage()
      } catch (error) {
        console.error('Error loading data:', error)
      }
    }

    // Watch for user changes and load data when admin user is detected
    watch(user, (newUser, oldUser) => {
      if (newUser && newUser.is_admin) {
        loadData()
      }
    }, { immediate: true })

    // Watch for system settings changes and auto-save
    watch(systemSettings, async (newSettings, oldSettings) => {
      // Skip the first call (initial load) and 2FA changes (handled manually)
      if (oldSettings && Object.keys(oldSettings).length > 0) {
        try {
          // Debounce the save operation to avoid too many API calls
          clearTimeout(autoSaveTimeout.value)
          autoSaveTimeout.value = setTimeout(async () => {
            // Save to localStorage for persistence
            saveSettingsToStorage(newSettings)
            
            // Also try to save to backend (optional)
            try {
              await adminService.saveSystemSettings(newSettings)
              console.log('Settings auto-saved successfully')
            } catch (error) {
              console.warn('Backend auto-save failed, but settings saved locally:', error)
            }
          }, 1000) // Wait 1 second after last change
        } catch (error) {
          console.warn('Auto-save failed:', error)
        }
      }
    }, { deep: true })

    // Filtered data computed properties
    const filteredTransactions = computed(() => {
      if (!transactionFilter.value) return transactions.value
      return transactions.value.filter(t => 
        t.transaction_type === transactionFilter.value
      )
    })

    const filteredTickets = computed(() => {
      if (!ticketFilter.value) return tickets.value
      return tickets.value.filter(t => 
        t.status === ticketFilter.value
      )
    })



    const getTransactionTypeColor = (type) => {
      switch (type) {
        case 'transfer':
          return 'bg-blue-100 text-blue-800'
        case 'deposit':
          return 'bg-green-100 text-green-800'
        case 'withdrawal':
          return 'bg-red-100 text-red-800'
        default:
          return 'bg-gray-100 text-gray-800'
      }
    }

    const getTicketStatusColor = (status) => {
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

    const getTicketPriorityColor = (priority) => {
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
      if (!dateString) return 'Unknown'
      try {
        return new Date(dateString).toLocaleDateString('tr-TR', {
          year: 'numeric',
          month: 'short',
          day: 'numeric',
          hour: '2-digit',
          minute: '2-digit'
        })
      } catch (error) {
        console.error('Error formatting date:', dateString, error)
        return 'Invalid Date'
      }
    }

    // Advanced Notification System
    const addNotification = (type, message, duration = 5000, options = {}) => {
      try {
        // Validate inputs
        if (!type || !message) {
          console.error('Invalid notification parameters:', { type, message })
          return null
        }

        const validTypes = ['success', 'error', 'warning', 'info']
        if (!validTypes.includes(type)) {
          console.error('Invalid notification type:', type)
          return null
        }

        const id = ++notificationId.value
        const notification = {
          id,
          type,
          message: String(message), // Ensure message is string
          duration: Math.max(2000, Math.min(10000, duration)), // Clamp between 2-10 seconds
          timestamp: Date.now(),
          persistent: options.persistent || false,
          actions: options.actions || []
        }
        
        // Check if we're at max capacity
        if (notifications.value.length >= maxNotifications) {
          // Remove oldest notification
          const oldest = notifications.value[0]
          removeNotification(oldest.id)
        }
        
        // Add to notifications
        notifications.value.push(notification)
        
        // Set auto-remove timeout (only if not persistent)
        if (!notification.persistent) {
          const timeoutId = setTimeout(() => {
            removeNotification(id)
          }, notification.duration)
          
          notificationTimeouts.value.set(id, timeoutId)
        }
        
        return id
      } catch (error) {
        console.error('Error adding notification:', error)
        return null
      }
    }

    const removeNotification = (id) => {
      try {
        // Clear timeout if exists
        const timeoutId = notificationTimeouts.value.get(id)
        if (timeoutId) {
          clearTimeout(timeoutId)
          notificationTimeouts.value.delete(id)
        }
        
        // Remove from notifications array
        const index = notifications.value.findIndex(n => n.id === id)
        if (index > -1) {
          notifications.value.splice(index, 1)
        }
        
        // Process queue if there are pending notifications
        processNotificationQueue()
      } catch (error) {
        console.error('Error removing notification:', error)
      }
    }

    const processNotificationQueue = () => {
      if (notificationQueue.value.length > 0 && notifications.value.length < maxNotifications) {
        const queuedNotification = notificationQueue.value.shift()
        addNotification(
          queuedNotification.type,
          queuedNotification.message,
          queuedNotification.duration,
          queuedNotification.options
        )
      }
    }

    const queueNotification = (type, message, duration = 5000, options = {}) => {
      notificationQueue.value.push({ type, message, duration, options })
      processNotificationQueue()
    }

    const clearAllNotifications = () => {
      try {
        // Clear all timeouts
        notificationTimeouts.value.forEach(timeoutId => clearTimeout(timeoutId))
        notificationTimeouts.value.clear()
        
        // Clear notifications and queue
        notifications.value = []
        notificationQueue.value = []
      } catch (error) {
        console.error('Error clearing notifications:', error)
      }
    }

    // Enhanced notification functions with error handling
    const showSuccess = (message, duration = 5000, options = {}) => {
      return addNotification('success', message, duration, options)
    }
    
    const showError = (message, duration = 7000, options = {}) => {
      return addNotification('error', message, duration, { ...options, persistent: true })
    }
    
    const showWarning = (message, duration = 6000, options = {}) => {
      return addNotification('warning', message, duration, options)
    }
    
    const showInfo = (message, duration = 5000, options = {}) => {
      return addNotification('info', message, duration, options)
    }

    // Queue-based notification functions for high-frequency scenarios
    const queueSuccess = (message, duration = 5000, options = {}) => {
      queueNotification('success', message, duration, options)
    }
    
    const queueError = (message, duration = 7000, options = {}) => {
      queueNotification('error', message, duration, { ...options, persistent: true })
    }

    // Confirmation System
    const showConfirmation = (message, options = {}) => {
      return new Promise((resolve) => {
        const id = ++confirmationId.value
        const confirmation = {
          id,
          message: String(message),
          type: options.type || 'warning',
          title: options.title || t('admin.confirmation'),
          confirmText: options.confirmText || t('admin.yes'),
          cancelText: options.cancelText || t('admin.no'),
          confirmClass: options.confirmClass || 'bg-red-500 hover:bg-red-600',
          cancelClass: options.cancelClass || 'bg-gray-500 hover:bg-gray-600',
          timestamp: Date.now(),
          resolve
        }
        
        pendingConfirmations.value.push(confirmation)
        isPageBlocked.value = true
        
        return confirmation
      })
    }

    const handleConfirmation = (id, confirmed) => {
      const index = pendingConfirmations.value.findIndex(c => c.id === id)
      if (index > -1) {
        const confirmation = pendingConfirmations.value[index]
        confirmation.resolve(confirmed)
        pendingConfirmations.value.splice(index, 1)
        
        // Unblock page if no more confirmations
        if (pendingConfirmations.value.length === 0) {
          isPageBlocked.value = false
        }
      }
    }

    const confirmAction = (message, options = {}) => {
      return showConfirmation(message, options)
    }

    const addUser = () => {
      showAddUserModal.value = true
      // Reset form
      newUser.value = {
        username: '',
        email: '',
        password: '',
        is_admin: false,
        is_active: true
      }
    }

    const submitAddUser = async () => {
      try {
        addUserLoading.value = true

        
        // Call API to add user
        await userService.createUser(newUser.value)
        
        // Show success message
        showSuccess(t('admin.userAddedSuccessfully'))

        
        // Close modal and reload users
        showAddUserModal.value = false
        await loadUsers()
        
        // Reset form
        newUser.value = {
          username: '',
          email: '',
          password: '',
          is_admin: false,
          is_active: true
        }
      } catch (error) {
        console.error('Error adding user:', error)
        
        // Show specific error message
        let errorMessage = 'Error adding user'
        if (error.response?.data?.error) {
          errorMessage = error.response.data.error
        } else if (error.message) {
          errorMessage = error.message
        }
        
        showError(errorMessage)
      } finally {
        addUserLoading.value = false
      }
    }

    const editUser = (userId) => {
      const user = users.value.find(u => u.id === userId)
      if (user) {
        selectedUser.value = user
        editUserData.value = {
          id: user.id,
          username: user.username,
          email: user.email,
          is_admin: user.is_admin || false,
          is_active: user.is_active !== false
        }
        showEditUserModal.value = true
      }
    }

    const submitEditUser = async () => {
      try {
        editUserLoading.value = true

        // Call API to update user
        await userService.updateUser(editUserData.value.id, editUserData.value)
        
        // Show success message
        showSuccess(t('admin.userUpdatedSuccessfully'))

        // Close modal and reload users
        showEditUserModal.value = false
        await loadUsers()
        
        // Reset form
        editUserData.value = {
          id: '',
          username: '',
          email: '',
          is_admin: false,
          is_active: true
        }
        selectedUser.value = null
        
      } catch (error) {
        console.error('Error updating user:', error)
        let errorMessage = t('admin.errorUpdatingUser')
        if (error.response?.data?.error) {
          errorMessage += ': ' + error.response.data.error
        } else if (error.message) {
          errorMessage += ': ' + error.message
        }
        
        showError(errorMessage)
      } finally {
        editUserLoading.value = false
      }
    }

    const deleteUser = async (userId) => {
      try {
        const confirmed = await confirmAction(t('admin.confirmDeleteUser'), {
          type: 'error',
          title: t('admin.deleteUser'),
          confirmText: t('admin.delete'),
          cancelText: t('common.cancel'),
          confirmClass: 'bg-red-500 hover:bg-red-600',
          cancelClass: 'bg-gray-500 hover:bg-gray-600'
        })
        
        if (confirmed) {
          // Call API to delete user
          await userService.deleteUser(userId)
          showSuccess(t('admin.userDeletedSuccessfully'))
          await loadUsers() // Reload users after deletion
        }
      } catch (error) {
        console.error('Error deleting user:', error)
        showError(t('admin.errorDeletingUser'))
      }
    }


    const toggleTwoFactor = async () => {
      try {
        const currentState = systemSettings.value.security.twoFactorEnabled
        const newState = !currentState
        
        const confirmed = await confirmAction(
          currentState 
            ? t('admin.confirmDisableTwoFactor') 
            : t('admin.confirmEnableTwoFactor'),
          {
            type: currentState ? 'warning' : 'info',
            title: currentState ? t('admin.disableTwoFactor') : t('admin.enableTwoFactor'),
            confirmText: currentState ? t('admin.disable') : t('admin.enable'),
            cancelText: t('common.cancel'),
            confirmClass: currentState 
              ? 'bg-red-500 hover:bg-red-600' 
              : 'bg-green-500 hover:bg-green-600',
            cancelClass: 'bg-gray-500 hover:bg-gray-600'
          }
        )
        
        if (confirmed) {
          systemSettings.value.security.twoFactorEnabled = newState
          
          // Save to localStorage for persistence
          saveSettingsToStorage(systemSettings.value)
          
          // Also try to save to backend (optional)
          try {
            await adminService.saveSystemSettings(systemSettings.value)
            showSuccess(
              newState 
                ? t('admin.twoFactorEnabledSuccessfully') 
                : t('admin.twoFactorDisabledSuccessfully')
            )
          } catch (saveError) {
            console.warn('Backend save failed, but setting saved locally:', saveError)
            // Still show success since localStorage save worked
            showSuccess(
              newState 
                ? t('admin.twoFactorEnabledSuccessfully') 
                : t('admin.twoFactorDisabledSuccessfully')
            )
          }
        }
      } catch (error) {
        console.error('Error toggling two-factor authentication:', error)
        showError(t('admin.errorTogglingTwoFactor'))
      }
    }

    const savePasswordPolicy = () => {
      showPasswordPolicyModal.value = false
      // TODO: Implement API call to save password policy
    }

    // Pagination functions
    const goToPage = (page) => {
      currentPage.value = page
    }

    const nextPage = () => {
      currentPage.value++
    }

    const prevPage = () => {
      currentPage.value--
    }

    const changeItemsPerPage = (newItemsPerPage) => {
      itemsPerPage.value = newItemsPerPage
      currentPage.value = 1 // Reset to first page
    }

    const exportTransactions = () => {
      try {
        const csvContent = [
          ['ID', 'Type', 'Amount', 'Description', 'Status', 'Date'],
          ...filteredTransactions.value.map(t => [
            t.id,
            t.transaction_type,
            t.amount,
            t.description || '',
            t.status,
            t.created_at ? formatDate(t.created_at) : 'Unknown'
          ])
        ].map(row => row.join(',')).join('\n')

        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
        const link = document.createElement('a')
        const url = URL.createObjectURL(blob)
        link.setAttribute('href', url)
        link.setAttribute('download', `admin-transactions-${new Date().toISOString().split('T')[0]}.csv`)
        link.style.visibility = 'hidden'
        document.body.appendChild(link)
        link.click()
        document.body.removeChild(link)
      } catch (error) {
        console.error('Error exporting transactions:', error)
      }
    }

    const replyToTicket = (ticketId) => {
      // Find the ticket to reply to
      const ticket = tickets.value.find(t => t.id == ticketId)
      if (!ticket) {
        alert('Ticket not found')
        return
      }
      
      // Set selected ticket and reset form
      selectedTicket.value = ticket
      replyMessage.value = ''
      showReplyModal.value = true
    }

    const submitReply = async () => {
      if (!replyMessage.value.trim()) {
        return
      }

      try {
        replyLoading.value = true
        const response = await adminService.replyToTicket(selectedTicket.value.id, replyMessage.value.trim())
        
        // Update ticket status in local data if it's a mock ticket
        const ticketIndex = tickets.value.findIndex(t => t.id === selectedTicket.value.id)
        if (ticketIndex !== -1 && response.status === 'in_progress') {
          tickets.value[ticketIndex].status = 'in_progress'
        }
        
        // Close modal and reload tickets
        showReplyModal.value = false
        selectedTicket.value = null
        replyMessage.value = ''
        
        showSuccess(t('admin.replySentSuccessfully'))
        await loadTickets() // Reload tickets after reply
      } catch (error) {
        console.error('Error replying to ticket:', error)
        showError(t('admin.errorSendingReply') + ': ' + (error.response?.data?.error || error.message))
      } finally {
        replyLoading.value = false
      }
    }

    const resolveTicket = async (ticketId) => {
      try {
        const confirmed = await confirmAction(t('admin.confirmResolveTicket'), {
          type: 'warning',
          title: t('admin.resolveTicket'),
          confirmText: t('admin.resolve'),
          cancelText: t('common.cancel'),
          confirmClass: 'bg-green-500 hover:bg-green-600',
          cancelClass: 'bg-gray-500 hover:bg-gray-600'
        })
        
        if (!confirmed) {
          return
        }

        const response = await adminService.resolveTicket(ticketId)
        
        // Update ticket status in local data if it's a mock ticket
        const ticketIndex = tickets.value.findIndex(t => t.id == ticketId)
        if (ticketIndex !== -1 && response.status === 'resolved') {
          tickets.value[ticketIndex].status = 'resolved'
        }
        
        showSuccess(t('admin.ticketResolvedSuccessfully'))
        await loadTickets() // Reload tickets after resolution
      } catch (error) {
        console.error('Error resolving ticket:', error)
        showError(t('admin.errorResolvingTicket') + ': ' + (error.response?.data?.error || error.message))
      }
    }

    onMounted(async () => {
      // If user is already loaded and is admin, load data immediately
      if (user.value?.is_admin) {
        loadData()
      }
    })

    return {
      user,
      activeTab,
      transactionFilter,
      ticketFilter,
      users,
      transactions,
      tickets,
      stats,
      tabs,
      getTransactionTypeColor,
      getTicketStatusColor,
      getTicketPriorityColor,
      formatDate,
      addUser,
      editUser,
      deleteUser,
      exportTransactions,
      replyToTicket,
      resolveTicket,
      toggleTwoFactor,
      savePasswordPolicy,
      systemSettings,
      showPasswordPolicyModal,
      loadSettingsFromStorage,
      saveSettingsToStorage,
      // Pagination
      currentPage,
      itemsPerPage,
      paginatedUsers,
      paginatedTransactions,
      paginatedTickets,
      goToPage,
      nextPage,
      prevPage,
      changeItemsPerPage,
      // User management
      newUser,
      addUserLoading,
      showAddUserModal,
      submitAddUser,
      editUser,
      editUserData,
      editUserLoading,
      showEditUserModal,
      submitEditUser,
      selectedUser,
      // Filtered data
      filteredTransactions,
      filteredTickets,
      // Auto-save
      autoSaveTimeout,
      // Reply modal
      showReplyModal,
      selectedTicket,
      replyMessage,
      replyLoading,
      submitReply,
      // Notifications
      notifications,
      notificationQueue,
      removeNotification,
      clearAllNotifications,
      showSuccess,
      showError,
      showWarning,
      showInfo,
      queueSuccess,
      queueError,
      // Confirmations
      pendingConfirmations,
      isPageBlocked,
      confirmAction,
      handleConfirmation
    }
  }
}
</script>

<style scoped>
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.notification-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

.notification-move {
  transition: transform 0.3s ease;
}

/* Ensure notifications are fully visible */
.notification-container {
  position: fixed;
  top: 1rem;
  right: 1rem;
  z-index: 9999;
  pointer-events: none;
  max-width: 24rem;
  width: 100%;
}

@media (min-width: 640px) {
  .notification-container {
    max-width: 28rem;
  }
}

/* Ensure notification content is fully visible */
.notification-item {
  pointer-events: auto;
  margin-bottom: 0.5rem;
}

/* Confirmation animations */
.confirmation-enter-active,
.confirmation-leave-active {
  transition: all 0.3s ease;
}

.confirmation-enter-from {
  opacity: 0;
  transform: scale(0.9) translateY(-20px);
}

.confirmation-leave-to {
  opacity: 0;
  transform: scale(0.9) translateY(-20px);
}

.confirmation-move {
  transition: transform 0.3s ease;
}

/* Page blocking styles */
.page-blocked {
  pointer-events: none;
  user-select: none;
}
</style>
