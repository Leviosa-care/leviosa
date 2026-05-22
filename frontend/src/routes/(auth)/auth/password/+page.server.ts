import type { Actions, RequestEvent } from "./$types"
import { fail, redirect } from "@sveltejs/kit"

import { env } from "$env/dynamic/private";

import { setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"
import { parseDate } from "@internationalized/date";

import { pad } from "$lib/utils/pad"
import { mapGenderToBackend, formatPhoneToE164 } from "$lib/utils/auth-helpers";
import { getCookieDomain } from "$lib/server/hostname";


const schema = type({
    password: "8 < string < 64",
    confirm: "string",
    email: "string",
    // Data from previous steps
    firstname: "string",
    lastname: "string",
    gender: "'' | 'man' | 'woman' | 'non_binary' | 'prefer_not_to_say' | 'custom' | 'not precised'",
    birthdate: "string",
    telephone: "string",
    address1: "string",
    address2: "string",
    city: "string",
    postalCode: "string == 5",
})

const defaults = {
    password: "",
    confirm: "",
    email: "",
    firstname: "",
    lastname: "",
    gender: '',
    birthdate: "",
    telephone: "",
    address1: "",
    address2: "",
    city: "",
    postalCode: "",
} as typeof schema.infer

export const load = async ({ cookies }: RequestEvent) => {
    // Get all data from registration cookies (set by previous steps)
    const regGeneralCookie = cookies.get("reg_general");
    const regAddressCookie = cookies.get("reg_address");

    let generalData = {
        firstname: defaults.firstname,
        lastname: defaults.lastname,
        gender: defaults.gender,
        birthdate: defaults.birthdate,
        telephone: defaults.telephone,
    };

    let addressData = {
        address1: defaults.address1,
        address2: defaults.address2,
        city: defaults.city,
        postalCode: defaults.postalCode,
    };

    if (regGeneralCookie) {
        try {
            generalData = JSON.parse(regGeneralCookie);
        } catch {
            // Invalid cookie, use defaults
        }
    }

    if (regAddressCookie) {
        try {
            addressData = JSON.parse(regAddressCookie);
        } catch {
            // Invalid cookie, use defaults
        }
    }

    const form = await superValidate(
        {
            ...defaults,
            ...generalData,
            ...addressData,
        },
        arktype(schema, { defaults })
    );
    return { form };
}


export const actions = {
    default: async ({ request, fetch, cookies, url }: RequestEvent) => {
        const form = await superValidate(request, arktype(schema, { defaults }))

        // Validate form
        if (!form.valid) {
            if (form.errors.password) {
                setError(form, "password", "Le mot de passe doit contenir au moins 8 caractères.");
            }
            if (form.errors.confirm) {
                setError(form, "confirm", "La confirmation du mot de passe est requise.");
            }
            return fail(400, { form })
        }

        // Verify passwords match
        if (form.data.password !== form.data.confirm) {
            return setError(form, "confirm", "Les mots de passe ne correspondent pas.");
        }

        // Format date to ISO 8601
        const date = parseDate(form.data.birthdate);
        const birth_date = `${date.year}-${pad(date.month)}-${pad(date.day)}T00:00:00Z`;

        // Map frontend gender to backend format
        const gender = mapGenderToBackend(form.data.gender);

        // Format phone to E.164
        const telephone = formatPhoneToE164(form.data.telephone);

        // Call /auth/complete endpoint
        const res = await fetch(`${env.API_URL}/auth/complete`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Cookie': `leviosa_access_token=${cookies.get("leviosa_access_token") ?? ""}`,
            },
            body: JSON.stringify({
                password: form.data.password,
                first_name: form.data.firstname,
                last_name: form.data.lastname,
                birth_date,
                gender,
                telephone,
                postal_code: form.data.postalCode,
                city: form.data.city,
                address1: form.data.address1,
                address2: form.data.address2,
            })
        });

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Erreur de validation. Veuillez vérifier vos informations.");
                case 401:
                    return setError(form, "Non authentifié. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Vous avez déjà complété votre inscription.");
                case 415:
                    return setError(form, "Type de contenu non supporté. Veuillez réessayer.");
                case 500:
                    return setError(form, "Une erreur serveur est survenue. Veuillez réessayer.");
                case 503:
                    return setError(form, "Le service est temporairement indisponible. Veuillez réessayer dans quelques instants.");
                default:
                    return setError(form, "Une erreur inattendue est survenue. Veuillez réessayer.");
            }
        }

        // Clear registration cookies on success
        const hostname = url.hostname;
        const cookieDomain = getCookieDomain(hostname);
        cookies.delete("reg_general", {
            path: "/auth",
            ...(cookieDomain && { domain: cookieDomain })
        });
        cookies.delete("reg_address", {
            path: "/auth",
            ...(cookieDomain && { domain: cookieDomain })
        });

        // Success (200 OK) - redirect to app
        redirect(302, "/");
    }
} satisfies Actions;
