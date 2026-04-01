<script setup>
import { ref, onMounted } from 'vue'

const emit = defineEmits(['complete'])
const bootLogs = ref([])
const showCursor = ref(true)

const sequence = [
  "Initializing Galaxies: Burn Rate...",
  "Memory allocation: OK",
  "Loading standard planetary assets: OK",
  "Establishing local simulation environment...",
  "Boot sequence complete."
]

onMounted(() => {
  let delay = 0
  // Fire lines with slight, random delays to simulate processing
  sequence.forEach((text, index) => {
    delay += Math.random() * 400 + 150 
    setTimeout(() => {
      bootLogs.value.push(text)
      if (index === sequence.length - 1) {
        setTimeout(() => emit('complete'), 800)
      }
    }, delay)
  })

  // Blinking cursor effect
  setInterval(() => {
    showCursor.value = !showCursor.value
  }, 500)
})
</script>

<template>
  <div class="splash-screen">
    <div v-for="(log, i) in bootLogs" :key="i" class="log-line">
      {{ log }}
    </div>
    <div class="cursor-line">
      <span v-if="showCursor" class="cursor">█</span>
    </div>
  </div>
</template>

<style scoped>
.splash-screen {
  display: flex;
  flex-direction: column;
  align-items: center;   /* Centers the text */
  justify-content: center; /* Vertically centers the boot sequence */
  width: 100%;
  flex-grow: 1;
  gap: 0.2rem;
}

.cursor {
  display: inline-block;
  width: 12px;
}
</style>