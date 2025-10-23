import type { PageServerLoad } from './$types';
import type { Actions } from "./$types"
import { arktype } from 'sveltekit-superforms/adapters';
import { superValidate, type SuperValidated } from 'sveltekit-superforms';

import {
    // product
    deleteProductSchema,
    deleteProductDefaults,
    type DeleteProduct,
    productSchema,
    productDefaults,
    type product,
    // category
    categorySchema,
    categoryDefaults,
    type category,
} from './schemas'
import { type CardType, cards, type Category } from "./products"
import {
    defaultStatus,
    defaultCategory,
    defaultAvailability,
} from "./default";

import { createProduct, updateProduct, deleteProduct, createCategory } from "./actions"

export const actions = {
    // products
    createProduct,
    deleteProduct,
    updateProduct,
    // categories
    createCategory,
} satisfies Actions

type Props = {
    cards: CardType[];
    statuses: Set<string>;
    // categories: Set<string>;
    categories: Category[];
    availabilities: Set<string>;
    deleteProductForm: SuperValidated<DeleteProduct>; createProductForm: SuperValidated<product>;
    updateProductForm: SuperValidated<product>;
    createCategoryForm: SuperValidated<category>;
}

export const load: PageServerLoad = async (): Promise<Props> => {
    const createProductForm = await superValidate(arktype(productSchema, { defaults: productDefaults }))
    // TODO: create as many forms for edit that I have cards or I will t
    const deleteProductForm = await superValidate({ id: "e3eb8aaa-a255-4059-8013-6fbfb97442c0" }, arktype(deleteProductSchema, { defaults: deleteProductDefaults }))
    const updateProductForm = await superValidate(arktype(productSchema, { defaults: productDefaults }))
    const createCategoryForm = await superValidate(arktype(categorySchema, { defaults: categoryDefaults }))

    // TODO: - get the cards server
    // - get the types server
    // That thing is constant but I can get it from the server

    // TODO: find how I am going to organise these products
    // try {
    // const res = await fetch(`${API_URL}/admin/categories`, {
    //     method: "GET",
    //     headers: { 'Content-Type': "application/json" },
    // })
    // if (!res.ok) {
    //    throw new Error(`Failed to fetch products: ${response.status}`);
    // }
    // const data = await res.json() as {
    //    cards: Card[];
    //    statuses: string[];
    //    categories: Category[];
    //    availabilities: string[];
    // }
    // } catch(error) {
    //    console.error('Error fetching product list:', error);
    //    throw error;
    // }

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

    // const categories: Set<string> = new Set([
    //     defaultCategory,
    //     "massage",
    //     "wellness",
    //     "mental coaching",
    // ]);

    const categories: Category[] = [
        // TODO: change IDs to fit what is going to be in the database
        // TODO: add the description into this since I might need it to be displayed on the category part
        { id: "default", name: defaultCategory }, { id: "massage", name: "Massage" },
        { id: "wellness", name: "Wellness" },
        { id: "mental coaching", name: "Mental Coaching" },
    ];

    return {
        // TODO: change the cards return for prodcuts that includes for each product an edit 
        cards,
        statuses,
        categories,
        availabilities,
        deleteProductForm,
        createProductForm,
        updateProductForm,
        createCategoryForm,
    }
}
