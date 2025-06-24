import { createRouter, createWebHistory } from 'vue-router';
   import Window from '../components/Window.vue';
   import SettingsPanel from '../components/SettingsPanel.vue';
   import RouletteSettings from '../components/RouletteSettings.vue';

   const routes = [
     { path: '/', component: Window },
     { path: '/settings', component: SettingsPanel },
     { path: '/roulette-settings', component: RouletteSettings },
   ];

   const router = createRouter({
     history: createWebHistory(),
     routes,
   });

   export default router;