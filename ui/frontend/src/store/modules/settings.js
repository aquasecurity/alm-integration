import api from "../../api"

export default {
    namespaced: true,
    state: {all: {}},
    actions: {
        update(context, payload) {
            api.saveConfig(api.toApiPayload(context, payload)).then( //entire config is saved
                context.commit("set", payload)
            ).catch((error) => {
                context.commit("error/set", error.response.data, {root: true})
            })

        },
    },
    mutations: {
        set(state, settings) {
            state.all = {...settings}
        },

    }
}