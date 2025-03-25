import type { PageServerLoad } from "./$types"
import type { Message } from '$lib/types';

import { messages } from '$lib/data';

type PageRes = { messages: Message[] };
export const load: PageServerLoad = ({ params }): PageRes => {
    // TODO: use that to do the fetching for the conversation
    // const id = params.id

    return { messages };
}
