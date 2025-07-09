import { createRouter, createWebHistory } from 'vue-router';
   import Window from '../components/Window.vue';
   import SettingsPanel from '../components/SettingsPanel.vue';
   import RouletteSettingsPanel from '../components/RouletteSettingsPanel.vue';
   import RouletteSettingsCardEditor from '../components/RouletteSettingsCardEditor.vue';

   const routes = [
     { path: '/', component: Window },
     { path: '/settings', component: SettingsPanel },
     { path: '/roulette-settings', component: RouletteSettingsPanel },
     //{ path: '/edit-category/:items', name: 'edit-category', component: RouletteSettingsCardEditor },
     //После подключения бд будем передавать индекс и делать по нему запрос
     { path: '/edit-category/:index', name: 'edit-category', component: RouletteSettingsCardEditor },
   ];

   const router = createRouter({
     history: createWebHistory(),
     routes,
   });

   export default router;