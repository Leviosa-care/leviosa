import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals, fetch, url }) => {
	const user = locals.user;

	if (!user) {
		throw redirect(302, '/auth?redirectTo=' + encodeURIComponent(url.pathname));
	}

	if (user.role !== 'standard' && user.role !== 'premium') {
		throw redirect(302, '/');
	}

	let unreadCount = 0;
	if (env.USE_MOCK_DATA !== 'true') {
		try {
			const res = await fetch(`${env.API_URL}/threads/unread-count`);
			if (res.ok) {
				const data = await res.json();
				unreadCount = data.unread_count ?? 0;
			}
		} catch {
			// unread count is non-critical
		}
	}

	return { user, unreadCount, profileIncomplete: user.profile_incomplete === true };
};
