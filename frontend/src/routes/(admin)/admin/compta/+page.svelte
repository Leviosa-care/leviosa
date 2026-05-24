<script lang="ts">
	import type { PageProps } from './$types';
	import { DollarSign, ArrowDown, TrendingUp, CreditCard, Calendar } from '@lucide/svelte';
	import { goto } from '$app/navigation';

	let { data }: PageProps = $props();

	function formatCents(cents: number): string {
		return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(cents / 100);
	}

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

	function getTypeBadge(type: string) {
		return type === 'payment'
			? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
			: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300';
	}

	function getTypeLabel(type: string): string {
		return type === 'payment' ? 'Paiement' : 'Remboursement';
	}

	function getBookingStatusBadge(status: string) {
		switch (status) {
			case 'completed':
				return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300';
			case 'cancelled':
				return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300';
			case 'confirmed':
				return 'bg-blue-100 text-blue-700 dark:bg-blue-900 dark:text-blue-300';
			case 'no_show':
				return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300';
			default:
				return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
		}
	}

	function getBookingStatusLabel(status: string): string {
		switch (status) {
			case 'completed': return 'Complété';
			case 'cancelled': return 'Annulé';
			case 'confirmed': return 'Confirmé';
			case 'no_show': return 'Absence';
			default: return status;
		}
	}

	function getPaymentStatusBadge(status: string) {
		switch (status) {
			case 'paid':
				return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300';
			case 'refunded':
				return 'bg-orange-100 text-orange-700 dark:bg-orange-900 dark:text-orange-300';
			default:
				return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
		}
	}

	function getPaymentStatusLabel(status: string): string {
		switch (status) {
			case 'paid': return 'Payé';
			case 'refunded': return 'Remboursé';
			default: return status;
		}
	}

	// Month selector: compute the selected month from the 'from' date
	let selectedMonth: string = $derived(data.from.substring(0, 7));

	function onMonthChange(e: Event) {
		const target = e.target as HTMLInputElement;
		const yearMonth = target.value; // "YYYY-MM"
		if (!yearMonth) return;
		const [year, month] = yearMonth.split('-').map(Number);
		const from = `${yearMonth}-01`;
		// Last day of selected month
		const lastDay = new Date(year, month, 0).getDate();
		const to = `${yearMonth}-${String(lastDay).padStart(2, '0')}`;
		goto(`/admin/compta?from=${from}&to=${to}`);
	}
</script>

<svelte:head>
	<title>Comptabilité | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
		<div>
			<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">
				Comptabilité
			</h1>
			<p class="text-muted-foreground">
				Suivi des revenus et des transactions
			</p>
		</div>
		<!-- Month selector -->
		<div class="flex items-center gap-2">
			<Calendar size={18} class="text-muted-foreground" />
			<input
				type="month"
				value={selectedMonth}
				onchange={onMonthChange}
				class="bg-card border border-border rounded-lg px-3 py-2 text-sm text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
			/>
		</div>
	</div>

	<!-- Summary KPIs -->
	<div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
		<!-- Gross Revenue -->
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-blue-100 dark:bg-blue-900/30 rounded-lg">
					<DollarSign size={24} class="text-blue-600 dark:text-blue-400" />
				</div>
			</div>
			<div class="text-3xl font-bold text-foreground mb-1">
				{formatCents(data.summary.grossRevenueInCents)}
			</div>
			<p class="text-sm text-muted-foreground">Chiffre d'affaires brut</p>
		</div>

		<!-- Refunds -->
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-red-100 dark:bg-red-900/30 rounded-lg">
					<ArrowDown size={24} class="text-red-600 dark:text-red-400" />
				</div>
			</div>
			<div class="text-3xl font-bold text-red-500 mb-1">
				-{formatCents(data.summary.refundsInCents)}
			</div>
			<p class="text-sm text-muted-foreground">Remboursements</p>
		</div>

		<!-- Net Revenue -->
		<div class="bg-card rounded-lg border border-border p-6">
			<div class="flex items-center justify-between mb-4">
				<div class="p-3 bg-green-100 dark:bg-green-900/30 rounded-lg">
					<TrendingUp size={24} class="text-green-600 dark:text-green-400" />
				</div>
			</div>
			<div class="text-3xl font-bold text-green-500 mb-1">
				{formatCents(data.summary.netRevenueInCents)}
			</div>
			<p class="text-sm text-muted-foreground">Chiffre d'affaires net</p>
		</div>
	</div>

	<!-- Transactions Table -->
	<div class="bg-card rounded-lg border border-border mb-8">
		<div class="p-6 border-b border-border">
			<h2 class="text-xl font-semibold text-foreground">Transactions</h2>
		</div>
		{#if data.transactions.length === 0}
			<div class="p-12 text-center">
				<CreditCard size={48} class="mx-auto mb-4 text-muted-foreground/50" />
				<p class="text-muted-foreground">Aucune transaction pour cette période.</p>
			</div>
		{:else}
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead>
						<tr class="bg-muted/50 border-b border-border">
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Date</th>
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Produit</th>
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Client</th>
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Praticien</th>
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Type</th>
							<th class="text-right py-4 px-6 text-sm font-medium text-muted-foreground">Montant</th>
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Paiement</th>
							<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Réservation</th>
						</tr>
					</thead>
					<tbody>
						{#each data.transactions as transaction}
							<tr class="border-b border-border hover:bg-muted/30 transition-colors">
								<td class="py-4 px-6 text-sm text-muted-foreground">{formatDateTime(transaction.date)}</td>
								<td class="py-4 px-6">
									<span class="text-sm font-medium text-foreground">{transaction.description}</span>
								</td>
								<td class="py-4 px-6 text-sm text-foreground">{transaction.clientName}</td>
								<td class="py-4 px-6 text-sm text-foreground">{transaction.partnerName}</td>
								<td class="py-4 px-6">
									<span class="px-3 py-1 rounded-full text-xs font-medium {getTypeBadge(transaction.type)}">
										{getTypeLabel(transaction.type)}
									</span>
								</td>
								<td class="py-4 px-6 text-right">
									<span class="text-sm font-medium {transaction.type === 'payment'
										? 'text-green-500'
										: 'text-red-500'}">
										{transaction.type === 'refund' ? '-' : ''}{formatCents(transaction.amountInCents)}
									</span>
								</td>
								<td class="py-4 px-6">
									<span class="px-3 py-1 rounded-full text-xs font-medium {getPaymentStatusBadge(transaction.paymentStatus)}">
										{getPaymentStatusLabel(transaction.paymentStatus)}
									</span>
								</td>
								<td class="py-4 px-6">
									<span class="px-3 py-1 rounded-full text-xs font-medium {getBookingStatusBadge(transaction.bookingStatus)}">
										{getBookingStatusLabel(transaction.bookingStatus)}
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
