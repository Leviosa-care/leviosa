<script lang="ts">
    import type { PageProps } from "./$types";

    import type { Event, BookingDTO } from "$lib/types/bookings";

    let { data }: PageProps = $props();

    interface Props {
        bookings: Promise<BookingDTO[]>;
        events?: Promise<Event[]>;
    }
    const { bookings, events }: Props = data;
</script>

{@render handleBookings()}
{@render handleEvent()}

{#snippet handleBookings()}
    {#await bookings}
        <p>Chargement...</p>
    {:then bookings}
        {#each bookings as booking}
            <div>{booking.product_id}</div>
        {/each}
    {/await}
{/snippet}

{#snippet handleEvent()}
    {#await events}
        <p>Chargement...</p>
    {:then events}
        {#if events}
            {#each events as event}
                <div>{event.title}</div>
            {/each}
        {/if}
    {/await}
{/snippet}
