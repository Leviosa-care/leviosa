export type Note = 1 | 2 | 3 | 4 | 5;

export type Review = {
    note: Note;
    count: number;
};

export type Prestation = {
    id: string;
    text: string;
    imageUrl: string;
};

export type Freelancer = {
    id: string;
    firstname: string;
    lastname: string;
    avatar: string;
    reviews?: Review[];
    ratings?: number;
    bio?: string;
};

// TODO: add the freelancers line and remove freelancers_names
export type Service = {
    name: string;
    label: string;
    description: string;
    image: string;
    positive_responses: number;
    clients_count: number;
    duration: number;
    freelancers: Freelancer[];
    bgurl: string;
    bgblur: number;
    isLight: boolean;
};

export type Offer = {
    type: string;
    services: Service[];
};

export type getOffersRes = {
    name: string;
    action: () => void;
};
