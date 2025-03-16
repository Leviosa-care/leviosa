type Card = {
    date: string;
};

const cards: Card[] = [
    {
        date: '12/12/2024'
    }
];

import type { Role } from '$lib/types'

type PageRes = {
    cards: Card[];
    role: Role
};

export function load({ locals }): PageRes {
    const role = locals.user.role as Role
    return {
        role,
        cards,
    };
}
