import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ params, fetch }) => {
	const { provider } = params;
	const res = await fetch(`${env.API_URL}/users/me/oauth/${provider}/link`, {
		method: 'POST',
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
