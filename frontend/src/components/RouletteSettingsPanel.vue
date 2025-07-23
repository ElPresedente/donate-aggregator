<template>
  <div class="container">
    <RouletteSettingsCard
      v-for="(category, index) in categories"
      :key="index"
      :title="category.title"
      :items="category.items"
      :percentage="category.percentage"
      :color="category.color"
      :index="index"
    />
  </div>
  <section class="card stretch" id="settings-panel">
    <button class="btn back" @click="goBack()">← Назад</button>
  </section>
</template>

<script>
import { useRouter } from 'vue-router';
import RouletteSettingsCard from './RouletteSettingsCard.vue'
import { onMounted, ref, onUnmounted } from "vue";
import { FrontendDispatcher } from '../../wailsjs/go/main/App'

export default {
  name: 'RouletteSettingsPanel',
  components: { RouletteSettingsCard },
  setup() {
    let unsubscribes = [];
    const router = useRouter();
    const goBack = () => router.go(-1);
    const categories = ref([])
    onMounted(() => {
      unsubscribes.push(
        window.runtime.EventsOn('groupsData', (data) => {
          console.log('📦 Группы:', data)
          categories.value = data 
        })
      );
      FrontendDispatcher("getGroups", "");
    });
    onUnmounted(() => {
      unsubscribes.forEach(unsub => unsub());
    });
    return { 
      categories,
      goBack
     }
  },  
  methods: {
    updateFormData(service, newData) {
      this.formData[service] = newData;
    },
    handleSave() {
    },
  },
};
</script>

<style scoped>
body {
  background-color: #1a1a1a;
  color: #fff;
  font-family: Arial, sans-serif;
}

.container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 15px;
  padding: 10px;
}

.card {
  background-color: #1e1e1e;
  border-radius: 12px;
  padding: 20px;
  margin: 20px;
  box-shadow: 0 0 10px #00000070;
  display: flex;
  flex-direction: column;
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
</style>
