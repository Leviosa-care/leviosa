import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ params, url, fetch }) => {
	const query = url.searchParams.toString();
	const upstream = `${env.API_URL}/threads/${params.id}/messages${query ? '?' + query : ''}`;
	const res = await fetch(upstream);
	if (!res.ok) return json({ messages: [] }, { status: res.status });
	return json(await res.json());
};

export const POST: RequestHandler = async ({ params, request, fetch }) => {
	const body = await request.json();
	const res = await fetch(`${env.API_URL}/threads/${params.id}/messages`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});
	const data = await res.json();
	return json(data, { status: res.status });
};
