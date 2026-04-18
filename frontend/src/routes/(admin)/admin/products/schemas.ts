import { type } from 'arktype'
import type { Infer } from 'sveltekit-superforms'

// NOTE: category
// create & update
export const categorySchema = type({
    id: "string",
    name: "string",
    description: "string",
})

export type category = Infer<typeof categorySchema>

export const categoryDefaults: category = {
    id: "",
    name: "massage",
    description: "Un petit massage pour tester l'application",
}

// NOTE: the delete thing is shared across ressources
export const deleteSchema = type({
    id: "string"
})

export type Delete = Infer<typeof deleteSchema>

export const deleteDefaults: Delete = {
    // NOTE: this is the one that is used to create the product so that I can try to destroy it
    id: "0143ec0d-1715-41de-8774-63821bce8d27"
}

// NOTE: product

// delete
export const deleteProductSchema = type({
    id: "string"
})

export type DeleteProduct = Infer<typeof deleteProductSchema>

export const deleteProductDefaults: DeleteProduct = {
    // NOTE: this is the one that is used to create the product so that I can try to destroy it
    id: "0143ec0d-1715-41de-8774-63821bce8d27"
}

// create & update
export const productSchema = type({
    id: "string",
    name: "string",
    price: "string",
    category: "string",
    // TODO: change the above 'category' for these three
    // categoryID: "string",
    // newCategoryName?: "string",
    // newCategoryDescription?: "string",
    description: "string",
    duration: "number",
    updatedAt: "string",
    published: "'published' | 'draft' | 'archived'",
    availability: "'online' | 'in-person' | 'hybrid'",
    bufferTime: "number > 0",
    cancellationHours: "number > 0"
})

export type product = Infer<typeof productSchema>

export const productDefaults: product = {
    id: "0143ec0d-1715-41de-8774-63821bce8d27",
    name: "Swedish Massage",
    price: "6000",
    category: "massage", // NOTE: missing
    // TODO: change the above 'category' for these three
    // categoryID: "93a2b59c-eba2-42f4-8a9f-02824723a36d",
    // newCategoryName?: "Accompagment physique",
    // newCategoryDescription?: "Une description que je veux tester",
    description: "A relaxing full-body massage using classic techniques.",
    duration: 60, // in minutes
    updatedAt: new Date().toISOString(),
    published: "draft", // or "draft" / "archived"
    availability: "hybrid",
    bufferTime: 15, // in minutes
    cancellationHours: 24, // in hours
};
