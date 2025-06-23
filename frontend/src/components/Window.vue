<template>
  <div id="body">
    <div v-if="!showSettingsPanel" class="dashboard" id="dashboard">
      <HistoryList :history="history" />
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
import LogList from './components/LogList.vue';
import ControlPanel from './components/ControlPanel.vue';
import SettingsPanel from './components/SettingsPanel.vue';

export default {
  components: {
    LogList,
    ControlPanel,
    SettingsPanel
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
      this.settingsTitle = titleMap[type];
      this.showSettingsPanel = true;
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
