import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// GET /api/partners/[partnerId]/availabilities/conflict — check conflict
export const GET: RequestHandler = async ({ params, url, fetch }) => {
	const qs = url.searchParams.toString();
	const res = await fetch(
		`${env.API_URL}/partners/${params.partnerId}/availabilities/conflict${qs ? `?${qs}` : ''}`
	);
	const data = await res.json();
	return json(data, { status: res.status });
};
