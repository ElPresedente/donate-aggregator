import { ref } from "vue"
import { defineStore } from 'pinia';

export const useLogStore = defineStore('log', {
  state: () => ({
    pinnedHistory: [], // Ваш массив для pinned
    rouletteHistory: [], // Ваш массив для roulette
    hiddensPins: [],
  }),
  actions: {
    pinRouletteItem(index) {
      const item = this.rouletteHistory[index];
      this.pinnedHistory.unshift(item); // Или .push
      this.hiddensPins.unshift(index);
    },
    unpinPinnedItem(index){
      const item = this.pinnedHistory[index];
      this.pinnedHistory.splice(index, 1);
      this.hiddensPins.splice(index, 1);
    },
    isPinned(index) {
      return this.hiddensPins.includes(index);
    }
  },
});
