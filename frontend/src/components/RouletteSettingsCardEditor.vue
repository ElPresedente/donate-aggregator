<template>
  <div>
    <h2>Редактирование категории</h2>
    <!-- Здесь будет форма редактирования -->
  </div>
</template>

<script>
import { useRoute } from 'vue-router';

export default {
  name: 'RouletteSettingsCardEditor',
  setup() {
    const route = useRoute();
    const index = Number(route.params.index);

    // Вывод аргумента в консоль
    console.log('Index из параметра маршрута:', index);

    return { index };
  }
};
</script>

<style scoped>
/* Стили для компонента */
div {
  padding: 20px;
}

h2 {
  color: #ccc;
}
</style>