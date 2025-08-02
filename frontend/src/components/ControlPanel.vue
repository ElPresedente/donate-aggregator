<template>
  <section class="card stretch" id="right-panel">
    <div class="card-block" id="main-controls">
      <header class="card-header">Управление рулеткой</header>
      <div class="status">
        <div class="status-row">
          <span v-if="connectionStore.donattyConnected === ConnectionStatus.CONNECTED" class="status-connected">✅ Donatty: Подключено</span>
          <span v-if="connectionStore.donattyConnected === ConnectionStatus.DISCONNECTED" class="status-disconnected">❌ Donatty: Не подключено</span>
          <span v-if="connectionStore.donattyConnected === ConnectionStatus.RECONNECTING" class="status-reconnecting">⚠️ Donatty: Попытка подключения...</span>
          <button v-if="connectionStore.isOnButtonDisabled" class="reload-btn" @click="reconnectDonatty">🔄</button>
        </div>
        
        <div class="status-row">
          <span v-if="connectionStore.donatepayConnected === ConnectionStatus.CONNECTED" class="status-connected">✅ Donatepay: Подключено</span>
          <span v-if="connectionStore.donatepayConnected === ConnectionStatus.DISCONNECTED" class="status-disconnected">❌ Donatepay: Не подключено</span>
          <span v-if="connectionStore.donatepayConnected === ConnectionStatus.RECONNECTING" class="status-reconnecting">⚠️ Donatepay: Попытка подключения...</span>
          <button v-if="connectionStore.isOnButtonDisabled" class="reload-btn" @click="reconnectDonatepay">🔄</button>
        </div>

        <div class="status-row">
          <span v-if="connectionStore.rouletteConnected === ConnectionStatus.CONNECTED" class="status-connected">✅ Виджет рулетки: Подключено</span>
          <span v-if="connectionStore.rouletteConnected === ConnectionStatus.DISCONNECTED" class="status-disconnected">❌ Виджет рулетки: Не подключено</span>
          <span v-if="connectionStore.rouletteConnected === ConnectionStatus.RECONNECTING" class="status-reconnecting">⚠️ Виджет рулетки: Попытка подключения...</span>
          <button v-if="connectionStore.rouletteConnected === ConnectionStatus.CONNECTED || connectionStore.rouletteConnected === ConnectionStatus.RECONNECTING" class="reload-btn" @click="reloadRoulette">🔄</button>
        </div>

        <br></br>

        <div class="status-row">
          <span>💲 Накоплено в рулетке: &nbsp;</span>
          <span id="current-amount">{{ connectionStore.currentAmount }}</span> <!--БЛЯ вынеси нахуй стили-->
          <span>  Донатов в очереди: &nbsp;</span>
          <span id="donate-queue-length">{{ connectionStore.donateQueueLength }}</span> <!--БЛЯ вынеси нахуй стили-->
        </div>
      </div>
      <div class="controls">
        <button id="onButton" class="btn green" @click="rouletteOn" :disabled="connectionStore.isOnButtonDisabled">Включить</button>
        <button id="offButton" class="btn red" @click="rouletteOff" :disabled="!connectionStore.isOnButtonDisabled">Выключить</button>
        <button class="btn gold" @click="rollRoulette" :disabled="connectionStore.rouletteConnected !== ConnectionStatus.CONNECTED">Крутить</button>
      </div>
    </div>
    <div class="card-block settings-buttons">
      <button class="btn gray" @click="showSettings">
        Настройка подключения
      </button>
      <button class="btn gray" @click="showRouletteSettings">
        Настройка рулетки
      </button>
    </div>
  </section>
</template>

<script>
import { useRouter } from 'vue-router';
import { onMounted, onUnmounted } from 'vue';
import { FrontendDispatcher } from '../../wailsjs/go/main/App'
import { useConnectionStore } from '../stores/connectionStore';
import { useToastStore } from '../stores/toastStore'

export default {
  name: 'ControlPanel',
  setup(){
    const connectionStore = useConnectionStore();
    const toast = useToastStore();
    const ConnectionStatus = Object.freeze({
      CONNECTED: 'connected',
      DISCONNECTED: 'disconnected',
      RECONNECTING: 'reconnecting',
    });
    const router = useRouter();
    onMounted(() => {

    });
    onUnmounted(() => {
      //unsubscribes.forEach(unsub => unsub());
    });
    const rollRoulette = () => {
      FrontendDispatcher("manualRouletteSpin", "");
      toast.showToast('Крутим рулетку', 'info', 3000)
    };
    const rouletteOn = () => {
      connectionStore.isOnButtonDisabled = true;
      FrontendDispatcher("startAllCollector", "");
      toast.showToast('Коллекторы включены', 'success', 3000)
    };
    const rouletteOff = () => {
      connectionStore.isOnButtonDisabled = false;
      FrontendDispatcher("stopAllCollector", "");
      toast.showToast('Коллекторы выключены', 'success', 3000)
    };
    const reconnectDonatty = () => {
      FrontendDispatcher("reconnectDonatty", "")
    };
    const reconnectDonatepay = () => {
      FrontendDispatcher("reconnectDonatepay", "")
    };
    const reloadRoulette = () => {
      FrontendDispatcher("reloadRoulette", "")
    };
    const rouletteReconnect = () => {
      FrontendDispatcher("reconnectAllCollector");
    };
    const showSettings = () => {
      router.push('/settings');
    };
    const showRouletteSettings = () => {
      router.push('/roulette-settings');
    };
    return {
      ConnectionStatus,
      connectionStore,
      rollRoulette, 
      rouletteOn, 
      rouletteOff, 
      reconnectDonatty,
      reconnectDonatepay,
      reloadRoulette,
      showSettings, 
      showRouletteSettings 
    }
  }
};
</script>

<style scoped>
.card {
  width: 50%;
  height: 100%;
  background-color: #1e1e1e;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 0 10px #00000070;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.card-header {
  font-size: 1.2rem;
  margin-bottom: 10px;
  font-weight: bold;
}

.status span {
  display: block;
  margin: 5px 0;
  font-size: 0.95rem;
  opacity: 0.9;
}

.status span.online {
  color: #22c55e;
}

.status-row {
  display: flex;
  align-items: center;
}

.reload-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 1.25rem;
  padding: 2px 6px;
  transition: transform 0.15s ease, filter 0.15s ease;
  color: inherit;
}

.reload-btn:hover {
  transform: scale(1.15);
  filter: brightness(0.8);
}

.reload-btn:active {
  transform: scale(1.1);
  filter: brightness(0.6);
}

.controls {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 15px;
}

.btn {
  flex: 1 1 45%;
  padding: 10px;
  border-radius: 8px;
  border: none;
  color: white;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s ease;
}

.btn.green {
  background-color: #22c55e;
  transition: background-color 0.2s ease; /* Плавный переход для фона */
}

.btn.green:hover {
  background-color: #16a34a; /* Темнее на 20% для наведения */
}

.btn.green:disabled {
  background-color: #15803d; /* Темный зелёный */
  color: #9ca3af;
  cursor: not-allowed;
  opacity: 0.7;
  pointer-events: none;
}

.btn.red {
  background-color: #ef4444;
  transition: background-color 0.2s ease;
}

.btn.red:hover {
  background-color: #dc2626; /* Темнее на 20% */
}

.btn.red:disabled {
  background-color: #991b1b; /* Темный красный */
  color: #9ca3af;
  cursor: not-allowed;
  opacity: 0.7;
  pointer-events: none;
}

.btn.gold {
  background-color: #fbbf24;
  transition: background-color 0.2s ease;
}

.btn.gold:hover {
  background-color: #d97706; /* Темнее на 20% */
}

.btn.gold:disabled {
  background-color: #b45309; /* Темный красный */
  color: #9ca3af;
  cursor: not-allowed;
  opacity: 0.7;
  pointer-events: none;
}

.btn.gray {
  background-color: #6b7280;
}

.btn.gray:hover {
  background-color: #4b5563;
}

.settings-buttons {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

#current-amount{
  color: rgb(28, 226, 28);
  font-weight: bold;
  font-size: 1.5rem;
  margin-right: 20px;
}

#donate-queue-length{
  color: rgb(28, 226, 28);
  font-weight: bold;
  font-size: 1.5rem;
}
</style>
