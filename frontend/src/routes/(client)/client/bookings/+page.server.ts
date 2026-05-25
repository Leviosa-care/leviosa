import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const clientId = locals.user!.id;
	let bookings: any[] = [];

	try {
		const res = await fetch(`${env.API_URL}/clients/${clientId}/bookings`);
		if (res.ok) {
			bookings = await res.json();
			if (!Array.isArray(bookings)) bookings = [];
		}
	} catch {
		// Non-critical
	}

	// Sort by start time descending (newest first)
	bookings.sort((a: any, b: any) => new Date(b.slot_start_time).getTime() - new Date(a.slot_start_time).getTime());

	return { bookings };
};
