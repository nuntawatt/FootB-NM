<template>
  <div class="card overflow-hidden">
    <div class="bg-white text-gray-900 px-4 py-3 text-center font-extrabold text-base border-b border-gray-200">
      รายชื่อนักเตะในทีม ({{ players.length }})
    </div>

    <div class="overflow-x-auto">
      <table class="w-full text-sm">
        <thead>
          <tr class="border-b text-xs" :class="isDark ? 'border-sport-border text-neutral-400' : 'border-gray-200 text-gray-500'">
            <th class="py-3 px-4 text-left">ตำแหน่ง</th>
            <th class="py-3 px-3 text-center">หมายเลข</th>
            <th class="py-3 px-4 text-left">ชื่อผู้เล่น</th>
            <th class="py-3 px-4 text-center hidden sm:table-cell">อายุ</th>
            <th class="py-3 px-4 text-left hidden sm:table-cell">สัญชาติ</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="p in players"
            :key="p.id"
            class="border-b last:border-b-0"
            :class="isDark ? 'border-sport-border/50 hover:bg-sport-hover/50' : 'border-gray-100 hover:bg-gray-50'"
          >
            <td class="py-2.5 px-4 font-bold text-sport-red">{{ p.position || '-' }}</td>
            <td class="py-2.5 px-3 text-center tabular-nums font-semibold">{{ p.number ?? '-' }}</td>
            <td class="py-2.5 px-4">
              <div class="flex items-center gap-2.5">
                <img
                  v-if="p.photo_url"
                  :src="p.photo_url"
                  :alt="p.name"
                  class="w-7 h-7 rounded-full object-cover bg-neutral-800 shrink-0"
                />
                <span class="font-medium">{{ p.name }}</span>
              </div>
            </td>
            <td class="py-2.5 px-4 text-center tabular-nums hidden sm:table-cell">{{ p.age ?? '-' }}</td>
            <td class="py-2.5 px-4 hidden sm:table-cell" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              {{ p.nationality || '-' }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { Player } from '@/types'

defineProps<{
  players: Player[]
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
</script>
