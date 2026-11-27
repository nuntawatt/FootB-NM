<template>
  <MainLayout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-xl sm:text-2xl font-extrabold">ทีมโปรดของคุณ</h1>
          <p class="text-xs mt-1" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            {{ authStore.isAuthenticated ? `บัญชี: ${authStore.user?.username}` : 'โหมดผู้เยี่ยมชม (บันทึกในฐานข้อมูลด้วย Guest ID หรือล็อกอินเพื่อซิงก์บัญชี)' }}
          </p>
        </div>
      </div>

      <div v-if="loading" class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        กำลังโหลดทีมโปรด...
      </div>

      <!-- 3-Column Grid -->
      <div v-else-if="favorites.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
        <FavoriteTeamCard
          v-for="item in favorites"
          :key="item.team.id"
          :item="item"
          @remove="handleRemoveFavorite"
        />
      </div>

      <div v-else class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        ยังไม่มีทีมโปรด กดปุ่มด้านล่างเพื่อเพิ่มทีมที่ชอบได้เลยครับ
      </div>

      <!-- Add Team Button -->
      <div class="flex justify-center pt-2">
        <button
          type="button"
          @click="showModal = true"
          class="flex items-center gap-2 px-8 py-3 rounded-full font-bold text-base shadow-lg transition-all
                 bg-white text-gray-900 hover:bg-sport-red hover:text-white"
        >
          <span class="text-xl leading-none">+</span>
          <span>เพิ่มทีมที่ชอบ</span>
        </button>
      </div>

      <!-- Modal เลือกทีม -->
      <div
        v-if="showModal"
        class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
        @click.self="showModal = false"
      >
        <div class="card w-full max-w-md p-5 space-y-4 max-h-[80vh] flex flex-col">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold">เลือกทีมที่ชอบ</h2>
            <button type="button" class="text-sm text-neutral-400 hover:text-white" @click="showModal = false">✕</button>
          </div>

          <div class="overflow-y-auto space-y-2 flex-1">
            <button
              v-for="team in availableTeams"
              :key="team.id"
              type="button"
              @click="handleAddFavorite(team)"
              class="w-full flex items-center gap-3 p-2.5 rounded-lg border transition-colors text-left"
              :class="isDark
                ? 'border-sport-border hover:bg-sport-hover hover:border-sport-red'
                : 'border-gray-200 hover:bg-gray-50 hover:border-sport-red'"
            >
              <img :src="team.logo_url" :alt="team.name" class="w-8 h-8 object-contain" />
              <span class="font-semibold text-sm">{{ team.name }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { useAuthStore } from '@/stores/auth'
import { fetchTeams } from '@/api/teams'
import { fetchFavorites, addFavorite, removeFavorite } from '@/api/favorites'
import { mockFavorites } from '@/mock/favorites'
import type { FavoriteTeamCard as FavCardType, Team } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import FavoriteTeamCard from '@/components/favorites/FavoriteTeamCard.vue'

const themeStore = useThemeStore()
const authStore = useAuthStore()
const { isDark } = storeToRefs(themeStore)

const favorites = ref<FavCardType[]>([])
const availableTeams = ref<Team[]>([])
const showModal = ref(false)
const loading = ref(false)

function toCard(team: Team): FavCardType {
  return {
    team,
    venue: `${team.name} Stadium`,
    next_match_date: 'ติดตามโปรแกรมแข่งขัน',
    broadcast: 'TrueID / AIS PLAY',
  }
}

async function loadFavorites() {
  loading.value = true
  try {
    const teams = await fetchFavorites()
    favorites.value = teams.map(toCard)
  } catch {
    favorites.value = [...mockFavorites]
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  await loadFavorites()
  try {
    availableTeams.value = await fetchTeams()
  } catch {
    availableTeams.value = mockFavorites.map((f) => f.team)
  }
})

watch(() => authStore.isAuthenticated, loadFavorites)

async function handleAddFavorite(team: Team) {
  if (!favorites.value.some((f) => f.team.id === team.id)) {
    try {
      await addFavorite(team.id)
    } catch {
      // still update UI if backend offline
    }
    favorites.value.push(toCard(team))
  }
  showModal.value = false
}

async function handleRemoveFavorite(teamId: string) {
  try {
    await removeFavorite(teamId)
  } catch {
    // ignore if offline
  }
  favorites.value = favorites.value.filter((f) => f.team.id !== teamId)
}
</script>
