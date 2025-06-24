<template>
  <div class="settings-card">
    <h3>{{ title }}</h3>
    <div v-for="(input, index) in inputsConfig" :key="index" class="input-group">
      <label :for="input.name">{{ input.label }}</label>
      <input
        :type="input.type"
        :id="input.name"
        :name="input.name"
        v-model="localFormData[input.name]"
        :placeholder="input.placeholder"
      />
    </div>
  </div>
</template>

<script>
export default {
  name: 'LoginSettingsCard',
  props: {
    title: {
      type: String,
      required: true,
    },
    inputsConfig: {
      type: Array,
      required: true,
      validator: (config) => {
        return config.every(
          (input) =>
            input.name &&
            input.label &&
            input.type &&
            typeof input.name === 'string' &&
            typeof input.label === 'string' &&
            typeof input.type === 'string'
        );
      },
    },
    formData: {
      type: Object,
      required: true,
    },
  },
  computed: {
    localFormData: {
      get() {
        return this.formData;
      },
      set(newValue) {
        this.$emit('update:formData', newValue);
      },
    },
  },
};
</script>

<style scoped>
.settings-card {
  border: 1px solid #ccc;
  padding: 16px;
  margin: 16px;
  border-radius: 8px;
}
.input-group {
  margin-bottom: 12px;
}
label {
  display: block;
  margin-bottom: 4px;
}
input {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}
button {
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}
button:hover {
  background-color: #0056b3;
}
</style>