<template>
  <div id="body">
    <div v-if="!showSettingsPanel" class="dashboard" id="dashboard">
      <LogList/>
      <ControlPanel/>
    </div>
  </div>
</template>

<script>
import LogList from './LogList.vue';
import ControlPanel from './ControlPanel.vue';

export default {
  components: {
    LogList,
    ControlPanel
  },
  setup() {
  },
  data() {
    return {
      showSettingsPanel: false,
      settingsTitle: '',
      donattyConnected: false,
      donatepayConnected: false
    };
  },
};
</script>

<style>
  #dashboard {
    height: 90vh;
    display: flex;
    gap: 20px;
    padding: 20px;
  }
</style>
