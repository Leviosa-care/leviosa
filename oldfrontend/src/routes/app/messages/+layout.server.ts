import type { LayoutServerLoad } from "./$types"
import type { Message, SessionNote } from '$lib/types';

import { messages, notes } from '$lib/data';

type LayoutRes = { messages: Message[]; notes: SessionNote[] };

export const load: LayoutServerLoad = ({ params }): LayoutRes => {
    // TODO: make a request to get the last conversation that I had
    return { messages, notes };
}
