import { ref } from "vue"
import { defineStore } from 'pinia';

export const useLogStore = defineStore('log', {
  state: () => ({
    rouletteHistory: [],
  })
});
