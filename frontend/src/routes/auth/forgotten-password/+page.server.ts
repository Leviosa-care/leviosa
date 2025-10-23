import type { PageServerLoad } from "./$types";
import { message, setError, superValidate } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import { type } from "arktype"

// TODO: find the right schema for this
const schema = type({
    prevPassword: "8 < string < 64",
    newPassword: "8 < string < 64",
})

const defaults = { prevPassword: "", newPassword: "" }

export const load: PageServerLoad = async () => {
    const form = await superValidate(arktype(schema, { defaults }))
    return { form }
}

export const actions = {
    default: async () => {

    }
}
