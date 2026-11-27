<template>
  <section class="space-y-3">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <span class="w-1.5 h-5 bg-sport-red rounded-full" />
        <h2 class="text-base sm:text-lg font-extrabold">ข่าวสาร &amp; ไฮไลท์ฟุตบอล</h2>
      </div>
      <RouterLink to="/news" class="text-xs font-semibold text-sport-red hover:underline">
        ดูข่าวทั้งหมด →
      </RouterLink>
    </div>

    <div v-if="loading" class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <div v-for="i in 3" :key="i" class="card p-4 space-y-3">
        <div class="skeleton h-40 w-full" />
        <div class="skeleton h-4 w-3/4" />
        <div class="skeleton h-3 w-full" />
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <article
        v-for="news in newsList"
        :key="news.id"
        @click="selectedArticle = news"
        class="card overflow-hidden group hover:border-sport-red transition-all duration-200 flex flex-col cursor-pointer"
      >
        <div class="relative h-40 overflow-hidden bg-black">
          <img
            :src="news.image_url"
            :alt="news.title"
            class="w-full h-full object-cover opacity-90 group-hover:scale-105 transition-transform duration-300"
          />
          <span class="absolute top-2.5 left-2.5 px-2 py-0.5 text-[10px] font-bold uppercase bg-sport-red text-white rounded">
            {{ news.category }}
          </span>
        </div>

        <div class="p-4 flex-1 flex flex-col justify-between space-y-2">
          <div class="space-y-1.5">
            <h3 class="text-sm font-bold leading-snug line-clamp-2 group-hover:text-sport-red transition-colors">
              {{ news.title }}
            </h3>
            <p v-if="news.summary" class="text-xs line-clamp-2" :class="isDark ? 'text-neutral-400' : 'text-gray-600'">
              {{ news.summary }}
            </p>
          </div>
          <span class="text-[11px] font-semibold text-sport-red">
            อ่านต่อ...
          </span>
        </div>
      </article>
    </div>

    <!-- News Preview Card Modal -->
    <NewsPreviewModal
      :article="selectedArticle"
      @close="selectedArticle = null"
    />
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchNews } from '@/api/news'
import { mockNews } from '@/mock/thaileague'
import type { NewsArticle } from '@/types'
import NewsPreviewModal from '@/components/news/NewsPreviewModal.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const fallbackNews: NewsArticle[] = mockNews.slice(0, 3).map((n) => ({
  id: String(n.id),
  title: n.title,
  summary: n.summary,
  content: n.content,
  image_url: n.image,
  category: n.category,
  published_at: n.date,
  author: n.author,
  redirect_url: n.redirect_url,
}))

const newsList = ref<NewsArticle[]>(fallbackNews)
const loading = ref(false)
const selectedArticle = ref<NewsArticle | null>(null)

onMounted(async () => {
  loading.value = true
  try {
    const apiNews = await fetchNews(3)
    if (apiNews.length > 0) {
      newsList.value = apiNews
    }
  } catch {
    // keep fallback mock
  } finally {
    loading.value = false
  }
})
</script>
