<template>
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <!-- Team Info Card -->
    <div class="md:col-span-2 card p-5 flex flex-wrap items-center justify-between gap-4">
      <div class="flex items-center gap-5 min-w-0 flex-1">
        <img
          :src="team.logo_url"
          :alt="team.name"
          class="w-20 h-20 sm:w-24 sm:h-24 object-contain shrink-0"
        />
        <div class="space-y-1 min-w-0">
          <h1 class="text-xl sm:text-2xl font-extrabold truncate" :class="isDark ? 'text-white' : 'text-gray-900'">
            {{ team.name }}
          </h1>
          <p class="text-sm" :class="isDark ? 'text-neutral-300' : 'text-gray-600'">
            {{ team.league?.name || 'Thai League 1' }}
          </p>
          <p class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            {{ team.league?.country || 'Thailand' }}
          </p>
        </div>
      </div>

      <!-- Add / Remove Favorite Button -->
      <button
        type="button"
        :disabled="favLoading"
        @click="toggleFavorite"
        class="shrink-0 inline-flex items-center gap-2 px-4 py-2 rounded-lg text-xs sm:text-sm font-bold border transition-all disabled:opacity-50"
        :class="isFavorite
          ? 'bg-sport-red border-sport-red text-white hover:bg-sport-red-hover shadow-lg shadow-sport-red/20'
          : (isDark
              ? 'border-sport-red text-sport-red hover:bg-sport-red hover:text-white'
              : 'border-sport-red text-sport-red hover:bg-sport-red hover:text-white')"
      >
        <svg class="w-4 h-4" :fill="isFavorite ? 'currentColor' : 'none'" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"
          />
        </svg>
        <span>{{ isFavorite ? 'ทีมโปรดแล้ว' : 'เพิ่มเป็นทีมโปรด' }}</span>
      </button>
    </div>

    <!-- Coach Card -->
    <div class="card p-4 flex flex-col items-center justify-center text-center">
      <img
        :src="coachImage"
        :alt="coachName"
        class="w-20 h-20 rounded-xl object-cover bg-neutral-800 mb-2 border border-sport-border"
        @error="onImgError"
      />
      <p class="font-bold text-sm" :class="isDark ? 'text-white' : 'text-gray-900'">{{ coachName }}</p>
      <p class="text-[11px]" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        หัวหน้าผู้ฝึกสอน{{ coachNationality ? ` (${coachNationality})` : '' }}
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchFavorites, addFavorite, removeFavorite } from '@/api/favorites'
import type { Team } from '@/types'

const props = defineProps<{
  team: Team
  coachName: string
  coachImage: string
  coachNationality?: string
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const isFavorite = ref(false)
const favLoading = ref(false)

async function checkFavorite() {
  try {
    const favs = await fetchFavorites()
    isFavorite.value = favs.some((t) => t.id === props.team.id)
  } catch {
    isFavorite.value = false
  }
}

async function toggleFavorite() {
  favLoading.value = true
  try {
    if (isFavorite.value) {
      await removeFavorite(props.team.id)
      isFavorite.value = false
    } else {
      await addFavorite(props.team.id)
      isFavorite.value = true
    }
  } catch {
    // fallback toggle UI if backend is offline
    isFavorite.value = !isFavorite.value
  } finally {
    favLoading.value = false
  }
}

watch(() => props.team.id, checkFavorite)
onMounted(checkFavorite)

function onImgError(e: Event) {
  (e.target as HTMLImageElement).src = 'https://ui-avatars.com/api/?name=Coach&background=0F172A&color=fff&size=256'
}
</script>
