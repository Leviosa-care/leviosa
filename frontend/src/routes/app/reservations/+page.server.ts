import type { PageServerLoad } from "./$types"
import type { Role } from '$lib/types'

type Card = {
    date: string;
};

const cards: Card[] = [
    {
        date: '12/12/2024'
    }
];


type PageRes = {
    cards: Card[];
    role: Role
};

export const load: PageServerLoad = ({ locals }): PageRes => {
    const role = locals.user?.role as Role
    return {
        role,
        cards,
    };
}
