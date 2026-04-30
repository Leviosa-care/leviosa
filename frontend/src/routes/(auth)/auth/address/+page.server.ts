
import type { Actions, RequestEvent } from "./$types"
import { fail, redirect } from "@sveltejs/kit"

import { setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"
import { getCookieDomain } from "$lib/server/hostname";

const schema = type({
    address1: "string",
    address2: "string",
    postalCode: "string == 5",
    city: "string",
})
const defaults = {
    address1: "",
    address2: "",
    postalCode: "",
    city: "",
}


export const load = async ({ cookies, url }: RequestEvent) => {
    // Try to get address data from registration cookie (if user is coming back from a later step)
    const regAddressCookie = cookies.get("reg_address");
    let initialData = { ...defaults };

    if (regAddressCookie) {
        try {
            initialData = JSON.parse(regAddressCookie) as typeof defaults;
        } catch {
            // Invalid cookie, use defaults
        }
    }

    const form = await superValidate(initialData, arktype(schema, { defaults }));
    return { form }
}

export const actions = {
    default: async ({ request, cookies, url }: RequestEvent) => {
        const form = await superValidate(request, arktype(schema, { defaults: defaults }))

        if (!form.valid) {
            if (form.errors.address1) {
                setError(form, "address1", "L'adresse est requise.");
            }
            if (form.errors.address2) {
                setError(form, "address2", "Complément d'adresse invalide.");
            }
            if (form.errors.city) {
                setError(form, "city", "La ville est requise.");
            }
            if (form.errors.postalCode) {
                setError(form, "postalCode", "Le code postal doit contenir exactement 5 chiffres.");
            }
            return fail(400, { form })
        }

        // Store address form data in HTTP-only cookie (30 min expiry)
        const hostname = url.hostname;
        const cookieDomain = getCookieDomain(hostname);
        cookies.set("reg_address", JSON.stringify(form.data), {
            path: "/auth",
            maxAge: 1800,
            httpOnly: true,
            secure: !hostname.startsWith('localhost'),
            sameSite: "strict",
            ...(cookieDomain && { domain: cookieDomain })
        });

        redirect(302, "/auth/password");
    },
} satisfies Actions
