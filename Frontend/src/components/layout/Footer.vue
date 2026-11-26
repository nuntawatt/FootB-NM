<template>
  <footer
    class="border-t mt-12 transition-colors duration-200"
    :class="isDark ? 'bg-[#050505] border-sport-border text-neutral-400' : 'bg-white border-gray-200 text-gray-600'"
  >
    <div class="h-1 bg-gradient-to-r from-sport-red-dark via-sport-red to-sport-red-dark" />
    <div class="max-w-[1400px] mx-auto px-4 py-8 flex flex-col md:flex-row items-center justify-between gap-4">
      <div class="flex items-center gap-3">
        <div class="w-8 h-8 rounded-lg bg-sport-red flex items-center justify-center">
          <span class="text-white font-black text-sm">FB</span>
        </div>
        <div>
          <p class="text-sm font-extrabold" :class="isDark ? 'text-white' : 'text-gray-900'">
            FOOTB — Thai League Hub
          </p>
          <p class="text-xs">ผลบอลสด โปรแกรมแข่งขัน ตารางคะแนน ฟุตบอลไทยลีก</p>
        </div>
      </div>

      <nav class="flex flex-wrap items-center gap-5 text-xs font-semibold">
        <RouterLink to="/" class="hover:text-sport-red transition-colors">หน้าแรก</RouterLink>
        <RouterLink to="/fixtures" class="hover:text-sport-red transition-colors">โปรแกรมแข่งขัน</RouterLink>
        <RouterLink to="/standings" class="hover:text-sport-red transition-colors">ตารางคะแนน</RouterLink>
        <RouterLink to="/favorites" class="hover:text-sport-red transition-colors">ทีมโปรด</RouterLink>
      </nav>
    </div>
  </footer>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
</script>
