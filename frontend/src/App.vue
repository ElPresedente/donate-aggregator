<script setup>
import Toast from './components/actions/Toast.vue'
import { onMounted, ref } from 'vue'
import { useToastStore } from './stores/toastStore'

const toastStore = useToastStore()
const toastInstance = ref(null)

onMounted(() => {
  toastStore.setRef(toastInstance.value)
})
</script>

<template>
  <router-view />
  <Toast ref="toastInstance" />
</template>


<style>
#app{
  height: auto;
}
</style>