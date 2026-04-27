<script lang="ts">
	import type { PageProps } from './$types';
	import { Calendar, Users, Clock, TrendingUp, Activity, Target } from '@lucide/svelte';

	let { data }: PageProps = $props();

	const sessionGrowth = $derived(
		data.stats.sessionsLastMonth > 0
			? Math.round(
					((data.stats.sessionsThisMonth - data.stats.sessionsLastMonth) /
						data.stats.sessionsLastMonth) *
						100,
				)
			: 0,
	);

	const maxWeekSessions = $derived(Math.max(...data.weeklyVolume.map((w) => w.sessions)));
</script>

<svelte:head>
	<title>Statistiques | Staff</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">Statistiques</h1>
		<p class="text-muted-foreground">Vos performances du mois en cours</p>
	</div>

	<!-- KPI Cards -->
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-blue-100 rounded-lg">
					<Calendar size={22} class="text-blue-600" />
				</div>
				<span
					class="text-xs font-medium px-2 py-1 rounded-full {sessionGrowth >= 0
						? 'bg-green-100 text-green-700'
						: 'bg-red-100 text-red-700'}"
				>
					{sessionGrowth >= 0 ? '+' : ''}{sessionGrowth}%
				</span>
			</div>
			<p class="text-3xl font-bold text-foreground mb-1">{data.stats.sessionsThisMonth}</p>
			<p class="text-sm text-muted-foreground">Séances ce mois</p>
		</div>

		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-purple-100 rounded-lg">
					<Users size={22} class="text-purple-600" />
				</div>
				<span class="text-xs text-muted-foreground">Ce mois</span>
			</div>
			<p class="text-3xl font-bold text-foreground mb-1">{data.stats.uniqueClientsThisMonth}</p>
			<p class="text-sm text-muted-foreground">Clients uniques</p>
		</div>

		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-green-100 rounded-lg">
					<Target size={22} class="text-green-600" />
				</div>
				<span class="text-xs text-muted-foreground">Taux</span>
			</div>
			<p class="text-3xl font-bold text-foreground mb-1">{data.stats.attendanceRate}%</p>
			<p class="text-sm text-muted-foreground">Taux de présence</p>
		</div>

		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-yellow-100 rounded-lg">
					<Activity size={22} class="text-yellow-600" />
				</div>
				<span class="text-xs text-muted-foreground">Taux</span>
			</div>
			<p class="text-3xl font-bold text-foreground mb-1">{data.stats.utilizationRate}%</p>
			<p class="text-sm text-muted-foreground">Taux d'occupation</p>
		</div>

		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-orange-100 rounded-lg">
					<Clock size={22} class="text-orange-600" />
				</div>
				<span class="text-xs text-muted-foreground">Moyenne</span>
			</div>
			<p class="text-3xl font-bold text-foreground mb-1">{data.stats.avgSessionDurationMin} min</p>
			<p class="text-sm text-muted-foreground">Durée moyenne</p>
		</div>

		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-pink-100 rounded-lg">
					<TrendingUp size={22} class="text-pink-600" />
				</div>
				<span class="text-xs text-muted-foreground">vs mois dernier</span>
			</div>
			<p class="text-3xl font-bold text-foreground mb-1">{data.stats.sessionsLastMonth}</p>
			<p class="text-sm text-muted-foreground">Séances mois dernier</p>
		</div>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
		<!-- Weekly Volume Chart -->
		<div class="bg-card rounded-lg border border-border p-6">
			<h2 class="text-lg font-semibold text-foreground mb-6">Volume hebdomadaire</h2>
			<div class="space-y-3">
				{#each data.weeklyVolume as week (week.week)}
					<div class="flex items-center gap-4">
						<span class="w-20 text-sm text-muted-foreground text-right">{week.week}</span>
						<div class="flex-1 h-7 bg-muted rounded-md overflow-hidden relative">
							<div
								class="h-full bg-foreground/80 rounded-md absolute top-0 left-0 transition-all"
								style="width: {maxWeekSessions > 0 ? (week.sessions / maxWeekSessions) * 100 : 0}%"
							></div>
						</div>
						<span class="w-8 text-right text-sm font-medium text-foreground">{week.sessions}</span>
					</div>
				{/each}
			</div>
		</div>

		<!-- Top Services -->
		<div class="bg-card rounded-lg border border-border p-6">
			<h2 class="text-lg font-semibold text-foreground mb-6">Services les plus demandés</h2>
			<div class="space-y-4">
				{#each data.topServices as service, i (service.name)}
					<div>
						<div class="flex items-center justify-between mb-1.5">
							<div class="flex items-center gap-2">
								<span
									class="w-6 h-6 rounded-full bg-muted flex items-center justify-center text-xs font-bold text-foreground"
								>
									{i + 1}
								</span>
								<span class="text-sm font-medium text-foreground">{service.name}</span>
							</div>
							<span class="text-sm text-muted-foreground">{service.sessions} séances</span>
						</div>
						<div class="h-2 bg-muted rounded-full overflow-hidden">
							<div
								class="h-full bg-foreground/70 rounded-full"
								style="width: {service.percentage}%"
							></div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>
