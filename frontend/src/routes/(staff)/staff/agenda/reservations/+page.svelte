<script lang="ts">
	import type { PageProps } from './$types';
	import type { Booking } from './+page.server';
	import { invalidateAll } from '$app/navigation';
	import { Clock, MapPin, CheckCircle, XCircle, FileText, Save, X, Pencil, AlertCircle } from '@lucide/svelte';

	let { data }: PageProps = $props();

	type Tab = 'upcoming' | 'completed' | 'no_show' | 'cancelled' | 'all';
	let activeTab = $state<Tab>((data.statusFilter as Tab) || 'upcoming');
	let editingNotesId: string | null = $state(null);
	let editingNotesValue = $state('');
	let savingNotes = $state(false);
	let notesError = $state('');
	let actionLoading: Record<string, boolean> = $state({});
	let actionError = $state('');
	let noteOverrides: Record<string, string> = $state({});

	function formatDateTime(iso: string): string {
		return new Date(iso).toLocaleDateString('fr-FR', {
			weekday: 'short',
			day: 'numeric',
			month: 'short',
			hour: '2-digit',
			minute: '2-digit',
		});
	}

	function formatTime(iso: string): string {
		return new Date(iso).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
	}

	function formatCents(cents: number): string {
		return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(
			cents / 100,
		);
	}

	function statusBadge(status: string): string {
		switch (status) {
			case 'confirmed':
				return 'bg-blue-100 text-blue-700';
			case 'completed':
				return 'bg-green-100 text-green-700';
			case 'no_show':
				return 'bg-orange-100 text-orange-700';
			case 'cancelled':
				return 'bg-red-100 text-red-700';
			default:
				return 'bg-gray-100 text-gray-700';
		}
	}

	function statusLabel(status: string): string {
		switch (status) {
			case 'confirmed':
				return 'À venir';
			case 'completed':
				return 'Terminé';
			case 'no_show':
				return 'Absent';
			case 'cancelled':
				return 'Annulé';
			default:
				return status;
		}
	}

	function paymentStatusBadge(status: string): string {
		switch (status) {
			case 'paid':
				return 'bg-green-50 text-green-600 border-green-200';
			case 'pending':
				return 'bg-yellow-50 text-yellow-600 border-yellow-200';
			case 'failed':
				return 'bg-red-50 text-red-600 border-red-200';
			case 'refunded':
				return 'bg-purple-50 text-purple-600 border-purple-200';
			default:
				return 'bg-gray-50 text-gray-600 border-gray-200';
		}
	}

	function paymentStatusLabel(status: string): string {
		switch (status) {
			case 'paid':
				return 'Payé';
			case 'pending':
				return 'En attente';
			case 'failed':
				return 'Échoué';
			case 'refunded':
				return 'Remboursé';
			default:
				return status;
		}
	}

	/** A confirmed booking is eligible for complete/no-show only if its end time has passed */
	function isActionable(booking: Booking): boolean {
		if (booking.status !== 'confirmed') return false;
		return new Date(booking.endTime) < new Date();
	}

	const filtered = $derived(
		activeTab === 'all'
			? data.bookings
			: activeTab === 'upcoming'
				? data.bookings.filter((b) => b.status === 'confirmed')
				: data.bookings.filter((b) => b.status === activeTab),
	);

	const upcomingCount = $derived(data.bookings.filter((b) => b.status === 'confirmed').length);
	const completedCount = $derived(data.bookings.filter((b) => b.status === 'completed').length);
	const completedEarnings = $derived(
		data.bookings
			.filter((b) => b.status === 'completed')
			.reduce((s, b) => s + b.amountInCents, 0),
	);

	const tabs: [Tab, string][] = [
		['upcoming', 'À venir'],
		['completed', 'Historique'],
		['no_show', 'Absents'],
		['cancelled', 'Annulés'],
		['all', 'Tout'],
	];

	function switchTab(tab: Tab) {
		activeTab = tab;
		const url = new URL(window.location.href);
		if (tab === 'all') {
			url.searchParams.delete('status');
		} else {
			url.searchParams.set('status', tab);
		}
		window.history.replaceState({}, '', url.toString());
	}

	async function completeBooking(bookingId: string) {
		actionLoading = { ...actionLoading, [bookingId + '_complete']: true };
		actionError = '';
		try {
			const res = await fetch(`/api/bookings/${bookingId}/complete`, { method: 'POST' });
			if (res.ok) {
				await invalidateAll();
			} else {
				actionError = 'Impossible de terminer la réservation. Veuillez réessayer.';
			}
		} catch {
			actionError = 'Erreur réseau. Veuillez réessayer.';
		} finally {
			actionLoading = { ...actionLoading, [bookingId + '_complete']: false };
		}
	}

	async function markNoShow(bookingId: string) {
		actionLoading = { ...actionLoading, [bookingId + '_noshow']: true };
		actionError = '';
		try {
			const res = await fetch(`/api/bookings/${bookingId}/no-show`, { method: 'POST' });
			if (res.ok) {
				await invalidateAll();
			} else {
				actionError = 'Impossible de marquer comme absent. Veuillez réessayer.';
			}
		} catch {
			actionError = 'Erreur réseau. Veuillez réessayer.';
		} finally {
			actionLoading = { ...actionLoading, [bookingId + '_noshow']: false };
		}
	}

	function startEditNotes(booking: Booking) {
		editingNotesId = booking.id;
		editingNotesValue = noteOverrides[booking.id] ?? booking.partnerNotes;
		notesError = '';
	}

	function cancelEditNotes() {
		editingNotesId = null;
		editingNotesValue = '';
		notesError = '';
	}

	async function saveNotes(bookingId: string) {
		savingNotes = true;
		notesError = '';
		try {
			const res = await fetch(`/api/bookings/${bookingId}/notes`, {
				method: 'PUT',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ partner_notes: editingNotesValue }),
			});
			if (res.ok) {
				noteOverrides = { ...noteOverrides, [bookingId]: editingNotesValue };
				editingNotesId = null;
				editingNotesValue = '';
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
	<title>Réservations | Staff</title>
</svelte:head>

<div class="p-6 lg:p-10">
	<div class="mb-10">
		<p class="text-[11px] font-semibold uppercase tracking-[0.2em] text-muted-foreground mb-3">Agenda</p>
		<h1 class="font-display text-4xl lg:text-5xl font-semibold tracking-tight text-foreground leading-[1.1]">Réservations</h1>
		<p class="text-muted-foreground mt-3 text-sm">Suivi de vos séances et clients</p>
		<div class="mt-4 h-px w-16 bg-foreground/20"></div>
	</div>

	<!-- Summary -->
	<div class="grid grid-cols-3 gap-4 mb-8">
		<div class="bg-background rounded-2xl border border-border p-5">
			<p class="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-2">À venir</p>
			<p class="text-2xl font-bold text-blue-600 dark:text-blue-400 tabular-nums">{upcomingCount}</p>
		</div>
		<div class="bg-background rounded-2xl border border-border p-5">
			<p class="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-2">Terminées</p>
			<p class="text-2xl font-bold text-green-600 dark:text-green-400 tabular-nums">{completedCount}</p>
		</div>
		<div class="bg-background rounded-2xl border border-border p-5">
			<p class="text-xs font-medium text-muted-foreground uppercase tracking-wide mb-2">Revenus</p>
			<p class="text-2xl font-bold text-foreground tabular-nums">{formatCents(completedEarnings)}</p>
		</div>
	</div>

	<!-- Tabs -->
	<div class="flex gap-1 mb-6 bg-muted p-1 rounded-lg w-fit flex-wrap">
		{#each tabs as [val, label]}
			<button
				class="px-4 py-2 rounded-md text-sm font-medium transition-colors {activeTab === val
					? 'bg-card shadow-mini text-foreground'
					: 'text-muted-foreground hover:text-foreground'}"
				onclick={() => switchTab(val)}
			>
				{label}
			</button>
		{/each}
	</div>

	<!-- Action error banner -->
	{#if actionError}
		<div class="flex items-center gap-2 mb-4 px-4 py-3 rounded-lg bg-red-50 border border-red-200 text-red-700 text-sm">
			<AlertCircle size={16} class="flex-shrink-0" />
			{actionError}
		</div>
	{/if}

	<!-- Bookings List -->
	<div class="space-y-3">
		{#each filtered as booking (booking.id)}
			<div class="bg-background rounded-xl border border-border p-4 sm:p-5 hover:shadow-mini transition-shadow">
				<div class="flex flex-col sm:flex-row sm:items-start gap-4">
					<!-- Client Avatar -->
					<div
						class="w-10 h-10 rounded-full bg-muted flex items-center justify-center flex-shrink-0"
					>
						<span class="text-sm font-semibold text-foreground">{booking.clientInitials}</span>
					</div>

					<!-- Main Info -->
					<div class="flex-1 min-w-0">
						<div class="flex flex-wrap items-center gap-2 mb-1">
							<p class="font-semibold text-foreground">{booking.clientName}</p>
							<span class="px-2 py-0.5 rounded-full text-xs font-medium {statusBadge(booking.status)}">
								{statusLabel(booking.status)}
							</span>
							<span class="px-2 py-0.5 rounded-full text-xs font-medium border {paymentStatusBadge(booking.paymentStatus)}">
								{paymentStatusLabel(booking.paymentStatus)}
							</span>
						</div>
						<p class="text-sm text-muted-foreground mb-2">{booking.productName}</p>

						<div class="flex flex-wrap gap-4 text-sm text-muted-foreground">
							<span class="flex items-center gap-1.5">
								<Clock size={13} />
								{formatDateTime(booking.startTime)} – {formatTime(booking.endTime)}
							</span>
							<span class="flex items-center gap-1.5">
								<MapPin size={13} />
								{booking.roomName}
							</span>
						</div>

						{#if booking.clientNotes}
							<div class="mt-2 flex items-start gap-1.5 text-sm text-muted-foreground">
								<FileText size={13} class="mt-0.5 flex-shrink-0" />
								<p class="italic">{booking.clientNotes}</p>
							</div>
						{/if}

						<!-- Partner Notes (inline editor) -->
						<div class="mt-2">
							{#if editingNotesId === booking.id}
								<div class="flex flex-col gap-2">
									<textarea
										class="w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring resize-none"
										rows="2"
										placeholder="Notes internes (visibles uniquement par vous)…"
										bind:value={editingNotesValue}
									></textarea>
									{#if notesError}
										<p class="text-xs text-red-600">{notesError}</p>
									{/if}
									<div class="flex gap-2">
										<button
											class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors disabled:opacity-50"
											onclick={() => saveNotes(booking.id)}
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
								{@const displayNote = noteOverrides[booking.id] ?? booking.partnerNotes}
								<button
									class="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors"
									onclick={() => startEditNotes(booking)}
								>
									{#if displayNote}
										<Pencil size={13} />
										<span class="italic">{displayNote}</span>
									{:else}
										<Pencil size={13} />
										<span>Ajouter une note</span>
									{/if}
								</button>
							{/if}
						</div>
					</div>

					<!-- Amount + Actions -->
					<div class="flex flex-col items-end gap-2 flex-shrink-0">
						<p class="font-semibold text-foreground">{formatCents(booking.amountInCents)}</p>
						{#if isActionable(booking)}
							<div class="flex gap-2">
								<button
									class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors disabled:opacity-50"
									title="Marquer comme terminé"
									onclick={() => completeBooking(booking.id)}
									disabled={actionLoading[booking.id + '_complete']}
								>
									<CheckCircle size={13} />
									{actionLoading[booking.id + '_complete'] ? '…' : 'Terminer'}
								</button>
								<button
									class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-orange-100 text-orange-700 rounded-md hover:bg-orange-200 transition-colors disabled:opacity-50"
									title="Marquer absent"
									onclick={() => markNoShow(booking.id)}
									disabled={actionLoading[booking.id + '_noshow']}
								>
									<XCircle size={13} />
									{actionLoading[booking.id + '_noshow'] ? '…' : 'Absent'}
								</button>
							</div>
						{/if}
					</div>
				</div>
			</div>
		{:else}
			<div class="text-center py-12 text-muted-foreground bg-background rounded-xl border border-border">
				Aucune réservation trouvée
			</div>
		{/each}
	</div>
</div>
