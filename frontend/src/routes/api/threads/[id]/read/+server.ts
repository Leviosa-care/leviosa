import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ params, fetch }) => {
	const res = await fetch(`${env.API_URL}/threads/${params.id}/read`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' }
	});
	return json(res.ok ? await res.json() : {}, { status: res.status });
};
