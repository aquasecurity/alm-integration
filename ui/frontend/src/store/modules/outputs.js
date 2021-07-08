import api from "../../api"
function updateOutputs(context, outputs) {
    api.saveConfig(api.toApiPayload(context, {outputs})).then( //entire config is saved
        context.commit("set", outputs)
    ).catch((error) => {
        context.commit("error/set", error.response.data, {root:true})
    })
}
export default {
    namespaced: true,
    state: {all: []},
    actions: {
        test(context, settings) {
            return new Promise((resolve, reject) => {
                api.test(settings).then(() => {
                    context.commit("error/clear", undefined, {root:true})
                    resolve();
                }).catch(error => {
                    if (error.response) {
                        context.commit("error/set", error.response.data, {root:true})
                        reject(error.response.data);
                    } else {
                        console.error(error)
                        reject(error);
                    }
                })
            })
        },

        update(context, payload) {
            const outputs = context.state.all;
            const { value, name } = payload

            for (let i = 0; i < outputs.length; i++) {
                if (outputs[i].name == name) {
                    outputs.splice(i, 1, value)
                }
            }
            updateOutputs(context, outputs)
        },
        remove(context, name) {
            const outputs = context.state.all.filter(item => item.name != name)

            updateOutputs(context, outputs)
        },
        add(context, settings) {
            const outputs = context.state.all
            outputs.push(settings)

            updateOutputs(context, outputs)
        },


    },
    mutations: {
        set(state, outputs) {
            state.all = [...outputs]
        },

    }
}