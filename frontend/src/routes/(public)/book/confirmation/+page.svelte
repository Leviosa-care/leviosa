<script lang="ts">
    import { reveal } from "$lib/actions/reveal";
    import Button from "$lib/ui/Button.svelte";
    import GuestClaimCard from "./GuestClaimCard.svelte";
    import type { PageProps } from "./$types";

    let { data }: PageProps = $props();

    const { booking, product, priceDisplay, user, guestInfo } = data;

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
        : booking.guest_first_name && booking.guest_last_name
          ? `${booking.guest_first_name} ${booking.guest_last_name}`
          : "";
</script>

<div class="bg-surface min-h-screen py-24 md:py-32 px-4 lg:px-8">
    <div class="max-w-2xl mx-auto" use:reveal={{ preset: "fade-up", delay: 100 }}>
        <!-- Success icon -->
        <div class="text-center mb-8">
            <div class="w-20 h-20 rounded-full bg-green-100 flex items-center justify-center mx-auto mb-4">
                <svg class="w-10 h-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
            </div>
            <h1 class="text-3xl md:text-4xl font-bold text-foreground mb-2">
                Réservation Confirmée !
            </h1>
            <p class="text-foreground-alt">
                Votre réservation a été enregistrée avec succès
            </p>
        </div>

        <!-- Booking details card -->
        <div class="bg-white rounded-3xl p-6 md:p-8 shadow-mini">
            <h2 class="text-xl font-bold text-foreground mb-6">Détails de la réservation</h2>

            <div class="grid gap-4">
                {#if product}
                    <div class="flex justify-between py-3 border-b border-border-input">
                        <span class="text-foreground-alt">Service</span>
                        <span class="font-semibold text-foreground">{product.name}</span>
                    </div>
                {/if}

                {#if dateDisplay}
                    <div class="flex justify-between py-3 border-b border-border-input">
                        <span class="text-foreground-alt">Date</span>
                        <span class="font-semibold text-foreground capitalize">{dateDisplay}</span>
                    </div>
                {/if}

                {#if timeDisplay}
                    <div class="flex justify-between py-3 border-b border-border-input">
                        <span class="text-foreground-alt">Horaire</span>
                        <span class="font-semibold text-foreground">{timeDisplay}</span>
                    </div>
                {/if}

                {#if durationMinutes}
                    <div class="flex justify-between py-3 border-b border-border-input">
                        <span class="text-foreground-alt">Durée</span>
                        <span class="font-semibold text-foreground">{durationMinutes} min.</span>
                    </div>
                {/if}

                {#if priceDisplay}
                    <div class="flex justify-between py-3 border-b border-border-input">
                        <span class="text-foreground-alt">Montant</span>
                        <span class="font-semibold text-foreground text-lg">{priceDisplay}€</span>
                    </div>
                {/if}

                {#if displayName}
                    <div class="flex justify-between py-3 border-b border-border-input">
                        <span class="text-foreground-alt">Réservé par</span>
                        <span class="font-semibold text-foreground">{displayName}</span>
                    </div>
                {/if}

                <div class="flex justify-between py-3">
                    <span class="text-foreground-alt">Référence</span>
                    <span class="font-mono text-sm text-foreground-alt">{booking.id}</span>
                </div>
            </div>
        </div>

        <!-- CTA section -->
        <div class="mt-8 grid gap-4">
            {#if !user && guestInfo}
                <!-- Guest: inline account creation card -->
                <GuestClaimCard {guestInfo} guestClaimForm={data.guestClaimForm} guestClaimVerifyForm={data.guestClaimVerifyForm} />

                <!-- Fallback link for no-JS or API-unreachable -->
                <noscript>
                    <div class="text-center">
                        <a href="/auth">
                            <Button class="text-white px-8 py-4 rounded-2xl cursor-pointer">
                                Créer un compte
                            </Button>
                        </a>
                    </div>
                </noscript>
            {:else if !user}
                <!-- Guest but no guest info cookie (e.g. page refreshed) — fallback link -->
                <div class="bg-white rounded-3xl p-6 md:p-8 shadow-mini text-center">
                    <h3 class="text-lg font-semibold text-foreground mb-2">
                        Créer un compte pour suivre vos réservations
                    </h3>
                    <p class="text-foreground-alt mb-6">
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
                <a href="/services" class="text-foreground-alt hover:text-foreground underline">
                    Retour aux services
                </a>
            </div>
        </div>
    </div>
</div>
