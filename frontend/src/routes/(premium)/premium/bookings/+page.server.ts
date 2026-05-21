import type { PageServerLoad, RequestEvent } from "./$types"
import { redirect } from "@sveltejs/kit"
import { env } from "$env/dynamic/private"

import type { BookingDTO } from "$lib/types/bookings"
import { handleLoginRedirect } from '$lib/utils/redirect';
import { type Role, ROLES } from "$lib/types/role"

export const load: PageServerLoad = async ({ locals, url, cookies }: RequestEvent) => {
    if (!locals.user) {
        throw redirect(302, handleLoginRedirect(url))
    }
    const clientId = locals.user.id
    const role = locals.user.role
    const sessionCookie = cookies.get('session')

    const getBookings = async () => {
        try {
            const res = await fetch(`${env.API_URL}/clients/${clientId}/bookings`, {
                method: "GET",
                headers: {
                    'Authorization': `Bearer ${sessionCookie}`,
                }
            })

            if (res.status === 401) {
                throw redirect(302, handleLoginRedirect(url))
            }

            if (!res.ok) {
                throw new Error(`Failed to fetch bookings: ${res.statusText}`)
            }

            return await res.json() as Promise<BookingDTO[]>
        } catch (error) {
            if (error && typeof error === 'object' && 'status' in error && error.status === 302) {
                throw error
            }
            console.error('Error fetching bookings:', error)
            return []
        }
    }

    const premiumRoles = [ROLES.premium, ROLES.partner, ROLES.admin] as Role[]
    if (premiumRoles.includes(role)) {
        return {
            bookings: getBookings(),
            events: Promise.resolve([]),
        }
    }

    return { bookings: getBookings() }
}
