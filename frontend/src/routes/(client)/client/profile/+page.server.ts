import { env } from '$env/dynamic/private';
import { fail, type Actions } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, locals }) => {
	if (env.USE_MOCK_DATA === 'true') {
		return { user: locals.user ?? null };
	}

	const res = await fetch(`${env.API_URL}/users/me`);

	if (!res.ok) {
		return { user: null };
	}

	const user = await res.json();
	return { user };
};

export const actions: Actions = {
	default: async ({ request, fetch }) => {
		const formData = await request.formData();

		const gender = formData.get('gender') as string;
		const birthdate = formData.get('birthdate') as string;
		const address1 = formData.get('address1') as string;
		const address2 = (formData.get('address2') as string) ?? '';
		const postalCode = formData.get('postalCode') as string;
		const city = formData.get('city') as string;

		const body: Record<string, unknown> = {};

		if (gender) body.gender = gender;
		if (birthdate) body.birthdate = birthdate + 'T00:00:00Z';
		if (address1) body.address1 = address1;
		if (address2) body.address2 = address2;
		if (postalCode) body.postal_code = postalCode;
		if (city) body.city = city;

		const res = await fetch(`${env.API_URL}/users/me`, {
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body),
		});

		if (!res.ok) {
			const err = await res.json().catch(() => ({}));
			return fail(res.status, { error: err.message ?? 'Impossible de mettre à jour votre profil.' });
		}

		const updatedUser = await res.json();
		return { success: true as const, user: updatedUser };
	}
};
