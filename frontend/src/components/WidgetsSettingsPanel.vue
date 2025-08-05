<template>
   <div>
        <button class="btn back" @click="showRouletteSettings()">Настройки рулетки</button>
        <button class="btn back" @click="goBack()">← Назад</button>
    </div>
</template>

<script>
import { useRouter } from 'vue-router';
export default {
    setup() {
        const router = useRouter();
        const goBack = () => router.go(-1);

        const showRouletteSettings = () => {
            router.push('/roulette-settings');
        };
        
        return { 
            goBack,
            showRouletteSettings,
        };
    }
}
</script>

<style>
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
</style>