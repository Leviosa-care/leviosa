import { env } from '$env/dynamic/private';
import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

export const POST: RequestHandler = async ({ fetch, url }) => {
	const origin = url.origin;
	try {
		const res = await fetch(`${env.API_URL}/partners/me/stripe/onboarding-link`, {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({
				return_url: `${origin}/staff/profile?stripe=return`,
				refresh_url: `${origin}/staff/profile?stripe=refresh`
			})
		});
		const data = await res.json();
		return json(data, { status: res.status });
	} catch {
		return json({ error: 'Failed to generate onboarding link' }, { status: 502 });
	}
};
