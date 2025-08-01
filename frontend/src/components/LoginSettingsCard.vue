<template>
  <div class="settings-card">
    <h2>{{ title }}</h2>

    <div v-for="(input, index) in inputsConfig" :key="index" class="input-group">
      <label :for="input.name">{{ input.label }}</label>
      <div class="input-wrapper">
        <component
          :is="resolveInputComponent(input.type)"
          v-model="localFormData[input.name]"
          :placeholder="input.placeholder"
          :id="input.name"
          :name="input.name"
          :index="index"
          :options="input.options"
          :domain="input.domain"
          :params="input.params"
          :callback="input.callback"
          @toggle-password="togglePassword(index)"
          :visible="showPassword[index]"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { ref } from 'vue';
import PasswordInput from './inputs/PasswordInput.vue';
import NumberInput from './inputs/NumberInput.vue';
import SelectInput from './inputs/SelectInput.vue';
import CheckboxInput from './inputs/CheckboxInput.vue';
import LinkInput from './inputs/LinkInput.vue';

export default {
  name: 'LoginSettingsCard',
  components: { PasswordInput, NumberInput, SelectInput, CheckboxInput, LinkInput },
  props: {
    title: String,
    inputsConfig: {
      type: Array,
      default: () => [],
    },
    formData: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props, { emit }) {
    const showPassword = ref(props.inputsConfig.map(() => false));

    const togglePassword = (index) => {
      showPassword.value[index] = !showPassword.value[index];
    };

    const resolveInputComponent = (type) => {
      switch (type) {
        case 'pass': return PasswordInput;
        case 'number': return NumberInput;
        case 'select': return SelectInput;
        case 'checkbox': return CheckboxInput;
        case 'link': return LinkInput;
        default: return 'input';
      }
    };

    return {
      showPassword,
      togglePassword,
      resolveInputComponent,
    };
  },
  computed: {
    localFormData: {
      get() {
        return this.formData;
      },
      set(val) {
        this.$emit('update:formData', val);
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
.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}
h2 {
  text-align: left;
  margin: 10px 0 0 10px;
}
label {
  left: 0px;
  top: 0px;
  display: block;
  margin-bottom: 4px;
}
input, select {
  width: calc(100% - 30px - 20px);
  padding: 5px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 20px;
}
input[type="checkbox"] {
  width: auto;
  margin-right: 8px;
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