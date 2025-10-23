import type { Actions, RequestEvent } from "./$types"
import { fail, redirect } from "@sveltejs/kit"

import { API_URL } from "$env/static/private";

import { setError, superValidate } from 'sveltekit-superforms'; import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"
import { parseDate } from "@internationalized/date";

import { pad } from "$lib/utils/pad"


const schema = type({
    password: "8 < string < 64",
    confirm: "string",
    email: "string.email",
    address1: "string",
    address2: "string",
    city: "string",
    postalCode: "string == 5",
    lastname: "string",
    firstname: "string",
    gender: "'' | 'man' | 'woman' | 'non_binary' | 'prefer_not_to_say' | 'custom' | 'not precised'",
    birthdate: "string",
    telephone: "string",
})



// TODO: change the default for this to try the API brother
const defaults = {
    password: "this_isnotaweakPassword123",
    confirm: "this_isnotaweakPassword123",
    email: "henry.gary@hotmail.com",
    address1: "01 Impasse Hoche",
    address2: "",
    city: "Ivry-Sur-Seine",
    postalCode: "94200",
    lastname: "HENRY",
    firstname: "Gary",
    gender: "man",
    // NOTE: not sure if this is the right formatting
    birthdate: "1997-08-12",
    // telephone: "0123456789",
    telephone: "0651919547",
} as typeof schema.infer

export const load = async () => {
    let form = await superValidate(arktype(schema, { defaults }))
    return { form }
}


export const actions = {
    default: async ({ request }: RequestEvent) => {
        const form = await superValidate(request, arktype(schema, { defaults }))
        console.log(form)
        if (!form.valid) {
            if (form.errors.password) {
                return setError(form, "password", "some error for the password field")
            }
            if (form.errors.confirm) {
                return setError(form, "confirm", "some error for the confirm field")
            }
            return fail(400, { form })
        }

        const date = parseDate(form.data.birthdate)
        const dateString = `${date.year}-${pad(date.month)}-${pad(date.day)}T00:00:00Z`;

        console.log("the date that I send to the backend is:", dateString)
        const res = await fetch(`${API_URL}/auth/register`, {
            method: "POST",
            headers: { 'Content-Type': "application/json" },
            body: JSON.stringify({
                email: form.data.email,
                password: form.data.password,
                address1: form.data.address1,
                address2: form.data.address2,
                telephone: form.data.telephone,
                city: form.data.city,
                postalCode: form.data.postalCode,
                lastname: form.data.lastname,
                firstname: form.data.firstname,
                // gender: form.data.gender,
                gender: {
                    gender: form.data.gender,
                    customGender: "",
                },
                birthdate: dateString
            })
        })

        if (!res.ok) {
            console.log("I get the status:", res.status)
            switch (res.status) {
                case 500:
                // some error with the server or service
                case 400:
                // bad request so an input failure
                // TODO: I might need to parse that value to give specific indication on what to change
                case 409:
                // the user exists
                case 422:
                // something went wrong with validating the data
                case 201:
                    // TODO: add the server error that you need brother
                    setError(form, "this is an error man, what is going on brother ?")
            }
        }
        throw redirect(302, "/")
    }
} satisfies Actions;
