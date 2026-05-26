import type { PageServerLoad } from './$types';
import type { Actions } from "./$types"
import { error } from '@sveltejs/kit';
import { arktype } from 'sveltekit-superforms/adapters';
import { superValidate, type SuperValidated } from 'sveltekit-superforms';
import { env } from "$env/dynamic/private";

import {
	deleteProductSchema,
	deleteProductDefaults,
	type DeleteProduct,
	productSchema,
	productDefaults,
	type product,
	categorySchema,
	categoryDefaults,
	type category,
} from './schemas'
import { type CardType, mockAdminCards, mockAdminCategories, type Category } from "$lib/data/mockData"
import {
	defaultStatus,
	defaultCategory,
	defaultAvailability,
} from "./default";

import { updateProduct, deleteProduct, createCategory } from "./actions"

// Backend returns CategoryWithImage: { category: Category, image?: Image }
interface BackendCategoryWithImage {
	category: {
		id: string;
		name: string;
		description: string;
		status: string;
		createdAt: string;
		updatedAt: string;
	};
	image?: {
		id: string;
		s3_key: string;
	};
}

// Backend returns ProductAggregator: { product: ProductRes, prices?: Price[], images?: Image }
interface BackendProductAggregator {
	product: {
		id: string;
		name: string;
		description: string;
		category: {
			id: string;
			name: string;
			description: string;
			status: string;
		};
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

function mapBackendCategoryToFrontend(item: BackendCategoryWithImage): Category {
	const cat = item.category;
	return {
		id: cat.id,
		name: cat.name,
		description: cat.description,
		status: cat.status as 'published' | 'draft' | 'archived',
	};
}

function mapBackendProductToFrontend(agg: BackendProductAggregator): CardType {
	const prod = agg.product;

	// Extract the first active one_time price
	const activePrice = (agg.prices ?? []).find(p => p.isActive && p.interval === 'one_time');
	const price = activePrice ? formatPriceFromCents(activePrice.amount) : '0.00';

	return {
		id: prod.id,
		name: prod.name,
		price,
		category: prod.category.name,
		description: prod.description,
		duration: prod.duration,
		image: '', // Images require presigned S3 URLs, not available in list endpoint
		updatedAt: prod.updatedAt,
		published: prod.publishedStatus as 'published' | 'draft' | 'archived',
		availability: prod.availability as 'online' | 'in-person' | 'hybrid',
		bufferTime: prod.bufferTime,
		cancellationHours: prod.cancellationHours,
	};
}

function computeProductCountsByCategory(
	categories: Category[],
	products: CardType[],
): Map<string, number> {
	const counts = new Map<string, number>();
	for (const cat of categories) {
		counts.set(cat.name, 0);
	}
	for (const prod of products) {
		const current = counts.get(prod.category) ?? 0;
		counts.set(prod.category, current + 1);
	}
	return counts;
}

export const actions = {
	deleteProduct,
	updateProduct,
	createCategory,
} satisfies Actions

type Props = {
	cards: CardType[];
	statuses: Set<string>;
	categories: Category[];
	availabilities: Set<string>;
	deleteProductForm: SuperValidated<DeleteProduct>;
	updateProductForm: SuperValidated<product>;
	createCategoryForm: SuperValidated<category>;
}

export const load: PageServerLoad = async ({ fetch }): Promise<Props> => {
	const deleteProductForm = await superValidate({ id: "e3eb8aaa-a255-4059-8013-6fbfb97442c0" }, arktype(deleteProductSchema, { defaults: deleteProductDefaults }))
	const updateProductForm = await superValidate(arktype(productSchema, { defaults: productDefaults }))
	const createCategoryForm = await superValidate(arktype(categorySchema, { defaults: categoryDefaults }))

	const availabilities: Set<string> = new Set([
		defaultAvailability,
		"online",
		"in-person",
		"hybrid",
	]);

	const statuses = new Set<string>([
		defaultStatus,
		"published",
		"draft",
		"archived",
	]);

	let categories: Category[];
	let cards: CardType[];

	if (env.USE_MOCK_DATA === "true") {
		categories = mockAdminCategories;
		cards = mockAdminCards;
	} else {
		try {
			const categoriesRes = await fetch(`${env.API_URL}/admin/categories`);
			if (!categoriesRes.ok) throw new Error(`Failed to fetch categories: ${categoriesRes.status} ${categoriesRes.statusText}`);
			const backendCategories: BackendCategoryWithImage[] = await categoriesRes.json();
			categories = backendCategories.map(mapBackendCategoryToFrontend);

			const productsRes = await fetch(`${env.API_URL}/admin/products`);
			if (!productsRes.ok) throw new Error(`Failed to fetch products: ${productsRes.status} ${productsRes.statusText}`);
			const backendProducts: BackendProductAggregator[] = await productsRes.json();
			cards = backendProducts.map(mapBackendProductToFrontend);

			// Compute product counts per category for the category view
			const countsByCategory = computeProductCountsByCategory(categories, cards);
			categories = categories.map(cat => ({
				...cat,
				productCount: countsByCategory.get(cat.name) ?? 0,
			}));
		} catch (err) {
			console.error("Error loading catalog data:", err);
			throw error(503, "Impossible de charger le catalogue. Veuillez réessayer.");
		}
	}

	// Combine the default category with the actual categories
	const allCategories: Category[] = [
		{ id: "default", name: defaultCategory },
		...categories
	];

	return {
		cards,
		statuses,
		categories: allCategories,
		availabilities,
		deleteProductForm,
		updateProductForm,
		createCategoryForm,
	}
}
