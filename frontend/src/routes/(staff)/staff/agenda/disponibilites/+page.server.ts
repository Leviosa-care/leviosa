import type { PageServerLoad } from './$types';

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

function getMockAvailabilities(): AvailabilityDay[] {
	const today = new Date();
	today.setHours(0, 0, 0, 0);

	function makeSlot(
		dayOffset: number,
		hour: number,
		durationMin: number,
		status: Availability['status'],
		room: string,
		clientName?: string,
		productName?: string,
	): Availability {
		const d = new Date(today);
		d.setDate(d.getDate() + dayOffset);
		const start = new Date(d);
		start.setHours(hour, 0, 0, 0);
		const end = new Date(start.getTime() + durationMin * 60 * 1000);
		return {
			id: `slot-${dayOffset}-${hour}`,
			date: d.toISOString().split('T')[0],
			startTime: start.toISOString(),
			endTime: end.toISOString(),
			status,
			roomName: room,
			clientName,
			productName,
		};
	}

	const days: AvailabilityDay[] = [];
	const dayNames = ['Dimanche', 'Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi'];

	for (let i = 0; i < 7; i++) {
		const d = new Date(today);
		d.setDate(d.getDate() + i);
		const dayName = dayNames[d.getDay()];
		const slots: Availability[] = [];

		if (d.getDay() !== 0) {
			// Skip Sundays
			slots.push(makeSlot(i, 9, 60, i === 0 ? 'booked' : 'available', 'Salle Sérénité', i === 0 ? 'Marie Dupont' : undefined, i === 0 ? 'Massage Relaxant 60min' : undefined));
			slots.push(makeSlot(i, 10, 90, i <= 1 ? 'booked' : 'available', 'Salle Sérénité', i <= 1 ? 'Jean Durand' : undefined, i <= 1 ? 'Drainage Lymphatique 90min' : undefined));
			slots.push(makeSlot(i, 14, 60, i === 2 ? 'cancelled' : 'available', 'Salle Harmonie', i === 2 ? 'Claire Bernard' : undefined, i === 2 ? 'Massage Relaxant 60min' : undefined));
			slots.push(makeSlot(i, 15, 60, 'available', 'Salle Harmonie'));
			slots.push(makeSlot(i, 17, 60, 'available', 'Salle Sérénité'));
		}

		if (slots.length > 0) {
			days.push({ date: d.toISOString().split('T')[0], dayName, slots });
		}
	}

	return days;
}

export const load: PageServerLoad = async () => {
	// TODO: Replace with GET /partners/{partnerId}/availabilities
	return {
		availabilities: getMockAvailabilities(),
	};
};
