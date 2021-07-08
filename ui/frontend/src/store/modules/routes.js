import api from "../../api"

function updateRoutes(context, routes) {
    api.saveConfig(api.toApiPayload(context, {routes})).then(
        context.commit("set", routes)
    ).catch((error) => {
        context.commit("error/set", error.response.data, {root:true})
    })
}

export default {
    namespaced: true,
    state: {all: []},
    actions: {
        update(context, payload) {
            const routes = context.state.all;
            const { value, name } = payload

            for (let i = 0; i < routes.length; i++) {
                if (routes[i].name == name) {
                    routes.splice(i, 1, value)
                }
            }
            updateRoutes(context, routes)
        },
        remove(context, name) {
            const routes = context.state.all.filter(item => item.name != name)

            updateRoutes(context, routes)
        },
        add(context, settings) {
            const routes = context.state.all
            routes.push(settings)

            updateRoutes(context, routes)
        },


    },
    mutations: {
        set(state, routes) {
            state.all = [...routes]
        },

    }
}