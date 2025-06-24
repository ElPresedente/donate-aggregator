<template>
  <div class="settings-card">
    <h2>{{ title }}</h2>
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
  background-color: #1e1e1e;
  border-radius: 12px;
  box-shadow: 0 0 10px #00000070;
  padding: 0 10px 10px 10px;
  margin: 16px;
  display: flex;
  flex-direction: column;
}
.input-group {
  margin-bottom: 12px;
}
h2{
  text-align: left;
  margin: 10px 0 0 10px;
}
label {
  left: 0px;
  top: 0px;
  display: block;
  margin-bottom: 4px;
}
input {
  width: calc(100% - 30px - 20px);
  padding: 5px;
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