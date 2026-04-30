import type { PageServerLoad, Actions, RequestEvent } from "./$types"

import { superValidate, setError } from "sveltekit-superforms"
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"
import { redirect } from "@sveltejs/kit";
import { isAdminDomain, isStaffDomain, getCookieDomain } from "$lib/server/hostname";

const schema = type({
    firstname: "string > 1",
    lastname: "string > 1",
    gender: "'' | 'man' | 'woman' | 'non_binary' | 'prefer_not_to_say' | 'custom' | 'not precised'",
    birthdate: "string",
    telephone: "string == 10"
})

const defaults = {
    firstname: "",
    lastname: "",
    gender: '',
    birthdate: "",
    telephone: ""
} as typeof schema.infer

export const load: PageServerLoad = async ({ cookies, url }: RequestEvent) => {
    // Try to get form data from registration cookie (if user is coming back from a later step)
    const regGeneralCookie = cookies.get("reg_general");
    let initialData = { ...defaults };

    if (regGeneralCookie) {
        try {
            initialData = JSON.parse(regGeneralCookie) as typeof defaults;
        } catch {
            // Invalid cookie, use defaults
        }
    }

    const form = await superValidate(initialData, arktype(schema, { defaults }));
    return { form };
}

export const actions: Actions = {
    default: async ({ request, cookies, url }: RequestEvent) => {
        const form = await superValidate(request, arktype(schema, { defaults }))

        if (!form.valid) {
            if (form.errors.firstname) {
                setError(form, "firstname", "Le prénom est requis (minimum 2 caractères).");
            }
            if (form.errors.lastname) {
                setError(form, "lastname", "Le nom de famille est requis (minimum 2 caractères).");
            }
            if (form.errors.gender) {
                setError(form, "gender", "Veuillez sélectionner un genre.");
            }
            if (form.errors.birthdate) {
                setError(form, "birthdate", "Veuillez renseigner votre date de naissance.");
            }
            if (form.errors.telephone) {
                setError(form, "telephone", "Le numéro de téléphone doit contenir exactement 10 chiffres.");
            }
            return { form };
        }

        // Store form data in HTTP-only cookie (30 min expiry)
        const hostname = url.hostname;
        const cookieDomain = getCookieDomain(hostname);
        cookies.set("reg_general", JSON.stringify(form.data), {
            path: "/auth",
            maxAge: 1800,
            httpOnly: true,
            secure: !hostname.startsWith('localhost'),
            sameSite: "strict",
            ...(cookieDomain && { domain: cookieDomain })
        });

        redirect(302, "/auth/address");
    }
}
