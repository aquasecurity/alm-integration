<template>
  <div class="col mb-4">
    <div class="card text-left">
      <div class="card-header">
        <b-icon-braces></b-icon-braces>
      </div>
      <div class="card-body">
        <div class="d-flex align-items-center">
          <div class="flex-grow-1">
            <h5 class="card-title">{{ name }}</h5>
          </div>
        </div>

      </div>
      <div class="card-footer text-center">
        <router-link
          :to="{ name: 'template', params: { name: name } }"
          class="btn btn-link"
          >Edit</router-link
        >
      </div>
    </div>
  </div>
</template>

<script>

export default {
  props: ["name"],
  data() {
    return {};
  }
};
</script>