<template>
  <MainLayout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div>
          <h1 class="text-xl sm:text-2xl font-extrabold">ทีมโปรดของคุณ</h1>
          <p class="text-xs mt-1" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
            แสดงตัวอย่างทีมโปรด (Mock Data — รอระบบ User Auth & Favorite API)
          </p>
        </div>
      </div>

      <!-- 3-Column Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
        <FavoriteTeamCard v-for="item in favorites" :key="item.team.id" :item="item" />
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

      <!-- Modal เลือกทีม (ดึงจาก /teams จริงถ้า backend ทำงาน) -->
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
              @click="addFavorite(team)"
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
import { ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchTeams } from '@/api/teams'
import { mockFavorites } from '@/mock/favorites'
import type { FavoriteTeamCard as FavCardType, Team } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import FavoriteTeamCard from '@/components/favorites/FavoriteTeamCard.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const favorites = ref<FavCardType[]>([...mockFavorites])
const availableTeams = ref<Team[]>([])
const showModal = ref(false)

onMounted(async () => {
  try {
    availableTeams.value = await fetchTeams()
  } catch {
    availableTeams.value = mockFavorites.map((f) => f.team)
  }
})

function addFavorite(team: Team) {
  if (!favorites.value.some((f) => f.team.id === team.id)) {
    favorites.value.push({
      team,
      venue: 'สนามเหย้า (Mock)',
      next_match_date: 'เร็วๆ นี้',
      broadcast: 'MONOMAX / AIS PLAY',
    })
  }
  showModal.value = false
}
</script>
