import { createRouter, createWebHistory } from 'vue-router';
   import Window from '../components/Window.vue';
   import SettingsPanel from '../components/SettingsPanel.vue';
   import RouletteSettingsPanel from '../components/RouletteSettingsPanel.vue';

   const routes = [
     { path: '/', component: Window },
     { path: '/settings', component: SettingsPanel },
     { path: '/roulette-settings', component: RouletteSettingsPanel },
   ];

   const router = createRouter({
     history: createWebHistory(),
     routes,
   });

   export default router;