<script lang="ts">
	import type { PageProps } from './$types';
	import { invalidateAll } from '$app/navigation';
	import {
		Clock,
		AlertCircle,
		ArrowRight,
		XCircle,
	} from '@lucide/svelte';

	let { data }: PageProps = $props();

	type Tab = 'all' | 'upcoming' | 'completed' | 'cancelled';
	let activeTab = $state<Tab>('all');
	let cancellingId: string | null = $state(null);
	let cancelError = $state('');

	function formatDate(iso: string): string {
		return new Date(iso).toLocaleDateString('fr-FR', {
			weekday: 'short',
			day: 'numeric',
			month: 'short',
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
			case 'confirmed': return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400';
			case 'completed': return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400';
			case 'cancelled': return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400';
			case 'no_show': return 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400';
			default: return 'bg-muted text-muted-foreground';
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

	const now = new Date();

	const filtered = $derived(
		activeTab === 'all'
			? data.bookings
			: activeTab === 'upcoming'
				? data.bookings.filter((b: any) => b.status === 'confirmed' && new Date(b.slot_start_time) > now)
				: data.bookings.filter((b: any) => b.status === activeTab)
	);

	const tabs: [Tab, string][] = [
		['all', 'Toutes'],
		['upcoming', 'À venir'],
		['completed', 'Passées'],
		['cancelled', 'Annulées'],
	];

	function isCancellable(booking: any): boolean {
		if (booking.status !== 'confirmed') return false;
		return new Date(booking.slot_start_time) > now;
	}

	async function cancelBooking(bookingId: string) {
		cancellingId = bookingId;
		cancelError = '';
		try {
			const res = await fetch(`/api/bookings/${bookingId}/cancel`, {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ reason: 'Annulé par le client' }),
			});
			if (res.ok) {
				await invalidateAll();
			} else {
				cancelError = 'Impossible d\'annuler la réservation. Veuillez réessayer.';
			}
		} catch {
			cancelError = 'Erreur réseau. Veuillez réessayer.';
		} finally {
			cancellingId = null;
		}
	}
</script>

<svelte:head>
	<title>Mes réservations | Leviosa</title>
</svelte:head>

<div class="space-y-6">
	<div>
		<h1 class="text-3xl lg:text-4xl font-bold text-foreground mb-1">Mes réservations</h1>
		<p class="text-muted-foreground">Historique de vos séances</p>
	</div>

	<!-- Tabs -->
	<div class="flex gap-1 bg-muted p-1 rounded-lg w-fit flex-wrap">
		{#each tabs as [val, label]}
			<button
				class="px-4 py-2 rounded-md text-sm font-medium transition-colors {activeTab === val
					? 'bg-background shadow-mini text-foreground'
					: 'text-muted-foreground hover:text-foreground'}"
				onclick={() => { activeTab = val; cancelError = ''; }}
			>
				{label}
			</button>
		{/each}
	</div>

	<!-- Error -->
	{#if cancelError}
		<div class="flex items-center gap-2 px-4 py-3 rounded-lg bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 text-sm">
			<AlertCircle size={16} class="flex-shrink-0" />
			{cancelError}
		</div>
	{/if}

	<!-- Bookings list -->
	<div class="space-y-3">
		{#each filtered as booking (booking.id)}
			<div class="bg-card rounded-lg border border-border p-4 sm:p-5 hover:bg-muted/30 transition-colors">
				<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
					<!-- Info -->
					<div class="flex-1 min-w-0">
						<div class="flex flex-wrap items-center gap-2 mb-1">
							<span class="px-2.5 py-0.5 rounded-full text-xs font-medium {statusBadge(booking.status)}">
								{statusLabel(booking.status)}
							</span>
							<span class="px-2 py-0.5 rounded-full text-xs border {booking.payment_status === 'paid'
								? 'bg-green-50 text-green-600 border-green-200 dark:bg-green-900/20 dark:text-green-400 dark:border-green-800'
								: 'bg-muted text-muted-foreground border-border'}">
								{paymentStatusLabel(booking.payment_status)}
							</span>
						</div>
						<p class="text-foreground font-medium">{formatDate(booking.slot_start_time)}</p>
						<div class="flex items-center gap-1.5 text-sm text-muted-foreground mt-0.5">
							<Clock size={13} />
							{formatTime(booking.slot_start_time)} – {formatTime(booking.slot_end_time)}
						</div>
					</div>

					<!-- Price + actions -->
					<div class="flex items-center gap-3 flex-shrink-0">
						{#if booking.total_price_cents}
							<span class="text-sm font-semibold text-foreground">{formatCents(booking.total_price_cents)}</span>
						{/if}
						<a
							href="/client/bookings/{booking.id}"
							class="inline-flex items-center gap-1 text-sm font-medium text-muted-foreground hover:text-foreground transition-colors"
						>
							Détail <ArrowRight size={14} />
						</a>
						{#if isCancellable(booking)}
							<button
								class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium text-red-600 dark:text-red-400 bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 rounded-md hover:bg-red-100 dark:hover:bg-red-900 transition-colors disabled:opacity-50"
								onclick={() => cancelBooking(booking.id)}
								disabled={cancellingId === booking.id}
							>
								<XCircle size={13} />
								{cancellingId === booking.id ? '…' : 'Annuler'}
							</button>
						{/if}
					</div>
				</div>
			</div>
		{:else}
			<div class="bg-card rounded-lg border border-border p-12 text-center text-muted-foreground">
				Aucune réservation trouvée
			</div>
		{/each}
	</div>
</div>
