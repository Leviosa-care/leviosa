type PageRes = { role: import('$lib/types').Role };

export function load({ locals }): PageRes {
    return { role: locals.user.role };
}
