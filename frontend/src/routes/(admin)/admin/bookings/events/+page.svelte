<script lang="ts">
	import type { PageProps } from './$types';
	import { Calendar, Clock, MapPin, Users, Euro } from '@lucide/svelte';

	let { data }: PageProps = $props();

	function formatDateTime(isoString: string): string {
		const date = new Date(isoString);
		return date.toLocaleDateString('fr-FR', {
			day: '2-digit',
			month: 'long',
			year: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatCents(cents: number): string {
		return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(cents / 100);
	}

	function getCapacityColor(registered: number, capacity: number): string {
		const percentage = (registered / capacity) * 100;
		if (percentage >= 90) return 'bg-green-500';
		if (percentage >= 70) return 'bg-blue-500';
		if (percentage >= 50) return 'bg-yellow-500';
		return 'bg-gray-400';
	}

	function getStatusBadge(status: string) {
		switch (status) {
			case 'upcoming':
				return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300';
			case 'ongoing':
				return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300';
			case 'completed':
				return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
			case 'cancelled':
				return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300';
			default:
				return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
		}
	}

	function getStatusLabel(status: string): string {
		switch (status) {
			case 'upcoming': return 'À venir';
			case 'ongoing': return 'En cours';
			case 'completed': return 'Terminé';
			case 'cancelled': return 'Annulé';
			default: return status;
		}
	}
</script>

<svelte:head>
	<title>Événements | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">
			Événements
		</h1>
		<p class="text-muted-foreground">
			Gérez les événements collectifs et ateliers
		</p>
	</div>

	<!-- Events Grid -->
	<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
		{#each data.events as event}
			<div class="bg-card rounded-lg border border-border p-6 hover:shadow-lg transition-shadow">
				<!-- Header -->
				<div class="flex items-start justify-between mb-4">
					<div class="flex-1">
						<h3 class="text-lg font-semibold text-foreground mb-1">{event.name}</h3>
						<p class="text-sm text-muted-foreground">{event.description}</p>
					</div>
					<span class="px-3 py-1 rounded-full text-xs font-medium {getStatusBadge(event.status)}">
						{getStatusLabel(event.status)}
					</span>
				</div>

				<!-- Date & Time -->
				<div class="flex items-center gap-2 text-sm text-muted-foreground mb-3">
					<Calendar size={16} />
					<span>{formatDateTime(event.date)}</span>
				</div>

				<!-- Duration -->
				<div class="flex items-center gap-2 text-sm text-muted-foreground mb-3">
					<Clock size={16} />
					<span>{Math.floor(event.duration / 60)}h {event.duration % 60 > 0 ? event.duration % 60 + 'min' : ''}</span>
				</div>

				<!-- Location -->
				<div class="flex items-center gap-2 text-sm text-muted-foreground mb-4">
					<MapPin size={16} />
					<span>{event.location}</span>
				</div>

				<!-- Capacity Progress -->
				<div class="mb-4">
					<div class="flex items-center justify-between text-sm mb-2">
						<div class="flex items-center gap-2 text-muted-foreground">
							<Users size={14} />
							<span>Inscriptions</span>
						</div>
						<span class="font-medium text-foreground">{event.registered} / {event.capacity}</span>
					</div>
					<div class="w-full h-2 bg-muted rounded-full overflow-hidden">
						<div
							class="h-full {getCapacityColor(event.registered, event.capacity)} rounded-full transition-all"
							style="width: {Math.min((event.registered / event.capacity) * 100, 100)}%"
						></div>
					</div>
				</div>

				<!-- Price -->
				<div class="flex items-center justify-between pt-4 border-t border-border">
					<div class="flex items-center gap-2 text-muted-foreground">
						<Euro size={14} />
						<span class="text-sm">Prix</span>
					</div>
					<span class="text-lg font-semibold text-foreground">
						{event.priceInCents === 0 ? 'Gratuit' : formatCents(event.priceInCents)}
					</span>
				</div>
			</div>
		{:else}
			<div class="col-span-1 md:col-span-2 text-center py-12 text-muted-foreground">
				Aucun événement trouvé
			</div>
		{/each}
	</div>
</div>
