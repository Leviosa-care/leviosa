import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// PUT /api/bookings/[id]/notes — update partner notes on a booking
export const PUT: RequestHandler = async ({ params, fetch, request }) => {
	const body = await request.json();
	const res = await fetch(`${env.API_URL}/bookings/${params.id}/notes`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
