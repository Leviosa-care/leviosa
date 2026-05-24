import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ fetch }) => {
	const res = await fetch(`${env.API_URL}/threads`);
	if (!res.ok) return json([], { status: res.status });
	return json(await res.json());
};

export const POST: RequestHandler = async ({ request, fetch }) => {
	const body = await request.json();
	const res = await fetch(`${env.API_URL}/threads`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
