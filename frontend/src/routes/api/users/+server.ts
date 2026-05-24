import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({ url, fetch }) => {
	const q = (url.searchParams.get('search') ?? '').toLowerCase().trim();

	const res = await fetch(`${env.API_URL}/admin/users`);
	if (!res.ok) return json([], { status: res.status });

	const users: Array<{ id: string; first_name?: string; last_name?: string; email?: string }> =
		await res.json();

	const filtered = q
		? users.filter(
				(u) =>
					`${u.first_name ?? ''} ${u.last_name ?? ''}`.toLowerCase().includes(q) ||
					(u.email ?? '').toLowerCase().includes(q)
			)
		: users;

	return json(filtered);
};
