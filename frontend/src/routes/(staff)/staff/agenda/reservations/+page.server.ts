import type { PageServerLoad } from './$types';
import { redirect, isRedirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export interface Booking {
	id: string;
	clientName: string;
	clientInitials: string;
	productName: string;
	startTime: string;
	endTime: string;
	roomName: string;
	status: 'upcoming' | 'completed' | 'no_show' | 'cancelled';
	amountInCents: number;
	notes?: string;
}

export const load: PageServerLoad = async ({ locals, url, fetch }) => {
	if (!locals.user?.id) {
		throw redirect(302, '/auth');
	}

	const partnerId = locals.user.id;
	const statusFilter = url.searchParams.get('status');

	try {
		const apiUrl = new URL(`${env.API_URL}/partners/${partnerId}/bookings`);
		if (statusFilter) {
			apiUrl.searchParams.set('status', statusFilter);
		}

		const res = await fetch(apiUrl.toString(), {
			headers: {
				'Content-Type': 'application/json',
			},
		});

		if (res.status === 401) {
			throw redirect(302, '/auth');
		}

		if (!res.ok) {
			console.error(`Failed to fetch bookings: ${res.status} ${res.statusText}`);
			return { bookings: [] };
		}

		const data = await res.json();

		if (!Array.isArray(data)) {
			console.error('Invalid response format: expected array');
			return { bookings: [] };
		}

		const bookings: Booking[] = data.map((booking: any) => {
			const clientInitials = booking.client_id
				? booking.client_id
						.split('-')
						.slice(-2)
						.map((s: string) => s[0]?.toUpperCase() || '')
						.join('')
				: '??';

			return {
				id: booking.id,
				clientName: booking.client_name || 'Client inconnu',
				clientInitials,
				productName: booking.product_name || 'Produit inconnu',
				startTime: booking.slot_start_time,
				endTime: booking.slot_end_time,
				roomName: booking.room_name || 'Salle inconnue',
				status: booking.status || 'upcoming',
				amountInCents: booking.total_price_cents || 0,
				notes: booking.client_notes || booking.partner_notes || undefined,
			};
		});

		bookings.sort((a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime());

		return { bookings };
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading bookings:', err);
		return { bookings: [] };
	}
};
