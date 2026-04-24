import type { PageServerLoad } from './$types';
import type { Actions } from "./$types"
import { arktype } from 'sveltekit-superforms/adapters';
import { superValidate, type SuperValidated } from 'sveltekit-superforms';

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
import { type CardType, cards, categories, type Category } from "./products"
import {
	defaultStatus,
	defaultCategory,
	defaultAvailability,
} from "./default";

import { createProduct, updateProduct, deleteProduct, createCategory } from "./actions"

export const actions = {
	createProduct,
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
	createProductForm: SuperValidated<product>;
	updateProductForm: SuperValidated<product>;
	createCategoryForm: SuperValidated<category>;
}

export const load: PageServerLoad = async (): Promise<Props> => {
	const createProductForm = await superValidate(arktype(productSchema, { defaults: productDefaults }))
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
		createProductForm,
		updateProductForm,
		createCategoryForm,
	}
}
