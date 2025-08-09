<template>
  <div class="input-wrapper">
    <input
      :type="visible ? 'text' : 'password'"
      :placeholder="placeholder"
      :id="id"
      :name="name"
      autocomplete="off"
      class="custom-input"
      v-model="inputValue"
    />
    <button
      type="button"
      class="toggle-password"
      @click="visible = !visible"
    >
      {{ visible ? '🙈' : '👁️' }}
    </button>
  </div>
</template>

<script setup>
import { ref, watch, defineProps, defineEmits } from 'vue'

const props = defineProps({
  modelValue: String,
  placeholder: String,
  id: String,
  name: String
})

const emit = defineEmits(['update:modelValue'])

const visible = ref(false)
const inputValue = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  inputValue.value = val
})

watch(inputValue, (val) => {
  emit('update:modelValue', val)
})
</script>

<style scoped>
input[type="password"]::-ms-reveal,
input[type="password"]::-ms-clear,
input[type="password"]::-webkit-credentials-auto-fill-button,
input[type="password"]::-webkit-password-toggle-button {
  display: none !important;
  appearance: none;
}
.custom-input {
    width: calc(100% - 30px - 20px);
    padding: 5px;
    border: 1px solid #ddd;
    border-radius: 4px;
    font-size: 20px;
}
.input-wrapper {
    width: 100%;
    position: relative;
    display: flex;
    align-items: center;
}
.toggle-password {
    position: absolute;
    right: 10px;
    background: none;
    border: none;
    cursor: pointer;
    font-size: 16px;
    padding: 0;
    margin: 0;
}
</style>
