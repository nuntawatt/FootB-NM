<template>
  <div class="min-h-screen flex flex-col transition-colors duration-200" :class="isDark ? 'bg-sport-bg text-white' : 'bg-gray-100 text-gray-900'">
    <Navbar />
    <ClubStrip />

    <!-- Optional Hero slot above sidebar + content -->
    <slot name="hero" />

    <!-- Main 2-column layout -->
    <main class="flex-1 max-w-[1400px] w-full mx-auto px-4 py-6">
      <div class="flex flex-col lg:flex-row gap-6 items-start">
        <LeagueSidebar />
        <div class="flex-1 min-w-0 w-full">
          <slot />
        </div>
      </div>
    </main>

    <Footer />
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import Navbar from './Navbar.vue'
import ClubStrip from './ClubStrip.vue'
import LeagueSidebar from './LeagueSidebar.vue'
import Footer from './Footer.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
</script>
