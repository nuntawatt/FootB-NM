<template>
  <header
    class="sticky top-0 z-50 border-b transition-colors duration-200"
    :class="isDark ? 'bg-sport-bg/95 border-sport-border backdrop-blur-md' : 'bg-white/95 border-gray-200 backdrop-blur-md'"
  >
    <div class="max-w-[1400px] mx-auto px-4 h-16 flex items-center justify-between gap-4">
      <!-- Logo & Search -->
      <div class="flex items-center gap-4 flex-1 min-w-0">
        <RouterLink to="/" class="flex items-center gap-2 shrink-0 group">
          <div class="w-9 h-9 rounded-lg bg-sport-red flex items-center justify-center shadow-lg shadow-sport-red/30 group-hover:bg-sport-red-hover transition-colors">
            <span class="text-white font-black text-lg tracking-tighter">FB</span>
          </div>
          <div class="hidden sm:flex flex-col leading-none">
            <span class="font-extrabold text-sm tracking-wide uppercase" :class="isDark ? 'text-white' : 'text-gray-900'">
              FOOTB
            </span>
            <span class="text-[10px] font-semibold text-sport-red tracking-widest uppercase">
              THAI LEAGUE
            </span>
          </div>
        </RouterLink>

        <!-- Search input -->
        <div class="relative flex-1 max-w-xs">
          <svg
            class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 pointer-events-none"
            :class="isDark ? 'text-neutral-400' : 'text-gray-500'"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="ค้นหาทีมที่ต้องการ..."
            class="w-full pl-9 pr-3 py-1.5 rounded-full text-sm outline-none transition-all border"
            :class="isDark
              ? 'bg-white text-gray-900 border-transparent placeholder-gray-500 focus:ring-2 focus:ring-sport-red'
              : 'bg-gray-100 text-gray-900 border-gray-300 placeholder-gray-500 focus:ring-2 focus:ring-sport-red'"
            @focus="showResults = true"
            @blur="handleBlur"
          />
          <!-- Search Dropdown -->
          <div
            v-if="showResults && filteredTeams.length > 0"
            class="absolute left-0 right-0 mt-2 rounded-xl shadow-2xl border max-h-72 overflow-y-auto z-50"
            :class="isDark ? 'bg-sport-card border-sport-border' : 'bg-white border-gray-200'"
          >
            <RouterLink
              v-for="team in filteredTeams"
              :key="team.id"
              :to="`/teams/${team.id}`"
              class="flex items-center gap-3 px-3 py-2.5 text-sm transition-colors"
              :class="isDark ? 'hover:bg-sport-hover text-white' : 'hover:bg-gray-100 text-gray-900'"
              @mousedown.prevent="goToTeam(team.id)"
            >
              <img :src="team.logo_url" :alt="team.name" class="w-6 h-6 object-contain shrink-0" />
              <span class="truncate font-medium">{{ team.name }}</span>
            </RouterLink>
          </div>
        </div>
      </div>

      <!-- Desktop Navigation -->
      <nav class="hidden md:flex items-center gap-1 lg:gap-2">
        <RouterLink
          v-for="item in navItems"
          :key="item.to"
          :to="item.to"
          class="px-3 py-2 text-sm font-semibold rounded-lg transition-colors relative"
          :class="isActive(item.to)
            ? (isDark ? 'text-white' : 'text-sport-red')
            : (isDark ? 'text-neutral-300 hover:text-white hover:bg-sport-card' : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100')"
        >
          {{ item.label }}
          <span
            v-if="isActive(item.to)"
            class="absolute bottom-0.5 left-3 right-3 h-0.5 bg-sport-red rounded-full"
          />
        </RouterLink>
      </nav>

      <!-- Right Controls: Theme Toggle + User + Mobile Menu -->
      <div class="flex items-center gap-3 shrink-0">
        <!-- Dark/Light Toggle -->
        <button
          type="button"
          @click="themeStore.toggle()"
          :title="isDark ? 'สลับเป็นโหมดสว่าง' : 'สลับเป็นโหมดมืด'"
          class="p-2 rounded-lg border transition-colors"
          :class="isDark
            ? 'border-sport-border bg-sport-card text-neutral-300 hover:text-white hover:border-sport-red'
            : 'border-gray-300 bg-gray-100 text-gray-700 hover:text-sport-red hover:border-sport-red'"
        >
          <svg v-if="isDark" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
          </svg>
        </button>

        <!-- User Profile (Mock) -->
        <div class="hidden sm:flex items-center gap-2">
          <span class="text-xs font-medium" :class="isDark ? 'text-neutral-400' : 'text-gray-600'">username</span>
          <div class="w-8 h-8 rounded-full flex items-center justify-center text-[10px] font-bold border"
            :class="isDark ? 'bg-neutral-700 text-neutral-300 border-sport-border' : 'bg-gray-200 text-gray-700 border-gray-300'">
            U
          </div>
        </div>

        <!-- Mobile Menu Button -->
        <button
          type="button"
          class="md:hidden p-2 rounded-lg border"
          :class="isDark ? 'border-sport-border text-neutral-300' : 'border-gray-300 text-gray-700'"
          @click="mobileMenuOpen = !mobileMenuOpen"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile Nav Drawer -->
    <div
      v-if="mobileMenuOpen"
      class="md:hidden border-t px-4 py-3 space-y-1"
      :class="isDark ? 'bg-sport-card border-sport-border' : 'bg-white border-gray-200'"
    >
      <RouterLink
        v-for="item in navItems"
        :key="item.to"
        :to="item.to"
        class="block px-3 py-2 rounded-lg text-sm font-semibold"
        :class="isActive(item.to)
          ? 'bg-sport-red text-white'
          : (isDark ? 'text-neutral-300 hover:bg-sport-hover' : 'text-gray-700 hover:bg-gray-100')"
        @click="mobileMenuOpen = false"
      >
        {{ item.label }}
      </RouterLink>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchTeams } from '@/api/teams'
import type { Team } from '@/types'

const route = useRoute()
const router = useRouter()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const navItems = [
  { label: 'หน้าแรก', to: '/' },
  { label: 'โปรแกรมแข่งขัน', to: '/fixtures' },
  { label: 'ตารางคะแนน', to: '/standings' },
  { label: 'ทีมโปรด', to: '/favorites' },
]

const searchQuery = ref('')
const showResults = ref(false)
const mobileMenuOpen = ref(false)
const allTeams = ref<Team[]>([])

onMounted(async () => {
  try {
    allTeams.value = await fetchTeams()
  } catch {
    // fallback quietly if backend isn't running
  }
})

const filteredTeams = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return []
  return allTeams.value
    .filter((t) => t.name.toLowerCase().includes(q) || t.short_name?.toLowerCase().includes(q))
    .slice(0, 8)
})

function isActive(path: string) {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}

function goToTeam(id: string) {
  showResults.value = false
  searchQuery.value = ''
  router.push(`/teams/${id}`)
}

function handleBlur() {
  setTimeout(() => {
    showResults.value = false
  }, 150)
}
</script>
