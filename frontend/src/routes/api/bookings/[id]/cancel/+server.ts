import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// POST /api/bookings/[id]/cancel — cancel a booking
export const POST: RequestHandler = async ({ params, fetch, request }) => {
	const body = await request.json();
	const res = await fetch(`${env.API_URL}/bookings/${params.id}/cancel`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
