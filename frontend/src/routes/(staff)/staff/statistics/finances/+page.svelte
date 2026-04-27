<script lang="ts">
	import type { PageProps } from './$types';
	import { Banknote, Clock, TrendingUp, ArrowDownToLine } from '@lucide/svelte';

	let { data }: PageProps = $props();

	function formatCents(cents: number): string {
		return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(
			cents / 100,
		);
	}

	function formatDate(iso: string): string {
		return new Date(iso).toLocaleDateString('fr-FR', {
			day: 'numeric',
			month: 'short',
			year: 'numeric',
		});
	}

	function formatDateTime(iso: string): string {
		return new Date(iso).toLocaleDateString('fr-FR', {
			day: 'numeric',
			month: 'short',
			hour: '2-digit',
			minute: '2-digit',
		});
	}

	function statusBadge(status: string): string {
		switch (status) {
			case 'paid': return 'bg-green-100 text-green-700';
			case 'pending': return 'bg-yellow-100 text-yellow-700';
			case 'refunded': return 'bg-red-100 text-red-700';
			default: return 'bg-gray-100 text-gray-700';
		}
	}

	function statusLabel(status: string): string {
		switch (status) {
			case 'paid': return 'Payé';
			case 'pending': return 'En attente';
			case 'refunded': return 'Remboursé';
			default: return status;
		}
	}

	const growthPct = $derived(
		data.summary.lastMonthInCents > 0
			? Math.round(
					((data.summary.currentMonthInCents - data.summary.lastMonthInCents) /
						data.summary.lastMonthInCents) *
						100,
				)
			: 0,
	);

	const maxEarning = $derived(Math.max(...data.monthlyEarnings.map((m) => m.amountInCents)));
</script>

<svelte:head>
	<title>Finances | Staff</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">Finances</h1>
		<p class="text-muted-foreground">Suivi de vos revenus et paiements</p>
	</div>

	<!-- Summary Cards -->
	<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-green-100 rounded-lg">
					<Banknote size={22} class="text-green-600" />
				</div>
				<span
					class="text-xs font-medium px-2 py-1 rounded-full {growthPct >= 0
						? 'bg-green-100 text-green-700'
						: 'bg-red-100 text-red-700'}"
				>
					{growthPct >= 0 ? '+' : ''}{growthPct}%
				</span>
			</div>
			<p class="text-2xl font-bold text-foreground mb-1">
				{formatCents(data.summary.currentMonthInCents)}
			</p>
			<p class="text-sm text-muted-foreground">Revenus ce mois</p>
		</div>

		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-blue-100 rounded-lg">
					<TrendingUp size={22} class="text-blue-600" />
				</div>
				<span class="text-xs text-muted-foreground">Mois dernier</span>
			</div>
			<p class="text-2xl font-bold text-foreground mb-1">
				{formatCents(data.summary.lastMonthInCents)}
			</p>
			<p class="text-sm text-muted-foreground">Revenus mois précédent</p>
		</div>

		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-yellow-100 rounded-lg">
					<Clock size={22} class="text-yellow-600" />
				</div>
				<span class="text-xs text-muted-foreground">En cours</span>
			</div>
			<p class="text-2xl font-bold text-foreground mb-1">
				{formatCents(data.summary.pendingInCents)}
			</p>
			<p class="text-sm text-muted-foreground">Paiements en attente</p>
		</div>

		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-purple-100 rounded-lg">
					<ArrowDownToLine size={22} class="text-purple-600" />
				</div>
				<span class="text-xs text-muted-foreground">
					{formatDate(data.summary.nextPayoutDate)}
				</span>
			</div>
			<p class="text-2xl font-bold text-foreground mb-1">
				{formatCents(data.summary.nextPayoutInCents)}
			</p>
			<p class="text-sm text-muted-foreground">Prochain virement</p>
		</div>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
		<!-- Monthly Chart -->
		<div class="bg-card rounded-lg border border-border p-6">
			<h2 class="text-lg font-semibold text-foreground mb-6">Revenus mensuels</h2>
			<div class="space-y-3">
				{#each data.monthlyEarnings as month (month.month)}
					<div class="flex items-center gap-4">
						<span class="w-10 text-sm text-muted-foreground">{month.month}</span>
						<div class="flex-1 h-7 bg-muted rounded-md overflow-hidden relative">
							<div
								class="h-full bg-foreground/80 rounded-md absolute top-0 left-0 transition-all"
								style="width: {maxEarning > 0 ? (month.amountInCents / maxEarning) * 100 : 0}%"
							></div>
						</div>
						<span class="w-24 text-right text-sm font-medium text-foreground">
							{formatCents(month.amountInCents)}
						</span>
					</div>
				{/each}
			</div>
		</div>

		<!-- Transaction History -->
		<div class="bg-card rounded-lg border border-border p-6">
			<h2 class="text-lg font-semibold text-foreground mb-6">Transactions récentes</h2>
			<div class="space-y-3">
				{#each data.transactions as tx (tx.id)}
					<div class="flex items-center gap-3">
						<div class="flex-1 min-w-0">
							<p class="text-sm font-medium text-foreground truncate">{tx.clientName}</p>
							<p class="text-xs text-muted-foreground truncate">{tx.productName}</p>
							<p class="text-xs text-muted-foreground">{formatDateTime(tx.date)}</p>
						</div>
						<div class="flex flex-col items-end gap-1 flex-shrink-0">
							<p
								class="text-sm font-semibold {tx.status === 'refunded'
									? 'text-red-600'
									: 'text-foreground'}"
							>
								{tx.status === 'refunded' ? '-' : ''}{formatCents(tx.amountInCents)}
							</p>
							<span class="px-2 py-0.5 rounded-full text-xs font-medium {statusBadge(tx.status)}">
								{statusLabel(tx.status)}
							</span>
						</div>
					</div>
				{/each}
			</div>
		</div>
	</div>
</div>
