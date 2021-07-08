<template>
  <div>
    <div class="row justify-content-end pb-3 pr-3">
      <button type="submit" class="btn btn-primary mr-2" @click="doSubmit">
        Submit
      </button>
    </div>
    <div class="card">
      <form @submit.prevent="doSubmit">
        <div class="card-body">
          <PropertyField
            :id="'tenant'"
            :label="'Name'"
            :value="formValues.name"
            :name="'name'"
            :description="'Tenant name'"
            :inputHandler="updateField"
          />
          <PropertyField
            :id="'aquaServer'"
            :label="'Aqua Server'"
            :value="formValues.AquaServer"
            :name="'AquaServer'"
            :description="'url of Aqua Server for links. E.g. https://myserver.aquasec.com'"
            :inputHandler="updateField"
          />
          <PropertyField
            :id="'maxDbSize'"
            :label="'Max Db size'"
            :inputType="'number'"
            :value="formValues.Max_DB_Size"
            :name="'Max_DB_Size'"
            :description="'Max size of DB. MB. if empty then unlimited'"
            :inputHandler="updateField"
          />
          <PropertyField
            :id="'deleteOldData'"
            :label="'Delete old data'"
            :inputType="'number'"
            :value="formValues.Delete_Old_Data"
            :name="'Delete_Old_Data'"
            :description="'delete data older than N day(s).  If empty then we do not delete.'"
            :inputHandler="updateField"
          />
          <PropertyField
            id="dbVerifyInterval"
            label="DB verify interval"
            inputType="number"
            :value="formValues.DbVerifyInterval"
            name="DbVerifyInterval"
            description="hours. an Interval between tests of DB. Default: 1 hour"
            :inputHandler="updateField"
          />
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

export default {
  data() {
    return {
      fields: {},
      errors: {},
    };
  },
  mixins: [FormFieldMixin, ValidationMixin],
  components: {
    PropertyField,
  },
  computed: {
    ...mapState({
      formValues(state) {
        return state.settings.all;
      },
    }),
  },
  methods: {
    doSubmit() {
        if (!this.isFormValid()) {
            return;
        }
        this.$store.dispatch("settings/update", this.formValues);
    },
  },
};
</script>
