<template>
  <section
    v-if="items.length > 0"
    @click="navigateToHero"
    class="relative w-full h-64 sm:h-80 md:h-[400px] overflow-hidden bg-black border-b-2 border-sport-red cursor-pointer group"
  >
    <img
      :src="currentItem.image_url"
      :alt="currentItem.title"
      class="w-full h-full object-cover opacity-75 group-hover:scale-[1.02] transition-all duration-500"
    />
    <div class="absolute inset-0 bg-gradient-to-t from-black via-black/50 to-black/20" />
    <div class="absolute inset-0 bg-gradient-to-r from-black/80 via-transparent to-transparent" />

    <!-- Prev / Next Arrows -->
    <button
      type="button"
      @click.stop="prevSlide"
      class="hidden sm:flex absolute left-4 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-black/50 hover:bg-sport-red text-white items-center justify-center transition-colors backdrop-blur-sm z-10"
      aria-label="ก่อนหน้า"
    >
      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
      </svg>
    </button>
    <button
      type="button"
      @click.stop="nextSlide"
      class="hidden sm:flex absolute right-4 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-black/50 hover:bg-sport-red text-white items-center justify-center transition-colors backdrop-blur-sm z-10"
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
          {{ currentItem.tag }}
        </span>
      </div>

      <h1 class="text-xl sm:text-3xl md:text-4xl font-black text-white leading-tight max-w-3xl drop-shadow-lg group-hover:text-sport-red-hover transition-colors">
        {{ currentItem.title }}
      </h1>
      <p v-if="currentItem.subtitle" class="text-xs sm:text-sm text-neutral-300 mt-1.5 max-w-2xl line-clamp-2">
        {{ currentItem.subtitle }}
      </p>

      <div class="flex flex-wrap items-center justify-between gap-4 mt-4">
        <!-- Dots -->
        <div class="flex items-center gap-2" @click.stop>
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

        <span
          class="px-4 py-1.5 rounded-lg text-xs font-bold bg-sport-red group-hover:bg-sport-red-hover text-white transition-colors shadow-lg"
        >
          {{ currentItem.redirect_url || currentItem.match_id ? 'อ่านเพิ่มเติม →' : 'ดูโปรแกรมแข่งขัน →' }}
        </span>
      </div>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { fetchHeroes } from '@/api/heroes'
import type { HeroHighlight } from '@/types'

const router = useRouter()

const items = ref<HeroHighlight[]>([])
const currentIndex = ref(0)
const currentItem = computed(() => items.value[currentIndex.value])

function navigateToHero() {
  const item = currentItem.value
  if (!item) return
  const url = item.redirect_url || (item.match_id ? `/matches/${item.match_id}` : '/fixtures')
  if (/^https?:\/\//i.test(url)) {
    window.open(url, '_blank', 'noopener,noreferrer')
  } else {
    router.push(url)
  }
}

function nextSlide() {
  if (items.value.length === 0) return
  currentIndex.value = (currentIndex.value + 1) % items.value.length
}

function prevSlide() {
  if (items.value.length === 0) return
  currentIndex.value = (currentIndex.value - 1 + items.value.length) % items.value.length
}

let timer: ReturnType<typeof setInterval> | null = null

onMounted(async () => {
  try {
    const apiHeroes = await fetchHeroes()
    if (apiHeroes.length > 0) {
      items.value = apiHeroes
      timer = setInterval(nextSlide, 6000)
    }
  } catch {
    items.value = []
  }
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>
