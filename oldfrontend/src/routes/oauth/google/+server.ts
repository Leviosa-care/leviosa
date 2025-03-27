import { google } from '$lib/server/oauth';
import { generateCodeVerifier, generateState } from 'arctic';

import type { RequestEvent } from './$types';

export function GET(event: RequestEvent): Response {
    const state = generateState();
    const codeVerifier = generateCodeVerifier();
    const scopes = [
        'openid',
        'profile',
        'email',
        'https://www.googleapis.com/auth/userinfo.profile', // Detailed profile
        'https://www.googleapis.com/auth/userinfo.email', // Detailed email
        'https://www.googleapis.com/auth/user.birthday.read', // Birthday
        'https://www.googleapis.com/auth/user.gender.read', // Gender
        'https://www.googleapis.com/auth/user.phonenumbers.read', // Phone numbers
        'https://www.googleapis.com/auth/user.addresses.read' // Addresses
    ];
    const url = google.createAuthorizationURL(state, codeVerifier, { scopes });

    event.cookies.set('google_oauth_state', state, {
        httpOnly: true,
        maxAge: 60 * 10,
        secure: import.meta.env.PROD,
        path: '/',
        sameSite: 'lax'
    });
    event.cookies.set('google_code_verifier', codeVerifier, {
        httpOnly: true,
        maxAge: 60 * 10,
        secure: import.meta.env.PROD,
        path: '/',
        sameSite: 'lax'
    });

    return new Response(null, {
        status: 302,
        headers: {
            Location: url.toString()
        }
    });
}
