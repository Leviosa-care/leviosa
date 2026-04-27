<script lang="ts">
	import type { PageProps } from './$types';
	import { TrendingUp, Users, Calendar, DollarSign, Repeat, Target, Award } from '@lucide/svelte';

	let { data }: PageProps = $props();

	function formatCents(cents: number): string {
		return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(cents / 100);
	}

	function formatNumber(num: number): string {
		return new Intl.NumberFormat('fr-FR').format(num);
	}

	const maxRevenue = $derived(Math.max(...data.monthlyRevenue.map(m => m.amountInCents)));
</script>

<svelte:head>
	<title>Analytics | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">
			Analytics
		</h1>
		<p class="text-muted-foreground">
			Tableau de bord des performances de votre activité
		</p>
	</div>

	<!-- KPI Cards -->
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
		<!-- Total Revenue -->
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-green-100 dark:bg-green-900/30 rounded-lg">
					<DollarSign size={24} class="text-green-600 dark:text-green-400" />
				</div>
				<span class="text-xs text-muted-foreground">Ce mois</span>
			</div>
			<div class="text-3xl font-bold text-foreground mb-1">
				{formatCents(data.stats.totalRevenueInCents)}
			</div>
			<p class="text-sm text-muted-foreground">Chiffre d'affaires total</p>
		</div>

		<!-- Bookings -->
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
					<Calendar size={24} class="text-blue-600 dark:text-blue-400" />
				</div>
				<span class="text-xs text-muted-foreground">Ce mois</span>
			</div>
			<div class="text-3xl font-bold text-foreground mb-1">
				{formatNumber(data.stats.bookingsThisMonth)}
			</div>
			<p class="text-sm text-muted-foreground">Réservations</p>
		</div>

		<!-- New Users -->
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-purple-100 dark:bg-purple-900/30 rounded-lg">
					<Users size={24} class="text-purple-600 dark:text-purple-400" />
				</div>
				<span class="text-xs text-muted-foreground">Ce mois</span>
			</div>
			<div class="text-3xl font-bold text-foreground mb-1">
				{formatNumber(data.stats.newUsersThisMonth)}
			</div>
			<p class="text-sm text-muted-foreground">Nouveaux clients</p>
		</div>

		<!-- Avg Booking Value -->
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-yellow-100 dark:bg-yellow-900/30 rounded-lg">
					<TrendingUp size={24} class="text-yellow-600 dark:text-yellow-400" />
				</div>
				<span class="text-xs text-muted-foreground">Moyenne</span>
			</div>
			<div class="text-3xl font-bold text-foreground mb-1">
				{formatCents(data.stats.avgBookingValue)}
			</div>
			<p class="text-sm text-muted-foreground">Panier moyen</p>
		</div>

		<!-- Repeat Rate -->
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-orange-100 dark:bg-orange-900/30 rounded-lg">
					<Repeat size={24} class="text-orange-600 dark:text-orange-400" />
				</div>
				<span class="text-xs text-muted-foreground">Taux</span>
			</div>
			<div class="text-3xl font-bold text-foreground mb-1">
				{data.stats.repeatRate}%
			</div>
			<p class="text-sm text-muted-foreground">Taux de fidélisation</p>
		</div>

		<!-- Conversion Rate -->
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-pink-100 dark:bg-pink-900/30 rounded-lg">
					<Target size={24} class="text-pink-600 dark:text-pink-400" />
				</div>
				<span class="text-xs text-muted-foreground">Taux</span>
			</div>
			<div class="text-3xl font-bold text-foreground mb-1">
				{data.stats.conversionRate}%
			</div>
			<p class="text-sm text-muted-foreground">Taux de conversion</p>
		</div>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
		<!-- Monthly Revenue Chart -->
		<div class="bg-card rounded-lg border border-border p-6">
			<h2 class="text-xl font-semibold text-foreground mb-6">Revenus mensuels</h2>
			<div class="space-y-4">
				{#each data.monthlyRevenue as month}
					<div class="flex items-center gap-4">
						<span class="w-12 text-sm font-medium text-muted-foreground">{month.month}</span>
						<div class="flex-1 h-8 bg-muted rounded-md overflow-hidden relative">
							<div
								class="h-full bg-primary rounded-md absolute top-0 left-0 transition-all"
								style="width: {(month.amountInCents / maxRevenue * 100)}%"
							></div>
						</div>
						<span class="w-24 text-right text-sm font-medium text-foreground">
							{formatCents(month.amountInCents)}
						</span>
					</div>
				{/each}
			</div>
		</div>

		<!-- Top Products -->
		<div class="bg-card rounded-lg border border-border p-6">
			<h2 class="text-xl font-semibold text-foreground mb-6 flex items-center gap-2">
				<Award size={20} class="text-yellow-500" />
				Top produits
			</h2>
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead>
						<tr class="border-b border-border">
							<th class="text-left py-3 px-2 text-sm font-medium text-muted-foreground">Produit</th>
							<th class="text-right py-3 px-2 text-sm font-medium text-muted-foreground">Réservations</th>
							<th class="text-right py-3 px-2 text-sm font-medium text-muted-foreground">Revenus</th>
						</tr>
					</thead>
					<tbody>
						{#each data.topProducts as product, index}
							<tr class="border-b border-border hover:bg-muted/30 transition-colors">
								<td class="py-3 px-2">
									<div class="flex items-center gap-2">
										<span class="flex items-center justify-center w-6 h-6 rounded-full bg-primary/10 text-primary text-xs font-bold">
											{index + 1}
										</span>
										<span class="text-sm font-medium text-foreground">{product.name}</span>
									</div>
								</td>
								<td class="text-right py-3 px-2 text-sm text-foreground">{formatNumber(product.bookings)}</td>
								<td class="text-right py-3 px-2 text-sm font-medium text-foreground">{formatCents(product.revenueInCents)}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>
