import { env } from '$env/dynamic/private';
import { error, redirect, isRedirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

interface WeeklyEvent {
	id: string;
	clientName: string;
	therapistName: string;
	productName: string;
	startTime: string;
	endTime: string;
	roomName: string;
	status: 'confirmed' | 'pending' | 'cancelled' | 'completed' | 'no_show';
}

interface DayEvents {
	date: string;
	dayName: string;
	events: WeeklyEvent[];
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

function getMonday(d: Date): Date {
	const date = new Date(d);
	const day = date.getUTCDay();
	const diff = date.getUTCDate() - day + (day === 0 ? -6 : 1);
	date.setUTCDate(diff);
	date.setUTCHours(0, 0, 0, 0);
	return date;
}

function getSunday(d: Date): Date {
	const monday = getMonday(d);
	const sunday = new Date(monday);
	sunday.setUTCDate(monday.getUTCDate() + 6);
	sunday.setUTCHours(23, 59, 59, 999);
	return sunday;
}

function formatDateISO(d: Date): string {
	return d.toISOString().split('T')[0];
}

const dayNames = ['Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi', 'Dimanche'];

function groupByDay(bookings: APIBooking[], monday: Date): DayEvents[] {
	const days: DayEvents[] = [];

	for (let i = 0; i < 7; i++) {
		const dayDate = new Date(monday);
		dayDate.setUTCDate(monday.getUTCDate() + i);
		const dateStr = formatDateISO(dayDate);

		const dayEvents: WeeklyEvent[] = bookings
			.filter(b => b.slot_start_time.startsWith(dateStr))
			.map(b => ({
				id: b.id,
				clientName: b.client_name,
				therapistName: b.partner_name,
				productName: b.product_name,
				startTime: b.slot_start_time,
				endTime: b.slot_end_time,
				roomName: b.room_name,
				status: b.status as WeeklyEvent['status']
			}));

		days.push({
			date: dateStr,
			dayName: dayNames[i],
			events: dayEvents
		});
	}

	return days;
}

async function getMockWeeklyEvents(): Promise<{ weekEvents: DayEvents[]; weekStart: string }> {
	const now = new Date();
	const monday = getMonday(now);

	const events: APIBooking[] = [
		{ id: '1', client_name: 'Marie Dupont', partner_name: 'Sophie Martin', product_name: 'Massage Relaxant 60min', room_name: 'Cabinet 1', slot_start_time: new Date(monday.getTime() + 9 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 10 * 60 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '2', client_name: 'Jean Durand', partner_name: 'Pierre Leroy', product_name: 'Consultation Kiné 45min', room_name: 'Cabinet 2', slot_start_time: new Date(monday.getTime() + 10 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 10 * 60 * 60 * 1000 + 45 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 4500, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '3', client_name: 'Claire Bernard', partner_name: 'Marie Dubois', product_name: 'Soin du Dos 60min', room_name: 'Cabinet 1', slot_start_time: new Date(monday.getTime() + 14 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 15 * 60 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'pending', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '4', client_name: 'Lucas Petit', partner_name: 'Sophie Martin', product_name: 'Massage Relaxant 90min', room_name: 'Cabinet 3', slot_start_time: new Date(monday.getTime() + 16 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 17 * 60 * 60 * 1000 + 30 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 9000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '5', client_name: 'Emma Moreau', partner_name: 'Pierre Leroy', product_name: 'Drainage Lymphatique', room_name: 'Cabinet 2', slot_start_time: new Date(monday.getTime() + 17 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 18 * 60 * 60 * 1000).toISOString(), status: 'cancelled', payment_status: 'refunded', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '6', client_name: 'Thomas Richard', partner_name: 'Sophie Martin', product_name: 'Massage Relaxant 60min', room_name: 'Cabinet 1', slot_start_time: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 9 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '7', client_name: 'Camille Simon', partner_name: 'Marie Dubois', product_name: 'Soin du Dos 60min', room_name: 'Cabinet 2', slot_start_time: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 12 * 60 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '8', client_name: 'Hugo Michel', partner_name: 'Pierre Leroy', product_name: 'Consultation Kiné 45min', room_name: 'Cabinet 1', slot_start_time: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 15 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 15 * 60 * 60 * 1000 + 45 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'pending', total_price_cents: 4500, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '9', client_name: 'Chloe Garcia', partner_name: 'Sophie Martin', product_name: 'Massage Relaxant 60min', room_name: 'Cabinet 3', slot_start_time: new Date(monday.getTime() + 48 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 48 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '10', client_name: 'Louis Laurent', partner_name: 'Marie Dubois', product_name: 'Soin du Dos 60min', room_name: 'Cabinet 1', slot_start_time: new Date(monday.getTime() + 72 * 60 * 60 * 1000 + 9 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 72 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '11', client_name: 'Jade Roux', partner_name: 'Pierre Leroy', product_name: 'Drainage Lymphatique', room_name: 'Cabinet 2', slot_start_time: new Date(monday.getTime() + 72 * 60 * 60 * 1000 + 14 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 72 * 60 * 60 * 1000 + 15 * 60 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '12', client_name: 'Nathan Girard', partner_name: 'Sophie Martin', product_name: 'Massage Relaxant 90min', room_name: 'Cabinet 3', slot_start_time: new Date(monday.getTime() + 96 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 96 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000 + 30 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'pending', total_price_cents: 9000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '13', client_name: 'Lola Bernard', partner_name: 'Marie Dubois', product_name: 'Soin du Dos 60min', room_name: 'Cabinet 1', slot_start_time: new Date(monday.getTime() + 120 * 60 * 60 * 1000 + 9 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 120 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '14', client_name: 'Enzo Dubois', partner_name: 'Pierre Leroy', product_name: 'Consultation Kiné 45min', room_name: 'Cabinet 2', slot_start_time: new Date(monday.getTime() + 120 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 120 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000 + 45 * 60 * 1000).toISOString(), status: 'confirmed', payment_status: 'paid', total_price_cents: 4500, currency: 'EUR', created_at: new Date().toISOString() },
		{ id: '15', client_name: 'Sarah Richard', partner_name: 'Sophie Martin', product_name: 'Massage Relaxant 60min', room_name: 'Cabinet 3', slot_start_time: new Date(monday.getTime() + 144 * 60 * 60 * 1000 + 14 * 60 * 60 * 1000).toISOString(), slot_end_time: new Date(monday.getTime() + 144 * 60 * 60 * 1000 + 15 * 60 * 60 * 1000).toISOString(), status: 'cancelled', payment_status: 'refunded', total_price_cents: 6000, currency: 'EUR', created_at: new Date().toISOString() },
	];

	return {
		weekEvents: groupByDay(events, monday),
		weekStart: formatDateISO(monday)
	};
}

export const load: PageServerLoad = async ({ fetch, url }) => {
	if (env.USE_MOCK_DATA === 'true') {
		return await getMockWeeklyEvents();
	}

	try {
		// Determine the week from the `week` query param (YYYY-MM-DD of Monday)
		const weekParam = url.searchParams.get('week');
		let monday: Date;
		if (weekParam) {
			monday = new Date(weekParam + 'T00:00:00Z');
			if (isNaN(monday.getTime())) {
				monday = getMonday(new Date());
			}
		} else {
			monday = getMonday(new Date());
		}
		const sunday = getSunday(monday);

		const params = new URLSearchParams({
			from: formatDateISO(monday),
			to: formatDateISO(sunday),
			limit: '500'
		});

		const res = await fetch(`${env.API_URL}/admin/bookings?${params.toString()}`);
		if (res.status === 401) {
			throw redirect(302, '/auth');
		}
		if (!res.ok) {
			throw new Error(`Failed to fetch bookings: ${res.status} ${res.statusText}`);
		}

		const data: APIBookingsResponse = await res.json();

		return {
			weekEvents: groupByDay(data.bookings, monday),
			weekStart: formatDateISO(monday)
		};
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading planning:', err);
		throw error(503, 'Impossible de charger le planning. Veuillez réessayer.');
	}
};
