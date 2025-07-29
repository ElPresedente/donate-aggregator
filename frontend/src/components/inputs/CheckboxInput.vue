<template>
    <div class="checkbox-wrapper">
      <input
        type="checkbox"
        :id="id"
        :name="name"
        v-model="localValue"
        class="checkbox-input"
      />
      <label :for="id">{{ label }}</label>
    </div>
  </template>
  
  <script>
  import { computed } from 'vue';
  
  export default {
    name: 'CheckboxInput',
    props: {
      modelValue: Boolean,
      id: String,
      name: String,
      label: String,
    },
    setup(props, { emit }) {
      const localValue = computed({
        get() {
          return props.modelValue;
        },
        set(value) {
          emit('update:modelValue', value);
        },
      });
  
      return {
        localValue,
      };
    },
  };
  </script>
  
  <style scoped>
  .checkbox-wrapper {
    display: flex;
    align-items: center;
  }
  .checkbox-input {
    width: auto;
    margin-right: 8px;
  }
  </style>