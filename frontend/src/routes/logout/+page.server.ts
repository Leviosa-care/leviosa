import { redirect } from '@sveltejs/kit';
import { env } from '$env/dynamic/private';

export const actions = {
	default: async ({ cookies, locals }) => {
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

			cookies.set(locals.sessionCookieName, '', {
				path: '/',
				domain: locals.cookieDomain,
				expires: new Date(0)
			});
		}

		throw redirect(302, '/auth');
	}
};
