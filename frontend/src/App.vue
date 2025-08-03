<script setup>
import Toast from './components/actions/Toast.vue'
import { onMounted, ref } from 'vue'
import { useToastStore } from './stores/toastStore'
import { useConnectionStore } from './stores/connectionStore';
import { useLogStore } from './stores/logStore';
import { FrontendDispatcher } from '../wailsjs/go/main/App'

const toastStore = useToastStore()
const toastInstance = ref(null)
const connectionStore = useConnectionStore();
const logStore = useLogStore();

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
    window.runtime.EventsOn('logNumData', (newData) => {
      console.log(newData)
      if(newData != null){
        logStore.rouletteHistory = [];
        newData.forEach(element => {
          logStore.rouletteHistory.push(element)
        });
      }
        logStore.rouletteHistory = newData;
    })
    window.runtime.EventsOn('logUpdated', (newData) => {
      try{
        const parsedData = JSON.parse( newData )
        parsedData.spins.forEach(element => {
          if (logStore.rouletteHistory.length > numLogs-1)
          {
            logStore.rouletteHistory.pop()
          }
          logStore.rouletteHistory.unshift({ time: parsedData.time, user: parsedData.user, value: element.sector })
        });
      } 
      catch( error ){
        console.error( error )
      }
    })
    FrontendDispatcher("getLogs", "");
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