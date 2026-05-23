import { env } from '$env/dynamic/private';
import { error, redirect, isRedirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

interface Consultation {
	id: string;
	clientName: string;
	therapistName: string;
	productName: string;
	date: string;
	duration: number;
	status: 'completed' | 'pending' | 'cancelled' | 'confirmed' | 'no_show';
	hasNotes: boolean;
}

interface APIBooking {
	id: string;
	client_name: string;
	partner_name: string;
	product_name: string;
	room_name: string;
	slot_start_time: string;
	slot_end_time: string;
	status: string;
	payment_status: string;
	total_price_cents: number;
	currency: string;
	created_at: string;
}

interface APIBookingsResponse {
	bookings: APIBooking[];
	total: number;
	page: number;
	limit: number;
}

async function getMockConsultations(): Promise<Consultation[]> {
	const now = new Date();

	return [
		{ id: '1', clientName: 'Marie Dupont', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 60min', date: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(), duration: 60, status: 'completed', hasNotes: true },
		{ id: '2', clientName: 'Jean Durand', therapistName: 'Pierre Leroy', productName: 'Consultation Kiné 45min', date: new Date(now.getTime() - 5 * 60 * 60 * 1000).toISOString(), duration: 45, status: 'completed', hasNotes: false },
		{ id: '3', clientName: 'Claire Bernard', therapistName: 'Marie Dubois', productName: 'Soin du Dos 60min', date: new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString(), duration: 60, status: 'completed', hasNotes: true },
		{ id: '4', clientName: 'Lucas Petit', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 90min', date: new Date(now.getTime() - 48 * 60 * 60 * 1000).toISOString(), duration: 90, status: 'cancelled', hasNotes: false },
		{ id: '5', clientName: 'Emma Moreau', therapistName: 'Pierre Leroy', productName: 'Drainage Lymphatique', date: new Date(now.getTime() + 3 * 60 * 60 * 1000).toISOString(), duration: 60, status: 'pending', hasNotes: false },
		{ id: '6', clientName: 'Thomas Richard', therapistName: 'Marie Dubois', productName: 'Soin du Dos 60min', date: new Date(now.getTime() + 6 * 60 * 60 * 1000).toISOString(), duration: 60, status: 'pending', hasNotes: false },
		{ id: '7', clientName: 'Camille Simon', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 60min', date: new Date(now.getTime() + 24 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), duration: 60, status: 'pending', hasNotes: false },
		{ id: '8', clientName: 'Hugo Michel', therapistName: 'Pierre Leroy', productName: 'Consultation Kiné 45min', date: new Date(now.getTime() - 72 * 60 * 60 * 1000).toISOString(), duration: 45, status: 'completed', hasNotes: true },
		{ id: '9', clientName: 'Chloe Garcia', therapistName: 'Marie Dubois', productName: 'Soin du Dos 60min', date: new Date(now.getTime() - 96 * 60 * 60 * 1000).toISOString(), duration: 60, status: 'completed', hasNotes: false },
		{ id: '10', clientName: 'Louis Laurent', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 90min', date: new Date(now.getTime() - 120 * 60 * 60 * 1000).toISOString(), duration: 90, status: 'completed', hasNotes: true },
		{ id: '11', clientName: 'Jade Roux', therapistName: 'Pierre Leroy', productName: 'Drainage Lymphatique', date: new Date(now.getTime() + 48 * 60 * 60 * 1000 + 14 * 60 * 60 * 1000).toISOString(), duration: 60, status: 'pending', hasNotes: false },
		{ id: '12', clientName: 'Nathan Girard', therapistName: 'Marie Dubois', productName: 'Soin du Dos 60min', date: new Date(now.getTime() - 144 * 60 * 60 * 1000).toISOString(), duration: 60, status: 'cancelled', hasNotes: false }
	];
}

function mapStatus(status: string): Consultation['status'] {
	// 'confirmed' maps to 'pending' because the UI vocabulary treats an upcoming
	// confirmed booking as "À venir" rather than showing raw API status names.
	if (status === 'confirmed') return 'pending';
	return status as Consultation['status'];
}

export const load: PageServerLoad = async ({ fetch, url }) => {
	if (env.USE_MOCK_DATA === 'true') {
		return { consultations: await getMockConsultations() };
	}

	try {
		const params = new URLSearchParams();
		const status = url.searchParams.get('status');
		const partnerId = url.searchParams.get('therapist');
		const page = url.searchParams.get('page') || '1';

		if (status) params.set('status', status);
		if (partnerId) params.set('partner_id', partnerId);
		params.set('page', page);
		params.set('limit', '20');

		const res = await fetch(`${env.API_URL}/admin/bookings?${params.toString()}`);
		if (res.status === 401) {
			throw redirect(302, '/auth');
		}
		if (!res.ok) {
			throw new Error(`Failed to fetch bookings: ${res.status} ${res.statusText}`);
		}

		const data: APIBookingsResponse = await res.json();

		const consultations: Consultation[] = data.bookings.map((b): Consultation => {
			const start = new Date(b.slot_start_time);
			const end = new Date(b.slot_end_time);
			const durationMin = Math.round((end.getTime() - start.getTime()) / 60000);

			return {
				id: b.id,
				clientName: b.client_name,
				therapistName: b.partner_name,
				productName: b.product_name,
				date: b.slot_start_time,
				duration: durationMin,
				status: mapStatus(b.status),
				hasNotes: false // notes not exposed in admin list for privacy
			};
		});

		return {
			consultations,
			total: data.total,
			page: data.page,
			limit: data.limit
		};
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading consultations:', err);
		throw error(503, 'Impossible de charger les consultations. Veuillez réessayer.');
	}
};
