import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// POST /api/availabilities/[id]/cancel — cancel availability
export const POST: RequestHandler = async ({ params, fetch }) => {
	const res = await fetch(`${env.API_URL}/availabilities/${params.id}/cancel`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' }
	});
	if (res.status === 204) return new Response(null, { status: 204 });
	const data = await res.json();
	return json(data, { status: res.status });
};
