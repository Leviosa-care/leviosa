// TODO: put these types somewhere else
import type { PageServerLoad } from "./$types"

import type { EventPhotos, EventVideos } from '$lib/types';

import { eventsVideos, eventsPhotos } from '$lib/data/media';

type PageRes = { eventsPhotos: EventPhotos[]; eventsVideos: EventVideos[] };

export const load: PageServerLoad = async (): Promise<PageRes> => {
    return { eventsPhotos, eventsVideos };
}
