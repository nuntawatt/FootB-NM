<template>
  <MainLayout>
    <div class="space-y-6">
      <!-- Page Header -->
      <div class="flex flex-wrap items-center justify-between gap-3">
        <div class="flex items-center gap-2.5">
          <span class="w-1.5 h-6 bg-sport-red rounded-full" />
          <h1 class="text-xl sm:text-2xl font-extrabold">ข่าวสารฟุตบอล</h1>
        </div>
      </div>

      <!-- Dynamic Category Filter Tabs -->
      <div class="flex flex-wrap items-center gap-2">
        <button
          v-for="cat in categories"
          :key="cat"
          type="button"
          @click="selectedCategory = cat"
          class="px-3.5 py-1.5 rounded-full text-xs font-bold transition-colors"
          :class="selectedCategory === cat
            ? 'bg-sport-red text-white'
            : (isDark ? 'bg-sport-card text-neutral-300 hover:bg-sport-hover border border-sport-border' : 'bg-white text-gray-700 hover:bg-gray-100 border border-gray-200')"
        >
          {{ cat }}
        </button>
      </div>

      <!-- Loading -->
      <div v-if="loading" class="card p-4 space-y-3">
        <div v-for="i in 4" :key="i" class="skeleton h-16 w-full" />
      </div>

      <div v-else-if="error" class="card p-6 text-center space-y-3">
        <p class="text-sm text-sport-red font-medium">{{ error }}</p>
        <button type="button" class="btn-primary text-xs" @click="loadNews">ลองใหม่อีกครั้ง</button>
      </div>

      <template v-else-if="filteredArticles.length > 0">
        <!-- Featured Top Article -->
        <article
          @click="activeArticle = filteredArticles[0]"
          class="card overflow-hidden grid grid-cols-1 md:grid-cols-2 cursor-pointer group hover:border-sport-red transition-all border-t-2 border-t-sport-red"
        >
          <div class="relative h-56 md:h-72 overflow-hidden bg-black">
            <img
              :src="filteredArticles[0].image_url"
              :alt="filteredArticles[0].title"
              class="w-full h-full object-cover opacity-90 group-hover:scale-105 transition-transform duration-300"
            />
            <span class="absolute top-3 left-3 px-2.5 py-0.5 text-xs font-bold bg-sport-red text-white rounded">
              {{ filteredArticles[0].category }}
            </span>
          </div>

          <div class="p-6 flex flex-col justify-between space-y-4">
            <div class="space-y-2.5">
              <span class="text-xs font-bold text-sport-red uppercase tracking-wider">ข่าวเด่น</span>
              <h2 class="text-lg sm:text-2xl font-black leading-snug group-hover:text-sport-red transition-colors">
                {{ filteredArticles[0].title }}
              </h2>
              <p class="text-sm leading-relaxed" :class="isDark ? 'text-neutral-300' : 'text-gray-600'">
                {{ filteredArticles[0].summary }}
              </p>
            </div>
            <div class="flex items-center justify-end text-xs">
              <span class="font-bold text-sport-red">อ่านต่อ →</span>
            </div>
          </div>
        </article>

        <!-- Remaining Articles Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5">
          <article
            v-for="item in filteredArticles.slice(1)"
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
                <p v-if="item.summary" class="text-xs line-clamp-2" :class="isDark ? 'text-neutral-400' : 'text-gray-600'">
                  {{ item.summary }}
                </p>
              </div>
              <div class="flex items-center justify-end text-[11px]">
                <span class="text-sport-red font-semibold">อ่านต่อ →</span>
              </div>
            </div>
          </article>
        </div>
      </template>

      <div v-else class="card p-8 text-center text-sm" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
        ไม่พบข่าวสารในหมวดหมู่ที่เลือก
      </div>

      <!-- News Preview Card Modal -->
      <NewsPreviewModal
        :article="activeArticle"
        expand-in-place
        @close="activeArticle = null"
      />
    </div>
  </MainLayout>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'
import { fetchNews } from '@/api/news'
import type { NewsArticle } from '@/types'
import MainLayout from '@/components/layout/MainLayout.vue'
import NewsPreviewModal from '@/components/news/NewsPreviewModal.vue'

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

const selectedCategory = ref('ทั้งหมด')
const allArticles = ref<NewsArticle[]>([])
const loading = ref(false)
const error = ref('')
const activeArticle = ref<NewsArticle | null>(null)

const categories = computed(() => {
  const tags = Array.from(new Set(allArticles.value.map((a) => a.category).filter(Boolean)))
  return ['ทั้งหมด', ...tags]
})

const filteredArticles = computed(() => {
  if (selectedCategory.value === 'ทั้งหมด') return allArticles.value
  return allArticles.value.filter((a) => a.category === selectedCategory.value)
})

async function loadNews() {
  loading.value = true
  error.value = ''
  try {
    allArticles.value = await fetchNews(20)
  } catch {
    allArticles.value = []
    error.value = 'ไม่สามารถโหลดข่าวสารจาก Backend API ได้'
  } finally {
    loading.value = false
  }
}

onMounted(loadNews)
</script>
