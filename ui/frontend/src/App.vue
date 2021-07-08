<template>
  <div id="app">
    <nav class="navbar navbar-expand-lg navbar-light bg-secondary mb-3">
      <a class="navbar-brand" href="#">Postee UI</a>
    </nav>

    <div class="container-fluid">
      <div class="row content">
        <div v-if="!isOnLogin" class="col-sm-2">
          <ul class="nav flex-column nav-pills">
            <li class="nav-item">
              <router-link
                active-class="active"
                :to="{ name: 'outputs' }"
                class="nav-link"
                >Outputs</router-link
              >
            </li>
            <li class="nav-item">
              <router-link
                active-class="active"
                :to="{ name: 'routes' }"
                class="nav-link"
                >Routes</router-link
              >
            </li>
            <li class="nav-item">
              <router-link
                active-class="active"
                :to="{ name: 'templates' }"
                class="nav-link"
                >Templates</router-link
              >
            </li>
            <li class="nav-item">
              <router-link
                active-class="active"
                :to="{ name: 'settings' }"
                class="nav-link"
                >Settings</router-link
              >
            </li>
            <li class="nav-item">
              <a class="nav-link" @click="doLogout" href="#">Logout</a>
            </li>
          </ul>
        </div>
        <div
          v-bind:class="[
            { 'col-sm-12': isOnLogin },
            { 'col-sm-9': !isOnLogin },
          ]"
        >
          <div v-if="error" class="alert alert-danger" role="alert">
            {{ error }}
          </div>
          <router-view></router-view>
        </div>
      </div>
    </div>
  </div>
</template>
<style>
</style>
<script>
import { mapState } from "vuex";

export default {
  name: "App",
  watch: {
    $route(to) {
      this.$store.commit("error/clear");
      if (["home", "outputs"].indexOf(to.name) >= 0 && !this.$store.state.flags.all.loaded) {
        this.startLoading();
      }
    },
  },
  data() {
    return {
    };
  },
  computed: {
    ...mapState({
      error(state) {
        return state.error.message;
      },
    }),
    isOnLogin() {
      return this.$route.name === 'login'
    }
  },
  methods: {
    doLogout() {
      this.$store.dispatch("account/logout");
    },
    startLoading() {
      this.$store.dispatch("load");
      this.$store.dispatch("stats/load");
    },
  },
  mounted() {
    if (this.$store.state.account.authenticated) {
      this.startLoading();
    } else {
      if (!this.isOnLogin) {
        this.$store
          .dispatch("account/login")
          .then(() => {
            this.startLoading();
          })
          .catch((error) => {
            if (!error) {
              //check failed - session is invalid
              this.$router.push({ name: "login" });
            }
          });
      }
    }
  },
};
</script>

