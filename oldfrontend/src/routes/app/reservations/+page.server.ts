import type { PageServerLoad } from "./$types"

import type { Role } from '$lib/types'

type EventCard = {
    date: string;
};

type ConsultationCard = {
    date: string;
};

const eventCards: EventCard[] = [
    {
        date: '12/12/2024'
    }
];

const consultationCards: ConsultationCard[] = [
    {
        date: '12/12/2024'
    }
];


type PageRes = {
    eventCards: EventCard[];
    consultationCards: ConsultationCard[];
    role: Role
};

export const load: PageServerLoad = ({ locals }): PageRes => {
    if (!locals.user?.role) {
        console.error('no role set from the local user, this is embarassing')
    }
    const role = locals.user?.role as Role
    console.log("the role is:", role);
    return {
        role,
        eventCards,
        consultationCards,
    };
}
