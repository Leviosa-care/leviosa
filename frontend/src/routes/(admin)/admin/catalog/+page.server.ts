import type { PageServerLoad } from './$types';
import type { Actions } from "./$types"
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
import { type CardType, cards as mockCards, categories as mockCategories, type Category } from "./products"
import {
	defaultStatus,
	defaultCategory,
	defaultAvailability,
} from "./default";

import { updateProduct, deleteProduct, createCategory } from "./actions"

interface BackendCategory {
	id: string;
	name: string;
	description: string;
	status: string;
	createdAt: string;
	updatedAt: string;
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
		price: "0.00", // Price is not in the base product response
		category: prod.category.name,
		description: prod.description,
		duration: prod.duration,
		image: "", // Image is not in the base product response
		updatedAt: prod.updatedAt,
		published: prod.publishedStatus as 'published' | 'draft' | 'archived',
		availability: prod.availability as 'online' | 'in-person' | 'hybrid',
		bufferTime: prod.bufferTime,
		cancellationHours: prod.cancellationHours,
	};
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
		categories = mockCategories;
		cards = mockCards;
	} else {
		try {
			const categoriesRes = await fetch(`${env.API_URL}/admin/categories`);
			if (!categoriesRes.ok) throw new Error(`Failed to fetch categories: ${categoriesRes.status} ${categoriesRes.statusText}`);
			const backendCategories: BackendCategory[] = await categoriesRes.json();
			categories = backendCategories.map(mapBackendCategoryToFrontend);

			const productsRes = await fetch(`${env.API_URL}/admin/products`);
			if (!productsRes.ok) throw new Error(`Failed to fetch products: ${productsRes.status} ${productsRes.statusText}`);
			const backendProducts: BackendProduct[] = await productsRes.json();
			cards = backendProducts.map(mapBackendProductToFrontend);
		} catch (error) {
			console.error("Error loading catalog data:", error);
			categories = [];
			cards = [];
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
