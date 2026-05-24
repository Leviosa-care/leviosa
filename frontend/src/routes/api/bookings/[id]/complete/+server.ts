import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// POST /api/bookings/[id]/complete — mark booking as completed
export const POST: RequestHandler = async ({ params, fetch }) => {
	const res = await fetch(`${env.API_URL}/bookings/${params.id}/complete`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' }
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
