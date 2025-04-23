import { setContext, getContext } from "svelte";

const USER_KEY = Symbol("user")

export class UserState {
    readonly user: App.User = $state({
        email: "",
        firstname: "",
        lastname: "",
        picture: "",
        role: "",
        birthdate: "",
        gender: "",
        telephone: "",
        postalCode: "",
        city: "",
        address1: "",
        address2: "",
    })

    constructor(user: App.User) {
        Object.assign(this.user, user);
        // Object.freeze(this.user);
    }
}

export function setUserContext(user: App.User) {
    return setContext(USER_KEY, new UserState(user))
}

export function getUserContext() {
    return getContext<ReturnType<typeof setUserContext>>(USER_KEY)
}
