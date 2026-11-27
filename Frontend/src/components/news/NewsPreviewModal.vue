<template>
  <div
    v-if="article"
    class="fixed inset-0 z-50 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4"
    @click.self="$emit('close')"
  >
    <div
      class="w-full max-w-md rounded-2xl p-4 sm:p-5 shadow-2xl border space-y-4 transition-all"
      :class="isDark ? 'bg-[#1a1a1a] border-sport-border text-white' : 'bg-white border-gray-200 text-gray-900'"
    >
      <!-- Top Rounded Image -->
      <div class="relative h-56 sm:h-64 rounded-xl overflow-hidden bg-neutral-900">
        <img
          :src="article.image_url"
          :alt="article.title"
          class="w-full h-full object-cover"
        />
        <span
          v-if="article.category"
          class="absolute top-3 left-3 px-2.5 py-0.5 text-[11px] font-bold bg-sport-red text-white rounded"
        >
          {{ article.category }}
        </span>
        <button
          type="button"
          @click="$emit('close')"
          class="absolute top-3 right-3 w-8 h-8 rounded-full bg-black/65 hover:bg-sport-red text-white flex items-center justify-center text-sm transition-colors"
          aria-label="ปิด"
        >
          ✕
        </button>
      </div>

      <!-- Title & Details -->
      <div class="space-y-2 px-1">
        <h3 class="text-lg sm:text-xl font-extrabold leading-snug">
          {{ article.title }}
        </h3>
        <p
          class="text-xs sm:text-sm leading-relaxed whitespace-pre-line"
          :class="isDark ? 'text-neutral-300' : 'text-gray-600'"
        >
          {{ showFull && article.content ? article.content : article.summary }}
        </p>
        <p
          v-if="article.published_at"
          class="text-[11px] pt-1"
          :class="isDark ? 'text-neutral-500' : 'text-gray-400'"
        >
          {{ article.author ? `${article.author} • ` : '' }}{{ article.published_at }}
        </p>
      </div>

      <!-- Bottom Right: อ่านต่อ... -->
      <div class="flex items-center justify-end pt-2 px-1">
        <button
          type="button"
          @click="handleReadMore"
          class="text-sm font-bold transition-colors hover:text-sport-red flex items-center gap-1"
          :class="isDark ? 'text-white' : 'text-gray-900'"
        >
          <span>{{ showFull && !article.redirect_url ? 'ปิดหน้าต่าง' : 'อ่านต่อ...' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import type { NewsArticle } from '@/types'

const props = defineProps<{
  article: NewsArticle | null
  expandInPlace?: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const router = useRouter()
const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const showFull = ref(false)

watch(() => props.article, () => {
  showFull.value = false
})

function handleReadMore() {
  if (!props.article) return

  if (props.article.redirect_url) {
    const url = props.article.redirect_url
    emit('close')
    if (/^https?:\/\//i.test(url)) {
      window.open(url, '_blank', 'noopener,noreferrer')
    } else {
      router.push(url)
    }
    return
  }

  if (props.expandInPlace) {
    if (showFull.value) {
      emit('close')
    } else {
      showFull.value = true
    }
    return
  }

  emit('close')
  router.push('/news')
}
</script>
