import api from "../../api"
export default {
    namespaced: true,
    state: {
        all: {}
    },
    actions: {
        load(context) {
            api.getStats().then((response) => {
                context.commit("set", response.data)
            }).catch((error) => {
                context.commit("error/set", error.response.data, {root: true})
            })
        },

    },
    mutations: {
        set(state, payload) {
            state.all = { ...payload }
        },

    }
}