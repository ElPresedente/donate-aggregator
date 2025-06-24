<template>
  <div class="card">
    <h3>{{ title }}</h3>
    <ul>
      <li v-for="item in items" :key="item">{{ item }}</li>
    </ul>
    <div class="stats">{{ percentage }}% <span>🎯</span></div>
  </div>
</template>

<script>
export default {
  name: 'Card',
  props: {
    title: String,
    items: Array,
    percentage: Number
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
  margin: 0 0 10px;
  color: #ccc;
}

.card ul {
  list-style: none;
  padding: 0;
  margin: 0 0 10px;
}

.card li {
  font-size: 12px;
  color: #999;
  margin-bottom: 5px;
}

.card .stats {
  font-size: 14px;
  text-align: right;
}

.card .stats span {
  margin-left: 5px;
}
</style>