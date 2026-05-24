import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const DELETE: RequestHandler = async ({ params, fetch }) => {
	const { provider } = params;
	const res = await fetch(`${env.API_URL}/users/me/oauth/${provider}/unlink`, {
		method: 'DELETE',
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
