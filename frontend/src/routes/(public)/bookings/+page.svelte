<script lang="ts">
	import { reveal } from '$lib/actions/reveal';
	import { enhance } from '$app/forms';
	import type { PageProps } from './$types';
	import { Eye, Clock, CheckCircle, XCircle, AlertCircle, Search, Phone as PhoneIcon } from '@lucide/svelte';

	let { data, form }: PageProps = $props();

	// Contact method toggle for the manual lookup form
	let contactMethod: 'email' | 'phone' = $state('email');

	// Booking detail (from either token path or manual form action)
	let booking = $derived(data.booking ?? form?.booking ?? null);
	let displayError = $derived(data.lookupError ?? (form && !form.success ? form.error : null));

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

<div class="min-h-screen bg-dark-50 py-24 md:py-32 px-4 lg:px-8">
	<div class="max-w-2xl mx-auto" use:reveal={{ preset: "fade-up", delay: 100 }}>

		{#if booking}
			<!-- ═══ Booking detail view ═══ -->
			<div class="text-center mb-8">
				<h1 class="text-3xl md:text-4xl font-bold text-dark-900 mb-2">
					Votre Réservation
				</h1>
			</div>

			<div class="bg-white rounded-3xl p-6 md:p-8 shadow-sm">
				<div class="flex items-center gap-3 mb-6">
					<span class="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-medium {statusColor(booking.status)}">
						{#if booking.status === 'confirmed' || booking.status === 'completed'}
							<CheckCircle size={16} />
						{:else if booking.status === 'cancelled'}
							<XCircle size={16} />
						{:else}
							<Clock size={16} />
						{/if}
						{statusLabel(booking.status)}
					</span>
				</div>

				<div class="grid gap-4">
					{#if booking.product_name}
						<div class="flex justify-between py-3 border-b border-dark-100">
							<span class="text-dark-600">Service</span>
							<span class="font-semibold text-dark-900">{booking.product_name}</span>
						</div>
					{/if}

					{#if booking.partner_name}
						<div class="flex justify-between py-3 border-b border-dark-100">
							<span class="text-dark-600">Praticien</span>
							<span class="font-semibold text-dark-900">{booking.partner_name}</span>
						</div>
					{/if}

					<div class="flex justify-between py-3 border-b border-dark-100">
						<span class="text-dark-600">Date</span>
						<span class="font-semibold text-dark-900 capitalize">{formatDate(booking.slot_start_time)}</span>
					</div>

					<div class="flex justify-between py-3 border-b border-dark-100">
						<span class="text-dark-600">Horaire</span>
						<span class="font-semibold text-dark-900">
							{formatTime(booking.slot_start_time)} — {formatTime(booking.slot_end_time)}
						</span>
					</div>

					{#if booking.total_price_cents}
						<div class="flex justify-between py-3 border-b border-dark-100">
							<span class="text-dark-600">Montant</span>
							<span class="font-semibold text-dark-900 text-lg">{formatCents(booking.total_price_cents)}</span>
						</div>
					{/if}

					<div class="flex justify-between py-3">
						<span class="text-dark-600">Référence</span>
						<span class="font-mono text-sm text-dark-700">{booking.id}</span>
					</div>
				</div>
			</div>

			<div class="mt-8 text-center">
				<a href="/services" class="text-dark-600 hover:text-dark-900 underline">
					Retour aux services
				</a>
			</div>

		{:else}
			<!-- ═══ Manual lookup form ═══ -->
			<div class="text-center mb-8">
				<h1 class="text-3xl md:text-4xl font-bold text-dark-900 mb-2">
					Retrouver ma réservation
				</h1>
				<p class="text-dark-600">
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
				class="bg-white rounded-3xl p-6 md:p-8 shadow-sm space-y-5"
				use:enhance
			>
				<!-- Booking reference -->
				<div>
					<label for="ref" class="block text-sm font-medium text-dark-700 mb-1.5">
						Référence de réservation
					</label>
					<input
						id="ref"
						name="ref"
						type="text"
						placeholder="ex. a1b2c3d4-e5f6-7890-abcd-ef1234567890"
						class="w-full px-4 py-3 rounded-xl border border-dark-200 text-dark-900 placeholder:text-dark-400 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent text-sm font-mono"
						required
					/>
				</div>

				<!-- Contact method toggle -->
				<div>
					<span class="block text-sm font-medium text-dark-700 mb-2">Vérification</span>
					<div class="flex gap-2 mb-3">
						<button
							type="button"
							class="flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition-colors {contactMethod === 'email'
								? 'bg-primary/10 text-primary border border-primary/30'
								: 'bg-dark-50 text-dark-600 border border-dark-200 hover:bg-dark-100'}"
							onclick={() => { contactMethod = 'email'; }}
						>
							<Eye size={14} />
							Email
						</button>
						<button
							type="button"
							class="flex items-center gap-1.5 px-4 py-2 rounded-lg text-sm font-medium transition-colors {contactMethod === 'phone'
								? 'bg-primary/10 text-primary border border-primary/30'
								: 'bg-dark-50 text-dark-600 border border-dark-200 hover:bg-dark-100'}"
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
							class="w-full px-4 py-3 rounded-xl border border-dark-200 text-dark-900 placeholder:text-dark-400 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent text-sm"
						/>
					{:else}
						<input
							name="phone"
							type="tel"
							placeholder="+33 6 12 34 56 78"
							class="w-full px-4 py-3 rounded-xl border border-dark-200 text-dark-900 placeholder:text-dark-400 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent text-sm"
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

			<div class="mt-6 text-center text-sm text-dark-500">
				<p>Un lien direct vous a été envoyé par email ou SMS lors de votre réservation.</p>
			</div>
		{/if}
	</div>
</div>
