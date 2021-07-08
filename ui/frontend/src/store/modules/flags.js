
export default {
    namespaced: true,
    state: {all: {
        loaded: false
    }},
    mutations: {
        set(state, flags) {
            state.all = {...state.all, ...flags}
        }
    }
}