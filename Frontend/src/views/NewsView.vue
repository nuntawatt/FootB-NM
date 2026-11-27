<template>
  <MainLayout>
    <div class="space-y-6">
      <!-- Page Header -->
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="flex items-center gap-2.5">
          <span class="w-1.5 h-6 bg-sport-red rounded-full" />
          <h1 class="text-xl sm:text-2xl font-extrabold">ข่าวสารฟุตบอลไทย</h1>
        </div>
        <span
          v-if="usingMock"
          class="text-xs px-2.5 py-1 rounded bg-amber-500/20 text-amber-400 font-semibold"
        >
          Mock Data (รอ API GET /news)
        </span>
      </div>

      <!-- Category Filter Tabs -->
      <div class="flex flex-wrap items-center gap-2">
        <button
          v-for="cat in categories"
          :key="cat"
          type="button"
          @click="selectedCategory = cat; loadNews()"
          class="px-3.5 py-1.5 rounded-full text-xs font-bold transition-colors"
          :class="selectedCategory === cat
            ? 'bg-sport-red text-white'
            : (isDark ? 'bg-sport-card text-neutral-300 hover:bg-sport-hover border border-sport-border' : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-200')"
        >
          {{ cat }}
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        กำลังโหลดข่าวสาร...
      </div>

      <template v-else-if="articles.length > 0">
        <!-- Featured Top Article -->
        <article
          @click="activeArticle = articles[0]"
          class="card overflow-hidden grid grid-cols-1 md:grid-cols-2 cursor-pointer group hover:border-sport-red transition-all border-t-2 border-t-sport-red"
        >
          <div class="relative h-56 md:h-72 overflow-hidden bg-black">
            <img
              :src="articles[0].image_url"
              :alt="articles[0].title"
              class="w-full h-full object-cover opacity-90 group-hover:scale-105 transition-transform duration-300"
            />
            <span class="absolute top-3 left-3 px-2.5 py-0.5 text-xs font-bold bg-sport-red text-white rounded">
              {{ articles[0].category }}
            </span>
          </div>

          <div class="p-6 flex flex-col justify-between space-y-4">
            <div class="space-y-2.5">
              <span class="text-xs font-bold text-sport-red uppercase tracking-wider">ข่าวเด่น</span>
              <h2 class="text-lg sm:text-2xl font-black leading-snug group-hover:text-sport-red transition-colors">
                {{ articles[0].title }}
              </h2>
              <p class="text-sm leading-relaxed" :class="isDark ? 'text-neutral-300' : 'text-gray-600'">
                {{ articles[0].summary }}
              </p>
            </div>
            <div class="flex items-center justify-between text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              <span>{{ articles[0].author || 'FOOTB' }} • {{ articles[0].published_at }}</span>
              <span class="font-bold text-sport-red">อ่านต่อ →</span>
            </div>
          </div>
        </article>

        <!-- Remaining Articles Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          <article
            v-for="item in articles.slice(1)"
            :key="item.id"
            @click="activeArticle = item"
            class="card overflow-hidden cursor-pointer group hover:border-sport-red transition-all flex flex-col"
          >
            <div class="relative h-44 overflow-hidden bg-black">
              <img
                :src="item.image_url"
                :alt="item.title"
                class="w-full h-full object-cover opacity-90 group-hover:scale-105 transition-transform duration-300"
              />
              <span class="absolute top-2.5 left-2.5 px-2 py-0.5 text-[10px] font-bold bg-sport-red text-white rounded">
                {{ item.category }}
              </span>
            </div>

            <div class="p-4 flex-1 flex flex-col justify-between space-y-3">
              <div class="space-y-1.5">
                <h3 class="text-sm font-bold leading-snug line-clamp-2 group-hover:text-sport-red transition-colors">
                  {{ item.title }}
                </h3>
                <p class="text-xs line-clamp-2" :class="isDark ? 'text-neutral-400' : 'text-gray-600'">
                  {{ item.summary }}
                </p>
              </div>
              <div class="flex items-center justify-between text-[11px]" :class="isDark ? 'text-neutral-500' : 'text-gray-400'">
                <span>{{ item.published_at }}</span>
                <span class="text-sport-red font-semibold">อ่านต่อ →</span>
              </div>
            </div>
          </article>
        </div>
      </template>

      <div v-else class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        ไม่พบข่าวสารในหมวดหมู่ที่เลือก
      </div>

      <!-- Article Reader Modal -->
      <div
        v-if="activeArticle"
        class="fixed inset-0 z-50 bg-black/80 backdrop-blur-sm flex items-center justify-center p-4"
        @click.self="activeArticle = null"
      >
        <div class="card w-full max-w-2xl max-h-[85vh] overflow-y-auto border-t-4 border-t-sport-red">
          <div class="relative h-56 sm:h-64 bg-black">
            <img :src="activeArticle.image_url" :alt="activeArticle.title" class="w-full h-full object-cover" />
            <button
              type="button"
              @click="activeArticle = null"
              class="absolute top-3 right-3 w-8 h-8 rounded-full bg-black/70 text-white hover:bg-sport-red flex items-center justify-center text-sm"
            >
              ✕
            </button>
            <span class="absolute bottom-3 left-4 px-2.5 py-0.5 text-xs font-bold bg-sport-red text-white rounded">
              {{ activeArticle.category }}
            </span>
          </div>

          <div class="p-6 space-y-4">
            <h2 class="text-xl sm:text-2xl font-black leading-snug">{{ activeArticle.title }}</h2>
            <p class="text-xs" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
              โดย {{ activeArticle.author || 'FOOTB' }} • {{ activeArticle.published_at }}
            </p>
            <p class="text-sm leading-relaxed whitespace-pre-line" :class="isDark ? 'text-neutral-200' : 'text-gray-700'">
              {{ activeArticle.content }}
            </p>
          </div>
        </div>
      </div>
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchNews } from '@/api/news'
import { mockNews } from '@/mock/thaileague'
import type { NewsArticle } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const categories = ['ทั้งหมด', 'Thai League 1', 'Thai League 2', 'FA CUP', 'พรีวิว']
const selectedCategory = ref('ทั้งหมด')
const articles = ref<NewsArticle[]>([])
const loading = ref(false)
const usingMock = ref(false)
const activeArticle = ref<NewsArticle | null>(null)

function mapMockNews(cat: string): NewsArticle[] {
  const list = mockNews.map((n) => ({
    id: String(n.id),
    title: n.title,
    summary: n.summary,
    content: n.content,
    image_url: n.image,
    category: n.category,
    published_at: n.date,
    author: n.author,
  }))
  if (cat === 'ทั้งหมด') return list
  return list.filter((a) => a.category === cat)
}

async function loadNews() {
  loading.value = true
  usingMock.value = false
  const catParam = selectedCategory.value === 'ทั้งหมด' ? undefined : selectedCategory.value
  try {
    const data = await fetchNews(catParam)
    if (Array.isArray(data) && data.length > 0) {
      articles.value = data
    } else {
      articles.value = mapMockNews(selectedCategory.value)
      usingMock.value = true
    }
  } catch {
    articles.value = mapMockNews(selectedCategory.value)
    usingMock.value = true
  } finally {
    loading.value = false
  }
}

onMounted(loadNews)
</script>
