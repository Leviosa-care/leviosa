import type { PageServerLoad } from "./$types";
import { env } from "$env/dynamic/private";
import { error } from "@sveltejs/kit";

export const load: PageServerLoad = async ({ fetch, url, locals }) => {
    const bookingId = url.searchParams.get("booking_id");
    if (!bookingId) {
        throw error(400, "Identifiant de réservation manquant");
    }

    // Try to fetch booking details.
    // For authenticated users, this will succeed via session cookie.
    // For guest bookings, we rely on URL search params for display data.
    let booking: any = null;
    let product: any = null;
    let priceDisplay: string | null = null;

    try {
        const bookingRes = await fetch(`${env.API_URL}/bookings/${bookingId}`);
        if (bookingRes.ok) {
            booking = await bookingRes.json();

            // Fetch product details
            if (booking.product_id) {
                const productRes = await fetch(`${env.API_URL}/products/${booking.product_id}`);
                if (productRes.ok) {
                    product = await productRes.json();
                }
            }

            if (booking.total_price_cents) {
                priceDisplay = (booking.total_price_cents / 100)
                    .toFixed(2)
                    .replace(/\.00$/, "");
            }
        }
    } catch {
        // Booking fetch may fail for guests — fall through to URL params
    }

    // Fallback: use URL search params for guest booking display
    if (!booking) {
        booking = {
            id: bookingId,
            product_id: url.searchParams.get("product_id"),
            slot_start_time: url.searchParams.get("slot_start_time"),
            slot_end_time: url.searchParams.get("slot_end_time"),
            total_price_cents: url.searchParams.get("price_cents")
                ? parseInt(url.searchParams.get("price_cents")!)
                : null,
            guest_first_name: url.searchParams.get("guest_first_name"),
            guest_last_name: url.searchParams.get("guest_last_name"),
        };

        // Try to fetch product
        if (booking.product_id) {
            const productRes = await fetch(`${env.API_URL}/products/${booking.product_id}`);
            if (productRes.ok) {
                product = await productRes.json();
            }
        }

        if (booking.total_price_cents) {
            priceDisplay = (booking.total_price_cents / 100)
                .toFixed(2)
                .replace(/\.00$/, "");
        }
    }

    return {
        booking,
        product,
        priceDisplay,
        user: locals.user ?? null,
    };
};
