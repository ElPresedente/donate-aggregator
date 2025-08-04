<template>
  <button class="link-button" @click="handleClick">
    {{ name }}
  </button>
</template>

<script>

export default {
  name: 'LinkButton',
  props: {
    name: {
      type: String,
      default: 'Twitch'
    },
    callback: {
      type: Function,
      required: true,
    },
  },
  methods: {
    handleClick() {
      this.callback();
    },
  },
};

</script>

<style scoped>
.link-button {
  padding: 10px 16px;
  font-size: 16px;
  background-color: #9146FF;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  transition: background-color 0.2s;
}
.link-button:hover {
  background-color: #6600ff;
}
</style>
