import { ref } from "vue"
import { defineStore } from 'pinia';

export const useConnectionStore = defineStore('connection', () => {
  const donattyConnected = ref('disconnected');
  const donatepayConnected = ref('disconnected');
  const twitchConnected = ref('disconnected');
  const rouletteConnected = ref('disconnected');
  const rewardConnected = ref('disconnected');
  const isOnButtonDisabled = ref(false);
  const currentAmount = ref(0); // В будущем перенести в нормальный стор
  const donateQueueLength = ref(0); // В будущем перенести в нормальный стор
  const subscribedStatus = ref(false); // В будущем перенести в нормальный стор


  return {
    donattyConnected,
    donatepayConnected,
    twitchConnected,
    rouletteConnected,
    rewardConnected,
    isOnButtonDisabled,
    currentAmount,
    donateQueueLength,
    subscribedStatus,
  };
});
