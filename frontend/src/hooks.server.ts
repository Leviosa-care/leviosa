import { redirect, type Handle } from '@sveltejs/kit';

import { mockUser } from '$lib/data/user';

import { env } from "$env/dynamic/private"
import { NODE_ENV } from "$env/static/private"

export const handle: Handle = async ({ event, resolve }) => {
    // enrich the fetch with the custom header with the client IP
    console.log("running the app in mode:", NODE_ENV)
    const client_IP_header = env.CLIENT_IP_HEADER || "31.111.93.187"
    if (!client_IP_header) {
        if (NODE_ENV === 'development') {
            console.warn("Client IP header missing for backend (ignored in development).");
        } else {
            throw new Error("Logging impossible: client IP header missing for backend.");
        }
    }
    event.fetch = async (input, init = {}) => {
        // Ensure headers object exists
        init.headers = {
            ...init.headers,
            client_IP_header: getClientIP(event.request)
        };

        // Proceed with the original fetch
        return fetch(input, init);
    };
    if (NODE_ENV === 'development') {
        event.locals.user = mockUser
        return await resolve(event)
    }
    // else {
    const signupProgress = getCookie(env.SIGNUP_COOKIE_NAME, event)
    const sessionID = getCookie(env.SESSION_COOKIE_NAME, event)
    if (signupProgress === "step1")
        throw redirect(302, "/signup/verify-email")
    if (signupProgress === "step2")
        throw redirect(302, "/signup/pending")
    if (!sessionID || !signupProgress) {
        throw redirect(302, "/")
    }
    const api_url = env.API_URL
    if (!api_url) {
        throw new Error("API URL not set, can not proceed")
    }
    if (sessionID) {
        try {
            const user = await validateSession(sessionID, event)
            event.locals.user = { ...user };
        } catch (err) {
            console.error("Error validating session:", err)
        }
    }
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

function getCookie(cookieName: string | undefined, event: any) {
    if (!cookieName) {
        return undefined
    }
    return event.cookies.get(cookieName)
}

async function validateSession(sessionID: string | undefined, event: any) {
    const res = await fetch(`${env.API_URL}/user/me`, {
        headers: {
            Authorization: `Bearer ${sessionID}`
        }
    });
    const sessionCookieName = env.SESSION_COOKIE_NAME
    if (!sessionCookieName) {
        throw redirect(302, '/');
    }
    if (res.status === 401) {
        console.log('the status is so that I should redirect the user and cancel the previous cookie.');
        event.cookies.set(sessionCookieName, '', {
            path: '/',
            expires: new Date(0)
        });
        throw redirect(302, '/');
    }
    return res.json();
}

// NOTE: the old thing used to redirect to the soon page
// if (event.url.pathname === '/soon') {
//     return resolve(event);
// }
// throw redirect(307, '/soon');

