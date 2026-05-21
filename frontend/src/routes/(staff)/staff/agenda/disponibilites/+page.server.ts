import type { PageServerLoad } from './$types';
import { redirect, isRedirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export interface Availability {
	id: string;
	date: string;
	startTime: string;
	endTime: string;
	status: 'available' | 'booked' | 'cancelled';
	roomName: string;
	clientName?: string;
	productName?: string;
}

export interface AvailabilityDay {
	date: string;
	dayName: string;
	slots: Availability[];
}

export const load: PageServerLoad = async ({ locals, url, fetch }) => {
	if (!locals.user?.id) {
		throw redirect(302, '/auth');
	}

	const partnerId = locals.user.id;

	try {
		const res = await fetch(`${env.API_URL}/partners/${partnerId}/availabilities`, {
			headers: {
				'Content-Type': 'application/json',
			},
		});

		if (res.status === 401) {
			throw redirect(302, '/auth');
		}

		if (!res.ok) {
			console.error(`Failed to fetch availabilities: ${res.status} ${res.statusText}`);
			return { availabilities: [] };
		}

		const data = await res.json();

		if (!Array.isArray(data)) {
			console.error('Invalid response format: expected array');
			return { availabilities: [] };
		}

		const dayNames = ['Dimanche', 'Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi'];
		const today = new Date();
		today.setHours(0, 0, 0, 0);

		const slotsByDate = new Map<string, Availability[]>();

		for (const slot of data) {
			const startDate = new Date(slot.start_time);
			const dateKey = startDate.toISOString().split('T')[0];

			const availability: Availability = {
				id: slot.id,
				date: dateKey,
				startTime: slot.start_time,
				endTime: slot.end_time,
				status: slot.status,
				roomName: slot.room_id || 'Salle inconnue',
				clientName: undefined,
				productName: undefined,
			};

			if (!slotsByDate.has(dateKey)) {
				slotsByDate.set(dateKey, []);
			}
			slotsByDate.get(dateKey)!.push(availability);
		}

		const availabilities: AvailabilityDay[] = [];
		for (let i = 0; i < 7; i++) {
			const d = new Date(today.getTime());
			d.setDate(d.getDate() + i);
			const dateKey = d.toISOString().split('T')[0];
			const dayName = dayNames[d.getDay()];

			const slots = slotsByDate.get(dateKey) || [];

			if (slots.length > 0 || d.getDay() !== 0) {
				availabilities.push({
					date: dateKey,
					dayName,
					slots,
				});
			}
		}

		return { availabilities };
	} catch (err) {
		if (isRedirect(err)) throw err;
		console.error('Error loading availabilities:', err);
		return { availabilities: [] };
	}
};
