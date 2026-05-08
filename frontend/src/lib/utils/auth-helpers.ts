import type { Cookies } from '@sveltejs/kit';

/**
 * Extracts and forwards Set-Cookie headers from a backend API response to the client
 *
 * In SvelteKit's SSR architecture:
 * 1. Client → SvelteKit server (frontend)
 * 2. SvelteKit server → Backend API
 * 3. Backend API returns with Set-Cookie headers
 * 4. This function forwards those cookies to the client
 *
 * @param response - The fetch Response object from the backend API
 * @param cookies - SvelteKit's cookies object to set cookies on the client
 */
// BACKEND_ACCESS_COOKIE is the cookie name the backend always sets and expects.
// The browser may store it under a different name (sessionCookieName) so that
// staging and production sessions don't collide on the shared .leviosa.care domain.
const BACKEND_ACCESS_COOKIE = 'leviosa_access_token';

export function forwardAuthCookies(response: Response, cookies: Cookies, sessionCookieName = BACKEND_ACCESS_COOKIE): void {
    const cookieStrings: string[] = response.headers.getSetCookie();

    for (const cookieString of cookieStrings) {
        // Parse cookie name and value
        const [nameValue, ...attributes] = cookieString.split(';').map(part => part.trim());
        const eqIdx = nameValue.indexOf('=');
        if (eqIdx === -1) continue;
        const backendName = nameValue.slice(0, eqIdx);
        const value = nameValue.slice(eqIdx + 1);
        // Rename the access token to the environment-specific browser cookie name.
        const name = backendName === BACKEND_ACCESS_COOKIE ? sessionCookieName : backendName;

        if (!name || !value) continue;

        // Parse cookie attributes
        const options: {
            path?: string;
            domain?: string;
            secure?: boolean;
            httpOnly?: boolean;
            sameSite?: 'strict' | 'lax' | 'none';
            maxAge?: number;
        } = {};

        for (const attr of attributes) {
            const [key, val] = attr.split('=').map(s => s.trim());
            const lowerKey = key.toLowerCase();

            if (lowerKey === 'path') options.path = val;
            else if (lowerKey === 'domain') options.domain = val;
            else if (lowerKey === 'secure') options.secure = true;
            else if (lowerKey === 'httponly') options.httpOnly = true;
            else if (lowerKey === 'samesite') {
                const sameSiteValue = val?.toLowerCase();
                if (sameSiteValue === 'strict' || sameSiteValue === 'lax' || sameSiteValue === 'none') {
                    options.sameSite = sameSiteValue;
                }
            }
            else if (lowerKey === 'max-age') {
                const maxAge = parseInt(val);
                if (!isNaN(maxAge)) options.maxAge = maxAge;
            }
        }

        if (!options.path) continue;

        cookies.set(name, value, { ...options, path: options.path });
    }
}

/**
 * Maps frontend gender values to backend-expected values
 *
 * Frontend uses: 'man' | 'woman' | 'non_binary' | 'prefer_not_to_say' | 'custom' | 'not precised'
 * Backend expects: 'male' | 'female' | 'other' | 'prefer_not_to_say'
 *
 * @param frontendGender - Gender value from frontend form
 * @returns Gender value expected by backend API
 */
export function mapGenderToBackend(frontendGender: string): string {
    const genderMap: Record<string, string> = {
        'man': 'male',
        'woman': 'female',
        'non_binary': 'other',
        'custom': 'other',
        'not precised': 'other',
        'prefer_not_to_say': 'prefer_not_to_say',
    };

    return genderMap[frontendGender] || 'other';
}

/**
 * Formats French phone number to E.164 international format
 *
 * Converts: 0612345678 → +33612345678
 * Already formatted numbers are left unchanged
 *
 * @param phone - Phone number in French format (10 digits starting with 0)
 * @returns Phone number in E.164 format (+33...)
 */
export function formatPhoneToE164(phone: string): string {
    // Remove any spaces, dashes, or dots
    const cleaned = phone.replace(/[\s\-\.]/g, '');

    // If already starts with +, assume it's already formatted
    if (cleaned.startsWith('+')) {
        return cleaned;
    }

    // If starts with 0 (French format), replace with +33
    if (cleaned.startsWith('0') && cleaned.length === 10) {
        return '+33' + cleaned.substring(1);
    }

    // If doesn't start with 0 but is 9 digits, assume it's missing the 0 prefix
    if (cleaned.length === 9 && !cleaned.startsWith('0')) {
        return '+33' + cleaned;
    }

    // Otherwise return as-is (may be invalid, but let backend validate)
    return cleaned;
}
