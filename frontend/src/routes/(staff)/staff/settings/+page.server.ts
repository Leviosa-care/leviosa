import type { Actions, PageServerLoad } from "./$types"
import { fail, redirect } from "@sveltejs/kit"
import { env } from "$env/dynamic/private"

import { setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';

import {
    companyNameSchema,
    companyNameDefaults,
    companyEmailSchema,
    companyEmailDefaults,
    companyPhoneSchema,
    companyPhoneDefaults,
    companyAddressSchema,
    companyAddressDefaults,
    companyInstagramSchema,
    companyInstagramDefaults,
    otpDurationSchema,
    otpDurationDefaults,
    otpLengthSchema,
    otpLengthDefaults,
    otpMaxAttemptsSchema,
    otpMaxAttemptsDefaults,
    accessTokenDurationSchema,
    accessTokenDurationDefaults,
    refreshTokenDurationSchema,
    refreshTokenDurationDefaults,
} from "./schemas"
// Enable mock mode for staging
const MOCK_MODE = env.USE_MOCK_DATA === 'true'

type Settings = {
    company: {
        name: string
        email: string
        telephone: string
        address: string
        instagram: string
        logo_url: string
        logo_content_type: string
    }
    otp: {
        duration: number
        length: number
        max_attempts: number
    }
    tokens: {
        access_duration: number
        refresh_duration: number
    }
}

export const load: PageServerLoad = async ({ parent, fetch, cookies }) => {
    // ⬅️ pulls data from (ops)/+layout.server.ts
    const { permissions } = await parent()

    if (!permissions.canAccessOps) {
        throw redirect(302, '/app')
    }

    // Initialize forms
    const companyNameForm = await superValidate(arktype(companyNameSchema, { defaults: companyNameDefaults }))
    const companyEmailForm = await superValidate(arktype(companyEmailSchema, { defaults: companyEmailDefaults }))
    const companyPhoneForm = await superValidate(arktype(companyPhoneSchema, { defaults: companyPhoneDefaults }))
    const companyAddressForm = await superValidate(arktype(companyAddressSchema, { defaults: companyAddressDefaults }))
    const companyInstagramForm = await superValidate(arktype(companyInstagramSchema, { defaults: companyInstagramDefaults }))
    const otpDurationForm = await superValidate(arktype(otpDurationSchema, { defaults: otpDurationDefaults }))
    const otpLengthForm = await superValidate(arktype(otpLengthSchema, { defaults: otpLengthDefaults }))
    const otpMaxAttemptsForm = await superValidate(arktype(otpMaxAttemptsSchema, { defaults: otpMaxAttemptsDefaults }))
    const accessTokenDurationForm = await superValidate(arktype(accessTokenDurationSchema, { defaults: accessTokenDurationDefaults }))
    const refreshTokenDurationForm = await superValidate(arktype(refreshTokenDurationSchema, { defaults: refreshTokenDurationDefaults }))

    let settings: Settings | null = null
    let settingsError = false

    const sessionCookie = cookies.get('session')

    const bulkKeys = [
        'company_name',
        'company_email',
        'company_phone',
        'company_address',
        'company_instagram',
        'otp_duration',
        'otp_length',
        'otp_max_attempts',
        'access_token_duration',
        'refresh_token_duration'
    ].join(',')

    try {
        const settingsRes = await fetch(`${env.API_URL}/admin/settings/bulk?keys=${bulkKeys}`, {
            headers: {
                'Authorization': `Bearer ${sessionCookie}`,
            }
        })

        if (settingsRes.ok) {
            const bulkData: { key: string; value: string }[] = await settingsRes.json()
            const find = (key: string) => bulkData.find((s) => s.key === key)?.value ?? ''

            settings = {
                company: {
                    name: find('company_name'),
                    email: find('company_email'),
                    telephone: find('company_phone'),
                    address: find('company_address'),
                    instagram: find('company_instagram'),
                    logo_url: '',
                    logo_content_type: '',
                },
                otp: {
                    duration: parseInt(find('otp_duration') || '300'),
                    length: parseInt(find('otp_length') || '6'),
                    max_attempts: parseInt(find('otp_max_attempts') || '5'),
                },
                tokens: {
                    access_duration: parseInt(find('access_token_duration') || '15'),
                    refresh_duration: parseInt(find('refresh_token_duration') || '168'),
                }
            }
        } else {
            settingsError = true
        }
    } catch {
        settingsError = true
    }

    return {
        companyNameForm,
        companyEmailForm,
        companyPhoneForm,
        companyAddressForm,
        companyInstagramForm,
        otpDurationForm,
        otpLengthForm,
        otpMaxAttemptsForm,
        accessTokenDurationForm,
        refreshTokenDurationForm,
        settings,
        settingsError,
    }
}

export const actions = {
    updateCompanyName: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(companyNameSchema, { defaults: companyNameDefaults }))

        if (!form.valid) {
            return fail(400, { companyNameForm: form })
        }

        if (MOCK_MODE) {
            console.log('🏢 [MOCK] Updating company name:', form.data.name)
            return { companyNameForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/name`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ name: form.data.name })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Nom invalide.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { companyNameForm: form }
    },

    updateCompanyEmail: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(companyEmailSchema, { defaults: companyEmailDefaults }))

        if (!form.valid) {
            return fail(400, { companyEmailForm: form })
        }

        if (MOCK_MODE) {
            console.log('📧 [MOCK] Updating company email:', form.data.email)
            return { companyEmailForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/email`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ email: form.data.email })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Email invalide.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { companyEmailForm: form }
    },

    updateCompanyPhone: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(companyPhoneSchema, { defaults: companyPhoneDefaults }))

        if (!form.valid) {
            return fail(400, { companyPhoneForm: form })
        }

        if (MOCK_MODE) {
            console.log('📞 [MOCK] Updating company phone:', form.data.telephone)
            return { companyPhoneForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/phone`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ telephone: form.data.telephone })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Téléphone invalide.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { companyPhoneForm: form }
    },

    updateCompanyAddress: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(companyAddressSchema, { defaults: companyAddressDefaults }))

        if (!form.valid) {
            return fail(400, { companyAddressForm: form })
        }

        if (MOCK_MODE) {
            console.log('🏠 [MOCK] Updating company address:', form.data.address)
            return { companyAddressForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/address`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ address: form.data.address })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Adresse invalide.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { companyAddressForm: form }
    },

    updateCompanyInstagram: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(companyInstagramSchema, { defaults: companyInstagramDefaults }))

        if (!form.valid) {
            return fail(400, { companyInstagramForm: form })
        }

        if (MOCK_MODE) {
            console.log('📸 [MOCK] Updating company Instagram:', form.data.instagram)
            return { companyInstagramForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/instagram`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ instagram: form.data.instagram })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "URL Instagram invalide.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { companyInstagramForm: form }
    },

    updateOtpDuration: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(otpDurationSchema, { defaults: otpDurationDefaults }))

        if (!form.valid) {
            return fail(400, { otpDurationForm: form })
        }

        if (MOCK_MODE) {
            console.log('⏱️ [MOCK] Updating OTP duration:', form.data.duration)
            return { otpDurationForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/otp/duration`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ duration: form.data.duration })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Durée invalide (60-3600 secondes).");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { otpDurationForm: form }
    },

    updateOtpLength: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(otpLengthSchema, { defaults: otpLengthDefaults }))

        if (!form.valid) {
            return fail(400, { otpLengthForm: form })
        }

        if (MOCK_MODE) {
            console.log('🔢 [MOCK] Updating OTP length:', form.data.length)
            return { otpLengthForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/otp/length`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ length: form.data.length })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Longueur invalide (4-10 chiffres).");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { otpLengthForm: form }
    },

    updateOtpMaxAttempts: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(otpMaxAttemptsSchema, { defaults: otpMaxAttemptsDefaults }))

        if (!form.valid) {
            return fail(400, { otpMaxAttemptsForm: form })
        }

        if (MOCK_MODE) {
            console.log('🚫 [MOCK] Updating OTP max attempts:', form.data.max_attempts)
            return { otpMaxAttemptsForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/otp/max-attempts`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ max_attempts: form.data.max_attempts })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Nombre de tentatives invalide (1-10).");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { otpMaxAttemptsForm: form }
    },

    updateAccessTokenDuration: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(accessTokenDurationSchema, { defaults: accessTokenDurationDefaults }))

        if (!form.valid) {
            return fail(400, { accessTokenDurationForm: form })
        }

        if (MOCK_MODE) {
            console.log('🔑 [MOCK] Updating access token duration:', form.data.duration)
            return { accessTokenDurationForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/tokens/access-duration`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ duration: form.data.duration })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Durée invalide (1-240 minutes).");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { accessTokenDurationForm: form }
    },

    updateRefreshTokenDuration: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(refreshTokenDurationSchema, { defaults: refreshTokenDurationDefaults }))

        if (!form.valid) {
            return fail(400, { refreshTokenDurationForm: form })
        }

        if (MOCK_MODE) {
            console.log('🔄 [MOCK] Updating refresh token duration:', form.data.duration)
            return { refreshTokenDurationForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${env.API_URL}/admin/settings/tokens/refresh-duration`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({ duration: form.data.duration })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Durée invalide (1-720 heures).");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { refreshTokenDurationForm: form }
    },
} satisfies Actions
