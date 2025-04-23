import { setContext, getContext } from "svelte"
// TODO: I should looking for the large icons and then to regroup them for the small but the state should be as extensive as possible
export const NAV_STATES = {
    Home: "Accueil",
    Profile: 'Profil',
    Services: 'Services',

    // TODO: finish the type here based on the routes that we have
    Messages: 'messages',
    Reservations: 'reservations',
    Conversations: 'conversations',
    NotesDeSeances: 'notes de seances',
} as const

export type NavState = typeof NAV_STATES[keyof typeof NAV_STATES]
export const NAV_STATES_ARRAY = Object.values(NAV_STATES)


export class NavStore {
    state = $state<NavState>(NAV_STATES.Home)
    constructor() { }
}

const NAV_KEY = Symbol("navigation")

export const setNavigationContext = () => {
    return setContext(NAV_KEY, new NavStore())
}

export const getNavigationContext = () => {
    return getContext<ReturnType<typeof setNavigationContext>>(NAV_KEY)
}
