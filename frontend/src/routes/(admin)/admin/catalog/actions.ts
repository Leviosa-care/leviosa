import type { RequestEvent } from "./$types"
import { superValidate, setError } from 'sveltekit-superforms';
import { fail } from "@sveltejs/kit";
import { arktype } from 'sveltekit-superforms/adapters';
import { env } from "$env/dynamic/private"

import {
    productSchema,
    productDefaults,
    deleteProductSchema,
    deleteProductDefaults,
    categorySchema,
    categoryDefaults,
    deleteSchema,
    deleteDefaults,
} from './schemas'

// PERF: Currently working on

export async function createCategory({ request, fetch }: RequestEvent) {
    const formData = await request.formData()
    const formValidated = await superValidate(formData, arktype(categorySchema, { defaults: categoryDefaults }))
    if (!formValidated.valid) {
        console.log("invalid from")
        if (formValidated.errors.name) return setError(formValidated, "name", "Le nom saisi n'est pas valide.")
        if (formValidated.errors.description) return setError(formValidated, "description", "La description est requise.")
    }
    formData.delete("id")
    const res = await fetch(`${env.API_URL}/admin/categories`, {
        method: "POST",
        body: formData,
    })
    console.log("fetch to golang backend done")
    if (!res.ok) {
        switch (res.status) {
            case 403:
                return setError(formValidated, "Vous n'avez pas les droits pour effectuer cette action.")
            case 404:
                return setError(formValidated, "La ressource demandée est introuvable.")
            case 409:
                return setError(formValidated, "Cette ressource est déjà utilisée et ne peut pas être supprimée.")
            case 500:
                return setError(formValidated, "Une erreur serveur est survenue. Le service semble momentanément indisponible.")
            case 400:
                const errorMsg = await res.json()
                return setError(formValidated, "La requête que vous avez soumise contient des données incorrectes ou incomplètes. Veuillez vérifier les champs saisis et réessayer.", errorMsg.msg);
            default:
                return setError(formValidated, "Une erreur est survenue. Veuillez réessayer.");
        }
    }
    console.log("Category successfully created. The status found is:", res.status)
    return {
        form: formValidated,
        success: true,
    }
}

export async function deleteCategory({ request, fetch }: RequestEvent) {
    const form = await superValidate(request, arktype(deleteSchema, { defaults: deleteDefaults }))
    if (!form.valid) {
        console.log("invalid from")
        return setError(form, "id", "L'ID saisi fourni n'est pas valide.")
    }
    const res = await fetch(`${env.API_URL}/admin/categories/${form.data.id}`, {
        method: "DELETE",
    })
    console.log("fetch to golang backend done")
    if (!res.ok) {
        switch (res.status) {
            case 403:
                return setError(form, "Vous n'avez pas les droits pour effectuer cette action.")
            case 404:
                return setError(form, "La ressource demandée est introuvable.")
            case 409:
                return setError(form, "Cette ressource est déjà utilisée et ne peut pas être supprimée.")
            case 500:
                return setError(form, "Une erreur serveur est survenue. Le service semble momentanément indisponible.")
            case 400:
                const errorMsg = await res.json()
                return setError(form, "La requête que vous avez soumise contient des données incorrectes ou incomplètes. Veuillez vérifier les champs saisis et réessayer.", errorMsg.msg);
            default:
                return setError(form, "Une erreur est survenue. Veuillez réessayer.");
        }
    }
    console.log("Category successfully removed. The status found is:", res.status)
    return {
        form,
        success: true,
    }
}

// NOTE: Done
export async function createProduct({ request, fetch }: RequestEvent) {
    const formData = await request.formData()

    const formValidated = await superValidate(formData, arktype(productSchema, { defaults: productDefaults }))

    if (!formValidated.valid) {
        if (formValidated.errors.name) return setError(formValidated, "name", "Le nom saisi n'est pas valide.")
        if (formValidated.errors.price) return setError(formValidated, "price", "Le prix saisi n'est pas valide.")
        if (formValidated.errors.category) return setError(formValidated, "category", "La catégorie est requise.")
        if (formValidated.errors.description) return setError(formValidated, "description", "La description est requise.")
        if (formValidated.errors.duration) return setError(formValidated, "duration", "La durée doit être un nombre valide.")
        if (formValidated.errors.bufferTime) return setError(formValidated, "bufferTime", "Le temps tampon doit être supérieur à 0.")
        if (formValidated.errors.cancellationHours) return setError(formValidated, "cancellationHours", "Le délai d'annulation doit être supérieur à 0.")
    }

    // image validation
    const image = formData.get("image")
    if (!image && formData.get("published") === "published") {
        console.log("the product should be published but there is no image")
        return setError(formValidated, "published", "Le status ne peut pas etre valide sans image publiee.")
    }
    if (image && image instanceof File && image.size > 0) {
        console.log("here this is the case where I check the type of image to see if it is a file")
        console.log("the image size:", image.size)
        console.log("the image type:", image.type)
        if (!image.type.startsWith('image/')) {
            // NOTE: this is where I went for some reason
            return fail(400, {
                form: formValidated,
                message: 'Uploaded file must be an image.'
            });
        }
    }

    // remove unecessary fields
    formData.delete("id")
    formData.delete("updatedAt")

    const res = await fetch(`${env.API_URL}/admin/products`, {
        method: "POST",
        body: formData,
    })

    if (!res.ok) {
        console.log("Something went wrong with the API response, status:", res.status)
        switch (res.status) {
            case 403:
                return setError(formValidated, "Vous n'avez pas les droits pour effectuer cette action.")
            case 404:
                return setError(formValidated, "La ressource demandée est introuvable.")
            case 409:
                return setError(formValidated, "Cette ressource est déjà utilisée et ne peut pas être supprimée.")
            case 500:
                return setError(formValidated, "Une erreur serveur est survenue. Le service semble momentanément indisponible.")
            case 400:
                const errorMsg = await res.json()
                return setError(formValidated, "La requête que vous avez soumise contient des données incorrectes ou incomplètes. Veuillez vérifier les champs saisis et réessayer.", errorMsg);
            default:
                return setError(formValidated, "Une erreur est survenue. Veuillez réessayer.");
        }
    }
    console.log("Product successfully created.")
    return {
        form: formValidated,
        success: true,
    }
}

export async function deleteProduct({ request, fetch }: RequestEvent) {
    const form = await superValidate(request, arktype(deleteProductSchema, { defaults: deleteProductDefaults }))
    console.log("here in the delete action")
    if (!form.valid) return setError(form, "id", "L'ID saisie n'est pas valide. Veuillez vérifier et réessayer.")

    const res = await fetch(`${env.API_URL}/admin/products/${form.data.id}`, {
        method: "DELETE",
    })

    if (!res.ok) {
        switch (res.status) {
            case 403:
                return setError(form, "Vous n'avez pas les droits pour effectuer cette action.")
            case 404:
                return setError(form, "La ressource demandée est introuvable.")
            case 409:
                return setError(form, "Cette ressource est déjà utilisée et ne peut pas être supprimée.")
            case 500:
                return setError(form, "Une erreur serveur est survenue. Le service semble momentanément indisponible.")
            case 400:
                const errorMsg = await res.json()
                return setError(form, "La requête que vous avez soumise contient des données incorrectes ou incomplètes. Veuillez vérifier les champs saisis et réessayer.", errorMsg.msg);
            default:
                return setError(form, "Une erreur est survenue. Veuillez réessayer.");
        }
    }
    return {
        form,
    }
}

export async function updateProduct() {
    console.log("here is the udpate action")
}
