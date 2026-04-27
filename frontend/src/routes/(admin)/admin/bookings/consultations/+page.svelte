<script lang="ts">
	import type { PageProps } from './$types';
	import { Search, Filter, FileText, Calendar, Clock, User } from '@lucide/svelte';

	let { data }: PageProps = $props();

	let searchQuery = $state('');
	let statusFilter = $state<'all' | 'completed' | 'pending' | 'cancelled'>('all');
	let therapistFilter = $state<string>('all');

	const therapists = $derived([...new Set(data.consultations.map(c => c.therapistName))]);

	function formatDateTime(isoString: string): string {
		const date = new Date(isoString);
		return date.toLocaleDateString('fr-FR', {
			day: '2-digit',
			month: '2-digit',
			year: '2-digit',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function getStatusBadge(status: string) {
		switch (status) {
			case 'completed':
				return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300';
			case 'pending':
				return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300';
			case 'cancelled':
				return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300';
			default:
				return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
		}
	}

	function getStatusLabel(status: string): string {
		switch (status) {
			case 'completed': return 'Terminé';
			case 'pending': return 'À venir';
			case 'cancelled': return 'Annulé';
			default: return status;
		}
	}

	const filteredConsultations = $derived(
		data.consultations.filter(c => {
			const matchesSearch =
				c.clientName.toLowerCase().includes(searchQuery.toLowerCase()) ||
				c.productName.toLowerCase().includes(searchQuery.toLowerCase()) ||
				c.therapistName.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesStatus = statusFilter === 'all' || c.status === statusFilter;
			const matchesTherapist = therapistFilter === 'all' || c.therapistName === therapistFilter;
			return matchesSearch && matchesStatus && matchesTherapist;
		})
	);
</script>

<svelte:head>
	<title>Consultations | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">
			Consultations
		</h1>
		<p class="text-muted-foreground">
			Gérez les consultations individuelles
		</p>
	</div>

	<!-- Search and Filters -->
	<div class="bg-card rounded-lg border border-border p-4 mb-6">
		<div class="flex flex-col md:flex-row gap-4">
			<!-- Search -->
			<div class="flex-1 relative">
				<Search size={18} class="absolute left-3 top-1/2 -translate-y-1/2 text-muted-foreground" />
				<input
					type="text"
					placeholder="Rechercher client, produit, thérapeute..."
					bind:value={searchQuery}
					class="w-full pl-10 pr-4 py-2 bg-background border border-border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary"
				/>
			</div>

			<!-- Status Filter -->
			<div class="flex items-center gap-2">
				<Filter size={18} class="text-muted-foreground" />
				<select
					bind:value={statusFilter}
					class="px-4 py-2 bg-background border border-border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary"
				>
					<option value="all">Tous les statuts</option>
					<option value="completed">Terminé</option>
					<option value="pending">À venir</option>
					<option value="cancelled">Annulé</option>
				</select>
			</div>

			<!-- Therapist Filter -->
			<div class="flex items-center gap-2">
				<User size={18} class="text-muted-foreground" />
				<select
					bind:value={therapistFilter}
					class="px-4 py-2 bg-background border border-border rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-primary"
				>
					<option value="all">Tous les thérapeutes</option>
					{#each therapists as therapist}
						<option value={therapist}>{therapist}</option>
					{/each}
				</select>
			</div>
		</div>
	</div>

	<!-- Consultations Table -->
	<div class="bg-card rounded-lg border border-border overflow-hidden">
		<div class="overflow-x-auto">
			<table class="w-full">
				<thead>
					<tr class="bg-muted/50 border-b border-border">
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">
							<div class="flex items-center gap-2">
								<Calendar size={14} />
								<span>Date</span>
							</div>
						</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Client</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Thérapeute</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Produit</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">
							<div class="flex items-center gap-2">
								<Clock size={14} />
								<span>Durée</span>
							</div>
						</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Statut</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Notes</th>
					</tr>
				</thead>
				<tbody>
					{#each filteredConsultations as consultation}
						<tr class="border-b border-border hover:bg-muted/30 transition-colors">
							<td class="py-4 px-6 text-sm text-foreground">
								{formatDateTime(consultation.date)}
							</td>
							<td class="py-4 px-6">
								<span class="text-sm font-medium text-foreground">{consultation.clientName}</span>
							</td>
							<td class="py-4 px-6 text-sm text-foreground">{consultation.therapistName}</td>
							<td class="py-4 px-6 text-sm text-foreground">{consultation.productName}</td>
							<td class="py-4 px-6 text-sm text-foreground">{consultation.duration} min</td>
							<td class="py-4 px-6">
								<span class="px-3 py-1 rounded-full text-xs font-medium {getStatusBadge(consultation.status)}">
									{getStatusLabel(consultation.status)}
								</span>
							</td>
							<td class="py-4 px-6">
								{#if consultation.hasNotes}
									<span class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300">
										<FileText size={12} />
										<span>Notes</span>
									</span>
								{:else}
									<span class="text-sm text-muted-foreground">—</span>
								{/if}
							</td>
						</tr>
					{:else}
						<tr>
							<td colspan="7" class="py-12 text-center text-muted-foreground">
								Aucune consultation trouvée
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
