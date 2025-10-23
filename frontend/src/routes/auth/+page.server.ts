import type { Actions, RequestEvent } from "./$types"
import { error, fail, redirect } from "@sveltejs/kit"
import { API_URL, SESSION_COOKIE_NAME } from "$env/static/private"

import { message, setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"

import { MESSAGES, type MessageType } from "$lib/utils/redirect";

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
    // const registerForm = await superValidate(arktype(registerSchema, { defaults: registerDefaults }))
    // TODO: just for testing purposes remove that later
    const registerForm = await superValidate({ registerEmail: "jean.dupont@leviosa.care" }, arktype(registerSchema, { defaults: registerDefaults }))
    const oauthForm = await superValidate(arktype(oauthSchema, { defaults: oauthDefaults }))

    // that part should be handled client side to be honest
    const redirectFrom = url.searchParams.get("redirectFrom");
    if (redirectFrom) {
        const queryMessage = url.searchParams.get("message") as MessageType;
        message(registerForm, MESSAGES[queryMessage])
    }

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
    console.log("the register form from the load function:", registerForm)
    return { loginForm, registerForm, oauthForm }
}

export const actions = {
    register: async ({ request, fetch }: RequestEvent) => {
        const form = await superValidate(request, arktype(registerSchema, { defaults: registerDefaults }))

        if (!form.valid) return setError(form, "registerEmail", "L'adresse email saisie n'est pas valide. Veuillez vérifier et réessayer.")

        const res = await fetch(`${API_URL}/auth/email`, {
            method: "POST",
            headers: { 'Content-Type': "application/json" },
            body: JSON.stringify({ email: form.data.registerEmail }),
        })

        if (!res.ok) {
            switch (res.status) {
                case 500:
                    return setError(form, "registerEmail", "Une erreur serveur est survenue. Le service semble momentanément indisponible.")
                case 422:
                // something went wrong with validating the data. I might need to parse that thing
                case 429:
                    return setError(form, "registerEmail", "Trop de tentatives. Veuillez réessayer plus tard.");
                case 409:
                // the user exists
                case 400:
                    return setError(form, "registerEmail", "Format d'adresse e-mail invalide. Veuillez vérifier.");
            }
        }
        // const exists = await res.json()
        // if (exists === true) {
        //     return setError(form, "registerEmail", "Cette adresse e-mail est déjà utilisée.");
        // }
        return {
            registerForm: form,
            registerSuccess: true,
        }
    },
    login: async ({ request, cookies }: RequestEvent) => {
        const form = await superValidate(request, arktype(loginSchema, { defaults: loginDefaults }))
        // TODO: remove that at the end
        console.log(form)
        if (!form.valid) {
            if (form.errors.email) {
                setError(form, "email", "L'adresse email saisie n'est pas valide. Veuillez vérifier et réessayer.")
            }
            if (form.errors.password) {
                setError(form, "password", "Le mot de passe saisi n'est pas valide. Veuillez vérifier et réessayer.")
            }
            return fail(400, { loginForm: form })
        }
        const res = await fetch(`${API_URL}/auth/login`, {
            method: "POST",
            headers: { 'Content-Type': "application/json" },
            body: JSON.stringify({
                email: form.data.email,
                password: form.data.password,
            })
        })
        // if (!res.ok) {
        if (res.status != 201) {
            const err = res.text() // from golang http.Error
            console.error("login attempt failed:", err)
            switch (res.status) {
                case 400:
                    return setError(form, "Format d'adresse e-mail ou de mot de passe invalide. Veuillez vérifier.");
                case 429: // too many requests
                    return setError(form, "Trop de tentatives. Veuillez réessayer plus tard.");
                case 500:
                    error(500, "Une erreur est survenue sur le serveur. Nous travaillons à résoudre le problème au plus vite.")
                default:
                    error(404, "Une erreur est survenue. Veuillez réessayer.")
            }
        }
        const sessionID = cookies.get(SESSION_COOKIE_NAME)
        if (!sessionID) {
            error(404, "Le serveur est en panne. Veuillez réessayer plus tard.")
        }
        cookies.set(SESSION_COOKIE_NAME, sessionID, {
            // TODO: set that properly with the right values
            httpOnly: true,
            maxAge: 60 * 10,
            secure: import.meta.env.PROD,
            path: '/',
            sameSite: 'lax'
        })

        return {
            form,
            success: true,
        }
    },
    oauth: async ({ request, cookies }: RequestEvent) => {
        const form = await superValidate(request, arktype(oauthSchema, { defaults: oauthDefaults }))
        console.log(form)
        if (!form.valid) {
            setError(form, "provider", "L'adresse provider saisie n'est pas valide. Veuillez vérifier et réessayer.")
            return fail(400, { oauthForm: form })
        }

        const res = await fetch(`${API_URL}/oauth/${form.data.provider}`, {
            method: "POST",
        })
        if (!res.ok) {
            const err = res.text() // from golang http.Error
            console.error("login attempt failed:", err)
            switch (res.status) {
                case 400:
                // return setError(form, "Format d'adresse e-mail ou de mot de passe invalide. Veuillez vérifier.");
                case 429: // too many requests
                // return setError(form, "Trop de tentatives. Veuillez réessayer plus tard.");
                case 500:
                    error(500, "Une erreur est survenue sur le serveur. Nous travaillons à résoudre le problème au plus vite.")
                default:
                    error(404, "Une erreur est survenue. Veuillez réessayer.")
            }
        }
        const sessionID = cookies.get(SESSION_COOKIE_NAME)
        if (!sessionID) {
            error(404, "Le serveur est en panne. Veuillez réessayer plus tard.")
        }
        cookies.set(SESSION_COOKIE_NAME, sessionID, {
            // TODO: set that properly with the right values
            httpOnly: true,
            maxAge: 60 * 10,
            secure: import.meta.env.PROD,
            path: '/',
            sameSite: 'lax'
        })
        throw redirect(302, "/")
    }
} satisfies Actions
