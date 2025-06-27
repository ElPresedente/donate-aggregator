<template>
   <div>
    <SettingsCard
      title="Donatty"
      :inputsConfig="Donatty"
      :formData="formData.Donatty"
      @update:formData="updateFormData('Donatty', $event)"
    />
    <SettingsCard
      title="DonatPay"
      :inputsConfig="Donatpay"
      :formData="formData.Donatpay"
      @update:formData="updateFormData('Donatpay', $event)"
    />
    <SettingsCard
      title="StreamElements"
      :inputsConfig="StreamElements"
      :formData="formData.StreamElements"
      @update:formData="updateFormData('StreamElements', $event)"
    />
  </div>
  <section class="card stretch" id="settings-panel">
    <button class="btn save" @click="handleSave">Сохранить все</button>
    <button class="btn back" @click="goBack">← Назад</button>
  </section>
</template>

<script>
import { useRouter } from 'vue-router';
import SettingsCard from './LoginSettingsCard.vue'
export default {
  setup() {
    const router = useRouter();
    const goBack = () => router.go(-1); //router.push('/');
    return { goBack };
  },
  components: { SettingsCard },
  data() {
    return {
      Donatty: [
        {
          name: 'token',
          label: 'Токен API',
          type: 'text',
          placeholder: 'Введите токен',
        },
        {
          name: 'url',
          label: 'URL сервиса',
          type: 'url',
          placeholder: 'ссылка',
        },
      ],
      Donatpay:[
        {
          name: 'key',
          label: 'Не помню, чо там надо',
          type: 'text',
          placeholder: 'Реально не помню',
        }
      ], 
      StreamElements: [
        {
          name: 'key',
          label: 'Не помню, чо там надо',
          type: 'text',
          placeholder: 'Реально не помню',
        }
      ],
      formData: {
        Donatty: { token: '', url: '' },
        Donatpay: { key: '' },
        StreamElements: { key: '' },
      },
    };
  },
  
  methods: {
    updateFormData(service, newData) {
      this.formData[service] = newData;
    },
    handleSave() {
      console.log('Сохранённые данные:', this.formData);
      // Отправка на сервер
    },
  },
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
