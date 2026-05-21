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

export type BookingDTO = {
    id: string
    availability_id: string
    client_id: string
    partner_id: string
    room_id: string
    product_id: string
    slot_start_time: string
    slot_end_time: string
    status: 'confirmed' | 'completed' | 'cancelled' | 'no_show'
    payment_status: 'pending' | 'paid' | 'failed' | 'refunded'
    total_price_cents: number
    currency: string
    payment_intent_id?: string
    client_notes?: string
    partner_notes?: string
    cancellation_reason?: string
    cancelled_at?: string
    completed_at?: string
}

export type Consultation = BookingDTO
