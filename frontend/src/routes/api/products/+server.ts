import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// GET /api/products — proxy to backend, forwards all query params
export const GET: RequestHandler = async ({ url, fetch }) => {
	const res = await fetch(`${env.API_URL}/products${url.search}`);
	const data = await res.json();
	return json(data, { status: res.status });
};
