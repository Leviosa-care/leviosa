import type { RequestEvent } from './$types';

export async function GET({ url }: RequestEvent): Promise<Response> {
    console.log('I am in the get function brother');
    const state = url.searchParams.get('state');
    console.log('the state that I get from this call is:', state);
    return new Response(null, {
        status: 302,
        headers: {
            Location: '/'
        }
    });
}
