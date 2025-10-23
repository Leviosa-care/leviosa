import type { PageServerLoad, RequestEvent } from "./$types"
import { redirect } from "@sveltejs/kit"
import { env } from "$env/dynamic/private"

import type { Event, Consultation } from "$lib/types/bookings"
import { handleLoginRedirect } from '$lib/utils/redirect';
import { type Role, ROLES } from "$lib/types/role";

export const load: PageServerLoad = async ({ locals, url }: RequestEvent) => {
    if (!locals.user) {
        throw redirect(302, handleLoginRedirect(url))
    }
    const role = locals.user.role
    const getConsultations = async () => {
        const res = await fetch(`${env.API_URL}/consultations`, { method: "GET", })
        return await res.json() as Promise<Consultation[]>
    }
    const premiumRoles = [ROLES.premium, ROLES.partner, ROLES.admin] as Role[];
    if (premiumRoles.includes(role)) {
        const getEvents = async () => {
            const res = await fetch(`${env.API_URL}/events`, { method: "GET", })
            return await res.json() as Promise<Event[]>;
        }
        return {
            consultations: getConsultations(),
            events: getEvents(),
        }
    }

    return { consultations: getConsultations() }
}
