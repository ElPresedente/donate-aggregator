import { ref } from "vue"
import { defineStore } from 'pinia';

export const useLogStore = defineStore('log', {
  state: () => ({
    pinnedHistory: [], // Ваш массив для pinned
    rouletteHistory: [], // Ваш массив для roulette
  }),
  actions: {
    pinRouletteItem(index) {
      const item = this.rouletteHistory[index];
      this.pinnedHistory.unshift(item); // Или .push
    },
    unpinPinnedItem(index){
      const item = this.pinnedHistory[index];
      this.pinnedHistory.splice(index, 1);
    }
  },
});
