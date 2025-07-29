<template>
<select
    :id="id"
    :name="name"
    v-model="localValue"
    :placeholder="placeholder"
    class="custom-input"
>
    <option v-if="placeholder" value="" disabled selected>{{ placeholder }}</option>
    <option v-for="(option, index) in options" :key="index" :value="option">
    {{ option }}
    </option>
</select>
</template>

<script>
import { computed } from 'vue';

export default {
name: 'SelectInput',
props: {
    modelValue: [String, Number],
    id: String,
    name: String,
    placeholder: String,
    options: {
    type: Array,
    default: () => [],
    },
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
.custom-input {
  width: calc(100% - 30px - 20px);
  padding: 5px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 20px;
}
</style>