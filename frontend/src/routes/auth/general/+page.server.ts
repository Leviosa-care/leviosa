import type { PageServerLoad, Actions, RequestEvent } from "./$types"

import { superValidate, setError } from "sveltekit-superforms"
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"
import { parseDate } from "@internationalized/date";

import { pad } from "$lib/utils/pad"


const schema = type({
    firstname: "string > 1",
    lastname: "string > 1",
    gender: "'' | 'man' | 'woman' | 'non_binary' | 'prefer_not_to_say' | 'custom' | 'not precised'",
    birthdate: "string",
    telephone: "string == 10"
})
// Month-Year-Day
const defaults = {
    firstname: "John",
    lastname: "DOE",
    gender: '',
    birthdate: "",
    telephone: "0123456789"
} as typeof schema.infer

export const load: PageServerLoad = async () => {
    const form = await superValidate(arktype(schema, { defaults }))
    return { form }
}

export const actions: Actions = {
    default: async ({ request }: RequestEvent) => {
        const form = await superValidate(request, arktype(schema, { defaults }))
        console.log("form:\n", form)

        const date = parseDate(form.data.birthdate)
        const dateString = `${date.year}-${pad(date.month)}-${pad(date.day)}T00:00:00Z`;

        console.log("the date string that I get is:", dateString)
        if (!form.valid) {
            // TODO: return some error to display to the user
            // TODO: check each field to see if there is an error
            if (form.errors.firstname) {
                setError(form, "firstname", "L'adresse email saisie n'est pas valide. Veuillez vérifier et réessayer.")
                setError(form, "lastname", "L'adresse email saisie n'est pas valide. Veuillez vérifier et réessayer.")
                setError(form, "gender", "Le genre n'est pas valide")
                setError(form, "birthdate", "Veuillez renseinge votre date d'anniversaire")
                setError(form, "telephone", "Wrong phone number brother")
            }
        }
        return { form }
    }
}
