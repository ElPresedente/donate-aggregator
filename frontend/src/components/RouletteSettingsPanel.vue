<template>
  <div class="container">
    <RouletteSettingsCard
      v-for="category in categories"
      :key="category.title"
      :title="category.title"
      :items="category.items"
      :percentage="category.percentage"
    />
  </div>
</template>

<script>
import { useRouter } from 'vue-router';
import RouletteSettingsCard from './RouletteSettingsCard.vue'
export default {
  setup() {
    const router = useRouter();
    const goBack = () => router.push('/');
    return { goBack };
  },
  components: { RouletteSettingsCard },
  data() {
    return {
      categories: [
        {
          title: 'Обычные',
          items: ['Попить воды', 'Один раз "Мяу"', 'Один раз "Ня"', 'Один раз "Фыр"', 'Сердечко в чатик', 'Поцелуй чатику'],
          percentage: 50
        },
        {
          title: 'Необычные',
          items: ['Слова поддержки', 'Сердечки в глазах на 5 мин', 'Звездочки в глазах на 5 мин', 'Размять ладошки'],
          percentage: 25
        },
        {
          title: 'Редкие',
          items: ['Детте слишком близко на 1 мин', 'Маленькая Детте на 1 мин', 'Потягушки', 'Чиби Детте на 20 мин'],
          percentage: 16
        },
        {
          title: 'Эпические',
          items: ['Ещё одна попытка покрутить', 'Модель терминатора на 20 мин'],
          percentage: 7
        },
        {
          title: 'Легендарные',
          items: ['Вип на месяц', 'Шепот режим наа 3 мин', 'Роспись в стиме'],
          percentage: 1.5
        },
        {
          title: 'Артифакты',
          items: ['Личная (милая) открытка от Детте', 'Личная открытка-подкат от Детте'],
          percentage: 0.5
        }
      ]
    };
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
  gap: 20px;
  padding: 20px;
}
</style>
