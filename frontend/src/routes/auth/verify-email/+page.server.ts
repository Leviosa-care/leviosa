import { type } from 'arktype';
import { superValidate, setError } from 'sveltekit-superforms';
import { arktype } from 'sveltekit-superforms/adapters';
import type { RequestEvent } from './$types';
import { API_URL } from "$env/static/private"
import { redirect } from '@sveltejs/kit';

// TODO: then add the otp thing brother
const schema = type({
    otp0: "/^\\d$/",
    otp1: "/^\\d$/",
    otp2: "/^\\d$/",
    otp3: "/^\\d$/",
    otp4: "/^\\d$/",
    otp5: "/^\\d$/",
});

// Defaults should also be defined outside the load function
const defaults = {
    // otp0: '',
    // otp1: '',
    // otp2: '',
    // otp3: '',
    // otp4: '',
    // otp5: '',

    otp0: '13',
    otp1: '4',
    otp2: '4',
    otp3: '29',
    otp4: '0',
    otp5: '2',
};

export const load = async () => {
    const form = await superValidate(arktype(schema, { defaults }));
    return { form };
};

export const actions = {
    default: async ({ request, fetch }: RequestEvent) => {
        const form = await superValidate(request, arktype(schema, { defaults }));
        console.log("in the form action:\n", form);

        // NOTE: this is the expected error
        // error(404, "Juste une erreur pour voir !")
        // NOTE: this is the not expected error
        // throw new Error(404, "Juste une erreur pour voir !")
        if (!form.valid) return setError(form, "Veuillez entrer des valeurs numériques à un seul chiffre (de 0 à 9).");

        const otp = [
            form.data.otp0,
            form.data.otp1,
            form.data.otp2,
            form.data.otp3,
            form.data.otp4,
            form.data.otp5,
        ]

        try {
            const res = await fetch(`${API_URL}/auth/otp`, {
                method: "POST",
                headers: { "Content-Type": "application/json", },
                body: JSON.stringify({ otp: otp.join("") })
            })
            if (!res.ok) {
                switch (res.status) {
                    case 500:
                    // server error
                    case 422:
                    // user not processed for some reason
                    case 409:
                    // user exists
                    case 400:
                        return setError(form, "")
                }
            }
        } catch (err) {
            console.error("here is the error:", err)
        }

        throw redirect(200, "/auth/general")
    }
};
