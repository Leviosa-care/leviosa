import type { PageServerLoad } from "./$types"

import type { EventPhotos, EventVideos } from '$lib/types';
import { eventsPhotos, eventsVideos } from '$lib/data';

type PageRes = { eventsPhotos: EventPhotos[]; eventsVideos: EventVideos[] };

export const load: PageServerLoad = (): PageRes => {
    return { eventsPhotos, eventsVideos };
}
