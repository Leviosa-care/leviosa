<script lang="ts">
	import type { PageProps } from './$types';
	import { invalidateAll } from '$app/navigation';
	import {
		Clock,
		MapPin,
		FileText,
		Save,
		X,
		Pencil,
		AlertCircle,
		XCircle,
		ArrowLeft,
		MessageSquare,
	} from '@lucide/svelte';

	let { data }: PageProps = $props();
	let booking = $derived(data.booking);

	const now = new Date();
	let cancelling = $state(false);
	let cancelError = $state('');

	let editingNotes = $state(false);
	let notesValue = $state(booking.client_notes ?? '');
	let notesOverride = $state(booking.client_notes ?? '');
	let savingNotes = $state(false);
	let notesError = $state('');

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

	function statusBadge(status: string): string {
		switch (status) {
			case 'confirmed': return 'bg-blue-100 text-blue-700';
			case 'completed': return 'bg-green-100 text-green-700';
			case 'cancelled': return 'bg-red-100 text-red-700';
			case 'no_show': return 'bg-orange-100 text-orange-700';
			default: return 'bg-gray-100 text-gray-700';
		}
	}

	function statusLabel(status: string): string {
		switch (status) {
			case 'confirmed': return 'À venir';
			case 'completed': return 'Terminé';
			case 'cancelled': return 'Annulé';
			case 'no_show': return 'Absent';
			default: return status;
		}
	}

	function paymentStatusLabel(status: string): string {
		switch (status) {
			case 'paid': return 'Payé';
			case 'pending': return 'En attente';
			case 'failed': return 'Échoué';
			case 'refunded': return 'Remboursé';
			default: return status;
		}
	}

	function isCancellable(): boolean {
		return booking.status === 'confirmed' && new Date(booking.slot_start_time) > now;
	}

	async function cancelBooking() {
		cancelling = true;
		cancelError = '';
		try {
			const res = await fetch(`/api/bookings/${booking.id}/cancel`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ reason: 'Annulé par le client' }),
			});
			if (res.ok) {
				await invalidateAll();
			} else {
				cancelError = 'Impossible d\'annuler. Veuillez réessayer.';
			}
		} catch {
			cancelError = 'Erreur réseau. Veuillez réessayer.';
		} finally {
			cancelling = false;
		}
	}

	function startEditNotes() {
		editingNotes = true;
		notesValue = notesOverride;
		notesError = '';
	}

	function cancelEditNotes() {
		editingNotes = false;
		notesValue = notesOverride;
		notesError = '';
	}

	async function saveNotes() {
		savingNotes = true;
		notesError = '';
		try {
			const res = await fetch(`/api/bookings/${booking.id}/notes`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ client_notes: notesValue }),
			});
			if (res.ok) {
				notesOverride = notesValue;
				editingNotes = false;
			} else {
				notesError = 'Impossible de sauvegarder. Veuillez réessayer.';
			}
		} catch {
			notesError = 'Erreur réseau. Veuillez réessayer.';
		} finally {
			savingNotes = false;
		}
	}
</script>

<svelte:head>
	<title>Réservation {booking.id.slice(0, 8)}… | Leviosa</title>
</svelte:head>

<div class="space-y-6">
	<!-- Back link -->
	<a
		href="/client/bookings"
		class="inline-flex items-center gap-1.5 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
	>
		<ArrowLeft size={16} />
		Retour aux réservations
	</a>

	<!-- Header -->
	<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
		<div>
			<h1 class="text-2xl font-bold text-foreground">Détail de la réservation</h1>
			<p class="text-sm text-muted-foreground mt-0.5">Réf. {booking.id}</p>
		</div>
		<div class="flex items-center gap-2">
			<span class="px-2.5 py-0.5 rounded-full text-xs font-medium {statusBadge(booking.status)}">
				{statusLabel(booking.status)}
			</span>
			<span class="px-2 py-0.5 rounded-full text-xs border {booking.payment_status === 'paid'
				? 'bg-green-50 text-green-600 border-green-200'
				: 'bg-gray-50 text-gray-600 border-gray-200'}">
				{paymentStatusLabel(booking.payment_status)}
			</span>
		</div>
	</div>

	<!-- Detail card -->
	<div class="bg-card rounded-lg border border-border-card p-5 sm:p-6 shadow-card space-y-5">
		<!-- Date/Time -->
		<div class="flex items-start gap-3">
			<Clock size={18} class="text-muted-foreground mt-0.5 flex-shrink-0" />
			<div>
				<p class="text-sm font-medium text-muted-foreground">Date et heure</p>
				<p class="text-foreground">{formatDate(booking.slot_start_time)}</p>
				<p class="text-sm text-muted-foreground">{formatTime(booking.slot_start_time)} – {formatTime(booking.slot_end_time)}</p>
			</div>
		</div>

		<!-- Price -->
		{#if booking.total_price_cents}
			<div class="flex items-start gap-3">
				<div class="w-[18px] mt-0.5 flex-shrink-0 text-center text-muted-foreground text-sm font-bold">€</div>
				<div>
					<p class="text-sm font-medium text-muted-foreground">Prix</p>
					<p class="text-foreground font-medium">{formatCents(booking.total_price_cents)}</p>
				</div>
			</div>
		{/if}

		<!-- Cancellation info -->
		{#if booking.status === 'cancelled' && booking.cancelled_at}
			<div class="flex items-start gap-3">
				<AlertCircle size={18} class="text-red-500 mt-0.5 flex-shrink-0" />
				<div>
					<p class="text-sm font-medium text-muted-foreground">Annulée le {new Date(booking.cancelled_at).toLocaleDateString('fr-FR')}</p>
					{#if booking.cancellation_reason}
						<p class="text-sm text-muted-foreground">{booking.cancellation_reason}</p>
					{/if}
				</div>
			</div>
		{/if}

		<!-- Completed info -->
		{#if booking.status === 'completed' && booking.completed_at}
			<div class="flex items-start gap-3">
				<div class="w-[18px] mt-0.5 flex-shrink-0 text-center text-green-600 text-sm">✓</div>
				<div>
					<p class="text-sm font-medium text-muted-foreground">Terminée le {new Date(booking.completed_at).toLocaleDateString('fr-FR')}</p>
				</div>
			</div>
		{/if}

		<!-- Partner notes (read-only) -->
		{#if booking.partner_notes}
			<div class="flex items-start gap-3">
				<FileText size={18} class="text-muted-foreground mt-0.5 flex-shrink-0" />
				<div>
					<p class="text-sm font-medium text-muted-foreground">Notes du praticien</p>
					<p class="text-foreground text-sm">{booking.partner_notes}</p>
				</div>
			</div>
		{/if}

		<!-- Client notes (editable) -->
		<div class="flex items-start gap-3">
			<FileText size={18} class="text-muted-foreground mt-0.5 flex-shrink-0" />
			<div class="flex-1">
				<p class="text-sm font-medium text-muted-foreground mb-2">Vos notes</p>
				{#if editingNotes}
					<div class="flex flex-col gap-2">
						<textarea
							class="w-full rounded-md border border-border-input bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring resize-none"
							rows="3"
							placeholder="Ajoutez une note pour le praticien…"
							bind:value={notesValue}
						></textarea>
						{#if notesError}
							<p class="text-xs text-red-600">{notesError}</p>
						{/if}
						<div class="flex gap-2">
							<button
								class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-foreground text-background rounded-md hover:bg-foreground/90 transition-colors disabled:opacity-50"
								onclick={saveNotes}
								disabled={savingNotes}
							>
								<Save size={13} />
								{savingNotes ? 'Enregistrement…' : 'Enregistrer'}
							</button>
							<button
								class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-muted text-muted-foreground rounded-md hover:bg-muted/80 transition-colors"
								onclick={cancelEditNotes}
							>
								<X size={13} />
								Annuler
							</button>
						</div>
					</div>
				{:else}
					<div class="flex items-start gap-2">
						{#if notesOverride}
							<p class="text-foreground text-sm flex-1">{notesOverride}</p>
						{:else}
							<p class="text-muted-foreground text-sm flex-1 italic">Aucune note</p>
						{/if}
						<button
							class="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground transition-colors flex-shrink-0"
							onclick={startEditNotes}
						>
							<Pencil size={13} />
							<span class="text-xs">{notesOverride ? 'Modifier' : 'Ajouter'}</span>
						</button>
					</div>
				{/if}
			</div>
		</div>

		<!-- Messaging shortcut (placeholder for Issue 008) -->
		{#if booking.status === 'confirmed' || booking.status === 'completed'}
			<div class="flex items-start gap-3">
				<MessageSquare size={18} class="text-muted-foreground mt-0.5 flex-shrink-0" />
				<div>
					<a
						href="/client/messages?partner={booking.partner_id}"
						class="text-sm font-medium text-foreground hover:underline"
					>
						Contacter le praticien
					</a>
				</div>
			</div>
		{/if}
	</div>

	<!-- Cancel action -->
	{#if isCancellable()}
		{#if cancelError}
			<div class="flex items-center gap-2 px-4 py-3 rounded-lg bg-red-50 border border-red-200 text-red-700 text-sm">
				<AlertCircle size={16} class="flex-shrink-0" />
				{cancelError}
			</div>
		{/if}
		<div class="pt-2">
			<button
				class="inline-flex items-center gap-1.5 px-4 py-2 text-sm font-medium text-red-600 bg-red-50 border border-red-200 rounded-lg hover:bg-red-100 transition-colors disabled:opacity-50"
				onclick={cancelBooking}
				disabled={cancelling}
			>
				<XCircle size={16} />
				{cancelling ? 'Annulation en cours…' : 'Annuler cette réservation'}
			</button>
		</div>
	{/if}
</div>
