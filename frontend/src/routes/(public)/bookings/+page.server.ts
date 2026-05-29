import { env } from '$env/dynamic/private';
import { redirect } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

function roleBookingsPath(role: string): string | null {
	switch (role) {
		case 'standard':
		case 'premium':
			return '/client/bookings';
		case 'administrator':
			return '/admin/bookings';
		case 'partner':
			return '/staff/agenda/reservations';
		default:
			return null;
	}
}

export const load: PageServerLoad = async ({ url, locals, fetch, cookies }) => {
	if (locals.user) {
		const dest = roleBookingsPath(locals.user.role);
		if (dest) throw redirect(302, dest);
	} else {
		// Production: hooks.server.ts skips session validation for non-protected routes,
		// so locals.user is not populated. Check the session cookie explicitly.
		const sessionID = cookies.get(locals.sessionCookieName);
		if (sessionID) {
			const res = await fetch(`${env.API_URL}/users/me`, {
				headers: { Cookie: `leviosa_access_token=${sessionID}` }
			}).catch(() => null);
			if (res?.ok) {
				const user = await res.json();
				const dest = roleBookingsPath(user.role);
				if (dest) throw redirect(302, dest);
			}
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

	return { booking, lookupError, token };
};

export const actions: Actions = {
	lookup: async ({ request, fetch }) => {
		const formData = await request.formData();
		const ref = (formData.get('ref') as string)?.trim() ?? '';
		const email = (formData.get('email') as string)?.trim() ?? '';
		const phone = (formData.get('phone') as string)?.trim() ?? '';

		if (!ref) {
			return { action: 'lookup' as const, success: false, error: 'Veuillez entrer la référence de réservation.' };
		}

		const params = new URLSearchParams({ ref });
		if (email) {
			params.set('email', email);
		} else if (phone) {
			params.set('phone', phone);
		} else {
			return { action: 'lookup' as const, success: false, error: 'Veuillez entrer votre email ou numéro de téléphone.' };
		}

		try {
			const res = await fetch(`${env.API_URL}/bookings/lookup?${params}`);
			if (res.ok) {
				return { action: 'lookup' as const, success: true, booking: await res.json() };
			}
			const body = await res.json().catch(() => ({ error: 'Erreur inconnue' }));
			return { action: 'lookup' as const, success: false, error: body.error ?? 'Aucune réservation trouvée.' };
		} catch {
			return { action: 'lookup' as const, success: false, error: 'Erreur réseau. Veuillez réessayer.' };
		}
	},

	cancel: async ({ request, fetch, url }) => {
		const formData = await request.formData();
		const bookingId = (formData.get('booking_id') as string)?.trim() ?? '';
		const token = (formData.get('token') as string)?.trim() ?? '';
		const reason = (formData.get('reason') as string)?.trim() ?? '';

		if (!bookingId || !token) {
			return { action: 'cancel' as const, success: false, error: 'Informations manquantes pour annuler.' };
		}

		try {
			const res = await fetch(
				`${env.API_URL}/bookings/${bookingId}/cancel-public?token=${encodeURIComponent(token)}`,
				{
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ reason: reason || 'Annulation par le client' })
				}
			);
			if (res.ok) {
				const booking = await res.json();
				return { action: 'cancel' as const, success: true, booking };
			}
			const body = await res.json().catch(() => ({ error: 'Erreur inconnue' }));
			return { action: 'cancel' as const, success: false, error: body.error ?? "Échec de l'annulation." };
		} catch {
			return { action: 'cancel' as const, success: false, error: 'Erreur réseau. Veuillez réessayer.' };
		}
	}
};
