import type { Actions, PageServerLoad } from "./$types"
import { fail, redirect } from "@sveltejs/kit"
import { API_URL } from "$env/static/private"

import { setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';

import {
    createCategorySchema,
    createCategoryDefaults,
    updateCategorySchema,
    updateCategoryDefaults,
    deleteCategorySchema,
    deleteCategoryDefaults,
    createProductSchema,
    createProductDefaults,
    updateProductSchema,
    updateProductDefaults,
    deleteProductSchema,
    deleteProductDefaults,
    createPriceSchema,
    createPriceDefaults,
    updatePriceSchema,
    updatePriceDefaults,
    createCouponSchema,
    createCouponDefaults,
    updateCouponSchema,
    updateCouponDefaults,
    deleteCouponSchema,
    deleteCouponDefaults,
} from "./schemas"
import { mockCategories, mockProducts, mockPrices, mockCoupons } from "./mockData"

// Enable mock mode for development
const MOCK_MODE = true

export const load: PageServerLoad = async ({ parent, fetch, cookies }) => {
    // ⬅️ pulls data from (ops)/+layout.server.ts
    const { permissions } = await parent()

    if (!permissions.canAccessOps) {
        throw redirect(302, '/app')
    }

    // Initialize category forms
    const createCategoryForm = await superValidate(arktype(createCategorySchema, { defaults: createCategoryDefaults }))
    const updateCategoryForm = await superValidate(arktype(updateCategorySchema, { defaults: updateCategoryDefaults }))
    const deleteCategoryForm = await superValidate(arktype(deleteCategorySchema, { defaults: deleteCategoryDefaults }))

    // Initialize product forms
    const createProductForm = await superValidate(arktype(createProductSchema, { defaults: createProductDefaults }))
    const updateProductForm = await superValidate(arktype(updateProductSchema, { defaults: updateProductDefaults }))
    const deleteProductForm = await superValidate(arktype(deleteProductSchema, { defaults: deleteProductDefaults }))

    // Initialize price forms
    const createPriceForm = await superValidate(arktype(createPriceSchema, { defaults: createPriceDefaults }))
    const updatePriceForm = await superValidate(arktype(updatePriceSchema, { defaults: updatePriceDefaults }))

    // Initialize coupon forms
    const createCouponForm = await superValidate(arktype(createCouponSchema, { defaults: createCouponDefaults }))
    const updateCouponForm = await superValidate(arktype(updateCouponSchema, { defaults: updateCouponDefaults }))
    const deleteCouponForm = await superValidate(arktype(deleteCouponSchema, { defaults: deleteCouponDefaults }))

    let categories = []
    let products = []
    let prices = []
    let coupons = []

    if (MOCK_MODE) {
        // Use mock data in development
        categories = mockCategories
        products = mockProducts
        prices = mockPrices
        coupons = mockCoupons
    } else {
        const sessionCookie = cookies.get('session');

        // Fetch all categories (admin endpoint to see all statuses)
        const categoriesRes = await fetch(`${API_URL}/admin/categories`, {
            headers: {
                'Authorization': `Bearer ${sessionCookie}`,
            }
        })

        if (categoriesRes.ok) {
            categories = await categoriesRes.json()
        }

        // Fetch all products (admin endpoint to see all statuses)
        const productsRes = await fetch(`${API_URL}/admin/products`, {
            headers: {
                'Authorization': `Bearer ${sessionCookie}`,
            }
        })

        if (productsRes.ok) {
            products = await productsRes.json()
        }

        // Fetch all prices across all products
        // Note: We could fetch prices per product, but for the catalog management page
        // it's more efficient to get all prices at once
        for (const product of products) {
            const pricesRes = await fetch(`${API_URL}/admin/products/${product.id}/prices`, {
                headers: {
                    'Authorization': `Bearer ${sessionCookie}`,
                }
            })

            if (pricesRes.ok) {
                const productPrices = await pricesRes.json()
                prices = [...prices, ...productPrices]
            }
        }

        // Fetch all coupons
        const couponsRes = await fetch(`${API_URL}/admin/coupons`, {
            headers: {
                'Authorization': `Bearer ${sessionCookie}`,
            }
        })

        if (couponsRes.ok) {
            coupons = await couponsRes.json()
        }
    }

    return {
        createCategoryForm,
        updateCategoryForm,
        deleteCategoryForm,
        createProductForm,
        updateProductForm,
        deleteProductForm,
        createPriceForm,
        updatePriceForm,
        createCouponForm,
        updateCouponForm,
        deleteCouponForm,
        categories,
        products,
        prices,
        coupons,
    }
}

export const actions = {
    createCategory: async ({ request, fetch, cookies }) => {
        const formData = await request.formData()

        // Extract and validate form fields
        const name = formData.get('name') as string
        const description = formData.get('description') as string
        const imageFile = formData.get('image') as File | null

        const form = await superValidate({ name, description }, arktype(createCategorySchema, { defaults: createCategoryDefaults }))

        if (!form.valid) {
            return fail(400, { createCategoryForm: form })
        }

        if (MOCK_MODE) {
            // Mock mode: just log and return success
            console.log('📝 [MOCK] Creating category:', {
                name: form.data.name,
                description: form.data.description,
                hasImage: imageFile ? imageFile.name : 'No image'
            })

            // Simulate success
            return { createCategoryForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Create category
        const res = await fetch(`${API_URL}/admin/categories`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({
                name: form.data.name,
                description: form.data.description,
                metadata: {},
            })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Données invalides. Veuillez vérifier les champs.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 409:
                    return setError(form, "name", "Une catégorie avec ce nom existe déjà.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        const result = await res.json()
        const categoryId = result.id

        // Upload image if provided
        if (imageFile && imageFile.size > 0) {
            const imageFormData = new FormData()
            imageFormData.append('file', imageFile)
            imageFormData.append('parent_id', categoryId)
            imageFormData.append('parent_type', 'category')
            imageFormData.append('title', `${form.data.name} - Image`)
            imageFormData.append('is_active', 'true')

            const imageRes = await fetch(`${API_URL}/admin/images`, {
                method: "POST",
                headers: {
                    'Authorization': `Bearer ${sessionCookie}`,
                },
                body: imageFormData
            })

            if (!imageRes.ok) {
                console.error('Failed to upload image:', await imageRes.text())
                // Don't fail the whole operation, just log the error
            }
        }

        return { createCategoryForm: form }
    },

    updateCategory: async ({ request, fetch, cookies }) => {
        const formData = await request.formData()

        const id = formData.get('id') as string
        const name = formData.get('name') as string
        const description = formData.get('description') as string
        const status = formData.get('status') as 'draft' | 'published' | 'archived'
        const imageFile = formData.get('image') as File | null

        const form = await superValidate({ id, name, description, status }, arktype(updateCategorySchema, { defaults: updateCategoryDefaults }))

        if (!form.valid) {
            return fail(400, { updateCategoryForm: form })
        }

        if (MOCK_MODE) {
            // Mock mode: just log and return success
            console.log('✏️ [MOCK] Updating category:', {
                id: form.data.id,
                name: form.data.name,
                description: form.data.description,
                status: form.data.status,
                hasImage: imageFile ? imageFile.name : 'No image'
            })

            // Simulate success
            return { updateCategoryForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Update category
        const res = await fetch(`${API_URL}/admin/categories/${form.data.id}`, {
            method: "PATCH",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({
                name: form.data.name,
                description: form.data.description,
                status: form.data.status,
            })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Données invalides ou aucun champ fourni.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Catégorie introuvable.");
                case 409:
                    return setError(form, "name", "Une catégorie avec ce nom existe déjà.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        // Upload new image if provided
        if (imageFile && imageFile.size > 0) {
            const imageFormData = new FormData()
            imageFormData.append('file', imageFile)
            imageFormData.append('parent_id', form.data.id)
            imageFormData.append('parent_type', 'category')
            imageFormData.append('title', `${form.data.name} - Image`)
            imageFormData.append('is_active', 'true')

            const imageRes = await fetch(`${API_URL}/admin/images`, {
                method: "POST",
                headers: {
                    'Authorization': `Bearer ${sessionCookie}`,
                },
                body: imageFormData
            })

            if (!imageRes.ok) {
                console.error('Failed to upload image:', await imageRes.text())
            }
        }

        return { updateCategoryForm: form }
    },

    deleteCategory: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(deleteCategorySchema, { defaults: deleteCategoryDefaults }))

        if (!form.valid) {
            return fail(400, { deleteCategoryForm: form })
        }

        if (MOCK_MODE) {
            // Mock mode: just log and return success
            console.log('🗑️ [MOCK] Deleting category:', {
                id: form.data.id
            })

            // Simulate success
            return { deleteCategoryForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${API_URL}/admin/categories/${form.data.id}`, {
            method: "DELETE",
            headers: {
                'Authorization': `Bearer ${sessionCookie}`,
            }
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Impossible de supprimer cette catégorie car des produits y sont associés.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Catégorie introuvable.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { deleteCategoryForm: form }
    },

    createProduct: async ({ request, fetch, cookies }) => {
        const formData = await request.formData()

        // Extract and validate form fields
        const name = formData.get('name') as string
        const description = formData.get('description') as string
        const category = formData.get('category') as string
        const duration = parseInt(formData.get('duration') as string)
        const bufferTime = parseInt(formData.get('bufferTime') as string)
        const cancellationHours = parseInt(formData.get('cancellationHours') as string)
        const availability = formData.get('availability') as 'online' | 'in-person' | 'hybrid'
        const stripeProductId = formData.get('stripeProductId') as string
        const imageFile = formData.get('image') as File | null

        const form = await superValidate({
            name,
            description,
            category,
            duration,
            bufferTime,
            cancellationHours,
            availability,
            stripeProductId,
        }, arktype(createProductSchema, { defaults: createProductDefaults }))

        if (!form.valid) {
            return fail(400, { createProductForm: form })
        }

        if (MOCK_MODE) {
            // Mock mode: just log and return success
            console.log('📝 [MOCK] Creating product:', {
                name: form.data.name,
                description: form.data.description,
                category: form.data.category,
                duration: form.data.duration,
                bufferTime: form.data.bufferTime,
                cancellationHours: form.data.cancellationHours,
                availability: form.data.availability,
                stripeProductId: form.data.stripeProductId,
                hasImage: imageFile ? imageFile.name : 'No image'
            })

            // Simulate success
            return { createProductForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Create product
        const res = await fetch(`${API_URL}/admin/products`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({
                name: form.data.name,
                description: form.data.description,
                category: form.data.category,
                duration: form.data.duration,
                buffer_time: form.data.bufferTime,
                cancellation_hours: form.data.cancellationHours,
                availability: form.data.availability,
                stripe_product_id: form.data.stripeProductId,
                metadata: {},
            })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Données invalides. Veuillez vérifier les champs.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 409:
                    return setError(form, "name", "Un produit avec ce nom existe déjà.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        const result = await res.json()
        const productId = result.id

        // Upload image if provided
        if (imageFile && imageFile.size > 0) {
            const imageFormData = new FormData()
            imageFormData.append('file', imageFile)
            imageFormData.append('parent_id', productId)
            imageFormData.append('parent_type', 'product')
            imageFormData.append('title', `${form.data.name} - Image`)
            imageFormData.append('is_active', 'true')

            const imageRes = await fetch(`${API_URL}/admin/images`, {
                method: "POST",
                headers: {
                    'Authorization': `Bearer ${sessionCookie}`,
                },
                body: imageFormData
            })

            if (!imageRes.ok) {
                console.error('Failed to upload image:', await imageRes.text())
                // Don't fail the whole operation, just log the error
            }
        }

        return { createProductForm: form }
    },

    updateProduct: async ({ request, fetch, cookies }) => {
        const formData = await request.formData()

        const id = formData.get('id') as string
        const name = formData.get('name') as string
        const description = formData.get('description') as string
        const category = formData.get('category') as string
        const duration = parseInt(formData.get('duration') as string)
        const bufferTime = parseInt(formData.get('bufferTime') as string)
        const cancellationHours = parseInt(formData.get('cancellationHours') as string)
        const status = formData.get('status') as 'draft' | 'published' | 'archived'
        const availability = formData.get('availability') as 'online' | 'in-person' | 'hybrid'
        const stripeProductId = formData.get('stripeProductId') as string
        const imageFile = formData.get('image') as File | null

        const form = await superValidate({
            id,
            name,
            description,
            category,
            duration,
            bufferTime,
            cancellationHours,
            status,
            availability,
            stripeProductId,
        }, arktype(updateProductSchema, { defaults: updateProductDefaults }))

        if (!form.valid) {
            return fail(400, { updateProductForm: form })
        }

        if (MOCK_MODE) {
            // Mock mode: just log and return success
            console.log('✏️ [MOCK] Updating product:', {
                id: form.data.id,
                name: form.data.name,
                description: form.data.description,
                category: form.data.category,
                duration: form.data.duration,
                bufferTime: form.data.bufferTime,
                cancellationHours: form.data.cancellationHours,
                status: form.data.status,
                availability: form.data.availability,
                stripeProductId: form.data.stripeProductId,
                hasImage: imageFile ? imageFile.name : 'No image'
            })

            // Simulate success
            return { updateProductForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Update product
        const res = await fetch(`${API_URL}/admin/products/${form.data.id}`, {
            method: "PATCH",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({
                name: form.data.name,
                description: form.data.description,
                category: form.data.category,
                duration: form.data.duration,
                buffer_time: form.data.bufferTime,
                cancellation_hours: form.data.cancellationHours,
                status: form.data.status,
                availability: form.data.availability,
                stripe_product_id: form.data.stripeProductId,
            })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Données invalides ou aucun champ fourni.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Produit introuvable.");
                case 409:
                    return setError(form, "name", "Un produit avec ce nom existe déjà.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        // Upload new image if provided
        if (imageFile && imageFile.size > 0) {
            const imageFormData = new FormData()
            imageFormData.append('file', imageFile)
            imageFormData.append('parent_id', form.data.id)
            imageFormData.append('parent_type', 'product')
            imageFormData.append('title', `${form.data.name} - Image`)
            imageFormData.append('is_active', 'true')

            const imageRes = await fetch(`${API_URL}/admin/images`, {
                method: "POST",
                headers: {
                    'Authorization': `Bearer ${sessionCookie}`,
                },
                body: imageFormData
            })

            if (!imageRes.ok) {
                console.error('Failed to upload image:', await imageRes.text())
            }
        }

        return { updateProductForm: form }
    },

    deleteProduct: async ({ request, fetch, cookies }) => {
        const form = await superValidate(request, arktype(deleteProductSchema, { defaults: deleteProductDefaults }))

        if (!form.valid) {
            return fail(400, { deleteProductForm: form })
        }

        if (MOCK_MODE) {
            // Mock mode: just log and return success
            console.log('🗑️ [MOCK] Deleting product:', {
                id: form.data.id
            })

            // Simulate success
            return { deleteProductForm: form }
        }

        const sessionCookie = cookies.get('session');

        const res = await fetch(`${API_URL}/admin/products/${form.data.id}`, {
            method: "DELETE",
            headers: {
                'Authorization': `Bearer ${sessionCookie}`,
            }
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Impossible de supprimer ce produit.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Produit introuvable.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { deleteProductForm: form }
    },

    createPrice: async ({ request, fetch, cookies }) => {
        const formData = await request.formData()

        const productId = formData.get('productId') as string
        const amount = parseInt(formData.get('amount') as string)
        const currency = formData.get('currency') as string
        const interval = formData.get('interval') as 'one_time' | 'month' | 'year'
        const nickname = formData.get('nickname') as string | null

        const form = await superValidate({
            productId,
            amount,
            currency,
            interval,
            nickname: nickname || undefined,
        }, arktype(createPriceSchema, { defaults: createPriceDefaults }))

        if (!form.valid) {
            return fail(400, { createPriceForm: form })
        }

        if (MOCK_MODE) {
            console.log('📝 [MOCK] Creating price:', {
                productId: form.data.productId,
                amount: form.data.amount,
                currency: form.data.currency,
                interval: form.data.interval,
                nickname: form.data.nickname,
            })

            return { createPriceForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Create price via product endpoint
        const res = await fetch(`${API_URL}/admin/products/${form.data.productId}/prices`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({
                amount: form.data.amount,
                currency: form.data.currency,
                interval: form.data.interval,
                nickname: form.data.nickname,
                metadata: {},
            })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Données invalides. Veuillez vérifier les champs.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Produit introuvable.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 502:
                    return setError(form, "Erreur Stripe. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { createPriceForm: form }
    },

    updatePrice: async ({ request, fetch, cookies }) => {
        const formData = await request.formData()

        const id = formData.get('id') as string
        const active = formData.get('active') === 'true'
        const nickname = formData.get('nickname') as string | null

        const form = await superValidate({
            id,
            active,
            nickname: nickname || undefined,
        }, arktype(updatePriceSchema, { defaults: updatePriceDefaults }))

        if (!form.valid) {
            return fail(400, { updatePriceForm: form })
        }

        if (MOCK_MODE) {
            console.log('✏️ [MOCK] Updating price:', {
                id: form.data.id,
                active: form.data.active,
                nickname: form.data.nickname,
            })

            return { updatePriceForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Update price
        const res = await fetch(`${API_URL}/admin/prices/${form.data.id}`, {
            method: "PATCH",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({
                active: form.data.active,
                nickname: form.data.nickname,
                metadata: {},
            })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Données invalides ou aucun champ fourni.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Prix introuvable.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 502:
                    return setError(form, "Erreur Stripe. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { updatePriceForm: form }
    },

    createCoupon: async ({ request, fetch, cookies }) => {
        const formData = await request.formData()

        const name = formData.get('name') as string
        const percentOff = formData.get('percentOff') ? parseFloat(formData.get('percentOff') as string) : undefined
        const amountOff = formData.get('amountOff') ? parseInt(formData.get('amountOff') as string) : undefined
        const currency = formData.get('currency') as string | null
        const duration = formData.get('duration') as string
        const durationInMonths = formData.get('durationInMonths') ? parseInt(formData.get('durationInMonths') as string) : undefined
        const maxRedemptions = formData.get('maxRedemptions') ? parseInt(formData.get('maxRedemptions') as string) : undefined
        const redeemBy = formData.get('redeemBy') as string | null

        const form = await superValidate({
            name,
            percentOff,
            amountOff,
            currency: currency || undefined,
            duration,
            durationInMonths,
            maxRedemptions,
            redeemBy: redeemBy || undefined,
        }, arktype(createCouponSchema, { defaults: createCouponDefaults }))

        if (!form.valid) {
            return fail(400, { createCouponForm: form })
        }

        // Validate mutually exclusive discount types
        if ((form.data.percentOff !== undefined && form.data.amountOff !== undefined) ||
            (form.data.percentOff === undefined && form.data.amountOff === undefined)) {
            return setError(form, "Vous devez fournir soit percentOff soit amountOff, mais pas les deux.");
        }

        // Validate currency is required if amountOff is set
        if (form.data.amountOff !== undefined && !form.data.currency) {
            return setError(form, "currency", "La devise est requise pour un montant fixe.");
        }

        // Validate durationInMonths is required if duration is repeating
        if (form.data.duration === "repeating" && !form.data.durationInMonths) {
            return setError(form, "durationInMonths", "La durée en mois est requise pour un coupon récurrent.");
        }

        if (MOCK_MODE) {
            console.log('📝 [MOCK] Creating coupon:', form.data)
            return { createCouponForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Prepare request body
        const requestBody: any = {
            name: form.data.name,
            duration: form.data.duration,
        }

        if (form.data.percentOff !== undefined) {
            requestBody.percentOff = form.data.percentOff
        }

        if (form.data.amountOff !== undefined) {
            requestBody.amountOff = form.data.amountOff
            requestBody.currency = form.data.currency
        }

        if (form.data.durationInMonths) {
            requestBody.durationInMonths = form.data.durationInMonths
        }

        if (form.data.maxRedemptions) {
            requestBody.maxRedemptions = form.data.maxRedemptions
        }

        if (form.data.redeemBy) {
            requestBody.redeemBy = form.data.redeemBy
        }

        requestBody.metadata = {}

        // Create coupon
        const res = await fetch(`${API_URL}/admin/coupons`, {
            method: "POST",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify(requestBody)
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Données invalides. Vérifiez que vous avez fourni soit un pourcentage soit un montant.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 502:
                    return setError(form, "Erreur Stripe. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { createCouponForm: form }
    },

    updateCoupon: async ({ request, fetch, cookies }) => {
        const formData = await request.formData()

        const id = formData.get('id') as string
        const name = formData.get('name') as string | null

        const form = await superValidate({
            id,
            name: name || undefined,
        }, arktype(updateCouponSchema, { defaults: updateCouponDefaults }))

        if (!form.valid) {
            return fail(400, { updateCouponForm: form })
        }

        // Validate at least one field is provided
        if (!form.data.name) {
            return setError(form, "Vous devez fournir au moins un champ à mettre à jour.");
        }

        if (MOCK_MODE) {
            console.log('✏️ [MOCK] Updating coupon:', form.data)
            return { updateCouponForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Update coupon
        const res = await fetch(`${API_URL}/admin/coupons/${form.data.id}`, {
            method: "PATCH",
            headers: {
                'Content-Type': "application/json",
                'Authorization': `Bearer ${sessionCookie}`,
            },
            body: JSON.stringify({
                name: form.data.name,
                metadata: {},
            })
        })

        if (!res.ok) {
            switch (res.status) {
                case 400:
                    return setError(form, "Données invalides ou aucun champ fourni.");
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Coupon introuvable.");
                case 415:
                    return setError(form, "Type de contenu non supporté.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 502:
                    return setError(form, "Erreur Stripe. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { updateCouponForm: form }
    },

    deleteCoupon: async ({ request, fetch, cookies }) => {
        const formData = await request.formData()

        const id = formData.get('id') as string

        const form = await superValidate({ id }, arktype(deleteCouponSchema, { defaults: deleteCouponDefaults }))

        if (!form.valid) {
            return fail(400, { deleteCouponForm: form })
        }

        if (MOCK_MODE) {
            console.log('🗑️ [MOCK] Deleting coupon:', form.data.id)
            return { deleteCouponForm: form }
        }

        const sessionCookie = cookies.get('session');

        // Delete coupon
        const res = await fetch(`${API_URL}/admin/coupons/${form.data.id}`, {
            method: "DELETE",
            headers: {
                'Authorization': `Bearer ${sessionCookie}`,
            },
        })

        if (!res.ok) {
            switch (res.status) {
                case 401:
                    return setError(form, "Non autorisé. Veuillez vous reconnecter.");
                case 403:
                    return setError(form, "Accès refusé. Vous devez être administrateur.");
                case 404:
                    return setError(form, "Coupon introuvable.");
                case 500:
                    return setError(form, "Erreur serveur. Veuillez réessayer.");
                case 502:
                    return setError(form, "Erreur Stripe. Veuillez réessayer.");
                case 503:
                    return setError(form, "Service temporairement indisponible.");
                default:
                    return setError(form, "Une erreur inattendue est survenue.");
            }
        }

        return { deleteCouponForm: form }
    },
} satisfies Actions
