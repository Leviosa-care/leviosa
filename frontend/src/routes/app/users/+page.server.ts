import type { PageServerLoad } from "./$types"
import { redirect } from '@sveltejs/kit';

type PageRes = { role: import('$lib/types').Role };

export const load: PageServerLoad = ({ locals }): PageRes => {
    if (locals.user?.role !== 'admin') throw redirect(302, '/app');
    return { role: locals.user?.role };
}
