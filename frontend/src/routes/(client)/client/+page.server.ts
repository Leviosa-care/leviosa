import { env } from '$env/dynamic/private';
import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals, fetch }) => {
	const user = locals.user!;
	const clientId = user.id;

	let bookings: any[] = [];

	try {
		const res = await fetch(`${env.API_URL}/clients/${clientId}/bookings`);
		if (res.ok) {
			bookings = await res.json();
			if (!Array.isArray(bookings)) bookings = [];
		}
	} catch {
		// Non-critical — will show empty state
	}

	const now = new Date();

	// Sort bookings by start time descending
	bookings.sort((a, b) => new Date(b.slot_start_time).getTime() - new Date(a.slot_start_time).getTime());

	// Next upcoming confirmed booking
	const upcoming = bookings
		.filter((b: any) => b.status === 'confirmed' && new Date(b.slot_start_time) > now)
		.sort((a: any, b: any) => new Date(a.slot_start_time).getTime() - new Date(b.slot_start_time).getTime());

	const nextBooking: any | null = upcoming[0] ?? null;

	// Last 3 completed bookings
	const recentCompleted = bookings
		.filter((b: any) => b.status === 'completed')
		.sort((a: any, b: any) => new Date(b.slot_start_time).getTime() - new Date(a.slot_start_time).getTime())
		.slice(0, 3);

	return { nextBooking, recentCompleted };
};
