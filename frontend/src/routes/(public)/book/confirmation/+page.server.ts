import type { PageServerLoad, Actions } from "./$types";
import { env } from "$env/dynamic/private";
import { error, fail, redirect } from "@sveltejs/kit";
import { superValidate, setError } from "sveltekit-superforms";
import { arktype } from "sveltekit-superforms/adapters";
import { type } from "arktype";
import { forwardAuthCookies, formatPhoneToE164 } from "$lib/utils/auth-helpers";
import { getCookieDomain } from "$lib/server/hostname";

// ── Phase 1: Guest claim (password + optional email/phone) ──

const guestClaimSchema = type({
    email: "string",
    phone: "string",
    password: "8 < string < 64",
});

const guestClaimDefaults = {
    email: "",
    phone: "",
    password: "",
};

// ── Phase 2: OTP verification ──

const guestClaimVerifySchema = type({
    otp0: "/^\\d$/",
    otp1: "/^\\d$/",
    otp2: "/^\\d$/",
    otp3: "/^\\d$/",
    otp4: "/^\\d$/",
    otp5: "/^\\d$/",
});

const guestClaimVerifyDefaults = {
    otp0: "",
    otp1: "",
    otp2: "",
    otp3: "",
    otp4: "",
    otp5: "",
} as unknown as typeof guestClaimVerifySchema.infer;

// ── Load ──

export const load: PageServerLoad = async ({ fetch, url, locals, cookies }) => {
    const bookingId = url.searchParams.get("booking_id");
    if (!bookingId) {
        throw error(400, "Identifiant de réservation manquant");
    }

    let booking: any = null;
    let product: any = null;
    let priceDisplay: string | null = null;

    try {
        const bookingRes = await fetch(`${env.API_URL}/bookings/${bookingId}`);
        if (bookingRes.ok) {
            booking = await bookingRes.json();

            if (booking.product_id) {
                const productRes = await fetch(`${env.API_URL}/products/${booking.product_id}`);
                if (productRes.ok) {
                    product = await productRes.json();
                }
            }

            if (booking.total_price_cents) {
                priceDisplay = (booking.total_price_cents / 100)
                    .toFixed(2)
                    .replace(/\.00$/, "");
            }
        }
    } catch {
        // Booking fetch may fail for guests — fall through to URL params
    }

    // Fallback: use URL search params for guest booking display
    if (!booking) {
        booking = {
            id: bookingId,
            product_id: url.searchParams.get("product_id"),
            slot_start_time: url.searchParams.get("slot_start_time"),
            slot_end_time: url.searchParams.get("slot_end_time"),
            total_price_cents: url.searchParams.get("price_cents")
                ? parseInt(url.searchParams.get("price_cents")!)
                : null,
            guest_first_name: url.searchParams.get("guest_first_name"),
            guest_last_name: url.searchParams.get("guest_last_name"),
        };

        if (booking.product_id) {
            const productRes = await fetch(`${env.API_URL}/products/${booking.product_id}`);
            if (productRes.ok) {
                product = await productRes.json();
            }
        }

        if (booking.total_price_cents) {
            priceDisplay = (booking.total_price_cents / 100)
                .toFixed(2)
                .replace(/\.00$/, "");
        }
    }

    // Read guest contact info from the short-lived cookie set by the booking action.
    // This avoids leaking PII into the URL while still pre-filling the inline card.
    let guestInfo: { guest_first_name: string; guest_last_name: string; guest_email: string; guest_phone: string } | null = null;
    const guestCookie = cookies.get("guest_booking_info");
    if (guestCookie) {
        try {
            guestInfo = JSON.parse(guestCookie);
        } catch {
            // Invalid cookie — ignore
        }
    }

    // Merge guest info into booking for display (e.g. "Réservé par")
    if (guestInfo && !booking.guest_first_name) {
        booking.guest_first_name = guestInfo.guest_first_name;
        booking.guest_last_name = guestInfo.guest_last_name;
    }

    // Build claim forms
    const claimFormDefaults = {
        email: guestInfo?.guest_email ?? "",
        phone: guestInfo?.guest_phone ?? "",
        password: "",
    };
    const guestClaimForm = await superValidate(claimFormDefaults, arktype(guestClaimSchema, { defaults: guestClaimDefaults }));
    const guestClaimVerifyForm = await superValidate(arktype(guestClaimVerifySchema, { defaults: guestClaimVerifyDefaults }));

    return {
        booking,
        product,
        priceDisplay,
        user: locals.user ?? null,
        guestInfo,
        guestClaimForm,
        guestClaimVerifyForm,
    };
};

// ── Actions ──

export const actions: Actions = {
    // Phase 1: Submit claim — sends OTP to guest's email
    guestClaim: async ({ request, fetch, cookies, url }) => {
        const form = await superValidate(request, arktype(guestClaimSchema, { defaults: guestClaimDefaults }));

        if (!form.valid) {
            if (form.errors.password) {
                return setError(form, "password", "Le mot de passe doit contenir au moins 8 caractères.");
            }
            return fail(400, { guestClaimForm: form });
        }

        // Read guest info from cookie for name data
        const guestCookie = cookies.get("guest_booking_info");
        let guestInfo: { guest_first_name: string; guest_last_name: string } | null = null;
        if (guestCookie) {
            try { guestInfo = JSON.parse(guestCookie); } catch { /* ignore */ }
        }

        const firstName = guestInfo?.guest_first_name ?? "";
        const lastName = guestInfo?.guest_last_name ?? "";
        const phone = formatPhoneToE164(form.data.phone);

        const res = await fetch(`${env.API_URL}/auth/guest-claim`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: form.data.email,
                phone,
                password: form.data.password,
                first_name: firstName,
                last_name: lastName,
            }),
        });

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Données invalides. Veuillez vérifier les informations saisies.");
                case 409:
                    return setError(form, "email", "Cette adresse email est déjà associée à un compte.");
                case 415:
                    return setError(form, "Format de requête non supporté. Veuillez réessayer.");
                case 429:
                    return setError(form, "Trop de tentatives. Veuillez réessayer plus tard.");
                case 500:
                    return setError(form, "Une erreur serveur est survenue. Veuillez réessayer.");
                case 503:
                    return setError(form, "Le service est temporairement indisponible. Veuillez réessayer dans quelques instants.");
                default:
                    return setError(form, "Une erreur inattendue est survenue. Veuillez réessayer.");
            }
        }

        // Store claim data in a short-lived cookie for Phase 2
        const cookieDomain = getCookieDomain(url.hostname);
        cookies.set("guest_claim_data", JSON.stringify({
            email: form.data.email,
            phone,
            password: form.data.password,
            first_name: firstName,
            last_name: lastName,
        }), {
            path: "/book/confirmation",
            maxAge: 300, // 5 minutes
            httpOnly: true,
            secure: !url.hostname.startsWith('localhost'),
            sameSite: "lax",
            ...(cookieDomain && { domain: cookieDomain }),
        });

        // Return success — the client will transition to Phase 2
        return { guestClaimForm: form, claimPhase: "otp" };
    },

    // Phase 2: Verify OTP — creates the account and logs in
    guestClaimVerify: async ({ request, fetch, cookies, locals, url }) => {
        const form = await superValidate(request, arktype(guestClaimVerifySchema, { defaults: guestClaimVerifyDefaults }));

        if (!form.valid) {
            return setError(form, "Veuillez entrer un code à 6 chiffres.");
        }

        // Read claim data from cookie
        const claimCookie = cookies.get("guest_claim_data");
        if (!claimCookie) {
            return setError(form, "Session expirée. Veuillez recommencer la création de compte.");
        }

        let claimData: { email: string; phone: string; password: string; first_name: string; last_name: string };
        try {
            claimData = JSON.parse(claimCookie);
        } catch {
            return setError(form, "Session expirée. Veuillez recommencer la création de compte.");
        }

        const code = [
            form.data.otp0,
            form.data.otp1,
            form.data.otp2,
            form.data.otp3,
            form.data.otp4,
            form.data.otp5,
        ].join("");

        const res = await fetch(`${env.API_URL}/auth/guest-claim/verify`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({
                email: claimData.email,
                code,
                password: claimData.password,
                first_name: claimData.first_name,
                last_name: claimData.last_name,
                phone: claimData.phone,
            }),
        });

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Code invalide. Veuillez vérifier et réessayer.");
                case 401:
                    return setError(form, "Code expiré. Veuillez recommencer la création de compte.");
                case 404:
                    return setError(form, "Aucune demande de vérification trouvée. Veuillez recommencer.");
                case 409:
                    return setError(form, "Un compte avec cette adresse email existe déjà.");
                case 415:
                    return setError(form, "Format de requête non supporté. Veuillez réessayer.");
                case 429:
                    return setError(form, "Trop de tentatives. Veuillez réessayer plus tard.");
                case 500:
                    return setError(form, "Une erreur serveur est survenue. Veuillez réessayer.");
                case 503:
                    return setError(form, "Le service est temporairement indisponible. Veuillez réessayer dans quelques instants.");
                default:
                    return setError(form, "Une erreur inattendue est survenue. Veuillez réessayer.");
            }
        }

        // Forward auth cookies from backend to client
        forwardAuthCookies(res, cookies, locals.sessionCookieName, locals.cookieDomain);

        // Clean up guest cookies
        const cookieDomain = getCookieDomain(url.hostname);
        cookies.delete("guest_booking_info", { path: "/book/confirmation", ...(cookieDomain && { domain: cookieDomain }) });
        cookies.delete("guest_claim_data", { path: "/book/confirmation", ...(cookieDomain && { domain: cookieDomain }) });

        // Success — redirect to client bookings
        redirect(302, "/client/bookings");
    },
};
