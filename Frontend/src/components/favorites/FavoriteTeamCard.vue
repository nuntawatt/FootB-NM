<template>
  <RouterLink
    :to="`/teams/${item.team.id}`"
    class="group relative rounded-2xl p-5 flex flex-col justify-between transition-all duration-200 border overflow-hidden
           bg-gradient-to-br from-[#0f172a] to-[#0a0f1f] border-sport-border hover:border-sport-red hover:-translate-y-0.5 shadow-lg"
  >
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
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import type { FavoriteTeamCard } from '@/types'

defineProps<{
  item: FavoriteTeamCard
}>()
</script>
