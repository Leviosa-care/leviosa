import type { PageServerLoad, Actions, RequestEvent } from "./$types"

import { superValidate, setError } from "sveltekit-superforms"
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"
import { redirect } from "@sveltejs/kit";

const schema = type({
    firstname: "string > 1",
    lastname: "string > 1",
    gender: "'' | 'man' | 'woman' | 'non_binary' | 'prefer_not_to_say' | 'custom' | 'not precised'",
    birthdate: "string",
    telephone: "string == 10"
})

const defaults = {
    firstname: "",
    lastname: "",
    gender: '',
    birthdate: "",
    telephone: ""
} as typeof schema.infer

export const load: PageServerLoad = async ({ url }: RequestEvent) => {
    // Get form data from URL search params (if user is coming back from a later step)
    const firstname = url.searchParams.get("firstname") || defaults.firstname;
    const lastname = url.searchParams.get("lastname") || defaults.lastname;
    const gender = url.searchParams.get("gender") || defaults.gender;
    const birthdate = url.searchParams.get("birthdate") || defaults.birthdate;
    const telephone = url.searchParams.get("telephone") || defaults.telephone;

    const form = await superValidate(
        { firstname, lastname, gender, birthdate, telephone },
        arktype(schema, { defaults })
    );
    return { form };
}

export const actions: Actions = {
    default: async ({ request }: RequestEvent) => {
        const form = await superValidate(request, arktype(schema, { defaults }))

        if (!form.valid) {
            if (form.errors.firstname) {
                setError(form, "firstname", "Le prénom est requis (minimum 2 caractères).");
            }
            if (form.errors.lastname) {
                setError(form, "lastname", "Le nom de famille est requis (minimum 2 caractères).");
            }
            if (form.errors.gender) {
                setError(form, "gender", "Veuillez sélectionner un genre.");
            }
            if (form.errors.birthdate) {
                setError(form, "birthdate", "Veuillez renseigner votre date de naissance.");
            }
            if (form.errors.telephone) {
                setError(form, "telephone", "Le numéro de téléphone doit contenir exactement 10 chiffres.");
            }
            return { form };
        }

        // Redirect to address page with form data as URL params
        const params = new URLSearchParams({
            firstname: form.data.firstname,
            lastname: form.data.lastname,
            gender: form.data.gender,
            birthdate: form.data.birthdate,
            telephone: form.data.telephone,
        });

        redirect(302, `/auth/address?${params.toString()}`);
    }
}
