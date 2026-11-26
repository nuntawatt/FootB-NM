<template>
  <section class="relative w-full h-64 sm:h-80 md:h-96 overflow-hidden bg-black">
    <img
      :src="currentItem.image"
      :alt="currentItem.title"
      class="w-full h-full object-cover opacity-80 transition-all duration-500"
    />
    <div class="absolute inset-0 bg-gradient-to-t from-black via-black/40 to-transparent" />

    <!-- Caption -->
    <div class="absolute bottom-0 left-0 right-0 max-w-[1400px] mx-auto px-4 pb-6 pt-12">
      <div class="flex items-center gap-2 mb-2">
        <span class="px-2.5 py-0.5 text-xs font-bold uppercase tracking-wider bg-sport-red text-white rounded">
          {{ currentItem.category }}
        </span>
        <span class="text-xs text-neutral-300">ตัวอย่างข้อมูล (Mock Highlight)</span>
      </div>
      <h1 class="text-xl sm:text-3xl md:text-4xl font-extrabold text-white leading-tight max-w-3xl drop-shadow-lg">
        {{ currentItem.title }}
      </h1>

      <!-- Dots -->
      <div class="flex items-center gap-2 mt-4">
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
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { mockHeroItems } from '@/mock/hero'

const items = mockHeroItems
const currentIndex = ref(0)
const currentItem = computed(() => items[currentIndex.value])

let timer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  timer = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % items.length
  }, 6000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>
