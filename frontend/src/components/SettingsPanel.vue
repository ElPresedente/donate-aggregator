<template>
  <section class="card stretch" id="settings-panel">
    <header class="card-header" id="settings-title">{{ settingsTitle }}</header>
    <div class="settings-content">
      <p>Настройки: <strong>{{ settingsTitle }}</strong></p>
    </div>
    <button class="btn back" @click="$emit('go-back')">← Назад</button>
  </section>
</template>

<script>
export default {
  name: 'SettingsPanel',
  props: {
    settingsTitle: {
      type: String,
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
  height: 100%;
}

.card-header {
  font-size: 1.2rem;
  margin-bottom: 10px;
  font-weight: bold;
}

.settings-content {
  flex: 1;
  padding-top: 10px;
}

.btn.back {
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
</style>
