<template>
  <div>
    <div class="row justify-content-end pb-3 pr-3">
      <button type="submit" class="btn btn-primary mr-2" @click="doSubmit">
        Submit
      </button>
      <button
        v-if="!!name"
        type="button"
        @click="doRemove"
        class="btn btn-outline-primary"
      >
        Remove
      </button>
    </div>
    <div class="card">
      <form @submit.prevent="doSubmit">
        <div class="card-body">
          <PropertyField
            :id="'name'"
            :label="'Name'"
            :value="formValues.name"
            :errorMsg="errors['name']"
            :inputHandler="updateField"
            :validator="v(uniqueName)"
          />
          <div class="form-group form-input">
            <b-tabs content-class="mt-3">
              <b-tab title="Inline" :active="!!formValues.body">
                <label class="form-label" for="input">REGO template:</label>
                <codemirror
                  :value="formValues.body"
                  :options="cmOptions"
                  id="body"
                  name="body"
                  @input="updateBody"
                >
                </codemirror>
                <small class="form-text text-muted">
                  REGO template to render received events
                </small>
              </b-tab>
              <b-tab title="Package" :active="!!formValues.regopackage">
                <PropertyField
                  id="rego-package"
                  label="Package"
                  :value="formValues['rego-package']"
                  description="Rego package with template"
                  :inputHandler="updateTemplateSource"
                />
              </b-tab>
              <b-tab title="Url" :active="!!formValues.url">
                <PropertyField
                  id="url"
                  label="Url"
                  :value="formValues.url"
                  description="Url to load rego from"
                  :inputHandler="updateTemplateSource"
                  :validator="v(url, true)"
                  :errorMsg="errors['url']"
                />
              </b-tab>
              <b-tab title="Legacy" :active="!!formValues.legacyScanRenderer">
                <div class="form-group form-input">
                  <label for="legacyScanRenderer">Legacy</label>
                  <select
                    class="form-select form-control"
                    :value="formValues['legacy-scan-renderer']"
                    id="legacyScanRenderer"
                    name="legacy-scan-renderer"
                    @input="updateTemplateSource"
                  >
                    <option value="html">Html</option>
                    <option value="slack">Slack</option>
                    <option value="jira">Jira</option>
                  </select>
                  <small id="aHelp" class="form-text text-muted"
                    >Use Postee v1 renderers</small
                  >
                </div>
              </b-tab>
            </b-tabs>


          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";
import ValidationMixin from "./validator";
import FormFieldMixin from "./form";
import PropertyField from "./PropertyField.vue";
import { codemirror } from "vue-codemirror";

import "codemirror-rego/mode";
import "codemirror/lib/codemirror.css";

export default {
  data() {
    return {
      fields: {},
      errors: {},
      name: "",
      cmOptions: {
        tabSize: 4,
        mode: "rego",
        lineNumbers: true,
        line: true,
      },
    };
  },
  mixins: [FormFieldMixin, ValidationMixin],
  components: {
    PropertyField,
    codemirror,
  },
  computed: {
    ...mapState({
      formValues(state) {
        //required for mixins
        const found = state.templates.all.filter(
          (item) => item.name === this.name
        );

        const result = found.length ? { ...found[0] } : {};

        return result;
      },
    }),
  },
  methods: {
    doSubmit() {
      if (!this.isFormValid()) {
        return;
      }

      if (this.name) {
        this.$store.dispatch("templates/update", {
          value: this.formValues,
          name: this.name,
        });
      } else {
        this.$store.dispatch("templates/add", this.formValues);
      }
      this.$router.push({ name: "templates" });
    },
    doRemove() {
      this.$store.dispatch("templates/remove", this.name);
      this.$router.push({ name: "templates" });
    },
    updateBody(v) {
      this.formValues.body = v;
      this.formValues.regopackage = undefined
      this.formValues.url = undefined
      this.formValues.legacyScanRenderer = undefined
    },
    updateTemplateSource(e) {
      const srcProperties = ["regopackage", "url", "legacyScanRenderer"] //body is not cleared
      const v = e.target.value;
      const propName = e.target.attributes["name"].value;

      this.formValues[propName] = v;

      srcProperties.filter(item=>item!=propName).forEach((item)=>{
          this.formValues[item]=undefined
      })
    },
    uniqueName(label, value) {
      if (!value) {
        return `${label} is required`;
      }
      const found = this.$store.state.templates.all.filter(
        (item) => item.name === value
      );

      if (found.length > 0 && found[0].name != this.name) {
        return `${value} is not unique`;
      }
      return false;
    },
  },
  mounted() {
    this.name = this.$route.params.name;
  },
};
</script>
