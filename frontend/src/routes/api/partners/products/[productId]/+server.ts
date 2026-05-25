import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// GET /api/partners/products/[productId] — partners offering a given product
export const GET: RequestHandler = async ({ params, fetch }) => {
	const res = await fetch(`${env.API_URL}/partners/products/${params.productId}`);
	const data = await res.json();
	return json(data, { status: res.status });
};
