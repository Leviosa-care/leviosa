import type { PageServerLoad } from "./$types"
import { values } from '$lib/data';

export const load: PageServerLoad = () => {
    return { values };
}
