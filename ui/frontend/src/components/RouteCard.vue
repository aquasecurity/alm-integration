<template>
  <div class="col mb-4">
    <div class="card text-left">
      <div class="card-header">
        <b-icon-shuffle></b-icon-shuffle>
      </div>
      <div class="card-body">
        <div class="d-flex align-items-center">
          <div class="flex-grow-1">
            <h5 class="card-title">{{ name }}</h5>
          </div>
        </div>

        <h6 class="card-subtitle text-muted">
          {{ outputCnt }} outputs configured
        </h6>
      </div>
      <div class="card-footer text-center">
        <router-link
          :to="{ name: 'route', params: { name: name } }"
          class="btn btn-link"
          >Edit</router-link
        >
      </div>
    </div>
  </div>
</template>

<script>

export default {
  props: ["name", "outputCnt"],
  data() {
    return {};
  }
};
</script>