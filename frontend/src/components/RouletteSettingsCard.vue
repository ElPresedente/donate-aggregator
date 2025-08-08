<template>
  <div class="card">
    <div class="header-panel">
      <h3 :style="{ color: color }">{{ title }}</h3>
      <button class="edit" @click="openEditor">
        <img src="./../assets/images/gear.svg" alt="Изменить" width="14px" height="14px">
      </button>
    </div>
    
    <ul>
      <template v-if="sectors.length > 0">
        <li v-for="item in sectors" :key="item">{{ item }}</li>
      </template>
      <template v-else>
        <li>...</li>
      </template>
    </ul>
    <div class="stats">{{ percentage }}% <span>🎯</span></div>
  </div>
</template>

<script>
import { useRouter } from "vue-router";

export default {
  name: 'Card',
  props: {
    title: String,
    sectors: Array,
    percentage: Number,
    color: {
      type: String,
      default: '#ccc'
    },
    index: Number
  },
  setup(props) {
    const router = useRouter();
    const openEditor = () => {
      router.push({ name: 'edit-category', params: { index: JSON.stringify(props.index) } });
    };
    return { openEditor };
  }
}
</script>

<style scoped>
.card {
  background-color: #2a2a2a;
  border-radius: 8px;
  padding: 15px;
  width: 200px;
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
  margin: 10px 0 10px;
}

.card li {
  font-size: 12px;
  color: #999;
  margin-bottom: 5px;
}

.card .stats {
  font-size: 14px;
  text-align: right;
  margin-top: auto;
}

.card .stats span {
  margin-left: 5px;
}

.edit{
  display: flex;
  justify-items: center;
  background-color: #303030;
  border: 1px solid transparent;
  align-self: center;
  border-radius: 8px;
  padding: 10px;
  font-weight: 600;
  color: #f57d07;
  position: absolute; /* Абсолютное позиционирование */
  right: 0; /* Прижимаем к правому краю */
  top: 50%; /* Центрируем по вертикали */
  transform: translateY(-50%);
}
.edit:hover{
  background-color: #464646;
  cursor: pointer;
}
.header-panel{
  display: flex;
  justify-content: center;
  position: relative;
  align-items: center;
  min-height: 34px; 
}
</style>