import type { PageServerLoad } from "./$types"
import type { Conversation, SessionNote } from '$lib/types';

import { conversations, notes } from '$lib/data';

type PageRes = { conversations: Conversation[]; notes: SessionNote[] };

export const load: PageServerLoad = (): PageRes => {
    return { conversations, notes };
}
