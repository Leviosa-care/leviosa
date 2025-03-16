export type NavState =
    | 'accueil'
    | 'messages'
    | 'services'
    | 'reservations'
    | 'profil'
    | 'conversations'
    | 'notes de seances';
export type ReservationState = 'consultations' | 'events';
export type EventState = 'Evenements a venir' | 'Reserve ta place' | 'Creer un evenement';
export type ConsultationState =
    | 'Consultations a venir'
    | 'Reserve ta consultation'
    | 'Creer une consulation';
export type MessageState = 'Conversations' | 'Notes de séances';

// =======================
// The service state
// =======================
export type ServiceState = 'A propos' | 'Deroule' | 'Prestataires';
// Convert the union type to an array to get the count
const serviceStates = ['A propos', 'Deroule', 'Prestataires'] as const;
type ServiceStateCount = typeof serviceStates.length;

// This will give the count as a constant value
export const numberOfServiceStates: ServiceStateCount = serviceStates.length;

// =======================
// The signup store (that is a link with the session storage)
// =======================
