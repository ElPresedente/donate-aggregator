<template>
  <div id="body">
    <div v-if="!showSettingsPanel" class="dashboard" id="dashboard">
      <LogList :history="history" />
      <ControlPanel
        :donattyConnected="donattyConnected"
        :donatepayConnected="donatepayConnected"
        @enable="enable"
        @disable="disable"
        @restart="restart"
        @spin="spin"
        @show-settings="showSettings"
      />
    </div>
    <SettingsPanel
      v-else
      :settingsTitle="settingsTitle"
      @go-back="goBack"
    />
  </div>
</template>

<script>
import { useRouter } from 'vue-router';
import LogList from './LogList.vue';
import ControlPanel from './ControlPanel.vue';
import SettingsPanel from './SettingsPanel.vue';

export default {
  components: {
    LogList,
    ControlPanel,
    SettingsPanel
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
      donatepayConnected: false,
      history: [
        'kamicute2 получает Пат-пат',
        'kamicute2 получает Пат-пат',
        'kamicute2 получает Пат-пат'
      ]
    };
  },
  methods: {
    showSettings(type) {
      const titleMap = {
        connection: 'Настройка подключения',
        roulette: 'Настройка рулетки'
      };
      this.router.push('/settings');
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
