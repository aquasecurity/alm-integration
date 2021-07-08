import axios from "axios";
import yaml from "js-yaml";

const transformYaml = (response) => {
    try {
        const json = yaml.load(response)
        return json
    } catch(e) {
        return response //this way text errors are handled
    }
}

export default {
    getConfig: function () {
        return axios.get("/api/config", { transformResponse: transformYaml })
    },
    getStats: function () {
        return axios.get("/api/outputs/stats")
    },
    saveConfig: function (settings) {
        const yamlObj = yaml.dump(settings)
        return axios.post("/api/config", yamlObj)

    },
    test: function (settings) {
        return axios.post("/api/test", settings)

    },
    login: function (username, password) {
        const bodyFormData = new FormData();

        bodyFormData.append('username', username ? username : "");
        bodyFormData.append('password', password ? password : "");

        return axios.post("/api/login", bodyFormData)

    },
    logout: function () {
        return axios.get("/api/logout")
    },
    toApiPayload: function (context, modification) {
        const rootState = context.rootGetters.getAppState
        return {
            outputs: rootState.outputs.all,
            routes: rootState.routes.all,
            templates: rootState.templates.all,
            ...rootState.settings.all,
            ...modification
        }
    }
}