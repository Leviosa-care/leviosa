import type { PageServerLoad } from './$types';

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

function getMockBookings(): Booking[] {
	const now = new Date();

	function booking(
		id: string,
		client: string,
		initials: string,
		product: string,
		offsetHours: number,
		durationMin: number,
		room: string,
		status: Booking['status'],
		cents: number,
		notes?: string,
	): Booking {
		const start = new Date(now.getTime() + offsetHours * 60 * 60 * 1000);
		const end = new Date(start.getTime() + durationMin * 60 * 1000);
		return {
			id,
			clientName: client,
			clientInitials: initials,
			productName: product,
			startTime: start.toISOString(),
			endTime: end.toISOString(),
			roomName: room,
			status,
			amountInCents: cents,
			notes,
		};
	}

	return [
		booking('b1', 'Marie Dupont', 'MD', 'Massage Relaxant 60min', 2, 60, 'Salle Sérénité', 'upcoming', 7500),
		booking('b2', 'Jean Durand', 'JD', 'Drainage Lymphatique 90min', 5, 90, 'Salle Harmonie', 'upcoming', 11000),
		booking('b3', 'Claire Bernard', 'CB', 'Massage Relaxant 60min', 8, 60, 'Salle Sérénité', 'upcoming', 7500),
		booking('b4', 'Lucas Petit', 'LP', 'Soin du Dos 60min', 24, 60, 'Salle Harmonie', 'upcoming', 9000),
		booking('b5', 'Emma Moreau', 'EM', 'Massage Sportif 60min', 26, 60, 'Salle Vitalité', 'upcoming', 8500),
		booking('b6', 'Thomas Richard', 'TR', 'Drainage Lymphatique 90min', -3, 90, 'Salle Sérénité', 'completed', 11000, 'Bonne séance, à revoir dans 3 semaines.'),
		booking('b7', 'Camille Simon', 'CS', 'Massage Relaxant 60min', -26, 60, 'Salle Harmonie', 'completed', 7500, 'Tension dans le dos et les épaules. Protocole standard.'),
		booking('b8', 'Hugo Michel', 'HM', 'Soin du Dos 60min', -50, 60, 'Salle Harmonie', 'no_show', 9000),
		booking('b9', 'Léa Fontaine', 'LF', 'Massage Relaxant 60min', -74, 60, 'Salle Sérénité', 'completed', 7500),
		booking('b10', 'Antoine Garnier', 'AG', 'Massage Sportif 60min', -98, 60, 'Salle Vitalité', 'cancelled', 8500),
	];
}

export const load: PageServerLoad = async () => {
	// TODO: Replace with GET /partners/{partnerId}/bookings
	const bookings = getMockBookings();
	return { bookings };
};
