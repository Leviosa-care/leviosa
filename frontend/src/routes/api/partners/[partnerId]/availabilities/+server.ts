import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

// GET /api/partners/[partnerId]/availabilities — forwards all query params (status, start_time, etc.)
export const GET: RequestHandler = async ({ params, url, fetch }) => {
	const qs = url.searchParams.toString();
	const res = await fetch(
		`${env.API_URL}/partners/${params.partnerId}/availabilities${qs ? `?${qs}` : ''}`
	);
	const data = await res.json();
	return json(data, { status: res.status });
};
