import { setContext, getContext } from "svelte"

export function setUserContext(user: App.User | undefined) {
    setContext<App.User | undefined>("user", user)
}


export function getUserContext(): App.User | undefined {
    return getContext<App.User | undefined>("user")
}
