<template>
  <section class="card stretch" id="right-panel">
    <div class="card-block" id="main-controls">
      <header class="card-header">Управление рулеткой</header>
      <div class="status">
        <span :class="{ online: donattyConnected }">
          {{ donattyConnected ? '✅ Donatty' : '❌ Donatty' }}
        </span>
        <span :class="{ online: donatepayConnected }">
          {{ donatepayConnected ? '✅ Donatpay' : '❌ Donatpay' }}
        </span>
      </div>
      <div class="controls">
        <button class="btn green" @click="$emit('enable')">Включить</button>
        <button class="btn red" @click="$emit('disable')">Выключить</button>
        <button class="btn blue" @click="$emit('restart')">Перезапустить</button>
        <button class="btn gold" @click="$emit('spin')">Крутить</button>
      </div>
    </div>
    <div class="card-block settings-buttons">
      <button class="btn gray" @click="$emit('show-settings')">
        Настройка подключения
      </button>
      <button class="btn gray" @click="$emit('show-roulette-settings')">
        Настройка рулетки
      </button>
    </div>
  </section>
</template>

<script>
export default {
  name: 'ControlPanel',
  props: {
    donattyConnected: Boolean,
    donatepayConnected: Boolean
  }
};
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
  height: 100%;
}

.card-header {
  font-size: 1.2rem;
  margin-bottom: 10px;
  font-weight: bold;
}

.status span {
  display: block;
  margin: 5px 0;
  font-size: 0.95rem;
  opacity: 0.9;
}

.status span.online {
  color: #22c55e;
}

.controls {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 15px;
}

.btn {
  flex: 1 1 45%;
  padding: 10px;
  border-radius: 8px;
  border: none;
  color: white;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s ease;
}

.btn.green {
  background-color: #22c55e;
  transition: background-color 0.2s ease; /* Плавный переход для фона */
}

.btn.green:hover {
  background-color: #16a34a; /* Темнее на 20% для наведения */
}

.btn.red {
  background-color: #ef4444;
  transition: background-color 0.2s ease;
}

.btn.red:hover {
  background-color: #dc2626; /* Темнее на 20% */
}

.btn.blue {
  background-color: #3b82f6;
  transition: background-color 0.2s ease;
}

.btn.blue:hover {
  background-color: #2563eb; /* Темнее на 20% */
}

.btn.gold {
  background-color: #fbbf24;
  transition: background-color 0.2s ease;
}

.btn.gold:hover {
  background-color: #d97706; /* Темнее на 20% */
}

.btn.gray {
  background-color: #6b7280;
}

.btn.gray:hover {
  background-color: #4b5563;
}

.settings-buttons {
  margin-top: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>
