<template>
  <div id="body">
    <div v-if="!showSettingsPanel" class="dashboard" id="dashboard">
      <LogList/>
      <ControlPanel
        :donattyConnected="donattyConnected"
        :donatepayConnected="donatepayConnected"
        @enable="enable"
        @disable="disable"
        @restart="restart"
        @spin="spin"
        @show-settings="showSettings"
        @show-roulette-settings="showRouletteSettings"
      />
    </div>
  </div>
</template>

<script>
import { useRouter } from 'vue-router';
import LogList from './LogList.vue';
import ControlPanel from './ControlPanel.vue';

export default {
  components: {
    LogList,
    ControlPanel
  },
  setup() {
    const router = useRouter();
    return { router };
  },
  data() {
    return {
      showSettingsPanel: false,
      settingsTitle: '',
      donattyConnected: false,
      donatepayConnected: false
    };
  },
  methods: {
    showSettings() {
      this.router.push('/settings');
    },
    showRouletteSettings() {
      this.router.push('/roulette-settings');
    },
    goBack() {
      this.showSettingsPanel = false;
    },
    enable() {
      this.donattyConnected = true;
    },
    disable() {
      this.donattyConnected = false;
    },
    restart() {
      alert('Рулетка перезапущена');
    },
    spin() {
      alert('Крутим рулетку!');
    }
  }
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
