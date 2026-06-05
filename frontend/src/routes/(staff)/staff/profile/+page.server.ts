import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export interface PartnerProfile {
	id: string;
	bio: string;
	experience: string;
	occupation: string;
	quote: string;
	tags: string[];
	isVerified: boolean;
	stripeAccountStatus: string;
	stripeOnboardingComplete: boolean;
	categories: { id: string; name: string }[];
	products: { id: string; name: string }[];
	joinedAt: string;
}

export interface LinkedProviders {
	google: boolean;
	apple: boolean;
}

interface PartnerResponse {
	id: string;
	bio: string;
	experience: string;
	occupation: string;
	quote: string;
	tags: string[];
	created_at: string;
	category_ids: string[];
	product_ids: string[];
	stripe_account_status: string;
	stripe_onboarding_complete: boolean;
}

interface Category {
	id: string;
	name: string;
}

interface Product {
	id: string;
	name: string;
}

export const load: PageServerLoad = async ({ fetch, url }) => {
	const partnerRes = await fetch(`${env.API_URL}/partners/me`, {
		headers: {
			'Content-Type': 'application/json',
		},
	});

	if (partnerRes.status === 401) {
		throw redirect(302, '/auth');
	}

	if (!partnerRes.ok) {
		if (partnerRes.status === 500) {
			return {
				profile: null,
				linkedProviders: { google: false, apple: false },
				stripeCallback: null,
				error: 'Erreur serveur. Veuillez réessayer dans quelques instants.',
			};
		}
		throw redirect(302, '/auth');
	}

	const partner: PartnerResponse = await partnerRes.json();

	// Fetch categories, products and user data in parallel — all independent of each other.
	const [categoriesRes, productsRes, userRes] = await Promise.all([
		fetch(`${env.API_URL}/categories`, { headers: { 'Content-Type': 'application/json' } }),
		fetch(`${env.API_URL}/products`, { headers: { 'Content-Type': 'application/json' } }),
		fetch(`${env.API_URL}/users/me`, { headers: { 'Content-Type': 'application/json' } }),
	]);

	let allCategories: Category[] = [];
	if (categoriesRes.ok) {
		allCategories = await categoriesRes.json();
	}

	let allProducts: Product[] = [];
	if (productsRes.ok) {
		allProducts = await productsRes.json();
	}

	// Map category IDs to category objects with names
	const partnerCategories = allCategories
		.filter((cat) => partner.category_ids.includes(cat.id))
		.map((cat) => ({
			id: cat.id,
			name: cat.name,
		}));

	// Map product IDs to product objects with names
	const partnerProducts = allProducts
		.filter((prod) => partner.product_ids.includes(prod.id))
		.map((prod) => ({
			id: prod.id,
			name: prod.name,
		}));

	// Transform to frontend's expected shape
	const profile: PartnerProfile = {
		id: partner.id,
		bio: partner.bio || '',
		experience: partner.experience || '',
		occupation: partner.occupation || '',
		quote: partner.quote || '',
		tags: partner.tags || [],
		isVerified: partner.stripe_account_status === 'active' && partner.stripe_onboarding_complete,
		stripeAccountStatus: partner.stripe_account_status,
		stripeOnboardingComplete: partner.stripe_onboarding_complete,
		categories: partnerCategories,
		products: partnerProducts,
		joinedAt: partner.created_at,
	};

	// Derive OAuth linking status from the user fetch done above in Promise.all.
	let linkedProviders: LinkedProviders = { google: false, apple: false };
	if (userRes.ok) {
		const userData = await userRes.json();
		linkedProviders = {
			google: !!userData.google_id,
			apple: !!userData.apple_id,
		};
	}

	// Derive the Stripe callback status from query parameter
	const stripeCallback = url.searchParams.get('stripe') || null;

	return {
		profile,
		linkedProviders,
		stripeCallback,
	};
};
