<script lang="ts">
    import type { PageProps } from "./$types";
    import type { BookingDTO } from "$lib/types/bookings";

    let { data }: PageProps = $props();

    interface Props {
        bookings: Promise<BookingDTO[]>;
    }
    const { bookings }: Props = data;

    const formatDate = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleDateString('fr-FR', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
        });
    };

    const formatTime = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleTimeString('fr-FR', {
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    const formatPrice = (cents: number, currency: string) => {
        return new Intl.NumberFormat('fr-FR', {
            style: 'currency',
            currency: currency,
        }).format(cents / 100);
    };

    const getStatusBadgeClass = (status: string) => {
        switch (status) {
            case 'confirmed':
                return 'bg-green-100 text-green-800';
            case 'completed':
                return 'bg-blue-100 text-blue-800';
            case 'cancelled':
                return 'bg-red-100 text-red-800';
            case 'no_show':
                return 'bg-yellow-100 text-yellow-800';
            default:
                return 'bg-gray-100 text-gray-800';
        }
    };

    const getStatusLabel = (status: string) => {
        switch (status) {
            case 'confirmed':
                return 'Confirmé';
            case 'completed':
                return 'Terminé';
            case 'cancelled':
                return 'Annulé';
            case 'no_show':
                return 'Absent';
            default:
                return status;
        }
    };
</script>

<div class="bookings-list">
    <h1>Mes consultations</h1>

    {#await bookings}
        <p>Chargement...</p>
    {:then bookings}
        {#if bookings.length === 0}
            <p>Aucune consultation trouvée.</p>
        {:else}
            <div class="bookings-grid">
                {#each bookings as booking}
                    <div class="booking-card">
                        <div class="booking-header">
                            <h3>Produit: {booking.product_id}</h3>
                            <span class="status-badge {getStatusBadgeClass(booking.status)}">
                                {getStatusLabel(booking.status)}
                            </span>
                        </div>

                        <div class="booking-details">
                            <div class="detail-row">
                                <span class="label">Date:</span>
                                <span class="value">{formatDate(booking.slot_start_time)}</span>
                            </div>

                            <div class="detail-row">
                                <span class="label">Heure:</span>
                                <span class="value">
                                    {formatTime(booking.slot_start_time)} - {formatTime(booking.slot_end_time)}
                                </span>
                            </div>

                            <div class="detail-row">
                                <span class="label">Praticien:</span>
                                <span class="value">{booking.partner_id}</span>
                            </div>

                            <div class="detail-row">
                                <span class="label">Prix:</span>
                                <span class="value">{formatPrice(booking.total_price_cents, booking.currency)}</span>
                            </div>

                            {#if booking.client_notes}
                                <div class="detail-row">
                                    <span class="label">Notes:</span>
                                    <span class="value">{booking.client_notes}</span>
                                </div>
                            {/if}
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    {:catch error}
        <p>Erreur lors du chargement des consultations.</p>
    {/await}
</div>

<style>
    .bookings-list {
        padding: 2rem;
        max-width: 1200px;
        margin: 0 auto;
    }

    h1 {
        font-size: 1.5rem;
        font-weight: 600;
        margin-bottom: 1.5rem;
    }

    .bookings-grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
        gap: 1.5rem;
    }

    .booking-card {
        border: 1px solid #e5e7eb;
        border-radius: 0.5rem;
        padding: 1.5rem;
        background-color: white;
    }

    .booking-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1rem;
    }

    .booking-header h3 {
        font-size: 1rem;
        font-weight: 600;
        color: #1f2937;
    }

    .status-badge {
        padding: 0.25rem 0.75rem;
        border-radius: 9999px;
        font-size: 0.75rem;
        font-weight: 500;
        text-transform: uppercase;
    }

    .booking-details {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
    }

    .detail-row {
        display: flex;
        justify-content: space-between;
        font-size: 0.875rem;
    }

    .label {
        color: #6b7280;
        font-weight: 500;
    }

    .value {
        color: #1f2937;
        font-weight: 400;
    }
</style>
