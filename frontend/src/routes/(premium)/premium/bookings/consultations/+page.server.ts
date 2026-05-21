import type { PageServerLoad, RequestEvent } from "./$types"
import { redirect } from "@sveltejs/kit"
import { handleLoginRedirect } from '$lib/utils/redirect';

export const load: PageServerLoad = async ({ parent, url }: RequestEvent & { parent: () => Promise<any> }) => {
    const parentData = await parent();

    if (!parentData.bookings) {
        throw redirect(302, handleLoginRedirect(url))
    }

    return {
        bookings: parentData.bookings,
    }
}
