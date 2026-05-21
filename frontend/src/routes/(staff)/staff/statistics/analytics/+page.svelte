<script lang="ts">
	import type { PageProps } from './$types';
	import { Activity, Clock, Target } from '@lucide/svelte';

	let { data }: PageProps = $props();

	const dailyMetrics = $derived(
		data.metrics?.roomMetrics.flatMap((rm) => rm.dailyMetrics) ?? []
	);

	const totalBookedMinutes = $derived(
		dailyMetrics.reduce((sum, d) => sum + d.totalMinutesBooked, 0)
	);
	const totalOpenMinutes = $derived(
		dailyMetrics.reduce((sum, d) => sum + d.totalMinutesOpen, 0)
	);
	const avgUtilization = $derived(
		totalOpenMinutes > 0 ? (totalBookedMinutes / totalOpenMinutes) * 100 : 0
	);
	const totalSessions = $derived(dailyMetrics.length);
	const avgEfficiency = $derived(
		dailyMetrics.length > 0
			? dailyMetrics.reduce((sum, d) => sum + d.efficiencyScore, 0) / dailyMetrics.length
			: 0
	);

	function formatDate(iso: string): string {
		return new Date(iso).toLocaleDateString('fr-FR', {
			day: 'numeric',
			month: 'short',
		});
	}

	const sortedMetrics = $derived(
		[...dailyMetrics].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime())
	);

	const last7Days = $derived(sortedMetrics.slice(-7));

	const maxMinutes = $derived(
		Math.max(...last7Days.map((d) => d.totalMinutesBooked), 1)
	);
</script>

<svelte:head>
	<title>Statistiques | Staff</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">Statistiques</h1>
		<p class="text-muted-foreground">Vos performances du mois en cours</p>
	</div>

	{#if data.error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-8">
			{data.error}
		</div>
	{:else if data.metrics}
		<!-- KPI Cards -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-green-100 rounded-lg">
						<Target size={22} class="text-green-600" />
					</div>
					<span class="text-xs text-muted-foreground">Taux</span>
				</div>
				<p class="text-3xl font-bold text-foreground mb-1">
					{avgUtilization.toFixed(1)}%
				</p>
				<p class="text-sm text-muted-foreground">Taux d'occupation moyen</p>
			</div>

			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-blue-100 rounded-lg">
						<Activity size={22} class="text-blue-600" />
					</div>
					<span class="text-xs text-muted-foreground">Total</span>
				</div>
				<p class="text-3xl font-bold text-foreground mb-1">{totalSessions}</p>
				<p class="text-sm text-muted-foreground">Jours analysés</p>
			</div>

			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-yellow-100 rounded-lg">
						<Activity size={22} class="text-yellow-600" />
					</div>
					<span class="text-xs text-muted-foreground">Score</span>
				</div>
				<p class="text-3xl font-bold text-foreground mb-1">{avgEfficiency.toFixed(1)}</p>
				<p class="text-sm text-muted-foreground">Score d'efficacité moyen</p>
			</div>

			<div class="bg-card rounded-lg border border-border p-6">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-orange-100 rounded-lg">
						<Clock size={22} class="text-orange-600" />
					</div>
					<span class="text-xs text-muted-foreground">Total</span>
				</div>
				<p class="text-3xl font-bold text-foreground mb-1">
					{Math.round(totalBookedMinutes / 60)}h
				</p>
				<p class="text-sm text-muted-foreground">Temps total réservé</p>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
			<!-- Daily Volume Chart (Last 7 days) -->
			<div class="bg-card rounded-lg border border-border p-6">
				<h2 class="text-lg font-semibold text-foreground mb-6">Volume quotidien (7 jours)</h2>
				{#if last7Days.length === 0}
					<p class="text-sm text-muted-foreground">Aucune données disponibles</p>
				{:else}
					<div class="space-y-3">
						{#each last7Days as day (day.date)}
							<div class="flex items-center gap-4">
								<span class="w-20 text-sm text-muted-foreground text-right">{formatDate(day.date)}</span>
								<div class="flex-1 h-7 bg-muted rounded-md overflow-hidden relative">
									<div
										class="h-full bg-foreground/80 rounded-md absolute top-0 left-0 transition-all"
										style="width: {(day.totalMinutesBooked / maxMinutes) * 100}%"
									></div>
								</div>
								<span class="w-8 text-right text-sm font-medium text-foreground">
									{Math.round(day.totalMinutesBooked / 60)}h
								</span>
							</div>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Summary Stats -->
			<div class="bg-card rounded-lg border border-border p-6">
				<h2 class="text-lg font-semibold text-foreground mb-6">Détails</h2>
				<div class="space-y-4">
					<div class="flex items-center justify-between">
						<span class="text-sm text-muted-foreground">Temps d'inactivité total</span>
						<span class="text-sm font-medium text-foreground">
							{data.metrics.summary.totalIdleMinutes} min
						</span>
					</div>
					<div class="flex items-center justify-between">
						<span class="text-sm text-muted-foreground">Fragmentations</span>
						<span class="text-sm font-medium text-foreground">
							{data.metrics.summary.totalFragmentation}
						</span>
					</div>
					<div class="flex items-center justify-between">
						<span class="text-sm text-muted-foreground">Efficacité moyenne</span>
						<span class="text-sm font-medium text-foreground">
							{data.metrics.summary.averageEfficiency.toFixed(2)}
						</span>
					</div>
					<div class="flex items-center justify-between">
						<span class="text-sm text-muted-foreground">Occupation moyenne</span>
						<span class="text-sm font-medium text-foreground">
							{data.metrics.summary.averageUtilization.toFixed(1)}%
						</span>
					</div>
					<div class="flex items-center justify-between">
						<span class="text-sm text-muted-foreground">Jours analysés</span>
						<span class="text-sm font-medium text-foreground">
							{data.metrics.summary.daysAnalyzed}
						</span>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
