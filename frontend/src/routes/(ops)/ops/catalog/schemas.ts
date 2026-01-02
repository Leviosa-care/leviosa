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
