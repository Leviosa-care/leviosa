<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
    import Button from "$lib/ui/Button.svelte";
    import type { PageProps } from "./$types";

    let { data }: PageProps = $props();

    const { booking, product, priceDisplay, user } = data;

    // Guard against null/invalid dates from the URL-param fallback path
    const rawStart = booking.slot_start_time;
    const rawEnd = booking.slot_end_time;
    const slotStart = rawStart ? new Date(rawStart) : null;
    const slotEnd = rawEnd ? new Date(rawEnd) : null;
    const datesValid = slotStart !== null && !isNaN(slotStart.getTime());

    const dateDisplay = datesValid
        ? slotStart!.toLocaleDateString("fr-FR", {
              weekday: "long",
              day: "numeric",
              month: "long",
              year: "numeric",
          })
        : null;

    const timeDisplay =
        datesValid && slotEnd
            ? `${slotStart!.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })} — ${slotEnd.toLocaleTimeString("fr-FR", { hour: "2-digit", minute: "2-digit" })}`
            : null;

    const durationMinutes =
        datesValid && slotEnd && product
            ? Math.round((slotEnd.getTime() - slotStart!.getTime()) / 60000)
            : null;

    const displayName = user
        ? `${user.firstname} ${user.lastname}`
        : "";
</script>

<div class="bg-dark-50 min-h-screen py-24 md:py-32 px-4 lg:px-8">
    <div class="max-w-2xl mx-auto" use:reveal={{ preset: "fade-up", delay: 100 }}>
        <!-- Success icon -->
        <div class="text-center mb-8">
            <div class="w-20 h-20 rounded-full bg-green-100 flex items-center justify-center mx-auto mb-4">
                <svg class="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
            </div>
            <h1 class="text-3xl md:text-4xl font-bold text-dark-900 mb-2">
                Réservation Confirmée !
            </h1>
            <p class="text-dark-600">
                Votre réservation a été enregistrée avec succès
            </p>
        </div>

        <!-- Booking details card -->
        <div class="bg-white rounded-3xl p-6 md:p-8 shadow-sm">
            <h2 class="text-xl font-bold text-dark-900 mb-6">Détails de la réservation</h2>

            <div class="grid gap-4">
                {#if product}
                    <div class="flex justify-between py-3 border-b border-dark-100">
                        <span class="text-dark-600">Service</span>
                        <span class="font-semibold text-dark-900">{product.name}</span>
                    </div>
                {/if}

                {#if dateDisplay}
                    <div class="flex justify-between py-3 border-b border-dark-100">
                        <span class="text-dark-600">Date</span>
                        <span class="font-semibold text-dark-900 capitalize">{dateDisplay}</span>
                    </div>
                {/if}

                {#if timeDisplay}
                    <div class="flex justify-between py-3 border-b border-dark-100">
                        <span class="text-dark-600">Horaire</span>
                        <span class="font-semibold text-dark-900">{timeDisplay}</span>
                    </div>
                {/if}

                {#if durationMinutes}
                    <div class="flex justify-between py-3 border-b border-dark-100">
                        <span class="text-dark-600">Durée</span>
                        <span class="font-semibold text-dark-900">{durationMinutes} min.</span>
                    </div>
                {/if}

                {#if priceDisplay}
                    <div class="flex justify-between py-3 border-b border-dark-100">
                        <span class="text-dark-600">Montant</span>
                        <span class="font-semibold text-dark-900 text-lg">{priceDisplay}€</span>
                    </div>
                {/if}

                {#if displayName}
                    <div class="flex justify-between py-3 border-b border-dark-100">
                        <span class="text-dark-600">Réservé par</span>
                        <span class="font-semibold text-dark-900">{displayName}</span>
                    </div>
                {/if}

                <div class="flex justify-between py-3">
                    <span class="text-dark-600">Référence</span>
                    <span class="font-mono text-sm text-dark-700">{booking.id}</span>
                </div>
            </div>
        </div>

        <!-- CTA section -->
        <div class="mt-8 grid gap-4">
            {#if !user}
                <!-- Guest: prompt to create account -->
                <div class="bg-white rounded-3xl p-6 md:p-8 shadow-sm text-center">
                    <h3 class="text-lg font-semibold text-dark-900 mb-2">
                        Créer un compte pour suivre vos réservations
                    </h3>
                    <p class="text-dark-600 mb-6">
                        Créez un compte pour accéder à l'historique de vos réservations, les modifier et recevoir des rappels.
                    </p>
                    <a href="/auth">
                        <Button class="text-white px-8 py-4 rounded-2xl cursor-pointer">
                            Créer un compte
                        </Button>
                    </a>
                </div>
            {:else}
                <!-- Authenticated: link to bookings -->
                <div class="text-center">
                    <a href="/client/bookings">
                        <Button class="text-white px-8 py-4 rounded-2xl cursor-pointer">
                            Voir mes réservations
                        </Button>
                    </a>
                </div>
            {/if}

            <div class="text-center">
                <a href="/services" class="text-dark-600 hover:text-dark-900 underline">
                    Retour aux services
                </a>
            </div>
        </div>
    </div>
</div>
