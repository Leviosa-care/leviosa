import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

interface Consultation {
	id: string;
	clientName: string;
	therapistName: string;
	productName: string;
	date: string;
	duration: number;
	status: 'completed' | 'pending' | 'cancelled';
	hasNotes: boolean;
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

export const load: PageServerLoad = async () => {
	if (env.USE_MOCK_DATA === 'true') {
		return { consultations: await getMockConsultations() };
	}
	return { consultations: [] };
};
