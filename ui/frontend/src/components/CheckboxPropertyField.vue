<template>
  <div class="form-group form-check" v-if="show">
    <input
      type="checkbox"
      class="form-check-input"
      :id="id"
      :checked="value"
      @input="inputHandler"
      :name="name || id"
    />
    <label class="form-check-label" :for="id">{{ label }}</label>
    <small v-show="!!description" class="form-text text-muted">{{
      description
    }}</small>
  </div>
</template>
<script>
export default {
  props: {
    name: String,
    label: String,
    value: Boolean,
    id: String,
    description: String,
    show: {
      type: Boolean,
      default: true,
    },
    inputHandler: Function,
  },
};
</script>