import { env } from '$env/dynamic/private';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ params, locals, fetch }) => {
	let booking: any = null;

	try {
		const res = await fetch(`${env.API_URL}/bookings/${params.id}`);
		if (res.status === 404) {
			throw error(404, 'Réservation introuvable');
		}
		if (!res.ok) {
			throw error(500, 'Erreur lors du chargement de la réservation');
		}
		booking = await res.json();
	} catch (err) {
		if (err && typeof err === 'object' && 'status' in err) throw err;
		throw error(500, 'Erreur lors du chargement de la réservation');
	}

	if (booking.client_id !== locals.user!.id) {
		throw error(403, 'Accès refusé');
	}

	return { booking };
};
