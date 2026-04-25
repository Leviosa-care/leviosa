import { redirect, type Handle, type RequestEvent } from '@sveltejs/kit';
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
        console.error("Error validating session:", err)
    }

    // enrich fetch with custom header with client IP
    event.fetch = async (input, init = {}) => {
        // Ensure headers object exists
        init.headers = {
            ...(init.headers ?? {}),
            [env.CLIENT_IP_HEADER ?? 'x-client-ip']: getClientIP(event.request),
            Authorization: `Bearer ${sessionID}`
        };
        return fetch(input, init);
    };

    // Store session cookie name for use in layouts/pages
    event.locals.sessionCookieName = sessionCookieName;
    event.locals.cookieDomain = cookieDomain ?? undefined;

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
        headers: { Authorization: `Bearer ${sessionID}` },
    });
    if (!res.ok) {
        event.cookies.set(event.locals.sessionCookieName, '', {
            path: '/',
            domain: event.locals.cookieDomain,
            expires: new Date(0)
        });
        switch (res.status) {
            case 404: // here the user is not found in the database for some reason
            // this is when I should send something to admin to solve that issue
            case 500:
                throw new Error("The server is currently unavailable. We apologize for the inconvenience and are working to resolve the issue as soon as possible.")
        }
        // TODO: maybe add some message to help the user understand what happened ?
        throw redirect(302, handleLoginRedirect(event.url));
    }
    return res.json();
}
