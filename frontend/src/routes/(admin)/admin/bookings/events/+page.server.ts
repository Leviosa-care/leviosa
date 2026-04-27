import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

interface Event {
	id: string;
	name: string;
	description: string;
	date: string;
	duration: number;
	capacity: number;
	registered: number;
	priceInCents: number;
	status: 'upcoming' | 'ongoing' | 'completed' | 'cancelled';
	location: string;
}

async function getMockEvents(): Promise<Event[]> {
	const now = new Date();

	return [
		{ id: '1', name: 'Atelier Respiration', description: 'Apprenez les techniques de respiration pour réduire le stress', date: new Date(now.getTime() + 7 * 24 * 60 * 60 * 1000).toISOString(), duration: 120, capacity: 15, registered: 12, priceInCents: 3500, status: 'upcoming', location: 'Salle A' },
		{ id: '2', name: 'Journée Portes Ouvertes', description: 'Venez découvrir nos nouveaux services et rencontrer nos thérapeutes', date: new Date(now.getTime() + 14 * 24 * 60 * 60 * 1000).toISOString(), duration: 240, capacity: 50, registered: 35, priceInCents: 0, status: 'upcoming', location: 'Toute la clinique' },
		{ id: '3', name: 'Conférence Nutrition', description: 'Les bases d\'une alimentation équilibrée pour la santé', date: new Date(now.getTime() + 21 * 24 * 60 * 60 * 1000).toISOString(), duration: 90, capacity: 20, registered: 8, priceInCents: 2000, status: 'upcoming', location: 'Salle B' },
		{ id: '4', name: 'Session de Yoga', description: 'Séance de yoga douce pour tous les niveaux', date: new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000).toISOString(), duration: 60, capacity: 10, registered: 10, priceInCents: 1500, status: 'completed', location: 'Salle A' },
		{ id: '5', name: 'Atelier Meditation', description: 'Introduction à la méditation pleine conscience', date: new Date(now.getTime() - 14 * 24 * 60 * 60 * 1000).toISOString(), duration: 90, capacity: 12, registered: 9, priceInCents: 2500, status: 'completed', location: 'Salle B' },
		{ id: '6', name: 'Weekend Bien-être', description: 'Un weekend complet de détente et de soins', date: new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000).toISOString(), duration: 720, capacity: 8, registered: 6, priceInCents: 25000, status: 'cancelled', location: 'Toute la clinique' }
	];
}

export const load: PageServerLoad = async () => {
	if (env.USE_MOCK_DATA === 'true') {
		return { events: await getMockEvents() };
	}
	return { events: [] };
};
