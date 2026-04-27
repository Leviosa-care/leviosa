import type { Actions, RequestEvent } from "./$types"
import { error, fail, redirect } from "@sveltejs/kit"
import { env } from "$env/dynamic/private"

import { message, setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"

import { MESSAGES, type MessageType } from "$lib/utils/redirect";
import { forwardAuthCookies } from "$lib/utils/auth-helpers";

const loginSchema = type({
    email: "string.email",
    password: "8 < string < 64",
})
const loginDefaults = { email: "", password: "" }

const registerSchema = type({ registerEmail: "string.email", })
const registerDefaults = { registerEmail: "" }


const oauthSchema = type({
    provider: "'google' | 'apple'", // it is either google or apple to be fair
})
const oauthDefaults = {
    provider: "google",
} as typeof oauthSchema.infer

export const load = async ({ url }: RequestEvent) => {
    const registerForm = await superValidate(arktype(registerSchema, { defaults: registerDefaults }))
    const oauthForm = await superValidate(arktype(oauthSchema, { defaults: oauthDefaults }))

    // that part should be handled client side to be honest
    const redirectFrom = url.searchParams.get("redirectFrom");
    if (redirectFrom) {
        const queryMessage = url.searchParams.get("message") as MessageType;
        message(registerForm, MESSAGES[queryMessage])
    }

    // Store the redirect target for use after successful login
    const redirectTo = url.searchParams.get("redirect");

    // NOTE: on the oauth implementation
    // provider
    // auth state
    // auth/{provider}/{action}
    //
    const email = url.searchParams.get("email") // get email is from redirect (via the user.locals if hooks)
    const init = email ? { email, password: "" } : undefined
    // TODO: remove that once the pipeline is done for the registration
    // const registerForm = await superValidate({ registerEmail: mockUser.email }, arktype(registerSchema, { defaults: registerDefaults }))
    // const registerForm = await superValidate(wrongUser, arktype(registerSchema, { defaults: registerDefaults }))
    // const registerForm = await superValidate(invalidUser, arktype(registerSchema, { defaults: registerDefaults }))

    const loginForm = await superValidate(init, arktype(loginSchema, { defaults: loginDefaults }))
    return { loginForm, registerForm, oauthForm, redirectTo }
}

export const actions = {
    register: async ({ request, fetch }: RequestEvent) => {
        const form = await superValidate(request, arktype(registerSchema, { defaults: registerDefaults }))

        if (!form.valid) return setError(form, "registerEmail", "L'adresse email saisie n'est pas valide. Veuillez vérifier et réessayer.")

        const res = await fetch(`${env.API_URL}/auth/email`, {
            method: "POST",
            headers: { 'Content-Type': "application/json" },
            body: JSON.stringify({ email: form.data.registerEmail }),
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "registerEmail", "Format d'adresse e-mail invalide. Veuillez vérifier.");
                case 409:
                    return setError(form, "registerEmail", "Cette adresse e-mail est déjà enregistrée.");
                case 415:
                    return setError(form, "registerEmail", "Type de contenu non supporté. Veuillez réessayer.");
                case 429:
                    return setError(form, "registerEmail", "Trop de tentatives. Veuillez réessayer plus tard.");
                case 500:
                    return setError(form, "registerEmail", "Une erreur serveur est survenue. Le service semble momentanément indisponible.");
                case 503:
                    return setError(form, "registerEmail", "Le service est temporairement indisponible. Veuillez réessayer dans quelques instants.");
                default:
                    return setError(form, "registerEmail", "Une erreur inattendue est survenue. Veuillez réessayer.");
            }
        }

        // Success - redirect to verify email page with email as query param
        redirect(302, `/auth/verify-email?email=${encodeURIComponent(form.data.registerEmail)}`)
    },
    login: async ({ request, cookies, url }: RequestEvent) => {
        const form = await superValidate(request, arktype(loginSchema, { defaults: loginDefaults }))

        if (!form.valid) {
            if (form.errors.email) {
                setError(form, "email", "L'adresse email saisie n'est pas valide. Veuillez vérifier et réessayer.")
            }
            if (form.errors.password) {
                setError(form, "password", "Le mot de passe saisi n'est pas valide. Veuillez vérifier et réessayer.")
            }
            return fail(400, { loginForm: form })
        }

        const res = await fetch(`${env.API_URL}/auth/login`, {
            method: "POST",
            headers: { 'Content-Type': "application/json" },
            body: JSON.stringify({
                email: form.data.email,
                password: form.data.password,
            })
        })

        if (res.status !== 201) {
            switch (res.status) {
                case 400:
                    return setError(form, "Format d'adresse e-mail ou de mot de passe invalide. Veuillez vérifier.");
                case 401:
                    return setError(form, "Identifiants incorrects. Veuillez vérifier votre email et mot de passe.");
                case 403:
                    return setError(form, "Votre compte n'est pas encore approuvé ou est inactif. Veuillez contacter l'administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté. Veuillez réessayer.");
                case 423:
                    return setError(form, "Votre compte est verrouillé suite à trop de tentatives échouées. Veuillez réessayer plus tard.");
                case 429:
                    return setError(form, "Trop de tentatives. Veuillez réessayer plus tard.");
                case 500:
                    return setError(form, "Une erreur est survenue sur le serveur. Nous travaillons à résoudre le problème au plus vite.");
                case 503:
                    return setError(form, "Le service est temporairement indisponible. Veuillez réessayer dans quelques instants.");
                default:
                    return setError(form, "Une erreur inattendue est survenue. Veuillez réessayer.");
            }
        }

        // Extract and forward authentication cookies from backend response to client
        forwardAuthCookies(res, cookies);

        // Get redirect target from URL params
        const redirectTo = url.searchParams.get("redirect");

        // Fetch user to determine role-based redirect
        // Note: We need to make a fresh request since we just received the session cookie
        let finalRedirect = "/";
        try {
            const userRes = await fetch(`${env.API_URL}/users/me`, {
                headers: {
                    'Cookie': `leviosa_access_token=${cookies.get("leviosa_access_token") || ""}`
                }
            });

            if (userRes.ok) {
                const user = await userRes.json();

                // Role-based redirect: always send users to the appropriate page for their role
                // This prevents staff users from getting 403 errors if they tried to access /admin
                if (user.role === "administrator") {
                    finalRedirect = "/admin";
                } else if (user.role === "partner") {
                    finalRedirect = "/staff";
                } else if (redirectTo) {
                    // For other roles, respect the redirect param if provided
                    finalRedirect = redirectTo;
                }
            } else if (redirectTo) {
                // Fallback: use redirect param if user fetch fails
                finalRedirect = redirectTo;
            }
        } catch (e) {
            // If user fetch fails, use redirect param or default to home
            if (redirectTo) {
                finalRedirect = redirectTo;
            }
        }

        // Success - redirect to appropriate destination
        redirect(302, finalRedirect);
    },
    oauth: async ({ request }: RequestEvent) => {
        const form = await superValidate(request, arktype(oauthSchema, { defaults: oauthDefaults }))

        if (!form.valid) {
            setError(form, "provider", "Le fournisseur OAuth sélectionné n'est pas valide.")
            return fail(400, { oauthForm: form })
        }

        // Redirect to OAuth provider
        // The backend GET endpoint will redirect to the provider's authorization screen
        // After authorization, the provider will redirect back to /auth/oauth/{provider}/callback
        // The backend callback handler will set cookies and redirect to the app
        redirect(302, `${env.API_URL}/auth/oauth/${form.data.provider}`)
    }
} satisfies Actions
