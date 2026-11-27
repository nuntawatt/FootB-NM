<template>
  <div class="min-h-screen flex flex-col transition-colors duration-200" :class="isDark ? 'bg-sport-bg text-white' : 'bg-gray-100 text-gray-900'">
    <Navbar />
    <ClubStrip />
    <LiveTicker />

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

    <!-- Scroll to Top Floating Button -->
    <button
      v-if="showScrollTop"
      type="button"
      @click="scrollToTop"
      title="กลับขึ้นด้านบน"
      class="fixed bottom-6 right-6 z-40 w-10 h-10 rounded-full bg-sport-red hover:bg-sport-red-hover text-white shadow-xl shadow-sport-red/30 flex items-center justify-center transition-all hover:-translate-y-0.5"
    >
      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 15l7-7 7 7" />
      </svg>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import Navbar from './Navbar.vue'
import ClubStrip from './ClubStrip.vue'
import LiveTicker from './LiveTicker.vue'
import LeagueSidebar from './LeagueSidebar.vue'
import Footer from './Footer.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const showScrollTop = ref(false)

function onScroll() {
  showScrollTop.value = window.scrollY > 320
}

function scrollToTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))
</script>
