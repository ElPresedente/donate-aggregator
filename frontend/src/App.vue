<script setup>
import Toast from './components/actions/Toast.vue'
import { onMounted, ref } from 'vue'
import { useToastStore } from './stores/toastStore'
import { useConnectionStore } from './stores/connectionStore';

const toastStore = useToastStore()
const toastInstance = ref(null)
const connectionStore = useConnectionStore();

onMounted(() => {
  toastStore.setRef(toastInstance.value)

  if( !connectionStore.subscribedStatus ){
    window.runtime.EventsOn('donattyConnectionUpdated', (connection) => {
      connectionStore.donattyConnected = connection;
    })
    window.runtime.EventsOn('donatepayConnectionUpdated', (connection) => {
      connectionStore.donatepayConnected = connection;
    })
    window.runtime.EventsOn('rouletteConnectionUpdated', (connection) => {
      connectionStore.rouletteConnected = connection;
    })
    window.runtime.EventsOn('currentAmountUpdate', (amount) => {
      connectionStore.currentAmount = amount;
    })
    window.runtime.EventsOn('donateQueueLengthUpdate', (amount) => {
      connectionStore.donateQueueLength = amount;
    })
    window.runtime.EventsOn('toastExec', (data) => {
      toastStore.showToast(data.message, data.type, 3000)
    })
    connectionStore.subscribedStatus = true
  }
})
</script>

<template>
  <router-view />
  <Toast ref="toastInstance" />
</template>


<style>
#app{
  height: auto;
}
</style>