import type { PageServerLoad } from './$types';

export interface PartnerProfile {
	id: string;
	bio: string;
	experience: string;
	isVerified: boolean;
	stripeOnboardingComplete: boolean;
	categories: { id: string; name: string }[];
	products: { id: string; name: string }[];
	joinedAt: string;
}

function getMockProfile(): PartnerProfile {
	return {
		id: 'partner-1',
		bio: 'Praticienne spécialisée en massage bien-être et drainage lymphatique, diplômée de l\'Institut Supérieur de Bien-Être de Paris. Passionnée par les thérapies manuelles et le soin du corps, j\'accompagne mes clients vers un équilibre physique et mental durable.',
		experience:
			'Plus de 8 ans d\'expérience en massage thérapeutique et bien-être. Formée aux techniques suédoises, californienne, et drainage lymphatique selon la méthode Vodder. Certifiée en aromathérapie appliquée. Ancienne praticienne au spa de l\'hôtel Le Royal Monceau (2016–2020), puis en cabinet libéral depuis 2020. Formation continue annuelle en techniques corporelles avancées.',
		isVerified: true,
		stripeOnboardingComplete: true,
		categories: [
			{ id: 'cat-1', name: 'Massage Bien-Être' },
			{ id: 'cat-2', name: 'Drainage Lymphatique' },
			{ id: 'cat-3', name: 'Thérapies Corporelles' },
		],
		products: [
			{ id: 'prod-1', name: 'Massage Relaxant 60min' },
			{ id: 'prod-2', name: 'Massage Relaxant 90min' },
			{ id: 'prod-3', name: 'Drainage Lymphatique 90min' },
			{ id: 'prod-4', name: 'Soin du Dos 60min' },
		],
		joinedAt: new Date(Date.now() - 420 * 24 * 3600 * 1000).toISOString(),
	};
}

export const load: PageServerLoad = async ({ locals }) => {
	// TODO: Replace with GET /partners/me
	return {
		user: locals.user,
		profile: getMockProfile(),
	};
};
