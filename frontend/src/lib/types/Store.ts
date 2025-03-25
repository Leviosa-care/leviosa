export const NAV_STATES = {
    Accueil: 'accueil',
    Messages: 'messages',
    Services: 'services',
    Reservations: 'reservations',
    Profil: 'profil',
    Conversations: 'conversations',
    NotesDeSeances: 'notes de seances',
} as const;
export type NavState = typeof NAV_STATES[keyof typeof NAV_STATES];

export const RESERVATION_STATES = {
    Consultations: 'consultations',
    Events: 'events',
} as const;
export type ReservationState = typeof RESERVATION_STATES[keyof typeof RESERVATION_STATES];

export const EVENT_STATES = {
    EvenementsAVenir: 'Evenements a venir',
    ReserveTaPlace: 'Reserve ta place',
    CreerUnEvenement: 'Creer un evenement',
} as const;
export type EventState = typeof EVENT_STATES[keyof typeof EVENT_STATES];

export const CONSULTATION_STATES = {
    ConsultationsAVenir: 'Consultations a venir',
    ReserveTaConsultation: 'Reserve ta consultation',
    CreerUneConsultation: 'Creer une consulation',
} as const;
export type ConsultationState = typeof CONSULTATION_STATES[keyof typeof CONSULTATION_STATES];

export const MESSAGE_STATES = {
    Conversations: 'Conversations',
    NotesDeSeances: 'Notes de séances',
} as const;
export type MessageState = typeof MESSAGE_STATES[keyof typeof MESSAGE_STATES];

export const SERVICE_STATES = {
    APropos: 'A propos',
    Deroule: 'Deroule',
    Prestataires: 'Prestataires',
} as const;
export type ServiceState = typeof SERVICE_STATES[keyof typeof SERVICE_STATES];
// Convert the union type to an array to get the count
const serviceStates = Object.values(SERVICE_STATES);
type ServiceStateCount = typeof serviceStates.length;

// This will give the count as a constant value
export const numberOfServiceStates: ServiceStateCount = serviceStates.length;
