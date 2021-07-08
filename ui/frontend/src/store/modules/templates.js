import api from "../../api"

function updateTemplates(context, templates) {
    api.saveConfig(api.toApiPayload(context, {templates})).then(
        context.commit("set", templates)
    ).catch((error) => {
        context.commit("error/set", error.response.data, {root:true})
    })
}

export default {
    namespaced: true,
    state: {all: []},
    actions: {
        update(context, payload) {
            const templates = context.state.all;
            const { value, name } = payload

            for (let i = 0; i < templates.length; i++) {
                if (templates[i].name == name) {
                    templates.splice(i, 1, value)
                }
            }
            updateTemplates(context, templates)
        },
        remove(context, name) {
            const templates = context.state.all.filter(item => item.name != name)

            updateTemplates(context, templates)
        },
        add(context, settings) {
            const templates = context.state.all
            templates.push(settings)

            updateTemplates(context, templates)
        },


    },
    mutations: {
        set(state, templates) {
            state.all = [...templates]
        },

    }
}