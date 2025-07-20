<template>
  <input
    type="text"
    :value="modelValue"
    @input="onInput"
    @keypress="onKeyPress"
    @paste="onPaste"
    class="w-full px-4 py-2 border rounded-lg outline-none focus:ring-2 focus:ring-blue-500"
    :placeholder="placeholder"
  />
</template>

<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  modelValue: String,
  placeholder: String,
})

const emit = defineEmits(['update:modelValue'])

// Запрещаем ввод нецифровых символов при вводе с клавиатуры
function onKeyPress(e) {
  const char = String.fromCharCode(e.charCode)
  if (!/^\d$/.test(char)) {
    e.preventDefault()
  }
}

// Обработка вставки — убираем все нецифры
function onPaste(e) {
  e.preventDefault()
  const pasted = (e.clipboardData || window.clipboardData).getData('text')
  const cleaned = pasted.replace(/\D/g, '')
  emit('update:modelValue', cleaned)
}

// Очищаем нецифры на случай, если что-то просочилось
function onInput(e) {
  const cleaned = e.target.value.replace(/\D/g, '')
  emit('update:modelValue', cleaned)
}
</script>


<style scoped>
/* В App.vue <style> или global.css */
input[type="password"]::-ms-reveal,
input[type="password"]::-ms-clear,
input[type="password"]::-webkit-credentials-auto-fill-button,
input[type="password"]::-webkit-password-toggle-button,
input[type="number"]::-webkit-outer-spin-button,
input[type="number"]::-webkit-inner-spin-button {
  display: none !important;
  appearance: none;
}

input[type="number"] {
  -moz-appearance: textfield;
}

.custom-input {
  width: calc(100% - 30px - 20px);
  padding: 5px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 20px;
}
</style>
