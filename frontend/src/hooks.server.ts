import { redirect, type Handle, type RequestEvent } from '@sveltejs/kit';
import { NODE_ENV, CLIENT_IP_HEADER, SESSION_COOKIE_NAME, API_URL, COMING_SOON } from "$env/static/private"

import { handleLoginRedirect } from '$lib/utils/redirect';
import { mockUser } from '$lib/data/user';

// TODO: for the session thing
// on each request, send the accesstoken or session_token
// when the access token is no longer valid, use will receive 401 status code
// then make a request to a dedicated refresh endpoint on the server and the server resend an access token
// and optionnaly a new refresh  token. This is the rotation step !

export const handle: Handle = async ({ event, resolve }) => {
    if (COMING_SOON === 'true' && event.url.pathname !== '/coming-soon') {
        throw redirect(302, '/coming-soon')
    }

    if (NODE_ENV === 'development' || NODE_ENV === 'staging') {
        event.locals.user = mockUser
        return await resolve(event)
    }

    // Protected routes are any routes that does not start with '/auth' or '/legal'.
    if (event.url.pathname.startsWith("/auth") || event.url.pathname.startsWith("/legal")) {
        return await resolve(event)
    }

    const sessionID = event.cookies.get(SESSION_COOKIE_NAME);
    if (!sessionID) {
        let path = handleLoginRedirect(event.url, "expiredSession");
        const email = event.locals.user.email;
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
            [CLIENT_IP_HEADER]: getClientIP(event.request),
            Authorization: `Bearer ${sessionID}`
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

async function validateSession({ cookies, url }: RequestEvent, sessionID: string): Promise<App.User> {
    const res = await fetch(`${API_URL}/users/me`, {
        method: "GET",
        headers: { Authorization: `Bearer ${sessionID}` },
    });
    if (!res.ok) {
        cookies.set(SESSION_COOKIE_NAME, '', {
            path: '/',
            expires: new Date(0)
        });
        switch (res.status) {
            case 404: // here the user is not found in the database for some reason
            // this is when I should send something to admin to solve that issue
            case 500:
                throw new Error("The server is currently unavailable. We apologize for the inconvenience and are working to resolve the issue as soon as possible.")
        }
        // TODO: maybe add some message to help the user understand what happened ?
        throw redirect(302, handleLoginRedirect(url));
    }
    return res.json();
}
