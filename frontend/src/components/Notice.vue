<template>
  <div v-if="message" :style="styleObj">
    {{ message }}
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{ message: string; type?: 'error' | 'success' | 'info' }>()

const styleObj = computed(() => {
  const base = {
    padding: '8px 12px',
    borderRadius: '6px',
    marginBottom: '8px',
  }
  if (props.type === 'success') return { ...base, color: '#14532d', background: '#dcfce7' }
  if (props.type === 'error') return { ...base, color: '#7f1d1d', background: '#fee2e2' }
  return { ...base, color: '#334155', background: '#e2e8f0' }
})
</script>
