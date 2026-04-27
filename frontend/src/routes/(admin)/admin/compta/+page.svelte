<script lang="ts">
	import type { PageProps } from './$types';
	import { DollarSign, ArrowDown, TrendingUp, CreditCard, Banknote, ArrowRightLeft } from '@lucide/svelte';

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

	function getPaymentMethodIcon(method: string) {
		switch (method) {
			case 'card':
				return CreditCard;
			case 'cash':
				return Banknote;
			case 'transfer':
				return ArrowRightLeft;
			default:
				return CreditCard;
		}
	}

	function getTypeBadge(type: string) {
		return type === 'payment'
			? 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300'
			: 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300';
	}

	function getTypeLabel(type: string): string {
		return type === 'payment' ? 'Paiement' : 'Remboursement';
	}

	function getStatusBadge(status: string) {
		switch (status) {
			case 'completed':
				return 'bg-green-100 text-green-700 dark:bg-green-900 dark:text-green-300';
			case 'pending':
				return 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900 dark:text-yellow-300';
			case 'failed':
				return 'bg-red-100 text-red-700 dark:bg-red-900 dark:text-red-300';
			default:
				return 'bg-gray-100 text-gray-700 dark:bg-gray-800 dark:text-gray-300';
		}
	}

	function getStatusLabel(status: string): string {
		switch (status) {
			case 'completed': return 'Complété';
			case 'pending': return 'En attente';
			case 'failed': return 'Échoué';
			default: return status;
		}
	}

	function getPaymentMethodLabel(method: string): string {
		switch (method) {
			case 'card': return 'Carte';
			case 'cash': return 'Espèces';
			case 'transfer': return 'Virement';
			default: return method;
		}
	}
</script>

<svelte:head>
	<title>Comptabilité | Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1 text-foreground">
			Comptabilité
		</h1>
		<p class="text-muted-foreground">
			Suivi des revenus et des transactions
		</p>
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
			<h2 class="text-xl font-semibold text-foreground">Transactions récentes</h2>
		</div>
		<div class="overflow-x-auto">
			<table class="w-full">
				<thead>
					<tr class="bg-muted/50 border-b border-border">
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Date</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Description</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Client</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Méthode</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Type</th>
						<th class="text-right py-4 px-6 text-sm font-medium text-muted-foreground">Montant</th>
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Statut</th>
					</tr>
				</thead>
				<tbody>
					{#each data.transactions as transaction}
						{@const MethodIcon = getPaymentMethodIcon(transaction.paymentMethod)}
						<tr class="border-b border-border hover:bg-muted/30 transition-colors">
							<td class="py-4 px-6 text-sm text-muted-foreground">{formatDateTime(transaction.date)}</td>
							<td class="py-4 px-6">
								<span class="text-sm font-medium text-foreground">{transaction.description}</span>
							</td>
							<td class="py-4 px-6 text-sm text-foreground">{transaction.clientName}</td>
							<td class="py-4 px-6">
								<div class="flex items-center gap-2">
									<MethodIcon size={14} class="text-muted-foreground" />
									<span class="text-sm text-muted-foreground">{getPaymentMethodLabel(transaction.paymentMethod)}</span>
								</div>
							</td>
							<td class="py-4 px-6">
								<span class="px-3 py-1 rounded-full text-xs font-medium {getTypeBadge(transaction.type)}">
									{getTypeLabel(transaction.type)}
								</span>
							</td>
							<td class="py-4 px-6 text-right">
								<span class="text-sm font-medium {transaction.amountInCents >= 0
									? 'text-green-500'
									: 'text-red-500'}">
									{formatCents(transaction.amountInCents)}
								</span>
							</td>
							<td class="py-4 px-6">
								<span class="px-3 py-1 rounded-full text-xs font-medium {getStatusBadge(transaction.status)}">
									{getStatusLabel(transaction.status)}
								</span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>

	<!-- Monthly Breakdown -->
	<div class="bg-card rounded-lg border border-border">
		<div class="p-6 border-b border-border">
			<h2 class="text-xl font-semibold text-foreground">Répartition mensuelle</h2>
		</div>
		<div class="overflow-x-auto">
			<table class="w-full">
				<thead>
					<tr class="bg-muted/50 border-b border-border">
						<th class="text-left py-4 px-6 text-sm font-medium text-muted-foreground">Mois</th>
						<th class="text-right py-4 px-6 text-sm font-medium text-muted-foreground">Paiements</th>
						<th class="text-right py-4 px-6 text-sm font-medium text-muted-foreground">Remboursements</th>
						<th class="text-right py-4 px-6 text-sm font-medium text-muted-foreground">Net</th>
					</tr>
				</thead>
				<tbody>
					{#each data.monthlyBreakdown as month}
						<tr class="border-b border-border hover:bg-muted/30 transition-colors">
							<td class="py-4 px-6">
								<span class="text-sm font-medium text-foreground">{month.month}</span>
							</td>
							<td class="py-4 px-6 text-right">
								<span class="text-sm font-medium text-green-500">{formatCents(month.payments)}</span>
							</td>
							<td class="py-4 px-6 text-right">
								<span class="text-sm font-medium text-red-500">{formatCents(month.refunds)}</span>
							</td>
							<td class="py-4 px-6 text-right">
								<span class="text-sm font-bold text-foreground">{formatCents(month.net)}</span>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
