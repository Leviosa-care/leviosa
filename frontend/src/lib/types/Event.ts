export const DAYS = {
    Lundi: 'Lundi',
    Mardi: 'Mardi',
    Mercredi: 'Mercredi',
    Jeudi: 'Jeudi',
    Vendredi: 'Vendredi',
    Samedi: 'Samedi',
    Dimanche: 'Dimanche',
} as const;
export type Day = typeof DAYS[keyof typeof DAYS];

export const MONTHS = {
    Janvier: 'Janvier',
    Fevrier: 'Fevrier',
    Mars: 'Mars',
    Avril: 'Avril',
    Mai: 'Mai',
    Juin: 'Juin',
    Juillet: 'Juillet',
    Aout: 'Aout',
    Septembre: 'Septembre',
    Octobre: 'Octobre',
    Novembre: 'Novembre',
    Decembre: 'Decembre',
} as const;
export type Month = typeof MONTHS[keyof typeof MONTHS];

// TODO: find a type for days that goes from 1 to 31
// TODO: find a type for hours that goes from 1 to 24
type EventDay = {
    day: Day;
    date: number;
    hours: number[];
};

export type EventPickerMonth = {
    month: Month;
    days: EventDay[];
};

type EventDescription = {
    id: string;
    date: string;
};
type Partner = {
    name: string;
    description: string;
    logo: string;
    website: string;
    instagram: string;
    facebook: string;
};
export type Event = {
    event: EventDescription;
    images: string[];
    partners: Partner[];
};

export type EventInformation = {
    name: string;
    address: string;
    postalCode: string;
    city: string;
    day: string;
    date: number;
    month: string;
    time: string;
    mapImg: string;
    headerImg: string;
};
