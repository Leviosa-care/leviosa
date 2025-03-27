import type { PageServerLoad } from './$types';

import type { Offer } from '$lib/types';
import { offers } from '$lib/data';

type PageRes = { offers: Offer[] }

export const load: PageServerLoad = (): PageRes => {
    return { offers };
};
