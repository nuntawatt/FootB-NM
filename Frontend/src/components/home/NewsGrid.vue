<template>
  <section class="space-y-3">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">
        <span class="w-1.5 h-5 bg-sport-red rounded-full" />
        <h2 class="text-base sm:text-lg font-extrabold">ข่าวสาร &amp; ไฮไลท์ไทยลีก</h2>
      </div>
      <RouterLink to="/news" class="text-xs font-semibold text-sport-red hover:underline">
        ดูข่าวทั้งหมด →
      </RouterLink>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <article
        v-for="news in newsList"
        :key="news.id"
        @click="openPreview(news)"
        class="card overflow-hidden group hover:border-sport-red transition-all duration-200 flex flex-col cursor-pointer"
      >
        <div class="relative h-40 overflow-hidden bg-black">
          <img
            :src="news.image"
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
            <p class="text-xs line-clamp-2" :class="isDark ? 'text-neutral-400' : 'text-gray-600'">
              {{ news.summary }}
            </p>
          </div>
          <span class="text-[11px]" :class="isDark ? 'text-neutral-500' : 'text-gray-400'">
            {{ news.date }}
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
import { ref } from 'vue'
import { RouterLink } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { mockNews, type NewsItem } from '@/mock/thaileague'
import type { NewsArticle } from '@/types'
import NewsPreviewModal from '@/components/news/NewsPreviewModal.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
const newsList = mockNews.slice(0, 3)

const selectedArticle = ref<NewsArticle | null>(null)

function openPreview(news: NewsItem) {
  selectedArticle.value = {
    id: String(news.id),
    title: news.title,
    summary: news.summary,
    content: news.content,
    image_url: news.image,
    category: news.category,
    published_at: news.date,
    author: news.author,
    redirect_url: news.redirect_url,
  }
}
</script>
