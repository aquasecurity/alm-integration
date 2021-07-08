<template>
  <div v-if="show" class="form-group form-input" v-bind:class="{ 'row pr-3' : inline }">
    <label v-bind:class="{ 'col-form-label col-sm-2': inline, 'form-label': !inline }" :for="id">{{ label }}</label>
    <input
      :type="inputType"
      :value="value"
      :name="nameOrId"
      @input="inputHandler"
      class="form-control"
      v-bind:class="{ 'is-invalid' : errorMsg , 'col-sm-10' : inline }"
      :id="id"
    />
    <small v-show="!!description" class="form-text text-muted">{{
      description
    }}</small>
    <div class="form-text invalid-feedback">{{ errorMsg }}</div>
  </div>
</template>
<script>
export default {
  data() {
    return {
    }
  },
  computed: {
    nameOrId(){
      return this.name || this.id
    },
  },
  props: {
    errorMsg : [String, Boolean],
    name: String,
    inputType: {
      type: String,
      default: "input",
    },
    label: String,
    value: [String, Number],
    id: String,
    description: String,
    show: {
      type: Boolean,
      default: true,
    },
    inline:{
      type: Boolean,
      default: false,
    },
    inputHandler: Function,
    validator: Object,
  },
  methods : {
  },
  mounted()  {
    if (this.validator) {
      this.validator.register(this.id, this.label, this.nameOrId)
    }
  }
};
</script>