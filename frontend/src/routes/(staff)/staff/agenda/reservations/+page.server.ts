import type { PageServerLoad } from './$types';
import { redirect, isRedirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export interface Booking {
	id: string;
	clientId: string;
	clientName: string;
	clientInitials: string;
	productName: string;
	startTime: string;
	endTime: string;
	roomName: string;
	status: 'confirmed' | 'completed' | 'no_show' | 'cancelled';
	paymentStatus: 'pending' | 'paid' | 'failed' | 'refunded';
	amountInCents: number;
	currency: string;
	clientNotes: string;
	partnerNotes: string;
}

export const load: PageServerLoad = async ({ locals, url, fetch }) => {
	if (!locals.user?.id) {
		throw redirect(302, '/auth');
	}

	const partnerId = locals.user.id;
	const statusFilter = url.searchParams.get('status') ?? '';

	try {
		const res = await fetch(`${env.API_URL}/partners/bookings/${partnerId}`);

		if (res.status === 401) {
			throw redirect(302, '/auth');
		}

		if (!res.ok) {
			console.error(`Failed to fetch bookings: ${res.status} ${res.statusText}`);
			return { bookings: [], statusFilter };
		}

		const data = await res.json();

		if (!Array.isArray(data)) {
			console.error('Invalid response format: expected array');
			return { bookings: [], statusFilter };
		}

		const bookings: Booking[] = data.map((booking: any) => {
			const clientInitials = (booking.client_name || '??')
				.split(' ')
				.map((s: string) => s[0]?.toUpperCase() || '')
				.slice(0, 2)
				.join('');

			return {
				id: booking.id,
				clientId: booking.client_id,
				clientName: booking.client_name || 'Client inconnu',
				clientInitials,
				productName: booking.product_name || 'Produit inconnu',
				startTime: booking.slot_start_time,
				endTime: booking.slot_end_time,
				roomName: booking.room_name || 'Salle inconnue',
				status: booking.status || 'confirmed',
				paymentStatus: booking.payment_status || 'pending',
				amountInCents: booking.total_price_cents || 0,
				currency: booking.currency || 'EUR',
				clientNotes: booking.client_notes || '',
				partnerNotes: booking.partner_notes || '',
			};
		});

		bookings.sort((a, b) => new Date(a.startTime).getTime() - new Date(b.startTime).getTime());

		return { bookings, statusFilter };
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading bookings:', err);
		return { bookings: [], statusFilter };
	}
};
