<template>
  <div class="card">
    <ul>
      <li 
      v-for="(item, idx) in localItems" :key="item.id || item.tempId"
      v-show="item.status !== 'delete'"
      >
        <input
          v-model="item.data"
          class="hidden-input"
          @input="markAsEdited(idx)"
        />
        <!--@input="updateItem(idx, $event.target.value)"-->
        <button class="del-btn" @click="deleteItem(idx)">Х</button>
      </li>
    </ul>
    <button class="btn add" @click="addItem()">Добавить</button>
  </div>
  <section class="card">
    <button class="btn save" @click="saveChanges()">Сохранить все</button>
    <button class="btn back" @click="goBack()">← Назад</button>
  </section>
</template>

<script>
import { useRoute, useRouter } from 'vue-router';
import { ref, onMounted, onUnmounted } from "vue";
import { FrontendDispatcher } from '../../wailsjs/go/main/App'
export default {
  name: 'RouletteSettingsCardEditor',
  setup() {
    let unsubscribes = [];
    const router = useRouter()
    const route = useRoute()
    const index = JSON.parse(route.params.index)+1
    const localItems = ref([]);
    onMounted(() => {
      unsubscribes.push(
        window.runtime.EventsOn('sectorsByCategoryIdData', (data) => {
          if(data)
            localItems.value = data // ← обновляем реактивно
        })
      );
      FrontendDispatcher("getSectorsByCategoryId", JSON.stringify({category_id: index }));
    });
    onUnmounted(() => {
      unsubscribes.forEach(unsub => unsub());
    });
    const markAsEdited = (idx) => {
      if (localItems.value[idx].id) {
        localItems.value[idx].status = 'edit';
      }
    }
    const addItem = () => {
      localItems.value.push({
        tempId: Date.now(),
        data: '',
        status: 'add'
      });
    }
    const deleteItem = (idx) => {
      if (localItems.value[idx].id) {
        localItems.value[idx].status = 'delete';

      } else {
        localItems.value.splice(idx, 1);
      }
    }
    const saveChanges = () => {
      const sectorsToSave = localItems.value.filter(
        item => item.data.trim() !== ''
      );
      const data = {
        id: index,
        sectors: sectorsToSave
      };
      FrontendDispatcher("sectorsToSave", JSON.stringify(data));
    }
    
    const goBack = () => router.go(-1);
    return { 
      localItems,
      markAsEdited,
      addItem,
      deleteItem,
      saveChanges,
      goBack
     };
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
  color: #ffffff;
  margin-bottom: 5px;
  display: flex;
  align-items: center;
  margin: 15px 10px;
}

.hidden-input {
  background: transparent;
  border: none;
  color: #ffffff;
  font-size: 20px;
  width: 100%;
  outline: none;
  padding: 5px 0;
  border-bottom: 1px solid #bd7e3f;
}

.hidden-input:focus {
  border-bottom: 1px solid #f77c00;
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

.btn.back:hover {
  background-color: #4b5563;
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

.btn.save:hover {
  background-color: #16a34a;
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
  margin-top: 40px;
}

.btn.add:hover {
  background-color: #16a34a;
}
</style>