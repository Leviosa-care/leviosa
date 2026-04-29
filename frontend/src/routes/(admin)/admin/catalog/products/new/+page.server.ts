import type { Actions, PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { arktype } from 'sveltekit-superforms/adapters';
import { superValidate } from 'sveltekit-superforms';
import { env } from '$env/dynamic/private';
import { productSchema, productDefaults } from '../../schemas';
import { categories as mockCategories, type Category } from '../../products';

interface BackendCategory {
	id: string;
	name: string;
	description: string;
	status: string;
}

function mapBackendCategoryToFrontend(cat: BackendCategory): Category {
	return {
		id: cat.id,
		name: cat.name,
		description: cat.description,
		status: cat.status as 'published' | 'draft' | 'archived',
	};
}

export const load: PageServerLoad = async ({ fetch }) => {
	const createProductForm = await superValidate(arktype(productSchema, { defaults: productDefaults }));

	let categories: Category[];

	if (env.USE_MOCK_DATA === 'true') {
		categories = mockCategories;
	} else {
		try {
			const res = await fetch(`${env.API_URL}/admin/categories`);
			if (!res.ok) throw new Error(`Failed to fetch categories: ${res.status} ${res.statusText}`);
			const backendCategories: BackendCategory[] = await res.json();
			categories = backendCategories.map(mapBackendCategoryToFrontend);
		} catch (error) {
			console.error('Error loading categories:', error);
			categories = [];
		}
	}

	return {
		categories: [{ id: "default", name: "Toutes les catégories" }, ...categories],
		createProductForm,
	};
};

export const actions: Actions = {
	default: async ({ request }) => {
		const formData = await request.formData();
		const name = formData.get('name');
		const description = formData.get('description');
		const categoryId = formData.get('categoryId');
		const duration = formData.get('duration');
		const price = formData.get('price');
		const status = formData.get('status');
		const availability = formData.get('availability');
		const bufferTime = formData.get('bufferTime');
		const cancellationHours = formData.get('cancellationHours');
		const imageUrl = formData.get('imageUrl');
		const stripeProductId = formData.get('stripeProductId');

		// TODO: Connect to backend API
		console.log('Creating product:', {
			name,
			description,
			categoryId,
			duration,
			price,
			status,
			availability,
			bufferTime,
			cancellationHours,
			imageUrl,
			stripeProductId
		});

		// For now, just redirect back to catalog
		throw redirect(303, '/admin/catalog?success=true');
	},
};
