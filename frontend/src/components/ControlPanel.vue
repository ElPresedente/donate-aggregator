<template>
  <section class="card stretch" id="right-panel">
    <div class="card-block" id="main-controls">
      <header class="card-header">Управление рулеткой</header>
      <div class="status">
        <span v-if="donattyConnected === ConnectionStatus.CONNECTED" class="status-connected">✅ Donatty: Подключено</span>
        <span v-if="donattyConnected === ConnectionStatus.DISCONNECTED" class="status-disconnected">❌ Donatty: Не подключено</span>
        <span v-if="donattyConnected === ConnectionStatus.RECONNECTING" class="status-reconnecting">⚠️ Donatty: Попытка переподключения...</span>
        <span v-if="donattyConnected === ConnectionStatus.CONNECTION_LOST" class="status-lost">⚠️ Donatty: Произошел разрыв соединения</span>
        
        <span v-if="donatepayConnected === 'connected'" class="status-connected">✅ Donatepay: Подключено</span>
        <span v-if="donatepayConnected === 'disconnected'" class="status-disconnected">❌ Donatepay: Не подключено</span>
        <span v-if="donatepayConnected === 'reconnecting'" class="status-reconnecting">⚠️ Donatepay: Попытка переподключения...</span>
        <span v-if="donatepayConnected === 'connection_lost'" class="status-lost">⚠️ Donatepay: Произошел разрыв соединения</span>
      </div>
      <div class="controls">
        <button class="btn green" @click="rouletteOn">Включить</button>
        <button class="btn red" @click="rouletteOff">Выключить</button>
        <button class="btn blue" @click="rouletteReconnect">Перезапустить</button>
        <button class="btn gold" @click="rollRoulette">Крутить</button>
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
import { ref, onMounted } from 'vue';
export default {
  name: 'ControlPanel',
  setup(){
    const ConnectionStatus = Object.freeze({
      CONNECTED: 'connected',
      DISCONNECTED: 'disconnected',
      RECONNECTING: 'reconnecting',
      CONNECTION_LOST: 'connectionLost'
    });
    const router = useRouter();
    const donattyConnected = ref("connectionLost");
    const donatepayConnected = ref("disconnected");
    //Я вот думаю, надо ли подобные флаги подключения вешать на виджеты...

    onMounted(() => {
      console.log(donatepayConnected === 'disconnected');
      window.runtime.EventsOn('donattyConnectionUpdated', (connection) => {
        /*
          для обоих сервисов сделать сообщения по типу
          disconnected - соединения нет
          connected - соединение есть
          reconnecting
          connectionLost
        */
        donattyConnected.value = connection;
      });
      window.runtime.EventsOn('donatepayConnectionUpdated', (connection) => {
        /*
          для обоих сервисов сделать сообщения по типу
          disconnected - соединения нет
          connected - соединение есть
          recconecting
          connectionLost
        */
        donatepayConnected.value = connection;
      });
    });
    
    const rollRoulette = () => {
      //Метод для прокрута рулетки без доната
      //window.go.main.App.SendMessageFromFrontend("Привет от кнопки!");
    };
    const rouletteOn = () => {
      //дёргаем из го
    };
    const rouletteOff = () => {
      //дёргаем из го
    };
    const rouletteReconnect = () => {
      //дёргаем из го
    };
    const showSettings = () => {
      router.push('/settings');
    };
    const showRouletteSettings = () => {
      router.push('/roulette-settings');
    };
    return {
      ConnectionStatus,
      donattyConnected,
      donatepayConnected,
      rollRoulette, 
      rouletteOn, 
      rouletteOff, 
      rouletteReconnect, 
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

.btn.red {
  background-color: #ef4444;
  transition: background-color 0.2s ease;
}

.btn.red:hover {
  background-color: #dc2626; /* Темнее на 20% */
}

.btn.blue {
  background-color: #3b82f6;
  transition: background-color 0.2s ease;
}

.btn.blue:hover {
  background-color: #2563eb; /* Темнее на 20% */
}

.btn.gold {
  background-color: #fbbf24;
  transition: background-color 0.2s ease;
}

.btn.gold:hover {
  background-color: #d97706; /* Темнее на 20% */
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
</style>
