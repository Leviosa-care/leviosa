import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';
type PageRes = { role: import('$lib/types').Role };

export const load: LayoutServerLoad = ({ locals }): PageRes => {
    if (!locals.user) throw redirect(302, '/');
    return { role: locals.user?.role };
}
