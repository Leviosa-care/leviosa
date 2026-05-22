import type { Actions, PageServerLoad } from './$types';
import { redirect, error, fail } from '@sveltejs/kit';
import { arktype } from 'sveltekit-superforms/adapters';
import { superValidate } from 'sveltekit-superforms';
import { env } from '$env/dynamic/private';
import { productSchema, productDefaults } from '../../schemas';
import { cards as mockCards, categories as mockCategories, type CardType, type Category } from '../../products';

interface BackendCategory {
	id: string;
	name: string;
	description: string;
	status: string;
}

interface BackendProduct {
	id: string;
	name: string;
	description: string;
	category: BackendCategory;
	duration: number;
	createdAt: string;
	updatedAt: string;
	publishedStatus: string;
	availability: string;
	bufferTime: number;
	cancellationHours: number;
}

function mapBackendCategoryToFrontend(cat: BackendCategory): Category {
	return {
		id: cat.id,
		name: cat.name,
		description: cat.description,
		status: cat.status as 'published' | 'draft' | 'archived',
	};
}

function mapBackendProductToFrontend(prod: BackendProduct): CardType {
	return {
		id: prod.id,
		name: prod.name,
		price: "0.00",
		category: prod.category.name,
		description: prod.description,
		duration: prod.duration,
		image: "",
		updatedAt: prod.updatedAt,
		published: prod.publishedStatus as 'published' | 'draft' | 'archived',
		availability: prod.availability as 'online' | 'in-person' | 'hybrid',
		bufferTime: prod.bufferTime,
		cancellationHours: prod.cancellationHours,
	};
}

export const load: PageServerLoad = async ({ params, fetch }) => {
	const productId = params.id;

	let product: CardType;
	let categories: Category[];

	if (env.USE_MOCK_DATA === 'true') {
		const found = mockCards.find(c => c.id === productId);
		if (!found) throw error(404, 'Product not found');
		product = found;
		categories = mockCategories;
	} else {
		try {
			const [productRes, categoriesRes] = await Promise.all([
				fetch(`${env.API_URL}/admin/products/${productId}`),
				fetch(`${env.API_URL}/admin/categories`),
			]);

			if (productRes.status === 404) throw error(404, 'Product not found');
			if (!productRes.ok) throw new Error(`Failed to fetch product: ${productRes.status} ${productRes.statusText}`);
			const backendProduct: BackendProduct = await productRes.json();
			product = mapBackendProductToFrontend(backendProduct);

			if (!categoriesRes.ok) throw new Error(`Failed to fetch categories: ${categoriesRes.status} ${categoriesRes.statusText}`);
			const backendCategories: BackendCategory[] = await categoriesRes.json();
			categories = backendCategories.map(mapBackendCategoryToFrontend);
		} catch (err) {
			// Re-throw SvelteKit errors (like 404) directly
			if (err && typeof err === 'object' && 'status' in err) throw err;
			console.error('Error loading product:', err);
			throw error(500, 'Failed to load product');
		}
	}

	const category = categories.find(c => c.name === product.category);
	const categoryId = category?.id || '';

	const updateProductForm = await superValidate(
		{
			id: product.id,
			name: product.name,
			description: product.description,
			categoryId,
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

	return {
		product,
		categories: [{ id: "default", name: "Toutes les catégories" }, ...categories],
		updateProductForm,
	};
};

export const actions: Actions = {
	uploadImage: async ({ request, params, fetch }) => {
		const productId = params.id;
		const formData = await request.formData();
		const imageFile = formData.get('image');

		if (!imageFile || !(imageFile instanceof File)) {
			return { success: false, error: 'Invalid image file' };
		}

		const uploadFormData = new FormData();
		uploadFormData.append('image', imageFile);

		const res = await fetch(`${env.API_URL}/admin/products/${productId}/images`, {
			method: 'POST',
			body: uploadFormData,
		});

		if (!res.ok) {
			if (res.status === 401) throw redirect(303, '/auth');
			const errorText = await res.text();
			console.error('Image upload failed:', res.status, errorText);
			return fail(res.status, { error: `Upload failed: ${res.status}` });
		}

		const result = await res.json();
		return { success: true, url: result.url };
	},
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
