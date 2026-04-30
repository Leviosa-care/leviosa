import type { Actions, PageServerLoad, RequestEvent } from "./$types";
import { fail, redirect } from "@sveltejs/kit";
import { setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"
import { env } from "$env/dynamic/private";
import { forwardAuthCookies } from "$lib/utils/auth-helpers";

// Step 1: Request password reset
const requestSchema = type({
    email: "string.email",
})

const requestDefaults = { email: "" }

// Step 2: Validate OTP
const validateSchema = type({
    email: "string.email",
    otp0: "/^\\d$/",
    otp1: "/^\\d$/",
    otp2: "/^\\d$/",
    otp3: "/^\\d$/",
    otp4: "/^\\d$/",
    otp5: "/^\\d$/",
})

const validateDefaults = {
    email: '',
    otp0: '',
    otp1: '',
    otp2: '',
    otp3: '',
    otp4: '',
    otp5: '',
}

// Step 3: Confirm new password
const confirmSchema = type({
    newPassword: "8 < string < 64",
    confirmPassword: "8 < string < 64",
})

const confirmDefaults = { newPassword: "", confirmPassword: "" }

export const load: PageServerLoad = async ({ url }: RequestEvent) => {
    const step = url.searchParams.get("step") || "request"; // request, validate, or confirm
    const email = url.searchParams.get("email") || "";

    let requestForm, validateForm, confirmForm;

    if (step === "validate") {
        validateForm = await superValidate({ ...validateDefaults, email }, arktype(validateSchema, { defaults: validateDefaults }));
    } else if (step === "confirm") {
        confirmForm = await superValidate(confirmDefaults, arktype(confirmSchema, { defaults: confirmDefaults }));
    } else {
        requestForm = await superValidate(arktype(requestSchema, { defaults: requestDefaults }));
    }

    return {
        step,
        requestForm,
        validateForm,
        confirmForm,
    };
}

export const actions = {
    // Step 1: Request password reset
    request: async ({ request, fetch }: RequestEvent) => {
        const form = await superValidate(request, arktype(requestSchema, { defaults: requestDefaults }));

        if (!form.valid) {
            return setError(form, "email", "L'adresse email n'est pas valide.");
        }

        const res = await fetch(`${env.API_URL}/auth/password/reset/request`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email: form.data.email })
        });

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "email", "Format d'email invalide.");
                case 415:
                    return setError(form, "email", "Type de contenu non supporté. Veuillez réessayer.");
                case 429:
                    return setError(form, "email", "Trop de demandes. Veuillez réessayer plus tard.");
                case 500:
                    return setError(form, "email", "Une erreur serveur est survenue. Veuillez réessayer.");
                case 503:
                    return setError(form, "email", "Le service d'email est temporairement indisponible. Veuillez réessayer dans quelques instants.");
                default:
                    return setError(form, "email", "Une erreur inattendue est survenue. Veuillez réessayer.");
            }
        }

        // Success - redirect to OTP validation step
        redirect(302, `/auth/forgotten-password?step=validate&email=${encodeURIComponent(form.data.email)}`);
    },

    // Step 2: Validate OTP
    validate: async ({ request, fetch, cookies }: RequestEvent) => {
        const form = await superValidate(request, arktype(validateSchema, { defaults: validateDefaults }));

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

        const res = await fetch(`${env.API_URL}/auth/password/reset/validate`, {
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
                    return setError(form, "Format d'email ou de code invalide.");
                case 401:
                    return setError(form, "Code OTP invalide ou expiré.");
                case 404:
                    return setError(form, "Aucune demande de réinitialisation trouvée pour cet email.");
                case 415:
                    return setError(form, "Type de contenu non supporté. Veuillez réessayer.");
                case 423:
                    return setError(form, "Trop de tentatives échouées. Veuillez réessayer plus tard.");
                case 500:
                    return setError(form, "Une erreur serveur est survenue. Veuillez réessayer.");
                case 503:
                    return setError(form, "Le service est temporairement indisponible. Veuillez réessayer dans quelques instants.");
                default:
                    return setError(form, "Une erreur inattendue est survenue. Veuillez réessayer.");
            }
        }

        // Forward auth cookies (leviosa_password_reset_token) from backend response
        forwardAuthCookies(res, cookies);

        // Success - redirect to password confirmation step
        redirect(302, `/auth/forgotten-password?step=confirm`);
    },

    // Step 3: Confirm new password
    confirm: async ({ request, fetch, cookies }: RequestEvent) => {
        const form = await superValidate(request, arktype(confirmSchema, { defaults: confirmDefaults }));

        if (!form.valid) {
            if (form.errors.newPassword) {
                setError(form, "newPassword", "Le mot de passe doit contenir au moins 8 caractères.");
            }
            if (form.errors.confirmPassword) {
                setError(form, "confirmPassword", "La confirmation du mot de passe est requise.");
            }
            return fail(400, { form });
        }

        // Verify passwords match
        if (form.data.newPassword !== form.data.confirmPassword) {
            return setError(form, "confirmPassword", "Les mots de passe ne correspondent pas.");
        }

        const res = await fetch(`${env.API_URL}/auth/password/reset/confirm`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Cookie": `leviosa_password_reset_token=${cookies.get("leviosa_password_reset_token") ?? ""}`,
            },
            body: JSON.stringify({
                new_password: form.data.newPassword
            })
        });

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "newPassword", "Le mot de passe ne respecte pas les exigences de sécurité.");
                case 401:
                    return setError(form, "Token de réinitialisation invalide ou expiré.");
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

        // Success - redirect to login page with success message
        redirect(302, "/auth?message=password_reset_success");
    }
} satisfies Actions;
