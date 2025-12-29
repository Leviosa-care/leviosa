
import type { Actions, RequestEvent } from "./$types"
import { fail, redirect } from "@sveltejs/kit"

import { setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"

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


export const load = async ({ url }: RequestEvent) => {
    // Get form data from URL search params (passed from general page)
    const address1 = url.searchParams.get("address1") || defaults.address1;
    const address2 = url.searchParams.get("address2") || defaults.address2;
    const postalCode = url.searchParams.get("postalCode") || defaults.postalCode;
    const city = url.searchParams.get("city") || defaults.city;

    const form = await superValidate(
        { address1, address2, postalCode, city },
        arktype(schema, { defaults })
    );

    return { form }
}

export const actions = {
    default: async ({ request, url }: RequestEvent) => {
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

        // Get data from general page (passed via URL search params)
        const firstname = url.searchParams.get("firstname") || "";
        const lastname = url.searchParams.get("lastname") || "";
        const gender = url.searchParams.get("gender") || "";
        const birthdate = url.searchParams.get("birthdate") || "";
        const telephone = url.searchParams.get("telephone") || "";

        // Combine general + address data and pass to password page
        const params = new URLSearchParams({
            firstname,
            lastname,
            gender,
            birthdate,
            telephone,
            address1: form.data.address1,
            address2: form.data.address2,
            postalCode: form.data.postalCode,
            city: form.data.city,
        });

        redirect(302, `/auth/password?${params.toString()}`);
    },
} satisfies Actions
