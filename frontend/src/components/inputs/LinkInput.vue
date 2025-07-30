<template>
    <input
        type="text"
        :id="id"
        :name="name"
        v-model="localValue"
        :placeholder="placeholder"
        class="link-input"
        @input="parseLink"
    />
</template>

<script>
import { computed } from 'vue';

export default {
name: 'LinkInput',
props: {
    modelValue: String,
    id: String,
    name: String,
    placeholder: String,
    domain: String,
    params: {
    type: Array,
    default: () => [],
    },
    callback: {
    type: Function,
    default: () => {},
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

    const parseLink = () => {
    const url = localValue.value;
    let result = {};

    try {
        console.log(url)
        console.log(props)
        const urlObj = new URL(url);
        console.log(urlObj)
        // Проверяем домен, если он указан
        if (props.domain && urlObj.hostname !== props.domain) {
        props.params.forEach((param) => {
            result[param] = '';
        });
        props.callback(result);
        return;
        }

        // Извлекаем указанные GET-параметры
        props.params.forEach((param) => {
            result[param] = urlObj.searchParams.get(param) || '';
            console.log('------------')
            console.log(urlObj.searchParams.get(param) || '')
            console.log(result[param])
        });

        // Вызываем callback с результатом
        console.log("result")
        console.log(result['ref'])
        props.callback(result);
    } catch (e) {
        // Если URL некорректен, возвращаем пустые значения для всех параметров
        props.params.forEach((param) => {
        result[param] = '';
        });
        props.callback(result);
    }
    };

    return {
    localValue,
    parseLink,
    };
},
};
</script>

<style scoped>
.link-input {
width: 100%;
padding: 5px;
border: 1px solid #ddd;
border-radius: 4px;
font-size: 20px;
}
</style>