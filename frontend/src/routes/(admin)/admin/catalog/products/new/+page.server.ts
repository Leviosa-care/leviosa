import type { Actions, PageServerLoad } from './$types';
import { fail } from '@sveltejs/kit';
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

interface CreateProductRequest {
	name: string;
	description: string;
	categoryId: string;
	duration: number;
	price: number;
	status: 'published' | 'draft' | 'archived';
	availability: 'online' | 'in-person' | 'hybrid';
	bufferTime: number;
	cancellationHours: number;
	stripeProductId?: string;
}

interface CreateProductResponse {
	id: string;
	name: string;
	description: string;
	categoryId: string;
	duration: number;
	price: number;
	status: string;
	availability: string;
	bufferTime: number;
	cancellationHours: number;
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
		} catch (err) {
			console.error('Error loading categories:', err);
			categories = [];
		}
	}

	return {
		categories: [{ id: "default", name: "Toutes les catégories" }, ...categories],
		createProductForm,
	};
};

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		const formData = await request.formData();

		const name = formData.get('name') as string;
		const description = formData.get('description') as string;
		const categoryId = formData.get('categoryId') as string;
		const duration = formData.get('duration') as string;
		const price = formData.get('price') as string;
		const status = formData.get('status') as 'published' | 'draft' | 'archived';
		const availability = formData.get('availability') as 'online' | 'in-person' | 'hybrid';
		const bufferTime = formData.get('bufferTime') as string;
		const cancellationHours = formData.get('cancellationHours') as string;
		const stripeProductId = formData.get('stripeProductId') as string | null;

		const requestBody: CreateProductRequest = {
			name,
			description,
			categoryId,
			duration: parseInt(duration, 10),
			price: parseFloat(price),
			status,
			availability,
			bufferTime: parseInt(bufferTime, 10),
			cancellationHours: parseInt(cancellationHours, 10),
		};

		if (stripeProductId) {
			requestBody.stripeProductId = stripeProductId;
		}

		try {
			const res = await fetch(`${env.API_URL}/admin/products`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify(requestBody),
			});

			if (!res.ok) {
				if (res.status === 401) return fail(401, { error: 'Session expirée. Veuillez vous reconnecter.' });
				if (res.status === 403) return fail(403, { error: "Vous n'avez pas les droits pour effectuer cette action." });
				if (res.status === 404) return fail(404, { error: "La ressource demandée est introuvable." });
				if (res.status === 409) return fail(409, { error: "Cette ressource est déjà utilisée." });
				if (res.status === 500) return fail(500, { error: "Une erreur serveur est survenue. Veuillez réessayer." });
				const errorData = await res.json().catch(() => ({ msg: 'Erreur inconnue' }));
				return fail(400, { error: `Erreur: ${errorData.msg || res.status}` });
			}

			const product: CreateProductResponse = await res.json();

			const imageFile = formData.get('imageFile') as File | null;
			if (imageFile && imageFile.size > 0) {
				const imageFormData = new FormData();
				imageFormData.append('image', imageFile);

				const imageRes = await fetch(`${env.API_URL}/admin/products/${product.id}/images`, {
					method: 'POST',
					body: imageFormData,
				});

				if (!imageRes.ok) {
					console.error('Image upload failed after product creation:', imageRes.status);
				}
			}

			return { success: true };
		} catch (e) {
			console.error('Product creation failed:', e);
			return fail(500, { error: 'Une erreur est survenue lors de la création du produit.' });
		}
	},
};
