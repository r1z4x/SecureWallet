<template>
  <div class="relative">
    <button
      @click="isOpen = !isOpen"
      class="flex items-center space-x-2 px-3 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
    >
      <span class="w-5 h-5 rounded-full overflow-hidden">
        <img
          :src="getFlagIcon(currentLocale)"
          :alt="getLanguageName(currentLocale)"
          class="w-full h-full object-cover"
        />
      </span>
      <span>{{ getLanguageName(currentLocale) }}</span>
      <svg
        class="w-4 h-4"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M19 9l-7 7-7-7"
        />
      </svg>
    </button>

    <div
      v-if="isOpen"
      class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-50 border border-gray-200"
    >
      <div class="py-1">
        <button
          v-for="locale in availableLocales"
          :key="locale.code"
          @click="changeLanguage(locale.code)"
          class="flex items-center space-x-3 w-full px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 hover:text-gray-900"
          :class="{ 'bg-indigo-50 text-indigo-700': currentLocale === locale.code }"
        >
          <span class="w-5 h-5 rounded-full overflow-hidden">
            <img
              :src="locale.flag"
              :alt="locale.name"
              class="w-full h-full object-cover"
            />
          </span>
          <span>{{ locale.name }}</span>
          <span v-if="currentLocale === locale.code" class="text-indigo-600">
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path
                fill-rule="evenodd"
                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                clip-rule="evenodd"
              />
            </svg>
          </span>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'

export default {
  name: 'LanguageSelector',
  setup() {
    const { locale } = useI18n()
    const isOpen = ref(false)
    const currentLocale = computed(() => locale.value)

    const availableLocales = [
      {
        code: 'en',
        name: 'English',
        flag: 'https://flagcdn.com/w40/gb.png'
      },
      {
        code: 'tr',
        name: 'Türkçe',
        flag: 'https://flagcdn.com/w40/tr.png'
      },
      {
        code: 'es',
        name: 'Español',
        flag: 'https://flagcdn.com/w40/es.png'
      }
    ]

    const getFlagIcon = (code) => {
      const locale = availableLocales.find(l => l.code === code)
      return locale ? locale.flag : availableLocales[0].flag
    }

    const getLanguageName = (code) => {
      const locale = availableLocales.find(l => l.code === code)
      return locale ? locale.name : availableLocales[0].name
    }

    const changeLanguage = (newLocale) => {
      locale.value = newLocale
      localStorage.setItem('locale', newLocale)
      isOpen.value = false
    }

    const handleClickOutside = (event) => {
      if (!event.target.closest('.relative')) {
        isOpen.value = false
      }
    }

    onMounted(() => {
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      isOpen,
      currentLocale,
      availableLocales,
      getFlagIcon,
      getLanguageName,
      changeLanguage
    }
  }
}
</script>

<style scoped>
.relative {
  position: relative;
}

.absolute {
  position: absolute;
}

.z-50 {
  z-index: 50;
}
</style>
