<script lang="ts">
	import type { PageProps } from './$types';
	import { Activity, Clock, Target, CalendarDays } from '@lucide/svelte';
	import { goto } from '$app/navigation';

	let { data }: PageProps = $props();

	const dailyMetrics = $derived(
		data.metrics?.roomMetrics.flatMap((rm) => rm.dailyMetrics) ?? []
	);

	const totalBookedMinutes = $derived(
		dailyMetrics.reduce((sum, d) => sum + d.totalMinutesBooked, 0)
	);
	const avgUtilization = $derived(data.metrics?.summary.averageUtilization ?? 0);
	const totalSessions = $derived(data.metrics?.summary.daysAnalyzed ?? 0);
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

	// Date range selector state
	let startDate = $state(data.startDate);
	let endDate = $state(data.endDate);

	function applyDateRange() {
		const params = new URLSearchParams();
		params.set('start_date', startDate);
		params.set('end_date', endDate);
		goto(`/staff/statistics/analytics?${params.toString()}`);
	}

	function setPresetRange(days: number) {
		const now = new Date();
		const start = new Date(now.getFullYear(), now.getMonth(), now.getDate() - (days - 1));
		startDate = start.toISOString().split('T')[0];
		endDate = now.toISOString().split('T')[0];
		applyDateRange();
	}

	const hasNoData = $derived(
		!data.metrics ||
		(data.metrics.roomMetrics.length === 0 && data.metrics.summary.daysAnalyzed === 0)
	);
</script>

<svelte:head>
	<title>Statistiques | Staff</title>
</svelte:head>

<div class="p-6 lg:p-10">
	<div class="mb-10">
		<p class="text-[11px] font-semibold uppercase tracking-[0.2em] text-muted-foreground mb-3">Statistiques</p>
		<h1 class="font-display text-4xl lg:text-5xl font-semibold tracking-tight text-foreground leading-[1.1]">
			Analytics
		</h1>
		<p class="text-muted-foreground mt-3 text-sm">
			Vos performances sur la période sélectionnée
		</p>
		<div class="mt-4 h-px w-16 bg-foreground/20"></div>
	</div>

	<!-- Date Range Selector -->
	<div class="bg-card rounded-lg border border-border-card p-4 mb-8">
		<div class="flex flex-col sm:flex-row items-start sm:items-end gap-4">
			<div class="flex flex-col gap-1.5">
				<label for="start_date" class="text-sm font-medium text-foreground flex items-center gap-1.5">
					<CalendarDays size={14} />
					Début
				</label>
				<input
					id="start_date"
					type="date"
					bind:value={startDate}
					class="px-3 py-2 rounded-lg border border-border bg-background text-foreground text-sm focus:ring-2 focus:ring-foreground focus:outline-none"
				/>
			</div>
			<div class="flex flex-col gap-1.5">
				<label for="end_date" class="text-sm font-medium text-foreground flex items-center gap-1.5">
					<CalendarDays size={14} />
					Fin
				</label>
				<input
					id="end_date"
					type="date"
					bind:value={endDate}
					class="px-3 py-2 rounded-lg border border-border bg-background text-foreground text-sm focus:ring-2 focus:ring-foreground focus:outline-none"
				/>
			</div>
			<button
				onclick={applyDateRange}
				class="px-4 py-2 rounded-lg bg-foreground text-background text-sm font-medium hover:bg-foreground/90 transition-colors cursor-pointer"
			>
				Appliquer
			</button>
			<div class="flex gap-2 sm:ml-auto">
				<button
					onclick={() => setPresetRange(7)}
					class="px-3 py-1.5 rounded-md text-xs font-medium border border-border text-foreground hover:bg-muted transition-colors cursor-pointer"
				>
					7 jours
				</button>
				<button
					onclick={() => setPresetRange(30)}
					class="px-3 py-1.5 rounded-md text-xs font-medium border border-border text-foreground hover:bg-muted transition-colors cursor-pointer"
				>
					30 jours
				</button>
				<button
					onclick={() => setPresetRange(90)}
					class="px-3 py-1.5 rounded-md text-xs font-medium border border-border text-foreground hover:bg-muted transition-colors cursor-pointer"
				>
					90 jours
				</button>
			</div>
		</div>
	</div>

	{#if data.error}
		<div class="bg-destructive/10 border border-destructive/30 text-destructive px-4 py-3 rounded-lg mb-8">
			{data.error}
		</div>
	{/if}

	{#if !data.error && hasNoData}
		<div class="bg-card rounded-lg border border-border-card p-12 text-center">
			<Target size={48} class="text-muted-foreground mx-auto mb-4" />
			<h2 class="text-xl font-semibold text-foreground mb-2">Aucune donnée disponible</h2>
			<p class="text-sm text-muted-foreground">Les statistiques d'utilisation apparaîtront ici dès que vous aurez des réservations sur la période sélectionnée.</p>
		</div>
	{:else if !data.error && data.metrics}
		<!-- KPI Cards -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
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

			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-blue-100 rounded-lg">
						<Activity size={22} class="text-blue-600" />
					</div>
					<span class="text-xs text-muted-foreground">Total</span>
				</div>
				<p class="text-3xl font-bold text-foreground mb-1">{totalSessions}</p>
				<p class="text-sm text-muted-foreground">Jours analysés</p>
			</div>

			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-yellow-100 rounded-lg">
						<Activity size={22} class="text-yellow-600" />
					</div>
					<span class="text-xs text-muted-foreground">Score</span>
				</div>
				<p class="text-3xl font-bold text-foreground mb-1">{avgEfficiency.toFixed(1)}</p>
				<p class="text-sm text-muted-foreground">Score d'efficacité moyen</p>
			</div>

			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
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
			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
				<h2 class="text-lg font-semibold text-foreground mb-6">Volume quotidien (7 jours)</h2>
				{#if last7Days.length === 0}
					<p class="text-sm text-muted-foreground">Aucune donnée disponible pour les 7 derniers jours</p>
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
			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
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
