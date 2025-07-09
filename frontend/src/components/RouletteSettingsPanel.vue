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
import { onMounted, ref } from "vue";
import { FrontendDispatcher } from '../../wailsjs/go/main/App'

export default {
  name: 'RouletteSettingsPanel',
  components: { RouletteSettingsCard },
  // data() {
  //   return {
  //     categories: [
  //       {
  //         title: 'Обычные',
  //         items: ['Попить воды', 'Один раз "Мяу"', 'Один раз "Ня"', 'Один раз "Фыр"', 'Сердечко в чатик', 'Поцелуй чатику'],
  //         percentage: 50
  //       },
  //       {
  //         title: 'Необычные',
  //         items: ['Слова поддержки', 'Сердечки в глазах на 5 мин', 'Звездочки в глазах на 5 мин', 'Размять ладошки'],
  //         percentage: 25,
  //         color: 'rgb(55, 255, 0)',
  //       },
  //       {
  //         title: 'Редкие',
  //         items: ['Детте слишком близко на 1 мин', 'Маленькая Детте на 1 мин', 'Потягушки', 'Чиби Детте на 20 мин'],
  //         percentage: 16,
  //         color: 'rgb(0, 200, 255)'
  //       },
  //       {
  //         title: 'Эпические',
  //         items: ['Ещё одна попытка покрутить', 'Модель терминатора на 20 мин'],
  //         percentage: 7,
  //         color: 'rgb(255, 0, 251)'
  //       },
  //       {
  //         title: 'Легендарные',
  //         items: ['Вип на месяц', 'Шепот режим наа 3 мин', 'Роспись в стиме'],
  //         percentage: 1.5,
  //         color: 'rgb(245, 117, 7)'
  //       },
  //       {
  //         title: 'Артифакты',
  //         items: ['Личная (милая) открытка от Детте', 'Личная открытка-подкат от Детте'],
  //         percentage: 0.5,
  //         color: 'rgb(229, 204, 128)'
  //       }
  //     ]
  //   };
  // },
  setup() {
    const router = useRouter();
    const goBack = () => router.go(-1); //router.push('/');
    const categories = ref([])
    onMounted(() => {
      // Вызываем функцию при монтировании компонента
      FrontendDispatcher("getGroups", ""); //Бля передача пустой строки выглядит как костыль
    
      window.runtime.EventsOn('groupsData', (data) => {
        console.log('📦 Группы:', data)
        categories.value = data // ← обновляем реактивно
      });
    })
    
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
      console.log('Сохранённые данные:', this.formData);
      // Отправка на сервер
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
