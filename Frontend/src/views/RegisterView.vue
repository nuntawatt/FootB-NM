<template>
  <div
    class="min-h-screen flex flex-col transition-colors duration-200"
    :class="isDark ? 'bg-sport-bg text-white' : 'bg-gray-100 text-gray-900'"
  >
    <Navbar />

    <main class="flex-1 flex items-center justify-center px-4 py-10 relative overflow-hidden">
      <div class="absolute w-96 h-96 bg-sport-red/10 rounded-full blur-3xl pointer-events-none" />
      <div class="relative z-10 w-full flex justify-center">
        <Register />
      </div>
    </main>

    <Footer />
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import Navbar from '@/components/layout/Navbar.vue'
import Footer from '@/components/layout/Footer.vue'
import Register from '@/components/login/Register.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
</script>
