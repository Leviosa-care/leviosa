import { redirect, } from "@sveltejs/kit"
import type { PageServerLoad, RequestEvent } from "./$types"
import { handleLoginRedirect } from "$lib/utils/redirect"

export const load: PageServerLoad = async ({ fetch, locals, url }: RequestEvent) => { // TODO: change the content of that page to get the information about the homepage once authencated
    // const { articleID }= event.params
    // TODO: make that thing better brother in terms of authentication
    if (!locals.user) {
        throw redirect(303, handleLoginRedirect(url))
    }
    // const res = await fetch("", {
    // })
    // const data = res.json()
    // console.log("data:", data)

    // I will have some information to get here brother
    return {
        blogPost: "This is an example of blog post",
    }
}

// TODO: I need to use that in the pages that redirect to login
// TODO: where to store that ?
