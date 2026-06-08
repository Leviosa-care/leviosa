<script lang="ts">
	import type { PageProps } from './$types';
	import { Banknote, Clock, TrendingUp, ArrowDownToLine } from '@lucide/svelte';

	let { data }: PageProps = $props();

	interface Transaction {
		id: string;
		slotStartTime: string;
		productId: string;
		productName: string;
		amountCents: number;
		paymentStatus: 'paid' | 'pending' | 'refunded';
	}

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
		data.summary?.lastMonthCents && data.summary.lastMonthCents > 0
			? Math.round(
					((data.summary.currentMonthCents - data.summary.lastMonthCents) /
						data.summary.lastMonthCents) *
						100,
				)
			: 0,
	);

	interface MonthlyEarning {
		key: string;
		label: string;
		amountInCents: number;
	}

	const monthlyEarnings = $derived.by((): MonthlyEarning[] => {
		if (!data.summary?.transactions) return [];
		const byMonth = new Map<string, MonthlyEarning>();
		data.summary.transactions.forEach((tx: Transaction) => {
			const date = new Date(tx.slotStartTime);
			const key = `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}`;
			const label = date.toLocaleDateString('fr-FR', { month: 'short', year: '2-digit' });
			const existing = byMonth.get(key);
			byMonth.set(key, { key, label, amountInCents: (existing?.amountInCents ?? 0) + tx.amountCents });
		});
		return Array.from(byMonth.values())
			.sort((a, b) => a.key.localeCompare(b.key))
			.slice(-6);
	});

	const maxEarning = $derived(
		monthlyEarnings.length > 0 ? Math.max(...monthlyEarnings.map((m) => m.amountInCents)) : 1
	);
</script>

<svelte:head>
	<title>Finances | Staff</title>
</svelte:head>

<div class="p-6 lg:p-10">
	<div class="mb-10">
		<p class="text-[11px] font-semibold uppercase tracking-[0.2em] text-muted-foreground mb-3">Finances</p>
		<h1 class="font-display text-4xl lg:text-5xl font-semibold tracking-tight text-foreground leading-[1.1]">Finances</h1>
		<p class="text-muted-foreground mt-3 text-sm">Suivi de vos revenus et paiements</p>
		<div class="mt-4 h-px w-16 bg-foreground/20"></div>
	</div>

	{#if data.error}
		<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-8">
			{data.error}
		</div>
	{:else if !data.summary}
		<div class="bg-card rounded-lg border border-border p-12 text-center">
			<Banknote size={48} class="text-muted-foreground mx-auto mb-4" />
			<h2 class="text-xl font-semibold text-foreground mb-2">Aucune donnée financière</h2>
			<p class="text-sm text-muted-foreground">Vos revenus apparaîtront ici dès que vous recevrez vos premières réservations.</p>
		</div>
	{:else if data.summary.transactions.length === 0 && data.summary.currentMonthCents === 0 && data.summary.lastMonthCents === 0}
		<div class="bg-card rounded-lg border border-border p-12 text-center">
			<Banknote size={48} class="text-muted-foreground mx-auto mb-4" />
			<h2 class="text-xl font-semibold text-foreground mb-2">Aucune donnée financière</h2>
			<p class="text-sm text-muted-foreground">Vos revenus apparaîtront ici dès que vous recevrez vos premières réservations.</p>
		</div>
	{:else if data.summary}
		<!-- Summary Cards -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
			<div class="bg-background rounded-2xl border border-border p-6 group/card hover:shadow-lg hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
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
					{formatCents(data.summary.currentMonthCents)}
				</p>
				<p class="text-sm text-muted-foreground">Revenus ce mois</p>
			</div>

			<div class="bg-background rounded-2xl border border-border p-6 group/card hover:shadow-lg hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-blue-100 rounded-lg">
						<TrendingUp size={22} class="text-blue-600" />
					</div>
					<span class="text-xs text-muted-foreground">Mois dernier</span>
				</div>
				<p class="text-2xl font-bold text-foreground mb-1">
					{formatCents(data.summary.lastMonthCents)}
				</p>
				<p class="text-sm text-muted-foreground">Revenus mois précédent</p>
			</div>

			<div class="bg-background rounded-2xl border border-border p-6 group/card hover:shadow-lg hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-yellow-100 rounded-lg">
						<Clock size={22} class="text-yellow-600" />
					</div>
					<span class="text-xs text-muted-foreground">En cours</span>
				</div>
				<p class="text-2xl font-bold text-foreground mb-1">
					{formatCents(data.summary.pendingCents)}
				</p>
				<p class="text-sm text-muted-foreground">Paiements en attente</p>
			</div>

			<div class="bg-background rounded-2xl border border-border p-6 group/card hover:shadow-lg hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-purple-100 rounded-lg">
						<ArrowDownToLine size={22} class="text-purple-600" />
					</div>
					<span class="text-xs text-muted-foreground">
						{formatDate(data.summary.nextPayoutDate)}
					</span>
				</div>
				<p class="text-2xl font-bold text-foreground mb-1">
					{formatCents(data.summary.nextPayoutCents)}
				</p>
				<p class="text-sm text-muted-foreground">Prochain virement</p>
			</div>
		</div>

		<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
			<!-- Monthly Chart (derived from transactions) -->
			{#if monthlyEarnings.length > 0}
				<div class="bg-background rounded-2xl border border-border p-6 group/card hover:shadow-lg hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
					<h2 class="text-lg font-semibold text-foreground mb-6">Revenus mensuels</h2>
					<div class="space-y-3">
						{#each monthlyEarnings as month (month.key)}
							<div class="flex items-center gap-4">
								<span class="w-10 text-sm text-muted-foreground">{month.label}</span>
								<div class="flex-1 h-7 bg-muted rounded-md overflow-hidden relative">
									<div
										class="h-full bg-foreground/80 rounded-md absolute top-0 left-0 transition-all"
										style="width: {(month.amountInCents / maxEarning) * 100}%"
									></div>
								</div>
								<span class="w-24 text-right text-sm font-medium text-foreground">
									{formatCents(month.amountInCents)}
								</span>
							</div>
						{/each}
					</div>
				</div>
			{/if}

			<!-- Transaction History -->
			<div class="bg-background rounded-2xl border border-border p-6 group/card hover:shadow-lg hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
				<h2 class="text-lg font-semibold text-foreground mb-6">Transactions récentes</h2>
				{#if data.summary.transactions.length === 0}
					<p class="text-sm text-muted-foreground">Aucune transaction</p>
				{:else}
					<div class="space-y-3">
						{#each data.summary.transactions as tx (tx.id)}
							<div class="flex items-center gap-3">
								<div class="flex-1 min-w-0">
									<p class="text-sm font-medium text-foreground truncate">
										Réservation #{tx.id.slice(0, 8)}
									</p>
									<p class="text-xs text-muted-foreground truncate">{tx.productName}</p>
									<p class="text-xs text-muted-foreground">{formatDateTime(tx.slotStartTime)}</p>
								</div>
								<div class="flex flex-col items-end gap-1 flex-shrink-0">
									<p
										class="text-sm font-semibold {tx.paymentStatus === 'refunded'
											? 'text-red-600'
											: 'text-foreground'}"
									>
										{tx.paymentStatus === 'refunded' ? '-' : ''}{formatCents(tx.amountCents)}
									</p>
									<span class="px-2 py-0.5 rounded-full text-xs font-medium {statusBadge(tx.paymentStatus)}">
										{statusLabel(tx.paymentStatus)}
									</span>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
