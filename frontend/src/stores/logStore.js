import { ref } from "vue"
import { defineStore } from 'pinia';

export const useLogStore = defineStore('log', () => {
  const rouletteHistory = ref([]);
  
  return {
    rouletteHistory
  };
});
