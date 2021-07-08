export default {
    namespaced: true,
    state: {
        message: undefined
    },
    actions: {

    },
    mutations: {
        set(state, error) {
            state.message = error
        },
        clear(state) {
            state.message = undefined
        },
    }
}