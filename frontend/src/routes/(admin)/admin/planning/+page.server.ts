import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

interface WeeklyEvent {
	id: string;
	clientName: string;
	therapistName: string;
	productName: string;
	startTime: string;
	endTime: string;
	roomName: string;
	status: 'confirmed' | 'pending' | 'cancelled';
}

interface DayEvents {
	date: string;
	dayName: string;
	events: WeeklyEvent[];
}

async function getMockWeeklyEvents(): Promise<DayEvents[]> {
	const now = new Date();
	const currentDay = now.getDay();
	const monday = new Date(now);
	const diff = now.getDate() - currentDay + (currentDay === 0 ? -6 : 1);
	monday.setDate(diff);
	monday.setHours(0, 0, 0, 0);

	const events: WeeklyEvent[] = [
		{ id: '1', clientName: 'Marie Dupont', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 60min', startTime: new Date(monday.getTime() + 9 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 10 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 1', status: 'confirmed' },
		{ id: '2', clientName: 'Jean Durand', therapistName: 'Pierre Leroy', productName: 'Consultation Kiné 45min', startTime: new Date(monday.getTime() + 10 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 10 * 60 * 60 * 1000 + 45 * 60 * 1000).toISOString(), roomName: 'Cabinet 2', status: 'confirmed' },
		{ id: '3', clientName: 'Claire Bernard', therapistName: 'Marie Dubois', productName: 'Soin du Dos 60min', startTime: new Date(monday.getTime() + 14 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 15 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 1', status: 'pending' },
		{ id: '4', clientName: 'Lucas Petit', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 90min', startTime: new Date(monday.getTime() + 16 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 17 * 60 * 60 * 1000 + 30 * 60 * 1000).toISOString(), roomName: 'Cabinet 3', status: 'confirmed' },
		{ id: '5', clientName: 'Emma Moreau', therapistName: 'Pierre Leroy', productName: 'Drainage Lymphatique', startTime: new Date(monday.getTime() + 17 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 18 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 2', status: 'cancelled' },
		{ id: '6', clientName: 'Thomas Richard', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 60min', startTime: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 9 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 1', status: 'confirmed' },
		{ id: '7', clientName: 'Camille Simon', therapistName: 'Marie Dubois', productName: 'Soin du Dos 60min', startTime: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 12 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 2', status: 'confirmed' },
		{ id: '8', clientName: 'Hugo Michel', therapistName: 'Pierre Leroy', productName: 'Consultation Kiné 45min', startTime: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 15 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 24 * 60 * 60 * 1000 + 15 * 60 * 60 * 1000 + 45 * 60 * 1000).toISOString(), roomName: 'Cabinet 1', status: 'pending' },
		{ id: '9', clientName: 'Chloe Garcia', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 60min', startTime: new Date(monday.getTime() + 48 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 48 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 3', status: 'confirmed' },
		{ id: '10', clientName: 'Louis Laurent', therapistName: 'Marie Dubois', productName: 'Soin du Dos 60min', startTime: new Date(monday.getTime() + 72 * 60 * 60 * 1000 + 9 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 72 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 1', status: 'confirmed' },
		{ id: '11', clientName: 'Jade Roux', therapistName: 'Pierre Leroy', productName: 'Drainage Lymphatique', startTime: new Date(monday.getTime() + 72 * 60 * 60 * 1000 + 14 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 72 * 60 * 60 * 1000 + 15 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 2', status: 'confirmed' },
		{ id: '12', clientName: 'Nathan Girard', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 90min', startTime: new Date(monday.getTime() + 96 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 96 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000 + 30 * 60 * 1000).toISOString(), roomName: 'Cabinet 3', status: 'pending' },
		{ id: '13', clientName: 'Lola Bernard', therapistName: 'Marie Dubois', productName: 'Soin du Dos 60min', startTime: new Date(monday.getTime() + 120 * 60 * 60 * 1000 + 9 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 120 * 60 * 60 * 1000 + 10 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 1', status: 'confirmed' },
		{ id: '14', clientName: 'Enzo Dubois', therapistName: 'Pierre Leroy', productName: 'Consultation Kiné 45min', startTime: new Date(monday.getTime() + 120 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 120 * 60 * 60 * 1000 + 11 * 60 * 60 * 1000 + 45 * 60 * 1000).toISOString(), roomName: 'Cabinet 2', status: 'confirmed' },
		{ id: '15', clientName: 'Sarah Richard', therapistName: 'Sophie Martin', productName: 'Massage Relaxant 60min', startTime: new Date(monday.getTime() + 144 * 60 * 60 * 1000 + 14 * 60 * 60 * 1000).toISOString(), endTime: new Date(monday.getTime() + 144 * 60 * 60 * 1000 + 15 * 60 * 60 * 1000).toISOString(), roomName: 'Cabinet 3', status: 'cancelled' },
	];

	const dayNames = ['Lundi', 'Mardi', 'Mercredi', 'Jeudi', 'Vendredi', 'Samedi', 'Dimanche'];
	const days: DayEvents[] = [];

	for (let i = 0; i < 7; i++) {
		const dayDate = new Date(monday.getTime() + i * 24 * 60 * 60 * 1000);
		const dateStr = dayDate.toISOString().split('T')[0];
		const dayEvents = events.filter(e => e.startTime.startsWith(dateStr));

		days.push({
			date: dateStr,
			dayName: dayNames[i],
			events: dayEvents
		});
	}

	return days;
}

export const load: PageServerLoad = async () => {
	const isDevelopment = env.NODE_ENV === 'development' || env.APP_ENV === 'development';

	if (isDevelopment) {
		return { weekEvents: await getMockWeeklyEvents() };
	}

	// In staging/production, return empty week structure
	// TODO: When backend is ready, fetch real events from API
	return { weekEvents: [] };
};
