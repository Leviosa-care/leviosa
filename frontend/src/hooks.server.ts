import { redirect, isRedirect, error, isHttpError, type Handle, type RequestEvent } from '@sveltejs/kit';
import { env } from "$env/dynamic/private"

import { handleLoginRedirect } from '$lib/utils/redirect';
import { mockUser } from '$lib/data/user';
import { isAdminDomain, isStaffDomain, getSessionCookieName, getCookieDomain } from '$lib/server/hostname';

// TODO: for the session thing
// on each request, send the accesstoken or session_token
// when the access token is no longer valid, use will receive 401 status code
// then make a request to a dedicated refresh endpoint on the server and the server resend an access token
// and optionnaly a new refresh  token. This is the rotation step !

export const handle: Handle = async ({ event, resolve }) => {
	// Detect hostname and set domain context
	const hostname = event.url.hostname;
	event.locals.isAdminDomain = isAdminDomain(hostname);
	event.locals.isStaffDomain = isStaffDomain(hostname);
	const cookieDomain = getCookieDomain(hostname);
	const sessionCookieName = getSessionCookieName(hostname);

	// On the admin subdomain (not localhost), redirect bare root to /admin so the auth flow runs
	const isLocalhost = hostname.startsWith('localhost') || hostname.startsWith('127.0.0.1');
	if (event.locals.isAdminDomain && !isLocalhost && event.url.pathname === '/') {
		throw redirect(302, '/admin');
	}

	// On the staff subdomain (not localhost), redirect bare root to /staff so the auth flow runs
	if (event.locals.isStaffDomain && !isLocalhost && event.url.pathname === '/') {
		throw redirect(302, '/staff');
	}

    if (env.COMING_SOON === 'true') {
        if (event.url.pathname !== '/coming-soon' && event.url.pathname !== '/healthz') {
            throw redirect(302, '/coming-soon')
        }
        return await resolve(event)
    }

    // USE_MOCK_DATA allows bypassing authentication for development/testing
    // Set via environment variable (e.g., in Ansible staging config)
    event.locals.sessionCookieName = sessionCookieName;
    event.locals.cookieDomain = cookieDomain ?? undefined;

    if (env.USE_MOCK_DATA === 'true') {
        event.locals.user = mockUser
        return await resolve(event)
    }

    // Routes that require authentication - all others are public
    const requiresAuth = (pathname: string) => {
        const protectedPrefixes = ['/staff', '/admin', '/premium'];
        return protectedPrefixes.some(prefix => pathname.startsWith(prefix));
    };

    if (!requiresAuth(event.url.pathname)) {
        return await resolve(event);
    }

    const sessionID = event.cookies.get(sessionCookieName);
    if (!sessionID) {
        let path = handleLoginRedirect(event.url, "expiredSession");
        const email = event.locals.user?.email;
        if (email) path += `&email=${email}`;
        throw redirect(302, path)
    }

    try {
        const user = await validateSession(event, sessionID)
        event.locals.user = user;
    } catch (err) {
        if (isRedirect(err)) throw err;
        if (isHttpError(err)) throw err;
        console.error("Error validating session:", err)
    }

    // enrich fetch with custom header with client IP
    event.fetch = async (input, init = {}) => {
        // Ensure headers object exists
        init.headers = {
            ...(init.headers ?? {}),
            [env.CLIENT_IP_HEADER ?? 'x-client-ip']: getClientIP(event.request),
            Cookie: `leviosa_access_token=${sessionID}`
        };
        return fetch(input, init);
    };

    return await resolve(event);
};

function getClientIP(request: Request): string {
    const headers = request.headers;
    return (
        headers.get('cf-connecting-ip') || // Cloudflare
        headers.get('x-real-ip') || // Nginx
        headers.get('x-forwarded-for')?.split(',')[0] || // Standard
        'unknown'
    );
}

async function validateSession(event: RequestEvent, sessionID: string): Promise<App.User> {
    const res = await fetch(`${env.API_URL}/users/me`, {
        method: "GET",
        headers: { Cookie: `leviosa_access_token=${sessionID}` },
    });
    if (!res.ok) {
        if (res.status === 401 || res.status === 404) {
            // Definitive rejection — wipe the cookie and send to login
            event.cookies.delete(event.locals.sessionCookieName, {
                path: '/',
                domain: event.locals.cookieDomain
            });
            throw redirect(302, handleLoginRedirect(event.url));
        }
        // Transient server error (5xx, etc.) — don't redirect to /auth, that
        // would mislead the user into thinking their session expired. Show an
        // error page instead so they can retry.
        throw error(503, "Le serveur est temporairement indisponible. Veuillez réessayer dans quelques instants.");
    }
    return res.json();
}
