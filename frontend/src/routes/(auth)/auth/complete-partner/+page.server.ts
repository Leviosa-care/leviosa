import type { Actions, RequestEvent } from "./$types"
import { fail, redirect } from "@sveltejs/kit"

import { API_URL } from "$env/static/private";

import { setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"
import { parseDate } from "@internationalized/date";

import { pad } from "$lib/utils/pad"
import { mapGenderToBackend, formatPhoneToE164 } from "$lib/utils/auth-helpers";


const schema = type({
    password: "8 < string < 64",
    confirm: "string",
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
    // Partner-specific fields
    bio: "string",
    experience: "string",
    category_ids: "string[]",
    product_ids: "string[]",
})

const defaults = {
    password: "",
    confirm: "",
    firstname: "",
    lastname: "",
    gender: '',
    birthdate: "",
    telephone: "",
    address1: "",
    address2: "",
    city: "",
    postalCode: "",
    bio: "",
    experience: "",
    category_ids: [],
    product_ids: [],
} as typeof schema.infer

export const load = async ({ url }: RequestEvent) => {
    // Get all data from URL search params (passed from address page)
    const firstname = url.searchParams.get("firstname") || defaults.firstname;
    const lastname = url.searchParams.get("lastname") || defaults.lastname;
    const gender = url.searchParams.get("gender") || defaults.gender;
    const birthdate = url.searchParams.get("birthdate") || defaults.birthdate;
    const telephone = url.searchParams.get("telephone") || defaults.telephone;
    const address1 = url.searchParams.get("address1") || defaults.address1;
    const address2 = url.searchParams.get("address2") || defaults.address2;
    const city = url.searchParams.get("city") || defaults.city;
    const postalCode = url.searchParams.get("postalCode") || defaults.postalCode;

    // Partner-specific data
    const bio = url.searchParams.get("bio") || defaults.bio;
    const experience = url.searchParams.get("experience") || defaults.experience;
    const categoryIdsParam = url.searchParams.get("category_ids");
    const productIdsParam = url.searchParams.get("product_ids");

    const category_ids = categoryIdsParam ? categoryIdsParam.split(',') : defaults.category_ids;
    const product_ids = productIdsParam ? productIdsParam.split(',') : defaults.product_ids;

    const form = await superValidate(
        {
            ...defaults,
            firstname,
            lastname,
            gender,
            birthdate,
            telephone,
            address1,
            address2,
            city,
            postalCode,
            bio,
            experience,
            category_ids,
            product_ids,
        },
        arktype(schema, { defaults })
    );
    return { form };
}


export const actions = {
    default: async ({ request, fetch }: RequestEvent) => {
        const form = await superValidate(request, arktype(schema, { defaults }))

        // Validate form
        if (!form.valid) {
            if (form.errors.password) {
                setError(form, "password", "Le mot de passe doit contenir au moins 8 caractères.");
            }
            if (form.errors.confirm) {
                setError(form, "confirm", "La confirmation du mot de passe est requise.");
            }
            if (form.errors.bio) {
                setError(form, "bio", "La biographie ne doit pas dépasser 1000 caractères.");
            }
            if (form.errors.experience) {
                setError(form, "experience", "L'expérience ne doit pas dépasser 2000 caractères.");
            }
            return fail(400, { form })
        }

        // Verify passwords match
        if (form.data.password !== form.data.confirm) {
            return setError(form, "confirm", "Les mots de passe ne correspondent pas.");
        }

        // Validate bio length (max 1000 chars)
        if (form.data.bio.length > 1000) {
            return setError(form, "bio", "La biographie ne doit pas dépasser 1000 caractères.");
        }

        // Validate experience length (max 2000 chars)
        if (form.data.experience.length > 2000) {
            return setError(form, "experience", "L'expérience ne doit pas dépasser 2000 caractères.");
        }

        // Format date to ISO 8601
        const date = parseDate(form.data.birthdate);
        const birth_date = `${date.year}-${pad(date.month)}-${pad(date.day)}T00:00:00Z`;

        // Map frontend gender to backend format
        const gender = mapGenderToBackend(form.data.gender);

        // Format phone to E.164
        const telephone = formatPhoneToE164(form.data.telephone);

        // Prepare request body
        const requestBody: Record<string, unknown> = {
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
        };

        // Add partner-specific fields only if they have values
        if (form.data.bio) {
            requestBody.bio = form.data.bio;
        }
        if (form.data.experience) {
            requestBody.experience = form.data.experience;
        }
        if (form.data.category_ids.length > 0) {
            requestBody.category_ids = form.data.category_ids;
        }
        if (form.data.product_ids.length > 0) {
            requestBody.product_ids = form.data.product_ids;
        }

        // Call /auth/complete/partner endpoint
        const res = await fetch(`${API_URL}/auth/complete/partner`, {
            method: "POST",
            headers: { 'Content-Type': "application/json" },
            body: JSON.stringify(requestBody)
        });

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Erreur de validation. Veuillez vérifier vos informations.");
                case 401:
                    return setError(form, "Non authentifié. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Vous avez déjà complété votre inscription.");
                case 404:
                    return setError(form, "Les catégories ou produits sélectionnés sont invalides.");
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

        // Success (200 OK) - redirect to app
        redirect(302, "/");
    }
} satisfies Actions;
