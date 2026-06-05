import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const actions = {
	default: async ({ cookies, locals }) => {
		if (env.USE_MOCK_DATA === 'true') {
			cookies.delete('mock_session', { path: '/' });
			throw redirect(302, '/auth');
		}

		const sessionID = cookies.get(locals.sessionCookieName);

		if (sessionID) {
			try {
				await fetch(`${env.API_URL}/auth/logout`, {
					method: 'POST',
					headers: {
						Cookie: `leviosa_access_token=${sessionID}`
					}
				});
			} catch (error) {
				console.error('Backend logout failed:', error);
			}

			cookies.delete(locals.sessionCookieName, {
				path: '/',
				domain: locals.cookieDomain
			});
			cookies.delete('leviosa_refresh_token', {
				path: '/',
				domain: locals.cookieDomain
			});
		}

		throw redirect(302, '/auth');
	}
};
