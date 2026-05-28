import { env } from '$env/dynamic/private';
import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ url, locals, fetch, cookies }) => {
	// Mock mode (dev/staging): locals.user is always set → redirect immediately
	if (locals.user) {
		throw redirect(302, '/client/bookings');
	}

	// Production: hooks.server.ts skips session validation for non-protected routes,
	// so locals.user is not populated. Check the session cookie explicitly.
	const sessionID = cookies.get(locals.sessionCookieName);
	if (sessionID) {
		const res = await fetch(`${env.API_URL}/users/me`, {
			headers: { Cookie: `leviosa_access_token=${sessionID}` }
		}).catch(() => null);
		if (res?.ok) {
			throw redirect(302, '/client/bookings');
		}
	}

	const token = url.searchParams.get('token');
	let booking = null;
	let lookupError: string | null = null;

	// Token path: server-side lookup so the guest sees the booking immediately
	if (token) {
		try {
			const res = await fetch(`${env.API_URL}/bookings/lookup?token=${encodeURIComponent(token)}`);
			if (res.ok) {
				booking = await res.json();
			} else {
				const body = await res.json().catch(() => ({ error: 'Erreur inconnue' }));
				lookupError = body.error ?? 'Lien invalide ou expiré';
			}
		} catch {
			lookupError = 'Erreur réseau. Veuillez réessayer.';
		}
	}

	return { booking, lookupError };
};

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		const formData = await request.formData();
		const ref = (formData.get('ref') as string)?.trim() ?? '';
		const email = (formData.get('email') as string)?.trim() ?? '';
		const phone = (formData.get('phone') as string)?.trim() ?? '';

		if (!ref) {
			return { success: false, error: 'Veuillez entrer la référence de réservation.' };
		}

		const params = new URLSearchParams({ ref });
		if (email) {
			params.set('email', email);
		} else if (phone) {
			params.set('phone', phone);
		} else {
			return { success: false, error: 'Veuillez entrer votre email ou numéro de téléphone.' };
		}

		try {
			const res = await fetch(`${env.API_URL}/bookings/lookup?${params}`);
			if (res.ok) {
				return { success: true, booking: await res.json() };
			}
			const body = await res.json().catch(() => ({ error: 'Erreur inconnue' }));
			return { success: false, error: body.error ?? 'Aucune réservation trouvée.' };
		} catch {
			return { success: false, error: 'Erreur réseau. Veuillez réessayer.' };
		}
	}
};
