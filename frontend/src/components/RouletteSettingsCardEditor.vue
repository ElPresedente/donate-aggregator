<template>
  <div class="card">
    <ul>
      <li v-for="(item, idx) in localItems" :key="idx">
        <input
          v-model="localItems[idx]"
          class="hidden-input"
          @input="updateItem(idx, $event.target.value)"
        />
        <button class="del-btn" @click="saveItem(idx)">Х</button>
      </li>
    </ul>
    <button class="btn add" @click="add()">Добавить</button>
  </div>
  <section class="card">
    <button class="btn save" @click="save()">Сохранить все</button>
    <button class="btn back" @click="goBack()">← Назад</button>
  </section>
</template>

<script>
import { useRoute, useRouter } from 'vue-router';
import { ref } from "vue";

export default {
  name: 'RouletteSettingsCardEditor',
  setup() {
    const router = useRouter();//Пока не понял, в чём прикол, но так работает
    const route = useRoute();
    const localItems = ref(JSON.parse(route.params.items || '[]'));

    // Вывод аргумента в консоль
    console.log('Index из параметра маршрута:', localItems);

    const save = () => {

      //Сохранение в бд
      
      router.go(-1);
    }

    const goBack = () => router.go(-1);
    return { localItems, save, goBack };
  },
};
</script>

<style scoped>
.card {
  width: 80%;
  margin: 20px auto;
  background-color: #2a2a2a;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.3);
}

.card h3 {
  font-size: 16px;
  color: #ccc;
  flex-grow: 1;
  text-align: center;
  margin: 0;
}

.card ul {
  list-style: none;
  padding: 0;
  margin: 10px 0;
}

.card li {
  font-size: 20px;
  color: #999;
  margin-bottom: 5px;
  display: flex;
  align-items: center;
  margin: 15px 10px;
}

.hidden-input {
  background: transparent;
  border: none;
  color: #999;
  font-size: 20px;
  width: 100%;
  outline: none;
  padding: 5px 0;
}

.hidden-input:focus {
  border-bottom: 1px solid #f57d07;
}

.del-btn {
  background-color: #790000;
  border: 1px solid transparent;
  border-radius: 4px;
  color: red;
  padding: 5px 10px;
  font-size: 10px;
  cursor: pointer;
}

.del-btn:hover {
  background-color: #b10000;
}

.btn.back {
  font-size: large;
  background-color: #6b7280;
  margin-top: 20px;
  width: 100%;
  padding: 10px;
  border-radius: 8px;
  border: none;
  color: white;
  font-weight: bold;
  cursor: pointer;
}
.btn.save {
  width: 100%;
  font-size: large;
  padding: 10px;
  border-radius: 8px;
  border: none;
  color: white;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s ease;
  background-color: #22c55e;
}
.btn.add {
  width: 100%;
  font-size: large;
  padding: 10px;
  border-radius: 8px;
  border: none;
  color: white;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s ease;
  background-color: #22c55e;
}
</style>