<template>
  <div class="col mb-4">
    <div class="card text-left">
      <div class="card-header">
        <b-icon-puzzle></b-icon-puzzle>
        {{ type }}
      </div>
      <div class="card-body">
        <div class="d-flex align-items-center">
          <div class="flex-grow-1">
            <h5 class="card-title">{{ isCommon ? "Defaults" : name }}</h5>
          </div>
        </div>

        <h6 v-if="!isCommon" class="card-subtitle text-muted">
          {{ scanCountMessage }}
        </h6>
      </div>
      <div class="card-footer text-center">
        <router-link
          :to="{ name: 'output', params: { name: name } }"
          class="btn btn-link"
          >Edit</router-link
        >
      </div>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  props: ["type", "name", "enable"],
  data() {
    return {};
  },
  computed: {
    ...mapState({
      scanCount(state) {
        return state.stats.all[this.name];
      },
    }),
    isCommon() {
      return this.type === "common";
    },
    scanCountMessage() {
      return this.scanCount === undefined
        ? "No scans received"
        : [
            this.scanCount < 10000 ? this.scanCount : "9999+",
            " scan",
            this.scanCount === 1 ? "" : "s",
            " received",
          ].join("");
    },
  },
};
</script>