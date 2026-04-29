import { redirect, error } from '@sveltejs/kit';
import { computePermissions } from '$lib/security/permissions';
import { env } from '$env/dynamic/private';

export const load = async ({ locals, parent }) => {
	const { user, permissions } = await parent();

	// Check if user is authenticated
	if (!user) {
		throw redirect(302, '/auth?redirect=' + encodeURIComponent('/admin'));
	}

	// Check if user has admin role
	if (user.role !== 'administrator') {
		throw error(403, 'Accès non autorisé');
	}

	return {
		user,
		permissions,
	};
};

export const actions = {
	logout: async ({ cookies, locals }) => {
		const sessionID = cookies.get(locals.sessionCookieName);

		if (sessionID) {
			try {
				// Call backend logout endpoint
				await fetch(`${env.API_URL}/auth/logout`, {
					method: 'POST',
					headers: {
						Cookie: `leviosa_access_token=${sessionID}`
					}
				});
			} catch (error) {
				console.error('Backend logout failed:', error);
			}

			// Clear the session cookie regardless of backend response
			cookies.set(locals.sessionCookieName, '', {
				path: '/',
				domain: locals.cookieDomain,
				expires: new Date(0)
			});
		}

		throw redirect(302, '/auth/login');
	}
};
