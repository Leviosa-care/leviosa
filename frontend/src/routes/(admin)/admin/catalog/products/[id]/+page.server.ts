import type { Actions, PageServerLoad } from './$types';
import { redirect, error } from '@sveltejs/kit';
import { arktype } from 'sveltekit-superforms/adapters';
import { superValidate } from 'sveltekit-superforms';
import { productSchema, productDefaults } from '../../schemas';
import type { product } from '../../schemas';
import { cards, type CardType, categories, type Category } from '../../products';

export const load: PageServerLoad = async ({ params }) => {
	const productId = params.id;

	// Find product by ID
	const product = cards.find(c => c.id === productId);

	if (!product) {
		throw error(404, 'Product not found');
	}

	// Get the category ID for this product
	const category = categories.find(c => c.name === product.category);
	const categoryId = category?.id || '';

	// Create form with existing product data
	const updateProductForm = await superValidate(
		{
			id: product.id,
			name: product.name,
			description: product.description,
			categoryId: categoryId,
			duration: product.duration,
			price: parseFloat(product.price),
			status: product.published,
			availability: product.availability,
			bufferTime: product.bufferTime,
			cancellationHours: product.cancellationHours,
			imageUrl: product.image
		},
		arktype(productSchema, { defaults: productDefaults })
	);

	// Combine default category with actual categories
	const allCategories: Category[] = [
		{ id: "default", name: "Toutes les catégories" },
		...categories
	];

	return {
		product,
		categories: allCategories,
		updateProductForm,
	};
};

export const actions: Actions = {
	default: async ({ request, params }) => {
		const productId = params.id;
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

		// TODO: Connect to backend API
		console.log('Updating product:', {
			id: productId,
			name,
			description,
			categoryId,
			duration,
			price,
			status,
			availability,
			bufferTime,
			cancellationHours,
			imageUrl
		});

		// For now, just redirect back to catalog
		throw redirect(303, '/admin/catalog?success=true');
	},
};
