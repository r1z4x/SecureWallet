import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import './style.css'

// Import views
import Landing from './views/Landing.vue'
import Login from './views/Login.vue'
import Register from './views/Register.vue'
import PasswordReset from './views/PasswordReset.vue'
import Dashboard from './views/Dashboard.vue'
import Wallet from './views/Wallet.vue'
import Transactions from './views/Transactions.vue'
import Transfer from './views/Transfer.vue'
import Profile from './views/Profile.vue'
import Support from './views/Support.vue'
import Admin from './views/Admin.vue'
import ResetDatabase from './views/ResetDatabase.vue'
import HelpCenter from './views/HelpCenter.vue'
import TermsOfService from './views/TermsOfService.vue'
import FAQ from './views/FAQ.vue'
import UserGuide from './views/UserGuide.vue'
import SecurityTips from './views/SecurityTips.vue'

// Import stores
import { createPinia } from 'pinia'

// Create router
const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Landing',
      component: Landing,
      meta: { requiresAuth: false }
    },
    {
      path: '/auth/login',
      name: 'Login',
      component: Login,
      meta: { requiresAuth: false }
    },
    {
      path: '/auth/register',
      name: 'Register',
      component: Register,
      meta: { requiresAuth: false }
    },
    {
      path: '/auth/password-reset',
      name: 'PasswordReset',
      component: PasswordReset,
      meta: { requiresAuth: false }
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: Dashboard,
      meta: { requiresAuth: true }
    },
    {
      path: '/wallet',
      name: 'Wallet',
      component: Wallet,
      meta: { requiresAuth: true }
    },
    {
      path: '/transactions',
      name: 'Transactions',
      component: Transactions,
      meta: { requiresAuth: true }
    },
    {
      path: '/transfer',
      name: 'Transfer',
      component: Transfer,
      meta: { requiresAuth: true }
    },
    {
      path: '/profile',
      name: 'Profile',
      component: Profile,
      meta: { requiresAuth: true }
    },
    {
      path: '/support',
      name: 'Support',
      component: Support,
      meta: { requiresAuth: true }
    },
    {
      path: '/admin',
      name: 'Admin',
      component: Admin,
      meta: { requiresAuth: true, requiresAdmin: true }
    },
    {
      path: '/reset-database',
      name: 'ResetDatabase',
      component: ResetDatabase,
      meta: { requiresAuth: true }
    },
    {
      path: '/help',
      name: 'HelpCenter',
      component: HelpCenter,
      meta: { requiresAuth: false }
    },
    {
      path: '/terms',
      name: 'TermsOfService',
      component: TermsOfService,
      meta: { requiresAuth: false }
    },
    {
      path: '/faq',
      name: 'FAQ',
      component: FAQ,
      meta: { requiresAuth: false }
    },
    {
      path: '/user-guide',
      name: 'UserGuide',
      component: UserGuide,
      meta: { requiresAuth: false }
    },
    {
      path: '/security-tips',
      name: 'SecurityTips',
      component: SecurityTips,
      meta: { requiresAuth: false }
    }
  ]
})

// Navigation guard
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  
  if (to.meta.requiresAuth && !token) {
    next('/auth/login')
  } else if (to.path === '/' && token) {
    next('/dashboard')
  } else {
    next()
  }
})

// Create app
const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

app.mount('#app')
