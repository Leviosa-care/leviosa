import { env } from "$env/dynamic/private"
import type { PageServerLoad } from './$types';
import type { Actions, RequestEvent } from "./$types"
import { arktype } from 'sveltekit-superforms/adapters';
import { superValidate, type SuperValidated, setError } from 'sveltekit-superforms';

import {
    productSchema,
    productDefaults,
    type product,
} from './schemas'
import { type CardType, mockAdminCards } from "$lib/data/mockData"
import {
    defaultStatus,
    defaultCategory,
    defaultAvailability,
} from "./default";

export const actions = {
    create: async ({ request, fetch }: RequestEvent) => {
        console.log("here in the create action")
        const form = await superValidate(request, arktype(productSchema, { defaults: productDefaults }))
        console.log("here in the create function after the superValidate") // NOTE: that case should never happen since I am the one providing the ID as a string and there is no reason for it to change
        if (!form.valid) {
            if (form.errors.name) {
                return setError(form, "name", "Le nom saisie n'est pas valide. Veuillez vérifier et réessayer.")
            }
            // TODO: do all remaining field
        }

        const { id, ...data } = form.data
        const res = await fetch(`${env.API_URL}/admin/products`, {
            method: "POST",
            headers: { 'Content-Type': "application/json" },
            body: JSON.stringify(data),
        })

        if (!res.ok) {
            console.log("something went wrong with the API response")
            switch (res.status) {
                // TODO: find the best status here because this is what I need
                case 500:
                    return setError(form, "id", "Une erreur serveur est survenue. Le service semble momentanément indisponible.")
                case 422:
                // something went wrong with validating the data. I might need to parse that thing
                case 429:
                    return setError(form, "id", "Trop de tentatives. Veuillez réessayer plus tard.");
                case 409:
                // the user exists
                case 400:
                    return setError(form, "id", "Format d'adresse e-mail invalide. Veuillez vérifier.");
            }
        }
        return {
            form,
            success: true,
        }
    },
} satisfies Actions

type Props = {
    cards: CardType[];
    statuses: Set<string>;
    categories: Set<string>;
    availabilities: Set<string>;
    createProductForm: SuperValidated<product>;
}

export const load: PageServerLoad = async (): Promise<Props> => {

    const createProductForm = await superValidate(arktype(productSchema, { defaults: productDefaults }))
    // const testForm = await superValidate(arktype(deleteProductSchema, { defaults: deleteProductDefaults }))
    // TODO: - get the cards server
    // - get the types server
    // That thing is constant but I can get it from the server

    // TODO: find how I am going to organise these products
    // try {
    // const res = await fetch(`${env.API_URL}/products`, {
    //     method: "GET",
    //     headers: { 'Content-Type': "application/json" },
    // })
    // if (!res.ok) {
    //    throw new Error(`Failed to fetch products: ${response.status}`);
    // }
    // const data = await res.json() as {
    //    cards: Card[];
    //    statuses: string[];
    //    categories: string[];
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

    const categories: Set<string> = new Set([
        defaultCategory,
        "massage",
        "wellness",
    ]);

    return {
        cards: mockAdminCards,
        statuses,
        categories,
        availabilities,
        createProductForm,
    }
}
