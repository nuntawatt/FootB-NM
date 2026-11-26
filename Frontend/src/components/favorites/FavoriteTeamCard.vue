<template>
  <div
    class="group relative rounded-2xl p-5 flex flex-col justify-between transition-all duration-200 border overflow-hidden
           bg-gradient-to-br from-[#0f172a] to-[#0a0f1f] border-sport-border hover:border-sport-red hover:-translate-y-0.5 shadow-lg"
  >
    <!-- Remove Favorite Button -->
    <button
      type="button"
      title="ลบออกจากทีมโปรด"
      @click.prevent="$emit('remove', item.team.id)"
      class="absolute top-3 right-3 w-7 h-7 rounded-full bg-black/50 hover:bg-sport-red text-neutral-300 hover:text-white flex items-center justify-center text-xs transition-colors z-10"
    >
      ✕
    </button>

    <RouterLink :to="`/teams/${item.team.id}`" class="block">
      <div class="flex justify-center py-3">
        <img
          :src="item.team.logo_url"
          :alt="item.team.name"
          class="w-24 h-24 object-contain group-hover:scale-105 transition-transform duration-200"
        />
      </div>

      <div class="space-y-1.5 mt-2">
        <h3 class="text-base font-bold text-white truncate">{{ item.team.name }}</h3>
        <div class="flex items-center gap-2 text-xs text-neutral-300">
          <span>🏟️</span>
          <span class="truncate">{{ item.venue }}</span>
        </div>
        <div class="flex items-center gap-2 text-xs text-neutral-300">
          <span>🕒</span>
          <span class="truncate">{{ item.next_match_date }}</span>
        </div>
        <div class="flex items-center gap-2 text-xs text-neutral-300">
          <span>📺</span>
          <span class="truncate">{{ item.broadcast }}</span>
        </div>
      </div>
    </RouterLink>
  </div>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { FavoriteTeamCard } from '@/types'

defineProps<{
  item: FavoriteTeamCard
}>()

defineEmits<{
  remove: [teamId: string]
}>()
</script>
