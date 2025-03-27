import { google } from '$lib/server/oauth';
import { ObjectParser } from '@pilcrowjs/object-parser';
// import { decodeIdToken } from 'arctic';
// import * as arctic from 'arctic';
import { decodeJWT } from "@oslojs/jwt";


import type { RequestEvent } from './$types';
// import type { OAuth2Tokens } from 'arctic';
import { redirect } from '@sveltejs/kit';

interface GoogleUserInfo {
    googleID: string;
    name: string;
    picture: string;
    email: string;
    birthday?: string;
    gender?: string;
    phoneNumber?: string;
    address?: {
        formatted?: string;
        streetAddress?: string;
        locality?: string;
        region?: string;
        postalCode?: string;
        country?: string;
    };
    given_name?: string;
    family_name?: string;
}

export async function GET({ url, cookies }: RequestEvent): Promise<Response> {
    const storedState = cookies.get('google_oauth_state') ?? null;
    const codeVerifier = cookies.get('google_code_verifier') ?? null;
    const code = url.searchParams.get('code');
    const state = url.searchParams.get('state');

    if (storedState === null || codeVerifier === null || code === null || state === null) {
        return new Response('Please restart the process.', {
            status: 400
        });
    }
    if (storedState !== state) {
        return new Response('Please restart the process.', {
            status: 400
        });
    }

    // let tokens: OAuth2Tokens;
    let tokens;
    try {
        tokens = await google.validateAuthorizationCode(code, codeVerifier);
    } catch (e) {
        return new Response('Please restart the process.', {
            status: 400
        });
    }

    // const claims = decodeIdToken(tokens.idToken());
    let claims: any
    try {
        claims = decodeJWT(String(tokens));
    }
    catch (e) {
        throw new Error("Invalid ID token", {
            cause: e
        });
    }
    const claimsParser = new ObjectParser(claims);

    const userInfo: GoogleUserInfo = {
        googleID: claimsParser.getString('sub'),
        name: claimsParser.getString('name'),
        picture: claimsParser.getString('picture'),
        email: claimsParser.getString('email'),
        given_name: claimsParser.getString('given_name'),
        family_name: claimsParser.getString('family_name')
    };

    // Fetch additional user info using the access token
    try {
        const userInfoResponse = await fetch('https://www.googleapis.com/oauth2/v3/userinfo', {
            headers: {
                Authorization: `Bearer ${tokens.accessToken}`
            }
        });
        const additionalInfo = await userInfoResponse.json();
        console.log("additionnal info:", additionalInfo)

        // Fetch additional profile information
        const peopleApiResponse = await fetch(
            'https://people.googleapis.com/v1/people/me?personFields=birthdays,genders,phoneNumbers,addresses,organizations',
            {
                headers: {
                    Authorization: `Bearer ${tokens.accessToken}`
                }
            }
        );
        const peopleData = await peopleApiResponse.json();

        // Update userInfo with additional data
        if (peopleData.birthdays?.[0]?.date) {
            const date = peopleData.birthdays[0].date;
            userInfo.birthday = `${date.year}-${date.month}-${date.day}`;
        }

        if (peopleData.genders?.[0]?.value) {
            userInfo.gender = peopleData.genders[0].value;
        }

        if (peopleData.phoneNumbers?.[0]?.value) {
            userInfo.phoneNumber = peopleData.phoneNumbers[0].value;
        }

        if (peopleData.addresses?.[0]) {
            userInfo.address = {
                formatted: peopleData.addresses[0].formattedValue,
                streetAddress: peopleData.addresses[0].streetAddress,
                locality: peopleData.addresses[0].city,
                region: peopleData.addresses[0].region,
                postalCode: peopleData.addresses[0].postalCode,
                country: peopleData.addresses[0].country
            };
        }
    } catch (error) {
        console.error('Error fetching additional user info:', error);
        // Continue with basic info even if additional info fetch fails
    }

    console.log('after all the fetching, the userInfo is:', userInfo);

    // TODO: Send to the backend
    // try {
    //     const res = await fetch("http://localhost:3500/api/v1/oauth/google", {
    //         method: 'POST',
    //         headers: {
    //             'Content-Type': 'application/json',
    //         },
    //         body: JSON.stringify(userInfo)
    //     })
    //
    //     const data = await res.json()
    //     console.log("the data received is:", data)
    // } catch (err) {
    //     return new Response("Please restart the process.", {
    //         status: 400
    //     });
    // }

    // TODO: I need to increase the scope, I do not have enough information about the user

    // TODO:
    // fetch the golang backend
    // if user does not exist
    //      - create the user
    //      - add user to the pending list that need admin for validation
    //      - redirect to page for pending users that can visit offline the page
    // if user exist
    //      - set the session for the user (golang backend)
    //      - create a cookie with the received information to send to the user
    //      - redirect to the '/app' page for the user

    // NOTE: the part from the code that I am going to delete
    // const existingUser = getUserFromGoogleId(googleId);
    // if (existingUser !== null) {
    //     const sessionToken = generateSessionToken();
    //     const session = createSession(sessionToken, existingUser.id);
    //     setSessionTokenCookie(event, sessionToken, session.expiresAt);
    //     return new Response(null, {
    //         status: 302,
    //         headers: {
    //             Location: "/"
    //         }
    //     });
    // }
    //
    // const user = createUser(googleId, email, name, picture);
    // const sessionToken = generateSessionToken();
    // const session = createSession(sessionToken, user.id);
    // setSessionTokenCookie(event, sessionToken, session.expiresAt);

    throw redirect(302, '/app');
}
