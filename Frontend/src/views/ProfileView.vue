<template>
  <div
    class="min-h-screen flex flex-col transition-colors duration-200"
    :class="isDark ? 'bg-sport-bg text-white' : 'bg-gray-100 text-gray-900'"
  >
    <Navbar />

    <main class="flex-1 max-w-[1400px] w-full mx-auto px-4 py-8">
      <Profile />
    </main>

    <Footer />
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import Navbar from '@/components/layout/Navbar.vue'
import Footer from '@/components/layout/Footer.vue'
import Profile from '@/components/profile/Profile.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
</script>
