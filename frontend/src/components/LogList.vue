<template>
  <section class="card stretch">
    <header class="card-header">История рулетки</header>
    <div class="scroll-container">
      <ul class="card-list">
        <li v-for="(item, index) in logStore.pinnedHistory" :key="index" class="log-item">
          <span class="log-time">{{ item.time }}</span>
          <span class="log-content">
            <strong style="color: rgb(245, 117, 7);">{{ item.user }}</strong> получает награду <strong style="color: rgb(245, 117, 7);">{{ item.value }}</strong>
          </span>
          <span @click="logStore.unpinPinnedItem(index); resetPinnedItem('ПОКА ЧТО ТУТ ХУЙ')" class="pin-button">🔓</span>
          <!--ТУТ ПОМЕНЯТЬ ------------------------------------------- ^^^^^^^^^^^^^^^^-->
        </li> 
      </ul>
      <ul class="card-list">
        <li class="log-item"></li>
      </ul>
      <ul class="card-list">
        <li v-for="(item, index) in logStore.rouletteHistory" :key="index" class="log-item" :class="{ hidden: logStore.isPinned(index) }">
          <span class="log-time">{{ item.time }}</span>
          <span class="log-content">
            <strong style="color: rgb(245, 117, 7);">{{ item.user }}</strong> получает награду <strong style="color: rgb(245, 117, 7);">{{ item.value }}</strong>
          </span>
          <span @click="logStore.pinRouletteItem(index); setPinnedItem('ПОКА ЧТО ТУТ ХУЙ')" class="pin-button">📌</span>
          <!--ТУТ ПОМЕНЯТЬ ------------------------------------------- ^^^^^^^^^^^^^^^^-->
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

    const setPinnedItem = (textOrIndexIDK) => {
      FrontendDispatcher("set-pinned-reward", textOrIndexIDK) //нет обработчика
    }

    const resetPinnedItem = (textOrIndexIDK) => {
      FrontendDispatcher("reset-pinned-reward", textOrIndexIDK)//нет обработчика
    }

    const updatePinned = (textOrIndexIDK) => {
      if( logStore.pinnedHistory.length == 0 ){
        //FrontendDispatcher("reset-pinned-rewards", "")
      }
      else{
        /*
        let resultStr = "Награды рулетки:"
        logStore.pinnedHistory.forEach((item) => {
          resultStr += '\n' + item.value 
        })*/
        FrontendDispatcher("update-pinned-rewards", textOrIndexIDK)
      }
    }
    return {
      logStore,
      updatePinned,
      resetPinnedItem,
      setPinnedItem,
    };
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

.hidden {
  display: none;
}
</style>