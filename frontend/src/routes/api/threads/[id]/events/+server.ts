import { env } from '$env/dynamic/private';
import type { RequestHandler } from './$types';

/**
 * Proxies the backend SSE endpoint. The browser's EventSource can only do GET
 * requests, and cookies are forwarded automatically by the SvelteKit fetch
 * (enriched in hooks.server.ts with the auth token).
 */
export const GET: RequestHandler = async ({ params, fetch }) => {
	const upstream = new Request(`${env.API_URL}/threads/${params.id}/events`, {
		headers: {
			Accept: 'text/event-stream'
		}
	});

	const res = await fetch(upstream);

	if (!res.ok) {
		return new Response(JSON.stringify({ error: 'upstream error' }), {
			status: res.status,
			headers: { 'Content-Type': 'application/json' }
		});
	}

	// Stream the SSE response through without buffering.
	const body = res.body!;
	return new Response(body, {
		status: 200,
		headers: {
			'Content-Type': 'text/event-stream',
			'Cache-Control': 'no-cache',
			Connection: 'keep-alive',
			'X-Accel-Buffering': 'no'
		}
	});
};
