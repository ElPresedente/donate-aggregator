// stores/toastStore.js
import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useToastStore = defineStore('toast', () => {
  const toastRef = ref(null)
  const toastsEnabled = ref(false)
  function setRef(ref) {
    toastRef.value = ref
  }

  function showToast(message, type = 'info', duration = 3000) {
    if(toastsEnabled.value){
      if (toastRef.value && typeof toastRef.value.showToast === 'function') {
        toastRef.value.showToast(message, type, duration)
      } else {
        console.warn('Toast компонент не инициализирован')
      }
    }
  }
  function setToastsEnabled (data){
    toastsEnabled.value = data
  }

  return { toastRef, toastsEnabled, setRef, showToast, setToastsEnabled }
})
