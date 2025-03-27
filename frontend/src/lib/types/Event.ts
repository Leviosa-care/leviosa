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
export const DAYS_ARRAY = Object.values(DAYS);

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
export const MONTHS_ARRAY = Object.values(MONTHS);

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

// 	Title            string`json:"title"`
// 	Description      string`json:"description"`
// 	City             string`json:"city"`
// 	PostalCode       string`json:"postal_code"`
// 	Address1         string`json:"address1"`
// 	Address2         string`json:"address2"`
// 	PlaceCount       int`json:"place_count"`
// 	FreePlace        int`json:"free_place"`
// 	BeginAt          time.Time`json:"begin_at"`
// 	EncryptedBeginAt string`json:"begin_at_formatted"`
// 	EndAt            time.Time`json:"end_at"`
// 	EncryptedEndAt   string`json:"end_at_formatted"`
// Products[]string`json:"products"`
// Offers[]string`json:"offers"`
// 	PriceID          string`json:"-"`
// 	Day              int`json:"day"`
// 	Month            int`json:"month"`
// 	Year             int`json:"year"`
