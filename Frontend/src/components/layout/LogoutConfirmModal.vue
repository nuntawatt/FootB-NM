<template>
  <div
    v-if="open"
    class="fixed inset-0 z-50 bg-black/75 backdrop-blur-sm flex items-center justify-center p-4"
    @click.self="$emit('cancel')"
  >
    <div class="card w-full max-w-sm p-6 space-y-5 border-t-4 border-t-sport-red shadow-2xl text-center">
      <div class="w-12 h-12 rounded-full bg-sport-red/15 text-sport-red flex items-center justify-center mx-auto">
        <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
          />
        </svg>
      </div>

      <div class="space-y-1.5">
        <h3 class="text-lg font-extrabold">ยืนยันการออกจากระบบ</h3>
        <p class="text-xs leading-relaxed" :class="isDark ? 'text-neutral-400' : 'text-gray-500'">
          คุณต้องการออกจากระบบใช่หรือไม่?
        </p>
      </div>

      <div class="flex items-center gap-3 pt-1">
        <button
          type="button"
          @click="$emit('cancel')"
          class="flex-1 px-4 py-2 rounded-lg text-xs font-bold border transition-colors"
          :class="isDark
            ? 'border-sport-border text-neutral-300 hover:bg-sport-hover'
            : 'border-gray-300 text-gray-700 hover:bg-gray-100'"
        >
          ยกเลิก
        </button>
        <button
          type="button"
          @click="$emit('confirm')"
          class="flex-1 btn-primary py-2 text-xs font-bold"
        >
          ออกจากระบบ
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'

defineProps<{
  open: boolean
}>()

defineEmits<{
  cancel: []
  confirm: []
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)
</script>
