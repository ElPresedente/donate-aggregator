<script setup>
import { ref, defineExpose } from 'vue'

const toasts = ref([])

function showToast(message, type = 'info', duration = 3000) {
  toasts.value.push({ message, type })
  setTimeout(() => {
    toasts.value.shift()
  }, duration)
}

defineExpose({ showToast })
</script>

<template>
  <div class="toast-container">
    <div v-for="(toast, i) in toasts" :key="i" :class="['toast', toast.type]">
      {{ toast.message }}
    </div>
  </div>
</template>

<style scoped>
.toast-container {
  position: fixed;
  bottom: 20px;
  right: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  z-index: 9999;
}

.toast {
  padding: 10px 20px;
  border-radius: 6px;
  color: white;
  font-weight: bold;
  min-width: 200px;
  box-shadow: 0 0 10px #00000044;
  animation: fadeIn 0.3s ease;
}

.toast.info {
  background-color: #007bff;
}

.toast.success {
  background-color: #28a745;
}

.toast.error {
  background-color: #dc3545;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
