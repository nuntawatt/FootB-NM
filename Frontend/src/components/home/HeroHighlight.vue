<template>
  <section class="relative w-full h-64 sm:h-80 md:h-[400px] overflow-hidden bg-black border-b-2 border-sport-red">
    <img
      :src="currentItem.image"
      :alt="currentItem.title"
      class="w-full h-full object-cover opacity-75 transition-all duration-500"
    />
    <div class="absolute inset-0 bg-gradient-to-t from-black via-black/50 to-black/20" />
    <div class="absolute inset-0 bg-gradient-to-r from-black/80 via-transparent to-transparent" />

    <!-- Prev / Next Arrows -->
    <button
      type="button"
      @click="prevSlide"
      class="hidden sm:flex absolute left-4 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-black/50 hover:bg-sport-red text-white items-center justify-center transition-colors backdrop-blur-sm"
      aria-label="ก่อนหน้า"
    >
      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
      </svg>
    </button>
    <button
      type="button"
      @click="nextSlide"
      class="hidden sm:flex absolute right-4 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-black/50 hover:bg-sport-red text-white items-center justify-center transition-colors backdrop-blur-sm"
      aria-label="ถัดไป"
    >
      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
      </svg>
    </button>

    <!-- Caption -->
    <div class="absolute bottom-0 left-0 right-0 max-w-[1400px] mx-auto px-4 pb-6 pt-12">
      <div class="flex items-center gap-2 mb-2.5">
        <span class="px-2.5 py-0.5 text-xs font-extrabold uppercase tracking-wider bg-sport-red text-white rounded">
          {{ currentItem.category }}
        </span>
        <span class="text-xs text-neutral-300 font-medium">THAI LEAGUE OFFICIAL STYLE (Mock)</span>
      </div>

      <h1 class="text-xl sm:text-3xl md:text-4xl font-black text-white leading-tight max-w-3xl drop-shadow-lg">
        {{ currentItem.title }}
      </h1>

      <div class="flex flex-wrap items-center justify-between gap-4 mt-4">
        <!-- Dots -->
        <div class="flex items-center gap-2">
          <button
            v-for="(item, idx) in items"
            :key="item.id"
            type="button"
            @click="currentIndex = idx"
            class="h-1.5 rounded-full transition-all"
            :class="idx === currentIndex ? 'w-8 bg-sport-red' : 'w-2 bg-white/40 hover:bg-white/70'"
            :aria-label="`Highlight ${idx + 1}`"
          />
        </div>

        <RouterLink
          to="/fixtures"
          class="px-4 py-1.5 rounded-lg text-xs font-bold bg-sport-red hover:bg-sport-red-hover text-white transition-colors shadow-lg"
        >
          ดูโปรแกรมแข่งขัน →
        </RouterLink>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import { mockHeroItems } from '@/mock/hero'

const items = mockHeroItems
const currentIndex = ref(0)
const currentItem = computed(() => items[currentIndex.value])

function nextSlide() {
  currentIndex.value = (currentIndex.value + 1) % items.length
}

function prevSlide() {
  currentIndex.value = (currentIndex.value - 1 + items.length) % items.length
}

let timer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  timer = setInterval(nextSlide, 6000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>
