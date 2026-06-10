<script lang="ts">
	import type { PageProps } from './$types';
	import { Banknote, Clock, TrendingUp, ArrowDownToLine, Calendar, Inbox } from '@lucide/svelte';
	import { goto } from '$app/navigation';

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

	function paymentStatusBadge(status: string): string {
		switch (status) {
			case 'paid': return 'bg-green-100 text-green-700';
			case 'pending': return 'bg-yellow-100 text-yellow-700';
			case 'refunded': return 'bg-red-100 text-red-700';
			default: return 'bg-gray-100 text-gray-700';
		}
	}

	function paymentStatusLabel(status: string): string {
		switch (status) {
			case 'paid': return 'Payé';
			case 'pending': return 'En attente';
			case 'refunded': return 'Remboursé';
			default: return status;
		}
	}

	function bookingStatusBadge(status: string): string {
		switch (status) {
			case 'completed': return 'bg-green-100 text-green-700';
			case 'confirmed': return 'bg-blue-100 text-blue-700';
			case 'cancelled': return 'bg-red-100 text-red-700';
			case 'no_show': return 'bg-yellow-100 text-yellow-700';
			default: return 'bg-gray-100 text-gray-700';
		}
	}

	function bookingStatusLabel(status: string): string {
		switch (status) {
			case 'completed': return 'Complétée';
			case 'confirmed': return 'Confirmée';
			case 'cancelled': return 'Annulée';
			case 'no_show': return 'Absence';
			default: return status;
		}
	}

	function getYearMonth(iso: string): string {
		const d = new Date(iso);
		return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
	}

	function prevYearMonth(ym: string): string {
		const [y, m] = ym.split('-').map(Number);
		const d = new Date(y, m - 2, 1);
		return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}`;
	}

	// Month selector
	let selectedMonth: string = $derived(data.selectedMonth);

	function onMonthChange(e: Event) {
		const target = e.target as HTMLInputElement;
		const value = target.value;
		if (!value) return;
		goto(`/staff/statistics/finances?month=${value}`);
	}

	// Filtered transactions for the selected month
	interface Transaction {
		id: string;
		slotStartTime: string;
		productId: string;
		productName: string;
		amountCents: number;
		paymentStatus: 'paid' | 'pending' | 'refunded';
		bookingStatus: 'confirmed' | 'cancelled' | 'completed' | 'no_show';
	}

	const filteredTransactions = $derived(
		(data.summary?.transactions ?? []).filter((tx: Transaction) =>
			getYearMonth(tx.slotStartTime) === selectedMonth,
		),
	);

	// Recompute KPIs for the selected month
	const selectedMonthPaid = $derived(
		filteredTransactions
			.filter((tx: Transaction) => tx.paymentStatus === 'paid')
			.reduce((sum: number, tx: Transaction) => sum + tx.amountCents, 0),
	);

	const prevMonthKey = $derived(prevYearMonth(selectedMonth));
	const prevMonthPaid = $derived(
		(data.summary?.transactions ?? [])
			.filter((tx: Transaction) => tx.paymentStatus === 'paid' && getYearMonth(tx.slotStartTime) === prevMonthKey)
			.reduce((sum: number, tx: Transaction) => sum + tx.amountCents, 0),
	);

	const selectedMonthPending = $derived(
		filteredTransactions
			.filter((tx: Transaction) => tx.paymentStatus === 'pending')
			.reduce((sum: number, tx: Transaction) => sum + tx.amountCents, 0),
	);

	const growthPct = $derived(
		prevMonthPaid > 0
			? Math.round(((selectedMonthPaid - prevMonthPaid) / prevMonthPaid) * 100)
			: 0,
	);

	// Monthly chart data (from all transactions)
	interface MonthlyEarning {
		key: string;
		label: string;
		amountInCents: number;
	}

	const monthlyEarnings = $derived.by((): MonthlyEarning[] => {
		if (!data.summary?.transactions) return [];
		const byMonth = new Map<string, MonthlyEarning>();
		data.summary.transactions.forEach((tx: Transaction) => {
			const key = getYearMonth(tx.slotStartTime);
			const d = new Date(tx.slotStartTime);
			const label = d.toLocaleDateString('fr-FR', { month: 'short', year: '2-digit' });
			const existing = byMonth.get(key);
			byMonth.set(key, { key, label, amountInCents: (existing?.amountInCents ?? 0) + tx.amountCents });
		});
		return Array.from(byMonth.values())
			.sort((a, b) => a.key.localeCompare(b.key))
			.slice(-6);
	});

	const maxEarning = $derived(
		monthlyEarnings.length > 0 ? Math.max(...monthlyEarnings.map((m) => m.amountInCents)) : 1,
	);
</script>

<svelte:head>
	<title>Finances | Staff</title>
</svelte:head>

<div class="p-6 lg:p-10">
	<div class="mb-10 flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4">
		<div>
			<p class="text-[11px] font-semibold uppercase tracking-[0.2em] text-muted-foreground mb-3">Revenus</p>
			<h1 class="font-display text-4xl lg:text-5xl font-semibold tracking-tight text-foreground leading-[1.1]">Finances</h1>
			<p class="text-muted-foreground mt-3 text-sm">Suivi de vos revenus et paiements</p>
			<div class="mt-4 h-px w-16 bg-foreground/20"></div>
		</div>
		<!-- Month picker -->
		<div class="flex items-center gap-2">
			<Calendar size={18} class="text-muted-foreground" />
			<input
				type="month"
				value={selectedMonth}
				onchange={onMonthChange}
				class="bg-background border border-border rounded-lg px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
			/>
		</div>
	</div>

	{#if data.error}
		<div class="bg-destructive/10 border border-destructive/30 text-destructive px-4 py-3 rounded-lg mb-8">
			{data.error}
		</div>
	{:else if !data.summary}
		<div class="bg-card rounded-lg border border-border-card p-12 text-center">
			<Banknote size={48} class="text-muted-foreground mx-auto mb-4" />
			<h2 class="text-xl font-semibold text-foreground mb-2">Aucune donnée financière</h2>
			<p class="text-sm text-muted-foreground">Vos revenus apparaîtront ici dès que vous recevrez vos premières réservations.</p>
		</div>
	{:else}
		<!-- Summary Cards (recomputed for selected month) -->
		<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
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
					{formatCents(selectedMonthPaid)}
				</p>
				<p class="text-sm text-muted-foreground">Revenus ce mois</p>
			</div>

			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-blue-100 rounded-lg">
						<TrendingUp size={22} class="text-blue-600" />
					</div>
					<span class="text-xs text-muted-foreground">Mois précédent</span>
				</div>
				<p class="text-2xl font-bold text-foreground mb-1">
					{formatCents(prevMonthPaid)}
				</p>
				<p class="text-sm text-muted-foreground">Revenus mois précédent</p>
			</div>

			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
				<div class="flex items-center justify-between mb-4">
					<div class="p-3 bg-yellow-100 rounded-lg">
						<Clock size={22} class="text-yellow-600" />
					</div>
					<span class="text-xs text-muted-foreground">En cours</span>
				</div>
				<p class="text-2xl font-bold text-foreground mb-1">
					{formatCents(selectedMonthPending)}
				</p>
				<p class="text-sm text-muted-foreground">Paiements en attente</p>
			</div>

			<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
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
			<!-- Monthly Chart (derived from all transactions) -->
			{#if monthlyEarnings.length > 0}
				<div class="bg-card rounded-2xl border border-border-card p-6 group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500">
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

			<!-- Transaction Table -->
			<div class="bg-card rounded-2xl border border-border-card group/card hover:shadow-card hover:shadow-foreground/[0.03] hover:-translate-y-0.5 transition-all duration-500 flex flex-col">
				<div class="p-6 border-b border-border">
					<h2 class="text-lg font-semibold text-foreground">Transactions</h2>
				</div>
				{#if filteredTransactions.length === 0}
					<div class="p-12 text-center flex-1 flex flex-col items-center justify-center">
						<Inbox size={48} class="text-muted-foreground mx-auto mb-4" />
						<h3 class="text-lg font-semibold text-foreground mb-2">Aucune transaction</h3>
						<p class="text-sm text-muted-foreground">Aucune transaction pour le mois sélectionné.</p>
					</div>
				{:else}
					<div class="overflow-x-auto">
						<table class="w-full">
							<thead>
								<tr class="bg-muted/50 border-b border-border">
									<th class="text-left py-3 px-6 text-xs font-medium text-muted-foreground uppercase tracking-wider">Date / Heure</th>
									<th class="text-left py-3 px-6 text-xs font-medium text-muted-foreground uppercase tracking-wider">Prestation</th>
									<th class="text-left py-3 px-6 text-xs font-medium text-muted-foreground uppercase tracking-wider">Statut paiement</th>
									<th class="text-left py-3 px-6 text-xs font-medium text-muted-foreground uppercase tracking-wider">Statut réservation</th>
									<th class="text-right py-3 px-6 text-xs font-medium text-muted-foreground uppercase tracking-wider">Montant</th>
								</tr>
							</thead>
							<tbody>
								{#each filteredTransactions as tx (tx.id)}
									<tr class="border-b border-border hover:bg-muted/30 transition-colors">
										<td class="py-3 px-6 text-sm text-muted-foreground whitespace-nowrap">
											{formatDateTime(tx.slotStartTime)}
										</td>
										<td class="py-3 px-6">
											<span class="text-sm font-medium text-foreground">{tx.productName}</span>
										</td>
										<td class="py-3 px-6">
											<span class="px-2.5 py-0.5 rounded-full text-xs font-medium {paymentStatusBadge(tx.paymentStatus)}">
												{paymentStatusLabel(tx.paymentStatus)}
											</span>
										</td>
										<td class="py-3 px-6">
											<span class="px-2.5 py-0.5 rounded-full text-xs font-medium {bookingStatusBadge(tx.bookingStatus)}">
												{bookingStatusLabel(tx.bookingStatus)}
											</span>
										</td>
										<td class="py-3 px-6 text-right">
											<span class="text-sm font-medium {tx.paymentStatus === 'refunded'
												? 'text-red-600'
												: 'text-foreground'}">
												{tx.paymentStatus === 'refunded' ? '-' : ''}{formatCents(tx.amountCents)}
											</span>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>
