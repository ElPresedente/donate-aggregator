<template>
  <section class="card stretch">
    <header class="card-header">История рулетки</header>
    <div class="scroll-container">
      <ul class="card-list">
        <li v-for="(log, index) in sortedLogs" :key="index" class="log-item" :class="{ 'pinned': log.pinned }">
          <span class="log-time" :class="{ 'pinned-time': log.pinned }">{{ log.time }}</span>
          <span class="log-content">
            <strong style="color: rgb(245, 117, 7);">{{ log.user }}</strong> получает награду <strong style="color: rgb(245, 117, 7);">{{ log.value }}</strong>
          </span>
          <button @click="togglePin(log)" class="pin-button">
            {{ log.pinned ? '🔓' : '📌' }}
          </button>
        </li> 
      </ul>
    </div>
  </section>
</template>

<script>
import { useLogStore } from '../stores/logStore';
import { FrontendDispatcher } from '../../wailsjs/go/main/App'

export default {
  name: 'LogList',
  setup() {
    const logStore = useLogStore();

    const togglePin = (log) => {
      const index = logStore.rouletteHistory.findIndex(l => l === log);
      if (index === -1) return;

      if(logStore.rouletteHistory[index].pinned){
        resetPinnedItem(logStore.rouletteHistory[index])
      } else {
        setPinnedItem(logStore.rouletteHistory[index])
      }
      logStore.rouletteHistory[index].pinned = !logStore.rouletteHistory[index].pinned;
    }

    const setPinnedItem = (textOrIndexIDK) => {
      const JSONtextOrIndexIDK = JSON.stringify(textOrIndexIDK)
      FrontendDispatcher("set-pinned-reward", JSONtextOrIndexIDK) //нет обработчика
    }

    const resetPinnedItem = (textOrIndexIDK) => {
      const JSONtextOrIndexIDK = JSON.stringify(textOrIndexIDK)
      FrontendDispatcher("reset-pinned-reward", JSONtextOrIndexIDK)//нет обработчика
    }
    return {
      logStore,
      togglePin
    };
  },
  computed: {
     sortedLogs() {
      // Разделяем логи на закрепленные и незакрепленные, сохраняя порядок
      const pinned = this.logStore.rouletteHistory.filter(log => log.pinned);
      const unpinned = this.logStore.rouletteHistory.filter(log => !log.pinned);
      return [...pinned, ...unpinned];
    }
  }
};

</script>
<style scoped>
.pin-button{
  display: inline-block;
  cursor: pointer;
}
.card {
  width: 50%;
  height: 100%;
  background-color: #1e1e1e;
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 0 10px #00000070;
  display: flex;
  flex-direction: column;
}

.card-header {
  font-size: 1.2rem;
  margin-bottom: 10px;
  font-weight: bold;
}

.card-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.scroll-container {
  overflow-y: scroll;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.scroll-container::-webkit-scrollbar {
  width: 10px;
  background: transparent;
}

.card-list li {
  padding: 8px 5px;
  border-bottom: 1px solid #2a2a2a;
}

.log-item {
  display: flex;
  align-items: flex-start; /* Выравнивание по началу для многострочного текста */
  font-size: 16px; /* Базовый шрифт */
  line-height: 1.5; /* Для читаемости */
  max-width: 100%; /* Ограничиваем контейнер */
}

.log-time {
  color: #888; /* Ненавязчивый цвет */
  font-size: 0.85em; /* Меньше основного */
  margin-right: 10px; /* Отступ от времени */
  white-space: nowrap; /* Время не переносится */
}

.log-content {
  text-align: left;
}

.user, .reward {
  font-size: 1.1em; /* Чуть крупнее */
  font-weight: bold; /* Жирный шрифт */
  margin: 0 4px; /* Отступы */
  word-wrap: break-word; /* Перенос длинных слов */
  overflow-wrap: break-word; /* Совместимость */
}
.pin-button{
  background-color: transparent;
  box-shadow: none;
  border: 1px solid transparent;
}
.pin-button:hover{
  border: 1px solid rgb(245, 117, 7);
}

.pinned{
  background: rgba(240, 255, 255, 0.073);
  border-radius: 5px;
}

.pinned-time{
  color:red;
  font-weight: 900;
}

.hidden {
  display: none;
}
</style>