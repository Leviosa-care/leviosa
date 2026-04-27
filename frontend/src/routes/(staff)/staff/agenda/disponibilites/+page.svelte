<script lang="ts">
	import type { PageProps } from './$types';
	import { CalendarClock, Clock, MapPin, Plus, X, ChevronLeft, ChevronRight } from '@lucide/svelte';

	let { data }: PageProps = $props();

	let filterStatus = $state<'all' | 'available' | 'booked' | 'cancelled'>('all');

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
	}

	function formatTime(iso: string): string {
		return new Date(iso).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
	}

	function statusBadge(status: string): string {
		switch (status) {
			case 'available': return 'bg-green-100 text-green-700';
			case 'booked': return 'bg-blue-100 text-blue-700';
			case 'cancelled': return 'bg-red-100 text-red-700';
			default: return 'bg-gray-100 text-gray-700';
		}
	}

	function statusLabel(status: string): string {
		switch (status) {
			case 'available': return 'Disponible';
			case 'booked': return 'Réservé';
			case 'cancelled': return 'Annulé';
			default: return status;
		}
	}

	const filteredDays = $derived(
		data.availabilities
			.map((day) => ({
				...day,
				slots: day.slots.filter((s) => filterStatus === 'all' || s.status === filterStatus),
			}))
			.filter((day) => day.slots.length > 0),
	);

	const totalAvailable = $derived(
		data.availabilities.flatMap((d) => d.slots).filter((s) => s.status === 'available').length,
	);
	const totalBooked = $derived(
		data.availabilities.flatMap((d) => d.slots).filter((s) => s.status === 'booked').length,
	);
</script>

<svelte:head>
	<title>Disponibilités | Staff</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
		<div>
			<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">Disponibilités</h1>
			<p class="text-muted-foreground">Gérez vos créneaux pour les 7 prochains jours</p>
		</div>
		<button
			class="inline-flex items-center gap-2 px-4 py-2 bg-foreground text-background rounded-lg text-sm font-medium hover:opacity-90 transition-opacity"
		>
			<Plus size={16} />
			Nouveau créneau
		</button>
	</div>

	<!-- Summary Cards -->
	<div class="grid grid-cols-2 sm:grid-cols-3 gap-4 mb-8">
		<div class="bg-card rounded-lg border border-border p-4">
			<p class="text-sm text-muted-foreground mb-1">Créneaux libres</p>
			<p class="text-2xl font-bold text-green-600">{totalAvailable}</p>
		</div>
		<div class="bg-card rounded-lg border border-border p-4">
			<p class="text-sm text-muted-foreground mb-1">Réservés</p>
			<p class="text-2xl font-bold text-blue-600">{totalBooked}</p>
		</div>
		<div class="bg-card rounded-lg border border-border p-4 col-span-2 sm:col-span-1">
			<p class="text-sm text-muted-foreground mb-1">Taux d'occupation</p>
			<p class="text-2xl font-bold text-foreground">
				{totalBooked + totalAvailable > 0
					? Math.round((totalBooked / (totalBooked + totalAvailable)) * 100)
					: 0}%
			</p>
		</div>
	</div>

	<!-- Week Navigation -->
	<div class="flex items-center justify-between mb-6 bg-card rounded-lg border border-border p-4">
		<button class="p-2 hover:bg-muted rounded-md transition-colors" disabled>
			<ChevronLeft size={20} class="text-muted-foreground" />
		</button>
		<div class="flex items-center gap-2">
			<CalendarClock size={18} class="text-muted-foreground" />
			<span class="font-semibold text-foreground">Cette semaine</span>
		</div>
		<button class="p-2 hover:bg-muted rounded-md transition-colors" disabled>
			<ChevronRight size={20} class="text-muted-foreground" />
		</button>
	</div>

	<!-- Status Filter -->
	<div class="flex gap-2 mb-6 flex-wrap">
		{#each [['all', 'Tous'], ['available', 'Libres'], ['booked', 'Réservés'], ['cancelled', 'Annulés']] as [val, label]}
			<button
				class="px-4 py-2 rounded-md text-sm font-medium transition-colors {filterStatus === val
					? 'bg-foreground text-background'
					: 'bg-card text-foreground border border-border hover:bg-muted'}"
				onclick={() => (filterStatus = val as typeof filterStatus)}
			>
				{label}
			</button>
		{/each}
	</div>

	<!-- Days and Slots -->
	<div class="space-y-6">
		{#each filteredDays as day (day.date)}
			<div class="bg-card rounded-lg border border-border overflow-hidden">
				<div class="bg-muted/50 px-5 py-3 border-b border-border">
					<h2 class="font-semibold text-foreground capitalize">
						{day.dayName} {formatDate(day.date)}
					</h2>
				</div>
				<div class="divide-y divide-border">
					{#each day.slots as slot (slot.id)}
						<div class="p-4 sm:p-5 hover:bg-muted/20 transition-colors">
							<div class="flex flex-col sm:flex-row sm:items-center gap-3">
								<div class="flex items-center gap-2 text-muted-foreground min-w-fit">
									<Clock size={15} />
									<span class="font-medium text-sm">
										{formatTime(slot.startTime)} – {formatTime(slot.endTime)}
									</span>
								</div>
								<div class="flex-1">
									<div class="flex items-center gap-1.5 text-sm text-muted-foreground mb-0.5">
										<MapPin size={13} />
										<span>{slot.roomName}</span>
									</div>
									{#if slot.clientName}
										<p class="text-sm font-medium text-foreground">
											{slot.clientName}
											{#if slot.productName}
												<span class="font-normal text-muted-foreground">— {slot.productName}</span>
											{/if}
										</p>
									{/if}
								</div>
								<div class="flex items-center gap-2">
									<span class="px-2.5 py-1 rounded-full text-xs font-medium {statusBadge(slot.status)}">
										{statusLabel(slot.status)}
									</span>
									{#if slot.status === 'available'}
										<button
											class="p-1.5 rounded-md text-muted-foreground hover:text-red-600 hover:bg-red-50 transition-colors"
											title="Annuler ce créneau"
										>
											<X size={14} />
										</button>
									{/if}
								</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		{:else}
			<div class="text-center py-12 text-muted-foreground">
				Aucun créneau pour cette période
			</div>
		{/each}
	</div>
</div>
