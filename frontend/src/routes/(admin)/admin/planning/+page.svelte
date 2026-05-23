<script lang="ts">
	import type { PageProps } from './$types';
	import { Calendar, Clock, User, MapPin, ChevronLeft, ChevronRight } from '@lucide/svelte';
	import { goto } from '$app/navigation';

	let { data }: PageProps = $props();

	let filterStatus = $state<'all' | 'confirmed' | 'completed' | 'cancelled' | 'no_show'>('all');

	function getWeekStartDate(): Date {
		if (data.weekStart) {
			return new Date(data.weekStart + 'T00:00:00');
		}
		return new Date();
	}

	function navigateWeek(deltaDays: number) {
		const current = getWeekStartDate();
		current.setDate(current.getDate() + deltaDays);
		const monday = current.toISOString().split('T')[0];
		goto(`/admin/planning?week=${monday}`);
	}

	function formatDate(dateStr: string): string {
		const date = new Date(dateStr);
		return date.toLocaleDateString('fr-FR', { day: 'numeric', month: 'short' });
	}

	function formatTime(isoString: string): string {
		const date = new Date(isoString);
		return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
	}

	function getStatusBadge(status: string) {
		switch (status) {
			case 'confirmed':
				return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300';
			case 'completed':
				return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300';
			case 'cancelled':
				return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300';
			case 'no_show':
				return 'bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300';
			default:
				return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
		}
	}

	function getStatusLabel(status: string): string {
		switch (status) {
			case 'confirmed': return 'Confirmé';
			case 'completed': return 'Terminé';
			case 'cancelled': return 'Annulé';
			case 'no_show': return 'Absent';
			default: return status;
		}
	}

	const filteredEvents = $derived(
		data.weekEvents.map(day => ({
			...day,
			events: day.events.filter(e =>
				filterStatus === 'all' || e.status === filterStatus
			)
		})).filter(day => day.events.length > 0)
	);
</script>

<svelte:head>
	<title>Planning | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">
			Planning
		</h1>
		<p class="text-muted-foreground">
			Vue d'ensemble des rendez-vous de la semaine
		</p>
	</div>

	<!-- Week Navigation -->
	<div class="flex items-center justify-between mb-6 bg-card rounded-lg border border-border p-4">
		<button class="p-2 hover:bg-muted rounded-md transition-colors" onclick={() => navigateWeek(-7)}>
			<ChevronLeft size={20} class="text-muted-foreground" />
		</button>
		<div class="flex items-center gap-2">
			<Calendar size={18} class="text-muted-foreground" />
			<span class="font-semibold">Semaine du {data.weekEvents.length > 0 ? formatDate(data.weekEvents[0].date) : formatDate(data.weekStart ?? new Date().toISOString())}</span>
		</div>
		<button class="p-2 hover:bg-muted rounded-md transition-colors" onclick={() => navigateWeek(7)}>
			<ChevronRight size={20} class="text-muted-foreground" />
		</button>
	</div>

	<!-- Status Filter -->
	{#if data.weekEvents.some(day => day.events.length > 0)}
		<div class="flex gap-2 mb-6">
			<button
				class="px-4 py-2 rounded-md text-sm font-medium transition-colors {filterStatus === 'all'
					? 'bg-primary text-primary-foreground'
					: 'bg-card text-foreground hover:bg-muted'}"
				onclick={() => filterStatus = 'all'}
			>
				Tous
			</button>
			<button
				class="px-4 py-2 rounded-md text-sm font-medium transition-colors {filterStatus === 'confirmed'
					? 'bg-primary text-primary-foreground'
					: 'bg-card text-foreground hover:bg-muted'}"
				onclick={() => filterStatus = 'confirmed'}
			>
				Confirmés
			</button>
			<button
				class="px-4 py-2 rounded-md text-sm font-medium transition-colors {filterStatus === 'completed'
					? 'bg-primary text-primary-foreground'
					: 'bg-card text-foreground hover:bg-muted'}"
				onclick={() => filterStatus = 'completed'}
			>
				Terminés
			</button>
			<button
				class="px-4 py-2 rounded-md text-sm font-medium transition-colors {filterStatus === 'cancelled'
					? 'bg-primary text-primary-foreground'
					: 'bg-card text-foreground hover:bg-muted'}"
				onclick={() => filterStatus = 'cancelled'}
			>
				Annulés
			</button>
			<button
				class="px-4 py-2 rounded-md text-sm font-medium transition-colors {filterStatus === 'no_show'
					? 'bg-primary text-primary-foreground'
					: 'bg-card text-foreground hover:bg-muted'}"
				onclick={() => filterStatus = 'no_show'}
			>
				Absents
			</button>
		</div>
	{/if}

	<!-- Events by Day -->
	{#if data.weekEvents.length === 0 || data.weekEvents.every(day => day.events.length === 0)}
		<div class="flex flex-col items-center justify-center py-16 text-center">
			<Calendar size={48} class="text-muted-foreground mb-4" />
			<h3 class="text-lg font-semibold text-foreground mb-2">Aucun rendez-vous cette semaine</h3>
			<p class="text-muted-foreground max-w-md">
				Aucun rendez-vous n'est prévu pour cette semaine.
				{#if filterStatus !== 'all'}
					Veuillez réinitialiser les filtres pour voir tous les événements.
				{/if}
			</p>
		</div>
	{:else}
		<div class="space-y-6">
			{#each filteredEvents as day (day.date)}
				<div class="bg-card rounded-lg border border-border overflow-hidden">
					<div class="bg-muted/50 px-5 py-3 border-b border-border">
						<h2 class="font-semibold text-foreground">{day.dayName} {formatDate(day.date)}</h2>
					</div>
					<div class="divide-y divide-border">
						{#each day.events as event}
							<div class="p-5 hover:bg-muted/30 transition-colors">
								<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
									<div class="flex items-start gap-3 flex-1">
										<div class="flex items-center gap-2 text-muted-foreground min-w-fit">
											<Clock size={16} />
											<span class="font-medium">
												{formatTime(event.startTime)} - {formatTime(event.endTime)}
											</span>
										</div>
									</div>
									<div class="flex-1">
										<div class="flex items-center gap-2 mb-1">
											<User size={14} class="text-muted-foreground" />
											<span class="font-medium text-foreground">{event.clientName}</span>
										</div>
										<div class="text-sm text-muted-foreground">
											{event.productName} avec {event.therapistName}
										</div>
									</div>
									<div class="flex items-center gap-3">
										<div class="flex items-center gap-1 text-sm text-muted-foreground">
											<MapPin size={14} />
											<span>{event.roomName}</span>
										</div>
										<span class="px-3 py-1 rounded-full text-xs font-medium {getStatusBadge(event.status)}">
											{getStatusLabel(event.status)}
										</span>
									</div>
								</div>
							</div>
						{/each}
					</div>
				</div>
			{:else}
				<div class="text-center py-12 text-muted-foreground">
					Aucun rendez-vous ne correspond aux filtres sélectionnés
				</div>
			{/each}
		</div>
	{/if}
</div>
