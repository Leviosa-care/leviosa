import type { Actions, PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { arktype } from 'sveltekit-superforms/adapters';
import { superValidate } from 'sveltekit-superforms';
import { productSchema, productDefaults } from '../../schemas';
import type { product } from '../../schemas';
import { categories } from '../../products';
import type { Category } from '../../products';

export const load: PageServerLoad = async () => {
	const createProductForm = await superValidate(arktype(productSchema, { defaults: productDefaults }));

	// Combine default category with actual categories
	const allCategories: Category[] = [
		{ id: "default", name: "Toutes les catégories" },
		...categories
	];

	return {
		categories: allCategories,
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
