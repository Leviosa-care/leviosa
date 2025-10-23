// import type { LayoutServerLoad, RequestEvent } from './$types';
import type { RequestEvent } from "@sveltejs/kit"
import { mockUser } from "$lib/data/user";


export const load = (event: RequestEvent) => {
    // TODO: get the redirectFrom, if it is from some valid url in the sign up I can use
    const redirectFrom = event.url.searchParams.get("redirectFrom")
    if (redirectFrom) {
        // TODO: create a visitor user user that has the role visitor or unknown I do nt remember
        event.locals.user = mockUser
        return { user: event.locals.user }
    }
    return { user: event.locals.user }
}

// https://leviosa.care/portal
