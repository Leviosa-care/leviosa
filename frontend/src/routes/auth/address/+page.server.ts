
import type { Actions, RequestEvent } from "./$types"
import { fail } from "@sveltejs/kit"

import { setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"

const schema = type({
    address1: "string",
    address2: "string",
    postalCode: "string == 5",
    city: "string",
})
const defaults = {
    address1: "",
    address2: "",
    postalCode: "",
    city: "",
}


export const load = async () => {
    const form = await superValidate(arktype(schema, { defaults }))

    return { form }
}

export const actions = {
    default: async ({ request }: RequestEvent) => {
        // TODO: here I should use throw new Error thing for some server error so that I can have the right path
        const form = await superValidate(request, arktype(schema, { defaults: defaults }))
        // TODO: remove that at the end
        console.log(form)
        if (!form.valid) {
            // TODO: find better messge for each field
            if (form.errors.address1) {
                setError(form, "address1", "L'adresse email saisie n'est pas valide. Veuillez vérifier et réessayer.")
            }
            if (form.errors.address2) {
                setError(form, "address2", "Le mot de passe saisi n'est pas valide. Veuillez vérifier et réessayer.")
            }
            if (form.errors.city) {
                setError(form, "city", "Le mot de passe saisi n'est pas valide. Veuillez vérifier et réessayer.")
            }
            if (form.errors.postalCode) {
                setError(form, "postalCode", "Le mot de passe saisi n'est pas valide. Veuillez vérifier et réessayer.")
            }
            return fail(400, { loginForm: form })
        }
        return { form }
    },
} satisfies Actions
