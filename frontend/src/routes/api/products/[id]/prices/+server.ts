import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// GET /api/products/[id]/prices — proxy to public price endpoint (no auth required)
export const GET: RequestHandler = async ({ params, fetch }) => {
	const res = await fetch(`${env.API_URL}/products/${params.id}/prices`);
	const data = await res.json();
	return json(data, { status: res.status });
};
