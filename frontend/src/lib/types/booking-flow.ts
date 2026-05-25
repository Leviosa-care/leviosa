export type Category = {
    id: string;
    name: string;
    description: string;
    status: string;
    createdAt: string;
    updatedAt: string;
};

export type Product = {
    id: string;
    name: string;
    description: string;
    duration: number;
    category: {
        id: string;
        name: string;
        description: string;
    };
    publishedStatus: string;
    availability: string;
    bufferTime: number;
    cancellationHours: number;
    metadata: Record<string, unknown>;
    createdAt: string;
    updatedAt: string;
};

export type Price = {
    id: string;
    amount: number;
    currency: string;
    interval: string;
};

export type Partner = {
    id: string;
    user_id: string;
    bio: string;
    experience: string;
    category_ids: string[];
    product_ids: string[];
    created_at: string;
    updated_at: string;
};

export type Availability = {
    id: string;
    partner_id: string;
    room_id: string;
    start_time: string;
    end_time: string;
    max_capacity: number;
    current_bookings: number;
    status: string;
    service_type: string;
    price_cents: number | null;
    notes: string;
    created_at: string;
    updated_at: string;
};

export type BookingRequest = {
    availability_id: string;
    product_id: string;
    slot_start_time: string;
    client_notes?: string;
    client_id?: string;
    guest_first_name?: string;
    guest_last_name?: string;
    guest_email?: string;
    guest_phone?: string;
};

export type BookingResponse = {
    id: string;
    availability_id: string;
    client_id?: string;
    partner_id: string;
    room_id: string;
    product_id: string;
    slot_start_time: string;
    slot_end_time: string;
    status: string;
    total_price_cents: number;
    currency: string;
    payment_status: string;
    client_notes?: string;
    guest_first_name?: string;
    guest_last_name?: string;
    guest_email?: string;
    guest_phone?: string;
};
