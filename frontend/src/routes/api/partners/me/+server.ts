import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const PUT: RequestHandler = async ({ request, fetch }) => {
	const body = await request.json();
	const res = await fetch(`${env.API_URL}/partners/me`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
