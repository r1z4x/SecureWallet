<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { watch } from 'vue'

export default {
  name: 'App',
  setup() {
    const authStore = useAuthStore()
    
    onMounted(async () => {
      try {
        // Initialize auth state on app mount
        await authStore.init()
        
        // Wait for user data to be loaded if token exists
        if (authStore.token && !authStore.isUserLoaded) {
          await new Promise(resolve => {
            const unwatch = watch(() => authStore.isUserLoaded, (loaded) => {
              if (loaded) {
                unwatch()
                resolve()
              }
            }, { immediate: true })
          })
        }
      } catch (error) {
        console.error('Failed to initialize auth store:', error)
      }
    })
  }
}
</script>

<style>
#app {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}
</style>
