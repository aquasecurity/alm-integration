import Vue from 'vue'
import Vuex from 'vuex'
import api from './../api'

import error from './modules/error.js'
import account from './modules/account.js'
import outputs from './modules/outputs.js'
import stats from './modules/stats.js'
import routes from './modules/routes.js'
import settings from './modules/settings.js'
import flags from './modules/flags.js'
import templates from './modules/templates.js'

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        error,
        outputs,
        account,
        stats,
        routes,
        settings,
        flags,
        templates
    },
    getters: {
        getAppState(state) {
            return state
        }
    },
    actions: {
        load(context) {
            api.getConfig().then((response) => {
                const data = response.data
                const settings = {
                    name: data.name,
                    AquaServer: data.AquaServer,
                    Delete_Old_Data: data.Delete_Old_Data,
                    DbVerifyInterval: data.DbVerifyInterval,
                    Max_DB_Size: data.Max_DB_Size
                }
                data.outputs && context.commit("outputs/set", data.outputs)
                data.routes && context.commit("routes/set", data.routes)
                data.templates && context.commit("templates/set", data.templates)
                context.commit("settings/set", settings)
                context.commit("flags/set", { loaded: true })
            }).catch((error) => {
                if (error.response) {
                    context.commit("error/set", error.response.data)
                } else {
                    console.error(error)
                }
            })
        }
    }

})