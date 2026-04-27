<script lang="ts">
	import type { PageProps } from './$types';
	import { Clock, MapPin, CheckCircle, XCircle, FileText } from '@lucide/svelte';

	let { data }: PageProps = $props();

	type Tab = 'upcoming' | 'completed' | 'all';
	let activeTab = $state<Tab>('upcoming');

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
			case 'upcoming': return 'bg-blue-100 text-blue-700';
			case 'completed': return 'bg-green-100 text-green-700';
			case 'no_show': return 'bg-orange-100 text-orange-700';
			case 'cancelled': return 'bg-red-100 text-red-700';
			default: return 'bg-gray-100 text-gray-700';
		}
	}

	function statusLabel(status: string): string {
		switch (status) {
			case 'upcoming': return 'À venir';
			case 'completed': return 'Terminé';
			case 'no_show': return 'Absent';
			case 'cancelled': return 'Annulé';
			default: return status;
		}
	}

	const filtered = $derived(
		activeTab === 'all'
			? data.bookings
			: activeTab === 'upcoming'
				? data.bookings.filter((b) => b.status === 'upcoming')
				: data.bookings.filter((b) => b.status !== 'upcoming'),
	);

	const upcomingCount = $derived(data.bookings.filter((b) => b.status === 'upcoming').length);
	const completedCount = $derived(data.bookings.filter((b) => b.status === 'completed').length);
	const totalEarnings = $derived(
		data.bookings.filter((b) => b.status === 'completed').reduce((s, b) => s + b.amountInCents, 0),
	);
</script>

<svelte:head>
	<title>Réservations | Staff</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">Réservations</h1>
		<p class="text-muted-foreground">Suivi de vos séances et clients</p>
	</div>

	<!-- Summary -->
	<div class="grid grid-cols-3 gap-4 mb-8">
		<div class="bg-card rounded-lg border border-border p-4">
			<p class="text-sm text-muted-foreground mb-1">À venir</p>
			<p class="text-2xl font-bold text-blue-600">{upcomingCount}</p>
		</div>
		<div class="bg-card rounded-lg border border-border p-4">
			<p class="text-sm text-muted-foreground mb-1">Terminées</p>
			<p class="text-2xl font-bold text-green-600">{completedCount}</p>
		</div>
		<div class="bg-card rounded-lg border border-border p-4">
			<p class="text-sm text-muted-foreground mb-1">Revenus (mois)</p>
			<p class="text-2xl font-bold text-foreground">{formatCents(totalEarnings)}</p>
		</div>
	</div>

	<!-- Tabs -->
	<div class="flex gap-1 mb-6 bg-muted p-1 rounded-lg w-fit">
		{#each [['upcoming', 'À venir'], ['completed', 'Historique'], ['all', 'Tout']] as [val, label]}
			<button
				class="px-4 py-2 rounded-md text-sm font-medium transition-colors {activeTab === val
					? 'bg-card shadow-sm text-foreground'
					: 'text-muted-foreground hover:text-foreground'}"
				onclick={() => (activeTab = val as Tab)}
			>
				{label}
			</button>
		{/each}
	</div>

	<!-- Bookings List -->
	<div class="space-y-3">
		{#each filtered as booking (booking.id)}
			<div class="bg-card rounded-lg border border-border p-4 sm:p-5 hover:shadow-sm transition-shadow">
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

						{#if booking.notes}
							<div class="mt-2 flex items-start gap-1.5 text-sm text-muted-foreground">
								<FileText size={13} class="mt-0.5 flex-shrink-0" />
								<p class="italic">{booking.notes}</p>
							</div>
						{/if}
					</div>

					<!-- Amount + Actions -->
					<div class="flex flex-col items-end gap-2 flex-shrink-0">
						<p class="font-semibold text-foreground">{formatCents(booking.amountInCents)}</p>
						{#if booking.status === 'upcoming'}
							<div class="flex gap-2">
								<button
									class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
									title="Marquer comme terminé"
								>
									<CheckCircle size={13} />
									Terminé
								</button>
								<button
									class="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-orange-100 text-orange-700 rounded-md hover:bg-orange-200 transition-colors"
									title="Marquer absent"
								>
									<XCircle size={13} />
									Absent
								</button>
							</div>
						{/if}
					</div>
				</div>
			</div>
		{:else}
			<div class="text-center py-12 text-muted-foreground bg-card rounded-lg border border-border">
				Aucune réservation trouvée
			</div>
		{/each}
	</div>
</div>
