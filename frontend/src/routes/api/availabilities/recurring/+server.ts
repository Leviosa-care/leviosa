import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// POST /api/availabilities/recurring — create recurring availability
export const POST: RequestHandler = async ({ request, fetch }) => {
	const body = await request.json();
	const res = await fetch(`${env.API_URL}/availabilities/recurring`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
