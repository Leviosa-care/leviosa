import type { Actions, PageServerLoad } from './$types';
import { redirect, error, fail } from '@sveltejs/kit';
import { arktype } from 'sveltekit-superforms/adapters';
import { superValidate } from 'sveltekit-superforms';
import { env } from '$env/dynamic/private';
import { productSchema, productDefaults } from '../../schemas';
import { mockAdminCards, mockAdminCategories, type CardType, type Category } from "$lib/data/mockData";

interface BackendCategory {
	id: string;
	name: string;
	description: string;
	status: string;
}

// GET /admin/categories returns CategoryWithImage: { category, image? }
interface BackendCategoryWithImage {
	category: BackendCategory;
}

// GET /products/{id} returns ProductAggregator: { product, prices?, images? }
interface BackendProductAggregator {
	product: {
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
	};
	prices?: {
		amount: number;
		currency: string;
		interval: string;
		isActive: boolean;
	}[];
	images?: {
		id: string;
		s3_key: string;
	};
}

function formatPriceFromCents(amount: number): string {
	return (amount / 100).toFixed(2);
}

function mapBackendCategoryToFrontend(cat: BackendCategory): Category {
	return {
		id: cat.id,
		name: cat.name,
		description: cat.description,
		status: cat.status as 'published' | 'draft' | 'archived',
	};
}

function mapBackendCategoryWithImageToFrontend(item: BackendCategoryWithImage): Category {
	return mapBackendCategoryToFrontend(item.category);
}

function mapBackendProductToFrontend(agg: BackendProductAggregator): CardType {
	const prod = agg.product;

	// Extract the first active one_time price
	const activePrice = (agg.prices ?? []).find(p => p.isActive && p.interval === 'one_time');
	const price = activePrice ? formatPriceFromCents(activePrice.amount) : '0.00';

	// Extract the active image S3 key
	const image = agg.images?.s3_key ?? '';

	return {
		id: prod.id,
		name: prod.name,
		price,
		category: prod.category.name,
		description: prod.description,
		duration: prod.duration,
		image,
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
		const found = mockAdminCards.find(c => c.id === productId);
		if (!found) throw error(404, 'Product not found');
		product = found;
		categories = mockAdminCategories;
	} else {
		try {
			const [productRes, categoriesRes] = await Promise.all([
				fetch(`${env.API_URL}/admin/products/${productId}`),
				fetch(`${env.API_URL}/admin/categories`),
			]);

			if (productRes.status === 404) throw error(404, 'Product not found');
			if (!productRes.ok) throw new Error(`Failed to fetch product: ${productRes.status} ${productRes.statusText}`);
			const backendProductAgg: BackendProductAggregator = await productRes.json();
			product = mapBackendProductToFrontend(backendProductAgg);

			if (!categoriesRes.ok) throw new Error(`Failed to fetch categories: ${categoriesRes.status} ${categoriesRes.statusText}`);
			const backendCategories: BackendCategoryWithImage[] = await categoriesRes.json();
			categories = backendCategories.map(mapBackendCategoryWithImageToFrontend);
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
			category: categoryId,
			duration: product.duration,
			price: product.price,
			updatedAt: product.updatedAt,
			published: product.published,
			availability: product.availability,
			bufferTime: product.bufferTime,
			cancellationHours: product.cancellationHours,
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
		const productName = formData.get('productName') as string | null;

		if (!imageFile || !(imageFile instanceof File)) {
			return fail(400, { error: 'Fichier image invalide' });
		}

		const uploadFormData = new FormData();
		uploadFormData.append('image', imageFile);
		uploadFormData.append('parent_id', productId);
		uploadFormData.append('parent_type', 'product');
		uploadFormData.append('title', productName ?? productId);
		uploadFormData.append('is_active', 'true');

		const res = await fetch(`${env.API_URL}/admin/images`, {
			method: 'POST',
			body: uploadFormData,
		});

		if (!res.ok) {
			if (res.status === 401) throw redirect(303, '/auth');
			const errorText = await res.text();
			console.error('Image upload failed:', res.status, errorText);
			return fail(res.status, { error: `Le téléchargement a échoué (${res.status})` });
		}

		// Endpoint returns { "id": "..." } — reload page to pick up new active image
		throw redirect(303, `/admin/catalog/products/${productId}`);
	},
	default: async ({ request, params, fetch }) => {
		const productId = params.id;
		const formData = await request.formData();

		const name = formData.get('name') as string | null;
		const description = formData.get('description') as string | null;
		const categoryId = formData.get('categoryId') as string | null;
		const durationStr = formData.get('duration') as string | null;
		const status = formData.get('status') as string | null;
		const availability = formData.get('availability') as string | null;
		const bufferTimeStr = formData.get('bufferTime') as string | null;
		const cancellationHoursStr = formData.get('cancellationHours') as string | null;

		// Build the PATCH body — only include fields that are present
		const body: Record<string, unknown> = {};
		if (name != null) body.name = name;
		if (description != null) body.description = description;
		if (categoryId != null) body.category = categoryId;
		if (durationStr != null) body.duration = parseInt(durationStr, 10);
		if (status != null) body.publishedStatus = status;
		if (availability != null) body.availability = availability;
		if (bufferTimeStr != null) body.bufferTime = parseInt(bufferTimeStr, 10);
		if (cancellationHoursStr != null) body.cancellationHours = parseInt(cancellationHoursStr, 10);

		const res = await fetch(`${env.API_URL}/admin/products/${productId}`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body),
		});

		if (res.status === 401) throw redirect(303, '/auth');

		if (!res.ok) {
			const errorText = await res.text();
			console.error('Product update failed:', res.status, errorText);
			return fail(res.status, { error: `La mise à jour a échoué (${res.status})` });
		}

		throw redirect(303, '/admin/catalog?success=true');
	},
};
