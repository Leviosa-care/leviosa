import { env } from '$env/dynamic/private';
import { redirect, error } from '@sveltejs/kit';
import { computePermissions } from '$lib/security/permissions';

export const load = async ({ locals, parent, fetch }) => {
	const { user, permissions } = await parent();

	if (!user) {
		throw redirect(302, '/auth?redirect=' + encodeURIComponent('/admin'));
	}

	if (user.role !== 'administrator') {
		throw error(403, 'Accès non autorisé');
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

	return { user, permissions, unreadCount };
};
