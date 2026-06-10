<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import { enhance } from '$app/forms';
	import type { PageProps } from './$types';
	import { Eye, Clock, CheckCircle, XCircle, AlertCircle, Search, Phone as PhoneIcon, Trash2, Loader2 } from '@lucide/svelte';

	let { data, form }: PageProps = $props();

	// Contact method toggle for the manual lookup form
	let contactMethod: 'email' | 'phone' = $state('email');

	// Cancel confirmation state
	let showCancelConfirm = $state(false);
	let cancelling = $state(false);

	// Booking detail (from either token path or manual form action)
	let booking = $derived(data.booking ?? (form?.action === 'lookup' && form.success ? form.booking : null));
	let displayError = $derived(data.lookupError ?? (form?.action === 'lookup' && !form.success ? form.error : null));
	let cancelError = $derived(form?.action === 'cancel' && !form.success ? form.error : null);

	// If the cancel action succeeded, update the booking in-place.
	// Spread original booking first so product_name / partner_name are preserved —
	// the cancel response is a minimal PublicBookingLookupResponse that omits those fields.
	let effectiveBooking = $derived(
		form?.action === 'cancel' && form.success && form.booking
			? { ...booking, ...form.booking }
			: booking
	);

	// Cancel requires a booking token — only available on the token URL path
	let canCancel = $derived(
		!!data.token &&
		effectiveBooking &&
		effectiveBooking.status === 'confirmed' &&
		new Date(effectiveBooking.slot_start_time) > new Date()
	);

	function formatDate(iso: string): string {
		return new Date(iso).toLocaleDateString('fr-FR', {
			weekday: 'long',
			day: 'numeric',
			month: 'long',
			year: 'numeric',
		});
	}

	function formatTime(iso: string): string {
		return new Date(iso).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
	}

	function formatCents(cents: number): string {
		return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(cents / 100);
	}

	function statusLabel(status: string): string {
		switch (status) {
			case 'confirmed': return 'Confirmée';
			case 'completed': return 'Terminée';
			case 'cancelled': return 'Annulée';
			case 'no_show': return 'Non présenté';
			default: return status;
		}
	}

	function statusColor(status: string): string {
		switch (status) {
			case 'confirmed': return 'text-blue-700 bg-blue-100';
			case 'completed': return 'text-green-700 bg-green-100';
			case 'cancelled': return 'text-red-700 bg-red-100';
			case 'no_show': return 'text-orange-700 bg-orange-100';
			default: return 'text-gray-700 bg-gray-100';
		}
	}
</script>

<svelte:head>
	<title>Mes Réservations | Leviosa</title>
</svelte:head>

<div class="min-h-screen bg-surface py-24 md:py-32 px-4 lg:px-8">
	<div class="max-w-2xl mx-auto" use:reveal={{ preset: "fade-up", delay: 100 }}>

		{#if effectiveBooking}
			<!-- ═══ Booking detail view ═══ -->
			<div class="text-center mb-8">
				<h1 class="text-3xl md:text-4xl font-bold text-foreground mb-2">
					Votre Réservation
				</h1>
			</div>

			<div class="bg-white rounded-3xl p-6 md:p-8 shadow-mini">
				<div class="flex items-center gap-3 mb-6">
					<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium {statusColor(effectiveBooking.status)}">
						{#if effectiveBooking.status === 'confirmed' || effectiveBooking.status === 'completed'}
							<CheckCircle size={16} />
						{:else if effectiveBooking.status === 'cancelled'}
							<XCircle size={16} />
						{:else}
							<Clock size={16} />
						{/if}
						{statusLabel(effectiveBooking.status)}
					</span>
				</div>

				<div class="grid gap-4">
					{#if effectiveBooking.product_name}
						<div class="flex justify-between py-3 border-b border-border-input">
							<span class="text-foreground-alt">Service</span>
							<span class="font-semibold text-foreground">{effectiveBooking.product_name}</span>
						</div>
					{/if}

					{#if effectiveBooking.partner_name}
						<div class="flex justify-between py-3 border-b border-border-input">
							<span class="text-foreground-alt">Praticien</span>
							<span class="font-semibold text-foreground">{effectiveBooking.partner_name}</span>
						</div>
					{/if}

					<div class="flex justify-between py-3 border-b border-border-input">
						<span class="text-foreground-alt">Date</span>
						<span class="font-semibold text-foreground capitalize">{formatDate(effectiveBooking.slot_start_time)}</span>
					</div>

					<div class="flex justify-between py-3 border-b border-border-input">
						<span class="text-foreground-alt">Horaire</span>
						<span class="font-semibold text-foreground">
							{formatTime(effectiveBooking.slot_start_time)} — {formatTime(effectiveBooking.slot_end_time)}
						</span>
					</div>

					{#if effectiveBooking.total_price_cents}
						<div class="flex justify-between py-3 border-b border-border-input">
							<span class="text-foreground-alt">Montant</span>
							<span class="font-semibold text-foreground text-lg">{formatCents(effectiveBooking.total_price_cents)}</span>
						</div>
					{/if}

					<div class="flex justify-between py-3">
						<span class="text-foreground-alt">Référence</span>
						<span class="font-mono text-sm text-foreground-alt">{effectiveBooking.id}</span>
					</div>
				</div>
			</div>

			<!-- ═══ Cancel section ═══ -->
			{#if canCancel}
				{#if showCancelConfirm}
					<div class="mt-6 bg-red-50 border border-red-200 rounded-2xl p-6">
						<p class="text-red-800 font-medium mb-4">Confirmer l'annulation ?</p>
						<p class="text-red-700 text-sm mb-5">Cette action est irréversible. Votre créneau sera libéré.</p>

						{#if cancelError}
							<div class="flex items-center gap-2 px-4 py-3 mb-4 rounded-lg bg-red-100 border border-red-300 text-red-800 text-sm">
								<AlertCircle size={16} class="flex-shrink-0" />
								{cancelError}
							</div>
						{/if}

						<form
							method="POST"
							action="?action=cancel"
							class="flex gap-3"
							use:enhance={() => {
								cancelling = true;
								return async ({ update }) => {
									await update({ reset: false });
									cancelling = false;
									showCancelConfirm = false;
								};
							}}
						>
							<input type="hidden" name="booking_id" value={effectiveBooking.id} />
							<input type="hidden" name="token" value={data.token ?? ''} />
							<input type="hidden" name="reason" value="Annulation par le client" />

							<button
								type="submit"
								disabled={cancelling}
								class="flex-1 flex items-center justify-center gap-2 bg-red-600 text-white py-2.5 rounded-xl font-medium hover:bg-red-700 transition-colors disabled:opacity-50 cursor-pointer"
							>
								{#if cancelling}
									<Loader2 size={18} class="animate-spin" />
									Annulation...
								{:else}
									<Trash2 size={18} />
									Confirmer l'annulation
								{/if}
							</button>

							<button
								type="button"
								disabled={cancelling}
								class="px-5 py-2.5 rounded-xl font-medium text-foreground-alt bg-white border border-border-input-hover hover:bg-surface transition-colors disabled:opacity-50 cursor-pointer"
								onclick={() => { showCancelConfirm = false; }}
							>
								Retour
							</button>
						</form>
					</div>
				{:else}
					<div class="mt-6 text-center">
						<button
							type="button"
							class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl font-medium text-red-600 bg-red-50 border border-red-200 hover:bg-red-100 transition-colors cursor-pointer"
							onclick={() => { showCancelConfirm = true; }}
						>
							<Trash2 size={18} />
							Annuler cette réservation
						</button>
					</div>
				{/if}
			{/if}

			<div class="mt-8 text-center">
				<a href="/services" class="text-foreground-alt hover:text-foreground underline">
					Retour aux services
				</a>
			</div>

		{:else}
			<!-- ═══ Manual lookup form ═══ -->
			<div class="text-center mb-8">
				<h1 class="text-3xl md:text-4xl font-bold text-foreground mb-2">
					Retrouver ma réservation
				</h1>
				<p class="text-foreground-alt">
					Entrez votre référence et vos coordonnées pour accéder à votre réservation
				</p>
			</div>

			{#if displayError}
				<div class="flex items-center gap-2 px-4 py-3 mb-6 rounded-lg bg-red-50 border border-red-200 text-red-700 text-sm">
					<AlertCircle size={16} class="flex-shrink-0" />
					{displayError}
				</div>
			{/if}

			<form
				method="POST"
				action="?action=lookup"
				class="bg-white rounded-3xl p-6 md:p-8 shadow-mini space-y-5"
				use:enhance
			>
				<!-- Booking reference -->
				<div>
					<label for="ref" class="block text-sm font-medium text-foreground-alt mb-1.5">
						Référence de réservation
					</label>
					<input
						id="ref"
						name="ref"
						type="text"
						placeholder="ex. a1b2c3d4-e5f6-7890-abcd-ef1234567890"
						class="w-full px-4 py-3 rounded-xl border border-border-input-hover text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent text-sm font-mono"
						required
					/>
				</div>

				<!-- Contact method toggle -->
				<div>
					<span class="block text-sm font-medium text-foreground-alt mb-2">Vérification</span>
					<div class="flex gap-2 mb-3">
						<button
							type="button"
							class="flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition-colors {contactMethod === 'email'
								? 'bg-primary/10 text-primary border border-primary/30'
								: 'bg-surface text-foreground-alt border border-border-input-hover hover:bg-surface-hover'}"
							onclick={() => { contactMethod = 'email'; }}
						>
							<Eye size={14} />
							Email
						</button>
						<button
							type="button"
							class="flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition-colors {contactMethod === 'phone'
								? 'bg-primary/10 text-primary border border-primary/30'
								: 'bg-surface text-foreground-alt border border-border-input-hover hover:bg-surface-hover'}"
							onclick={() => { contactMethod = 'phone'; }}
						>
							<PhoneIcon size={14} />
							Téléphone
						</button>
					</div>

					{#if contactMethod === 'email'}
						<input
							name="email"
							type="email"
							placeholder="votre@email.com"
							class="w-full px-4 py-3 rounded-xl border border-border-input-hover text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent text-sm"
						/>
					{:else}
						<input
							name="phone"
							type="tel"
							placeholder="+33 6 12 34 56 78"
							class="w-full px-4 py-3 rounded-xl border border-border-input-hover text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent text-sm"
						/>
					{/if}
				</div>

				<button
					type="submit"
					class="w-full flex items-center justify-center gap-2 bg-primary text-white py-3 rounded-xl font-medium hover:bg-primary/90 transition-colors cursor-pointer"
				>
					<Search size={18} />
					Rechercher
				</button>
			</form>

			<div class="mt-6 text-center text-sm text-muted-foreground">
				<p>Un lien direct vous a été envoyé par email ou SMS lors de votre réservation.</p>
			</div>
		{/if}
	</div>
</div>
