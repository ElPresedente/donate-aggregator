<template>
  <section class="card stretch">
    <header class="card-header">История рулетки</header>
    <ul class="card-list" id="log-list">
      <li v-for="(item, index) in data" :key="index">{{ item }}</li>
    </ul>
  </section>
</template>

<script>
import { ref, onMounted } from 'vue';

export default {
  name: 'LogList',
  setup() {
    const data = ref([]);

    onMounted(() => {
      window.runtime.EventsOn('db_updated', (newData) => {
        data.value = newData;
      });
    });

    return {
      data,
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
</style>