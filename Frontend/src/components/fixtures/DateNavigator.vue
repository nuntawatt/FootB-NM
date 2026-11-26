<template>
  <div class="card px-4 py-2.5 flex items-center justify-between gap-3">
    <button
      type="button"
      @click="$emit('prev')"
      class="p-1.5 rounded-lg transition-colors"
      :class="isDark ? 'hover:bg-sport-hover text-neutral-300' : 'hover:bg-gray-100 text-gray-700'"
      aria-label="ก่อนหน้า"
    >
      <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
        <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
      </svg>
    </button>

    <div class="flex items-center gap-3">
      <span class="text-sm sm:text-base font-bold tracking-wide" :class="isDark ? 'text-white' : 'text-gray-900'">
        {{ label }}
      </span>
    </div>

    <div class="flex items-center gap-1">
      <button
        type="button"
        @click="$emit('next')"
        class="p-1.5 rounded-lg transition-colors"
        :class="isDark ? 'hover:bg-sport-hover text-neutral-300' : 'hover:bg-gray-100 text-gray-700'"
        aria-label="ถัดไป"
      >
        <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
        </svg>
      </button>

      <input
        v-if="showDatePicker"
        type="date"
        :value="selectedDate"
        @change="onDateChange"
        class="ml-1 px-2 py-1 text-xs rounded-lg border outline-none cursor-pointer"
        :class="isDark
          ? 'bg-sport-bg border-sport-border text-white'
          : 'bg-gray-50 border-gray-300 text-gray-900'"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useThemeStore } from '@/stores/theme'

defineProps<{
  label: string
  selectedDate?: string
  showDatePicker?: boolean
}>()

const emit = defineEmits<{
  prev: []
  next: []
  'update:date': [date: string]
}>()

const themeStore = useThemeStore()
const { isDark } = storeToRefs(themeStore)

function onDateChange(e: Event) {
  const val = (e.target as HTMLInputElement).value
  emit('update:date', val)
}
</script>
