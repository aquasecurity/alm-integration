import api from "../../api"
import { router } from './../../main'

export default {
    namespaced: true,
    state: {
        authenticated: false
    },
    actions: {
        login(context, payload) {
            const { username, password } = payload || {}
            return new Promise((resolve, reject) => {
                api.login(username, password).then(() => {
                    context.commit("update", { authenticated: true })
                    context.commit("error/clear", undefined, {root: true})
                    resolve()
                }).catch(error => {
                    if (username && password) {
                        const errorMsg = error.response.status === 401 ? "Invalid credentials" : error.response.data;
                        context.commit("error/set", errorMsg, {root: true})
                        reject(errorMsg)
                    } else {
                        reject() //just checking
                    }
                })
            })
        },
        logout(context) {
            api.logout().then(() => {
                context.commit("update", { authenticated: false })
                router.push({ name: "login" });
            }).catch(error => {
                context.commit("error/set", error.response.data, {root: true})
            })
        },


    },
    mutations: {
        update(state, info) {
            state.userInfo = { ...info }
        },

    }
}