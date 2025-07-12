<template>
  <section class="card stretch">
    <header class="card-header">История рулетки</header>
    <ul class="card-list" id="log-list">
      <li v-for="(item, index) in rouletteHistory" :key="index" class="log-item">
        <span class="log-time">{{ item.time }}</span>
        <span class="log-content">
          <strong style="color: rgb(245, 117, 7);">{{ item.user }}</strong> получает награду <strong style="color: rgb(245, 117, 7);">{{ item.value }}</strong>
          <!--<strong style="color: rgb(255, 0, 251);">{{ item.user }}</strong> получает награду <strong style="color: rgb(255, 0, 251);">{{ item.data }}</strong>-->
        </span>
      </li> 
    </ul>
  </section>
</template>

<script>
import { ref, onMounted } from 'vue';

export default {
  name: 'LogList',
  setup() {
    const rouletteHistory = ref([
      {time: "24.06 12:46", user: "kamicute2", value:"какая-то награда"},
      {time: "24.06 15:12", user: "ElPresedente", value:"магнит из пятигорска"},
      {time: "25.03 19:00", user: "moonseere", value:"очень длинное название награды чтобы перешло на другую строку"},
    ]);
    onMounted(() => {
      /*
        Допустим к нам будут приходить массив [...] данных вида
        {
          user: пользователь, для которого активировалась рулетка
          time: время активации рулетки DD.MM HH.MM
          spins: [
            {
              winnerItem: выпавший итем
              winnerSector: выпавший сектор //не используется
            }
          ]
          value: сектор, выпавший на рулетке
        }
      */
      window.runtime.EventsOn('logUpdated2', (newData) => {
        try{
          const parsedData = JSON.parse( newData )
          parsedData.spins.forEach(element => {
            rouletteHistory.push({ time: parsedData.time, user: parsedData.user, value: parsedData.winnerItem })
          });
        } 
        catch( error ){
          console.error( error )
        }
      });

      window.runtime.EventsOn('logUpdated', (newData) => {
        //Вот так тянуть данные
        //window.go.main.App.FrontendDispatcher("getGroupById", [1]);
        rouletteHistory.value = newData;
        try {
        // Парсим JSON-строку в массив объектов
        const parsedData = JSON.parse(newData);
        rouletteHistory.value = parsedData.map(item => ({
          time: item.time, // Время активации
          user: item.user, // Пользователь
          value: item.value // Выпавший сектор
        }));
        } catch (error) {
          console.error('Ошибка парсинга JSON:', error);
          // Запасные данные на случай ошибки
          rouletteHistory.value = [
            { time: '2023-10-20T00:00:00Z', user: 'default_user', value: 'Сектор по умолчанию' }
          ];
        }
      });
    });
    return {
      rouletteHistory,
    };
  }
};
  //В теории вот так надо будет подписываться на ивенты с го для обновления интерфейса
  /*
  mounted() {
    window.wails.runtime.events.on('logUpdate', (data) => {
      this.history = data;
    });
  },
  beforeUnmount() {
    // Отписка от события при уничтожении компонента
    window.wails.runtime.events.off('logUpdate');
  }
    
};*/

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