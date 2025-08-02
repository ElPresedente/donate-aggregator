<template>
  <section class="card stretch">
    <header class="card-header">История рулетки</header>
    <div class="scroll-container">
      <ul class="card-list" id="log-list">
        <li v-for="(item, index) in rouletteHistory" :key="index" class="log-item">
          <span class="log-time">{{ item.time }}</span>
          <span class="log-content">
            <strong style="color: rgb(245, 117, 7);">{{ item.user }}</strong> получает награду <strong style="color: rgb(245, 117, 7);">{{ item.value }}</strong>
          </span>
        </li> 
      </ul>
    </div>
  </section>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue';
import { FrontendDispatcher } from '../../wailsjs/go/main/App'
const numLogs = 100;

export default {
  name: 'LogList',
  setup() {
    let unsubscribes = [];
    const rouletteHistory = ref([]);
    onMounted(() => {
      unsubscribes.push(
        window.runtime.EventsOn('logNumData', (newData) => {
          if(newData != null)
            rouletteHistory.value = newData;
        })
      );
      unsubscribes.push(
        window.runtime.EventsOn('logUpdated', (newData) => {
          try{
            const parsedData = JSON.parse( newData )
            parsedData.spins.forEach(element => {
              if (rouletteHistory.value.length > numLogs-1)
              {
                rouletteHistory.value.pop()
              }
              rouletteHistory.value.unshift({ time: parsedData.time, user: parsedData.user, value: element.sector })
            });
          } 
          catch( error ){
            console.error( error )
          }
        })
      );
      FrontendDispatcher("getNumLogs", String(numLogs));
    });
    onUnmounted(() => {
      unsubscribes.forEach(unsub => unsub());
    });
    return {
      rouletteHistory,
    };
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
  padding: 8px 0;
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
</style>