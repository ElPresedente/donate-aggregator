<template>
  <div class="container">
    <RouletteSettingsCard
      v-for="(category, index) in categories"
      :key="index"
      :title="category.title"
      :sectors="category.sectors"
      :percentage="category.percentage"
      :color="category.color"
      :index="index"
    />
  </div>
  <SettingsCard
      title="Настройки рулетки"
      :inputsConfig="rouletteSettingsCfg"
      :formData="rouletteSettings"
      @update:formData="updateFormData('rouletteSettings', $event)"
  >
    <button class="btn save" @click="handleSave()">Сохранить настройки</button>
  </SettingsCard>


  <section class="card stretch" id="settings-panel">
    <button class="btn back" @click="goBack()">← Назад</button>
  </section>
</template>

<script>
import { useRouter } from 'vue-router';
import RouletteSettingsCard from './RouletteSettingsCard.vue'
import SettingsCard from './SettingsCard.vue'
import { onMounted, ref, onUnmounted } from "vue";
import { FrontendDispatcher } from '../../wailsjs/go/main/App'

export default {
  name: 'RouletteSettingsPanel',
  components: { 
    RouletteSettingsCard,
    SettingsCard
   },
  setup() {
    let unsubscribes = [];
    const router = useRouter();
    const goBack = () => router.go(-1);
    const categories = ref([])
    const rouletteSettings = ref([{ rollPrice: '', rollPriceIncrease: '' }])
    const rouletteSettingsCfg = [
      {
        name: 'rollPrice',
        label: 'Цена прокрутки рулетки в рублях',
        type: 'number',
        placeholder: 'Введите цену',
      },
      {
        name: 'rollPriceIncrease',
        label: 'Увеличение цены прокрутки рулетки в рублях',
        type: 'number',
        placeholder: 'Введите цену',
      },
    ]
    const updateFormData = (target, newData) => {
      if (target === 'rouletteSettings')   rouletteSettings.value = newData
    }
    const handleSave = () => {
      const settingsToSave = {
        settings:  [
          {name: "rollPrice",        value: String(rouletteSettings.value.rollPrice)},
          {name: "rollPriceIncrease",value: String(rouletteSettings.value.rollPriceIncrease)},
        ]
      }
      FrontendDispatcher("updateRouletteSettings", JSON.stringify(settingsToSave));
    }
    onMounted(() => {
      unsubscribes.push(
        window.runtime.EventsOn('sectorsData', (data) => {
          categories.value = data 
        })
      );
      window.runtime.EventsOn('rouletteSettingsData', (data) => {
        data.forEach(setting => {
          switch (setting.name) {
            case 'rollPrice':
              rouletteSettings.value.rollPrice = setting.value;
              break;
            case 'rollPriceIncrease':
              rouletteSettings.value.rollPriceIncrease = setting.value;
              break;
            default:
              console.warn(`⚠️ Неизвестная настройка: ${setting.name}`);
          }
        });
      });
      FrontendDispatcher("getRouletteSectors", "");
      FrontendDispatcher("getRouletteSettings", "");
    });
    onUnmounted(() => {
      unsubscribes.forEach(unsub => unsub());
    });
    return { 
      handleSave,
      rouletteSettings,
      categories,
      rouletteSettingsCfg,
      updateFormData,
      goBack
     }
  },  
};
</script>

<style scoped>
body {
  background-color: #1a1a1a;
  color: #fff;
  font-family: Arial, sans-serif;
}

.container {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 15px;
  padding: 10px;
}

.card {
  background-color: #1e1e1e;
  border-radius: 12px;
  padding: 20px;
  margin: 20px;
  box-shadow: 0 0 10px #00000070;
  display: flex;
  flex-direction: column;
}
.btn.back {
  font-size: large;
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

.btn.back:hover {
  background-color: #4b5563;
}

.btn.save {
  margin-top: 20px;
  font-size: large;
  padding: 10px;
  border-radius: 8px;
  border: none;
  color: white;
  font-weight: bold;
  cursor: pointer;
  transition: background 0.2s ease;
  background-color: #22c55e;
}

.btn.save:hover {
  background-color: #16a34a;
}
</style>
