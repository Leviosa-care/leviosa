import { type } from 'arktype'
import type { Infer } from 'sveltekit-superforms'

// Category schemas
export const createCategorySchema = type({
    name: "string",
    description: "string",
    // Note: image file will be handled separately via FormData
})

export const updateCategorySchema = type({
    id: "string",
    name: "string",
    description: "string",
    status: "'draft' | 'published' | 'archived'",
})

export const deleteCategorySchema = type({
    id: "string"
})

// Types
export type CreateCategory = Infer<typeof createCategorySchema>
export type UpdateCategory = Infer<typeof updateCategorySchema>
export type DeleteCategory = Infer<typeof deleteCategorySchema>

// Defaults
export const createCategoryDefaults: CreateCategory = {
    name: "",
    description: "",
}

export const updateCategoryDefaults: UpdateCategory = {
    id: "",
    name: "",
    description: "",
    status: "draft",
}

export const deleteCategoryDefaults: DeleteCategory = {
    id: ""
}

// Product schemas
export const createProductSchema = type({
    name: "string",
    description: "string",
    category: "string",  // category ID
    duration: "number",  // in minutes
    bufferTime: "number",  // in minutes
    cancellationHours: "number",
    availability: "'online' | 'in-person' | 'hybrid'",
    stripeProductId: "string",
    // Note: image file will be handled separately via FormData
})

export const updateProductSchema = type({
    id: "string",
    name: "string",
    description: "string",
    category: "string",
    duration: "number",
    bufferTime: "number",
    cancellationHours: "number",
    status: "'draft' | 'published' | 'archived'",
    availability: "'online' | 'in-person' | 'hybrid'",
    stripeProductId: "string",
})

export const deleteProductSchema = type({
    id: "string"
})

// Product types
export type CreateProduct = Infer<typeof createProductSchema>
export type UpdateProduct = Infer<typeof updateProductSchema>
export type DeleteProduct = Infer<typeof deleteProductSchema>

// Product defaults
export const createProductDefaults: CreateProduct = {
    name: "",
    description: "",
    category: "",
    duration: 60,
    bufferTime: 0,
    cancellationHours: 24,
    availability: "in-person",
    stripeProductId: "",
}

export const updateProductDefaults: UpdateProduct = {
    id: "",
    name: "",
    description: "",
    category: "",
    duration: 60,
    bufferTime: 0,
    cancellationHours: 24,
    status: "draft",
    availability: "in-person",
    stripeProductId: "",
}

export const deleteProductDefaults: DeleteProduct = {
    id: ""
}

// Price schemas
export const createPriceSchema = type({
    productId: "string",
    amount: "number",  // in cents
    currency: "string",  // ISO currency code (eur, usd)
    interval: "'one_time' | 'month' | 'year'",
    nickname: "string?",  // optional
})

export const updatePriceSchema = type({
    id: "string",
    active: "boolean",
    nickname: "string?",  // optional
})

// Price types
export type CreatePrice = Infer<typeof createPriceSchema>
export type UpdatePrice = Infer<typeof updatePriceSchema>

// Price defaults
export const createPriceDefaults: CreatePrice = {
    productId: "",
    amount: 0,
    currency: "eur",
    interval: "one_time",
    nickname: undefined,
}

export const updatePriceDefaults: UpdatePrice = {
    id: "",
    active: true,
    nickname: undefined,
}
