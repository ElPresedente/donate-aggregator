<template>
  <section class="card stretch">
    <header class="card-header">История рулетки</header>
    <ul class="card-list">
      <li v-for="(item, index) in history" :key="index">{{ item }}</li>
    </ul>
  </section>
</template>

<script>
export default {
  name: 'LogList',
  props: {
    history: {
      type: Array,
      required: true
    }
  }
};
</script>

<style scoped>
.card {
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