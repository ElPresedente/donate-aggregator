import { ref } from "vue"
import { defineStore } from 'pinia';

export const useConnectionStore = defineStore('connection', () => {
  const donattyConnected = ref('disconnected');
  const donatepayConnected = ref('disconnected');
  const rouletteConnected = ref('disconnected');
  const isOnButtonDisabled = ref(false);
  const isOffButtonDisabled = ref(true);

  return {
    donattyConnected,
    donatepayConnected,
    rouletteConnected,
    isOnButtonDisabled,
    isOffButtonDisabled,
  };
});
