<template>
   <div>
    <SettingsCard
      title="Donatty"
      :inputsConfig="donattyCfg"
      :formData="donatty"
      @update:formData="updateFormData('donatty', $event)"
    />
    <SettingsCard
      title="DonatPay"
      :inputsConfig="donatpayCfg"
      :formData="donatpay"
      @update:formData="updateFormData('donatpay', $event)"
    />
  </div>
  <section class="card stretch" id="settings-panel">
    <button class="btn save" @click="handleSave()">Сохранить все</button>
    <button class="btn back" @click="goBack()">← Назад</button>
  </section>
</template>

<script>
import { useRouter } from 'vue-router';
import SettingsCard from './LoginSettingsCard.vue'
import { onMounted, ref } from 'vue';
import { FrontendDispatcher } from '../../wailsjs/go/main/App'
export default {
  setup() {
    const router = useRouter();
    const goBack = () => router.go(-1); //router.push('/');
    const donattyCfg = [
      {
        name: 'donattyToken',
        label: 'Токен Донатти',
        type: 'text',
        placeholder: 'Введите токен',
      },
      {
        name: 'donattyUrl',
        label: 'URL Донатти',
        type: 'url',
        placeholder: 'Введите ссылка',
      },
    ]
    const donatpayCfg = [
      {
        name: 'donatpayToken',
        label: 'Токен Донатпей',
        type: 'text',
        placeholder: 'Введите токен',
      },
      {
        name: 'donatpayUserId',
        label: 'Пользовательский ID',
        type: 'text',
        placeholder: 'Введите ID пользователя',
      }
    ]
    const donatty = ref([{donattyToken: '', donattyUrl: ''}])
    const donatpay = ref([{donatpayToken: '', donatpayUserId: ''}])

    const updateFormData = (target, newData) => {
      if (target === 'donatty')   donatty.value = newData
      if (target === 'donatpay')  donatpay.value = newData
    }
        
    const handleSave = () => {
      const settingsToSave = {
        settings:  [
          {name: "donattyToken",   value: donatty.value.donattyToken},
          {name: "donattyUrl",     value: donatty.value.donattyUrl},
          {name: "donatpayToken",  value: donatpay.value.donatpayToken},
          {name: "donatpayUserId", value: donatpay.value.donatpayUserId}
        ]
      }
      FrontendDispatcher("updateSettings", JSON.stringify(settingsToSave));
      // Отправка на сервер
    }
    onMounted(() =>{
      FrontendDispatcher("getSettings", "");
      window.runtime.EventsOn('SettingsData', (data) => {
        console.log('📦 Группы:', data)
        data.forEach(setting => {
        switch (setting.name) {
          case 'donattyToken':
            donatty.value.donattyToken = setting.value;
            break;
          case 'donattyUrl':
            donatty.value.donattyUrl = setting.value;
            break;
          case 'donatpayToken':
            donatpay.value.donatpayToken = setting.value;
            break;
          case 'donatpayUserId':
            donatpay.value.donatpayUserId = setting.value;
            break;
          default:
            console.warn(`⚠️ Неизвестная настройка: ${setting.name}`);
          }
        });
    });
  });
    return { 
      goBack,
      donattyCfg,
      donatpayCfg,
      donatty,
      donatpay,
      updateFormData,
      handleSave
    };
  },
  components: { SettingsCard },
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

.settings-content {
  flex: 1;
  padding-top: 10px;
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
