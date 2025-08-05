<template>
  <div class="section-wrapper">
    <span class="section-title">{{ title }}</span>
    <div class="section-content">
      <slot />
    </div>
  </div>
</template>

<script>
export default {
  name: 'ControlPanelSection',
  props: {
    title: {
      type: String,
      required: true,
    },
  },
};
</script>

<style scoped>
.section-wrapper {
  position: relative;
  padding: 24px;
  padding-top: 32px;
  border: 2px solid white;
  border-radius: 12px;
  background-color: transparent;
  margin-bottom: 20px;
}

.section-title {
  position: absolute;
  top: 0;
  left: 20px;
  transform: translateY(-60%);
  background-color: #1e1e1e; /* ← подложка под заголовком, под цвет фона */
  padding: 0 12px;
  color: white;
  font-weight: 600;
  font-size: 18px;
  z-index: 1;
}

.section-content {
  color: white;
}
</style>
