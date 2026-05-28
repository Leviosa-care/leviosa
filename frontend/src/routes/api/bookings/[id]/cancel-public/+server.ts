import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// POST /api/bookings/[id]/cancel-public — cancel a booking using a booking token (no auth required)
export const POST: RequestHandler = async ({ params, fetch, request, url }) => {
	const token = url.searchParams.get('token');
	if (!token) {
		return json({ error: 'Token manquant' }, { status: 400 });
	}

	const body = await request.json();
	const res = await fetch(`${env.API_URL}/bookings/${params.id}/cancel-public?token=${encodeURIComponent(token)}`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
