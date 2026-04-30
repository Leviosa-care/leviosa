import { type } from 'arktype';
import { superValidate, setError } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import type { RequestEvent } from './$types';
import { env } from "$env/dynamic/private"
import { redirect } from '@sveltejs/kit';
import { forwardAuthCookies } from '$lib/utils/auth-helpers';

const schema = type({
    email: "string.email",
    otp0: "/^\\d$/",
    otp1: "/^\\d$/",
    otp2: "/^\\d$/",
    otp3: "/^\\d$/",
    otp4: "/^\\d$/",
    otp5: "/^\\d$/",
});

const defaults = {
    email: '',
    otp0: '',
    otp1: '',
    otp2: '',
    otp3: '',
    otp4: '',
    otp5: '',
};

export const load = async ({ url }: RequestEvent) => {
    // Get email from URL search params (passed from register action)
    const email = url.searchParams.get("email") || '';
    const form = await superValidate({ ...defaults, email }, arktype(schema, { defaults }));
    return { form };
};

export const actions = {
    default: async ({ request, fetch, cookies }: RequestEvent) => {
        const form = await superValidate(request, arktype(schema, { defaults }));

        if (!form.valid) {
            if (form.errors.email) {
                return setError(form, "email", "L'adresse email n'est pas valide.");
            }
            return setError(form, "Veuillez entrer des valeurs numériques à un seul chiffre (de 0 à 9).");
        }

        const code = [
            form.data.otp0,
            form.data.otp1,
            form.data.otp2,
            form.data.otp3,
            form.data.otp4,
            form.data.otp5,
        ].join("");

        const res = await fetch(`${env.API_URL}/auth/otp`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: form.data.email,
                code: code
            })
        });

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Format d'email ou de code OTP invalide.");
                case 401:
                    return setError(form, "Code OTP invalide ou expiré. Veuillez demander un nouveau code.");
                case 404:
                    return setError(form, "Aucune demande de vérification trouvée pour cet email.");
                case 409:
                    return setError(form, "Un utilisateur avec cet email existe déjà.");
                case 415:
                    return setError(form, "Type de contenu non supporté. Veuillez réessayer.");
                case 423:
                    return setError(form, "Trop de tentatives échouées. Votre compte est temporairement verrouillé.");
                case 500:
                    return setError(form, "Une erreur serveur est survenue. Veuillez réessayer.");
                case 503:
                    return setError(form, "Le service est temporairement indisponible. Veuillez réessayer dans quelques instants.");
                default:
                    return setError(form, "Une erreur inattendue est survenue. Veuillez réessayer.");
            }
        }

        // Forward authentication cookies from backend response to client
        forwardAuthCookies(res, cookies);

        // Success - redirect to general info page
        redirect(302, "/auth/general");
    }
};
