<template>
  <div>
    <div class="row justify-content-end pb-3 pr-3">
      <button type="submit" class="btn btn-primary mr-2" @click="doSubmit">
        Submit
      </button>
      <button
        type="button"
        @click="doTest"
        class="btn btn-outline-primary mr-2"
      >
        Test config
      </button>
      <b-spinner
        v-if="isTestingInProgress"
        variant="primary"
        label="Spinning"
        class="mr-2"
      ></b-spinner>
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
          <div class="row">
            <div class="col">
              <div class="form-group form-input">
                <label for="outputType">Type</label>
                <select
                  class="form-select form-control"
                  :value="formValues.type"
                  id="outputType"
                  name="type"
                  @input="updateOutputType"
                >
                  <option value="email">Email</option>
                  <option value="jira">Jira</option>
                  <option value="slack">Slack</option>
                  <option value="teams">Teams</option>
                  <option value="webhook">Webhook</option>
                  <option value="splunk">Splunk</option>
                  <option value="serviceNow">ServiceNow</option>
                </select>
                <small id="aHelp" class="form-text text-muted"
                  >The output type</small
                >
              </div>
            </div>
            <div class="col">
              <PropertyField
                id="name"
                label="Name"
                :value="formValues.name"
                :errorMsg="errors['name']"
                :inputHandler="updateField"
                :validator="v(uniqueName)"
              />
            </div>
          </div>
          <!-- enable is not required here-->
          <CheckboxPropertyField
            id="enable"
            label="Enable output"
            :value="formValues.enable"
            :inputHandler="updateField"
          />
          <div class="row">
            <div class="col">
              <PropertyField
                id="user"
                label="User"
                :value="formValues.user"
                name="user"
                :show="showCredentials"
                :inputHandler="updateField"
                :errorMsg="errors['user']"
                :validator="v(required)"
              />
            </div>
            <div class="col">
              <PropertyField
                id="password"
                label="Password"
                inputType="password"
                :errorMsg="errors['password']"
                :value="formValues.password"
                name="password"
                :show="showCredentials"
                :inputHandler="updateField"
                :validator="v(required)"
              />
            </div>
          </div>
          <PropertyField
            id="url"
            label="Url"
            :value="formValues.url"
            :errorMsg="errors['url']"
            name="url"
            :description="getUrlDescription"
            :show="showUrl"
            :inputHandler="updateField"
            :validator="v([url, required])"
          />

          <!-- email custom properties start -->
          <div class="row">
            <div class="col">
              <PropertyField
                id="host"
                label="Host"
                :errorMsg="errors['host']"
                :value="formValues.host"
                description="Mandatory: SMTP host name (e.g. smtp.gmail.com)"
                :show="isEmail"
                :inputHandler="updateField"
                :validator="v(required)"
              />
            </div>
            <div class="col">
              <PropertyField
                id="port"
                label="Port"
                inputType="number"
                :errorMsg="errors['port']"
                :value="formValues.port"
                description="Mandatory: SMTP server port (e.g. 587)"
                :show="isEmail"
                :inputHandler="updateField"
                :validator="v(required)"
              />
            </div>
          </div>

          <div class="row">
            <div class="col">
              <PropertyField
                id="sender"
                label="Sender"
                :value="formValues.sender"
                description="The email address to use as a sender"
                :errorMsg="errors['sender']"
                :show="isEmail"
                :inputHandler="updateField"
                :validator="v(email)"
              />
            </div>
            <div class="col">
              <PropertyField
                id="recipients"
                label="Recipients"
                :value="formValues.recipients | toString"
                description="Mandatory: comma separated list of recipients"
                :errorMsg="errors['recipients']"
                :show="isEmail"
                :inputHandler="updateCollectionField"
                :validator="v(recipients)"
              />
            </div>
          </div>

          <CheckboxPropertyField
            id="use-mx"
            label="Use MX"
            :value="formValues['use-mx']"
            :show="isEmail"
            :inputHandler="updateField"
          />
          <!-- email custom properties end -->

          <!-- jira custom properties start -->
          <div class="row">
            <div class="col">
              <PropertyField
                id="projectKey"
                label="Project Key"
                name="project-key"
                :value="formValues['project-key']"
                :errorMsg="errors['project-key']"
                :description="'Mandatory. Specify the JIRA project key'"
                :show="isJira"
                :inputHandler="updateField"
                :validator="v(required)"
              />
            </div>
            <div class="col">
              <PropertyField
                id="board"
                label="Board"
                :value="formValues.board"
                description="Optional. Specify the Jira board name to open tickets on"
                :show="isJira"
                :inputHandler="updateField"
              />
            </div>
          </div>

          <CheckboxPropertyField
            id="tlsVerify"
            label="TLS verify"
            :value="formValues['tls-verify']"
            name="tls-verify"
            :show="isJira"
            :inputHandler="updateField"
          />


          <div class="row">
            <div class="col">
              <PropertyField
                id="fixVersions"
                name="fix-versions"
                label="Fix Versions"
                :value="formValues['fix-versions'] | toString"
                description="Optional, specify comma separated list of Fix versions to add to Ticket"
                :show="isJira"
                :inputHandler="updateCollectionField"
              />
            </div>
            <div class="col">
              <PropertyField
                id="affectsVersions"
                name="affects-versions"
                label="Affects Versions"
                :value="formValues['affects-versions'] | toString"
                description="Optional, specify comma separated list of Affects versions to add to Ticket"
                :show="isJira"
                :inputHandler="updateCollectionField"
              />
            </div>
          </div>

          <PropertyField
            id="labels"
            label="Labels"
            :value="formValues.labels | toString"
            description="Optional, specify comma separated list of labels to add to Ticket"
            :show="isJira"
            :inputHandler="updateCollectionField"
          />

          <div class="row">
            <div class="col">
              <PropertyField
                id="issuetype"
                label="Issue Type"
                :value="formValues.issuetype"
                description="Optional. Specify the issue type to open (Bug, Task, etc.). Default is Task"
                :show="isJira"
                :inputHandler="updateField"
              />
            </div>
            <div class="col">
              <PropertyField
                id="priority"
                label="Priority"
                :value="formValues.priority"
                description="Optional. Specify the issues severity. Default is High"
                :show="isJira"
                :inputHandler="updateField"
              />
            </div>
          </div>
          <div class="row">
            <div class="col">
              <PropertyField
                id="assignee"
                label="Assignee"
                :value="formValues.assignee | toString"
                :description="jiraAssigneeDescription"
                :show="isJira"
                :inputHandler="updateCollectionField"
              />
            </div>
            <div class="col">
              <PropertyField
                id="sprint"
                label="Sprint"
                :value="formValues.sprint"
                description="Optional Sprint name, e.g., '3.5 Sprint 8'"
                :show="isJira"
                :inputHandler="updateField"
              />
            </div>
          </div>
          <div class="row mb-3" v-show="isJira">
            <div class="col-sm-7">
              <h5>Jira custom fields</h5>
            </div>
            <div class="col-sm-4 ">
                <b-form-input
                  id="cf-adder"
                  v-model="addedControl"
                  type="text"
                  placeholder="Add new custom field"
                ></b-form-input>
            </div>
            <div class="col-sm-1 text-right">
               <b-button variant="primary" @click="addCf" title="Add field">Add</b-button>
            </div>
          </div>

          <b-row class="form-group" v-for="(cfValue, cfName) in unknowns" v-bind:key="cfName" :show="isJira">
            <b-col sm="2">
              <label :for="'cf-' + cfName">{{cfName}}</label>
            </b-col>
            <b-col sm="9">
              <b-form-input :id="'cf-' + cfName" :name="cfName" v-model="unknowns[cfName]"></b-form-input>
            </b-col>
            <b-col sm="1" class="text-right">
              <b-button variant="primary" title="Remove field" @click="removeCf(cfName)" ><span aria-hidden="true">&times;</span></b-button>
            </b-col>
          </b-row>


          <!-- jira custom properties end -->
          <!-- serviceNow custom properties start -->
          <div class="row">
            <div class="col">
              <PropertyField
                id="instance"
                label="Instance"
                :value="formValues.instance"
                description="Mandatory. Name of ServiceNow  or Instance"
                :errorMsg="errors['instance']"
                :show="isServiceNow"
                :inputHandler="updateField"
                :validator="v(required)"
              />
            </div>
            <div class="col">
              <PropertyField
                id="board"
                label="Board"
                :value="formValues.board"
                description="Specify the ServiceNow board name to open tickets on. Default is incident"
                :show="isServiceNow"
                :inputHandler="updateField"
              />
            </div>
          </div>
          <!-- serviceNow custom properties end -->
          <!-- splunk custom properties start -->
          <div class="row">
            <div class="col">
              <PropertyField
                id="token"
                label="Token"
                :value="formValues.token"
                :errorMsg="errors['token']"
                description="Mandatory. a HTTP Event Collector Token"
                :show="isSplunk"
                :inputHandler="updateField"
                :validator="v(required)"
              />
            </div>
            <div class="col">
              <PropertyField
                id="sizeLimit"
                name="size-limit"
                label="Size Limit"
                :value="formValues['size-limit']"
                inputType="number"
                description="Optional. Maximum scan length, in bytes. Default: 10000"
                :show="isSplunk"
                :inputHandler="updateField"
              />
            </div>
          </div>
          <!-- splunk custom properties end -->
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
import CheckboxPropertyField from "./CheckboxPropertyField.vue";

const urlDescriptionByType = {
  splunk: "Mandatory. Url of a Splunk server",
  webhook: "Webhook's url",
  teams: "Webhook's url",
  jira: 'Mandatory. E.g "https://johndoe.atlassian.net"',
  slack: "",
};
const typesWithCredentials = ["serviceNow", "jira", "email"]; //TODO add description strings

export default {
  data() {
    return {
      name: "",
      addedControl: "",
      unknowns: [],
      isTestingInProgress: false,
      fields: {}, //required for mixins
      errors: {}, //required for mixins
      outputType: "", //stored separately to track dependencies
      jiraAssigneeDescription:
        'Optional: comma separated list of users (emails) that will be assigned to ticket, e.g., ["john@yahoo.com"]. To assign a ticket to the Application Owner email address (as defined in Aqua Application Scope, owner email field), specify ["<%application_scope_owner%>"] as the assignee value',
    };
  },
  mixins: [FormFieldMixin, ValidationMixin],
  components: {
    PropertyField,
    CheckboxPropertyField,
  },
  computed: {
    ...mapState({
      formValues(state) { //required for mixins
        const found = state.outputs.all.filter(
          (item) => item.name === this.name
        );

        const result = found.length ? { ...found[0] } : { type: "email" };

        this.outputType = result.type;
        this.unknowns = {...result.unknowns}

        return result;
      },
    }),
    showUrl() {
      return urlDescriptionByType[this.outputType] !== undefined;
    },
    getUrlDescription() {
      return urlDescriptionByType[this.outputType];
    },
    isServiceNow() {
      return this.outputType === "serviceNow";
    },
    isSplunk() {
      return this.outputType === "splunk";
    },
    isEmail() {
      return this.outputType === "email";
    },
    isJira() {
      return this.outputType === "jira";
    },
    showCredentials() {
      return typesWithCredentials.indexOf(this.outputType) >= 0;
    },
  },
  filters: {
    toString(col) {
      return col ? col.join(", ") : undefined;
    },
  },
  methods: {
    doTest() {
      if (!this.isFormValid()) {
        return;
      }

      this.isTestingInProgress = true;

      this.$store
        .dispatch("outputs/test", this.formValues)
        .then(() => {
          this.$bvToast.toast("Output is configured correctly", {
            title: "Success",
            variant: "success",
            autoHideDelay: 5000,
          });
          this.isTestingInProgress = false;
        })
        .catch((error) => {
          this.$bvToast.toast(error, {
            title: "Connection error",
            variant: "danger",
            autoHideDelay: 15000,
          });
          this.isTestingInProgress = false;
        });
    },
    doSubmit() {
      if (!this.isFormValid()) {
        return;
      }
      //apply unknowns
      if (Object.keys(this.unknowns).length > 0) {
        this.formValues.unknowns={...this.unknowns}
      }
      if (this.name) {
        this.$store.dispatch("outputs/update", {
          value: this.formValues,
          name: this.name,
        });
      } else {
        this.$store.dispatch("outputs/add", this.formValues);
      }
      this.$router.push({ name: "home" });
    },
    updateOutputType(e) {
      this.outputType = e.target.value;
      this.updateField(e);
    },
    doRemove() {
      this.$store.dispatch("outputs/remove", this.name);
      this.$router.push({ name: "home" });
    },
    uniqueName(label, value) {
      if  (!value) {
        return `${label} is required`
      }
      const found = this.$store.state.outputs.all.filter(
        (item) => item.name === value
      );

      if (found.length > 0 && found[0].name != this.name) {
        return `${value} is not unique`
      }
      return false
    },
    addCf() {
      const unknowns = {...this.unknowns}
      unknowns[this.addedControl] = ""
      this.addedControl = ""
      this.unknowns = {...unknowns}
    },
    removeCf(propertyName) {
      const unknowns = {...this.unknowns}
      delete unknowns[propertyName]
      this.unknowns = {...unknowns}
    }

  },
  mounted() {
    this.name = this.$route.params.name;
  },
};
</script>