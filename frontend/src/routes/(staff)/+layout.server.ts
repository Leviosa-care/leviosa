import { redirect, error } from '@sveltejs/kit';
import { computePermissions } from '$lib/security/permissions';
import { env } from '$env/dynamic/private';

export const load = async ({ locals }) => {
	const user = locals.user;

	// Check if user is authenticated
	if (!user) {
		throw redirect(302, '/auth?redirect=' + encodeURIComponent('/staff'));
	}

	// Compute permissions to check access
	const permissions = computePermissions(user.role);

	// Check if user has staff or admin role
	if (!permissions.canAccessOps) {
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
