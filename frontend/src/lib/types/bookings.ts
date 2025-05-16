export type Product = {
    id: string;
}

export type Offer = {
    id: string;
}

export type Event = {
    id: string;
    title: string;
    description: string;
    city: string;
    postalCode: string;
    address1: string;
    address2: string;
    beginAt: string;
    endAt: string;
    products: Product[];
    offers: Offer[];
    day: number;
    month: number;
    year: number;
}

export type Consultation = {
    name: string;
    date: string;
}
