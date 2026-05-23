<script lang="ts">
	import {
		Calendar,
		TrendingUp,
		Clock,
		Package,
		ArrowRight,
		Users,
		BarChart3,
		Plus,
		CheckCircle2,
		Circle
	} from '@lucide/svelte';
	import type { PageData } from './$types';
	import type { RecentBookingUI } from './+page.server';

	let { data }: { data: PageData } = $props();

	function formatCurrency(cents: number, currency = 'EUR'): string {
		return new Intl.NumberFormat('fr-FR', { style: 'currency', currency }).format(cents / 100);
	}

	function formatDate(dateStr: string): string {
		return new Date(dateStr).toLocaleDateString('fr-FR', {
			day: 'numeric',
			month: 'short',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function formatTime(dateStr: string): string {
		return new Date(dateStr).toLocaleTimeString('fr-FR', {
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	function timeAgo(dateStr: string): string {
		const diff = Date.now() - new Date(dateStr).getTime();
		const minutes = Math.floor(diff / 60000);
		const hours = Math.floor(diff / 3600000);
		const days = Math.floor(diff / 86400000);
		if (minutes < 1) return "à l'instant";
		if (minutes < 60) return `il y a ${minutes}m`;
		if (hours < 24) return `il y a ${hours}h`;
		return `il y a ${days}j`;
	}

	function getStatusBadge(status: RecentBookingUI['status']) {
		switch (status) {
			case 'confirmed':
				return { label: 'Confirmé', class: 'text-green-700 bg-green-50', icon: CheckCircle2 };
			case 'completed':
				return { label: 'Terminé', class: 'text-blue-700 bg-blue-50', icon: CheckCircle2 };
			case 'pending':
				return { label: 'En attente', class: 'text-yellow-700 bg-yellow-50', icon: Clock };
			case 'cancelled':
				return { label: 'Annulé', class: 'text-red-700 bg-red-50', icon: Circle };
			case 'no_show':
				return { label: 'Absent', class: 'text-gray-700 bg-gray-100', icon: Circle };
		}
	}
</script>

<svelte:head>
	<title>Tableau de bord | Leviosa Admin</title>
</svelte:head>

<div class="container mx-auto px-4 py-8 lg:py-12">
	<div class="mb-8">
		<h1 class="text-3xl lg:text-4xl font-bold mb-1">Tableau de bord</h1>
		<p class="text-muted-foreground">Activité récente et aperçu de votre activité</p>
	</div>

	<!-- Stats -->
	<div class="grid grid-cols-2 lg:grid-cols-5 gap-4 mb-8">
		<div class="bg-background rounded-lg border border-border-card p-5">
			<div class="flex items-center gap-2 mb-3">
				<TrendingUp size={16} class="text-muted-foreground" />
				<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
					>Revenus (7j)</span
				>
			</div>
			<div class="text-2xl font-bold text-foreground">{formatCurrency(data.stats.revenueThisWeek)}</div>
		</div>

		<div class="bg-background rounded-lg border border-border-card p-5">
			<div class="flex items-center gap-2 mb-3">
				<Calendar size={16} class="text-muted-foreground" />
				<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
					>Réservations (7j)</span
				>
			</div>
			<div class="text-2xl font-bold text-foreground">{data.stats.bookingsThisWeek}</div>
		</div>

		<div class="bg-background rounded-lg border border-border-card p-5">
			<div class="flex items-center gap-2 mb-3">
				<Users size={16} class="text-muted-foreground" />
				<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
					>À venir</span
				>
			</div>
			<div class="text-2xl font-bold text-foreground">{data.stats.upcomingBookingsCount}</div>
		</div>

		<div class="bg-background rounded-lg border border-border-card p-5">
			<div class="flex items-center gap-2 mb-3">
				<Clock size={16} class="text-muted-foreground" />
				<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
					>En attente</span
				>
			</div>
			<div class="text-2xl font-bold text-foreground">{data.stats.pendingBookingsCount}</div>
		</div>

		<div class="bg-background rounded-lg border border-border-card p-5">
			<div class="flex items-center gap-2 mb-3">
				<Package size={16} class="text-muted-foreground" />
				<span class="text-xs font-medium text-muted-foreground uppercase tracking-wider"
					>Produits actifs</span
				>
			</div>
			<div class="text-2xl font-bold text-foreground">{data.stats.activeProductsCount}</div>
		</div>
	</div>

	<div class="grid grid-cols-1 lg:grid-cols-5 gap-6 mb-8">
		<!-- Recent Bookings -->
		<div class="lg:col-span-3 bg-background rounded-lg border border-border-card">
			<div class="flex items-center justify-between px-6 py-4 border-b border-border-card">
				<div>
					<h2 class="text-base font-semibold text-foreground">Dernières réservations</h2>
					<p class="text-xs text-muted-foreground mt-0.5">Réservations récentes</p>
				</div>
				<a
					href="/admin/bookings"
					class="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors font-medium"
				>
					Voir tout
					<ArrowRight size={12} />
				</a>
			</div>

			{#if data.recentBookings.length === 0}
				<div class="flex items-center justify-center h-48 text-muted-foreground">
					<div class="text-center">
						<Users size={36} class="mx-auto mb-3 opacity-25" />
						<p class="text-sm">Aucune réservation</p>
					</div>
				</div>
			{:else}
				<div class="divide-y divide-border-card">
					{#each data.recentBookings as booking}
						{@const statusBadge = getStatusBadge(booking.status)}
						{@const StatusIcon = statusBadge.icon}
						<div class="px-6 py-4 hover:bg-muted/40 transition-colors">
							<div class="flex items-start justify-between gap-4">
								<div class="min-w-0 flex-1">
									<div class="flex items-center gap-2 mb-1">
										<span class="font-medium text-foreground text-sm">{booking.clientName}</span>
										<span
											class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-medium {statusBadge.class}"
										>
											<StatusIcon size={10} />
											{statusBadge.label}
										</span>
									</div>
									<div class="text-xs text-muted-foreground">
										{booking.productName} · {booking.therapistName}
									</div>
								</div>
								<div class="flex-shrink-0 text-right">
									<div class="text-sm text-muted-foreground">{timeAgo(booking.startTime)}</div>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>

		<!-- Upcoming Bookings -->
		<div class="lg:col-span-2 bg-background rounded-lg border border-border-card">
			<div class="flex items-center justify-between px-6 py-4 border-b border-border-card">
				<div>
					<h2 class="text-base font-semibold text-foreground">Prochaines séances</h2>
					<p class="text-xs text-muted-foreground mt-0.5">Réservations à venir</p>
				</div>
				<a
					href="/admin/bookings"
					class="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors font-medium"
				>
					Voir tout
					<ArrowRight size={12} />
				</a>
			</div>

			{#if data.upcomingBookings.length === 0}
				<div class="flex items-center justify-center h-48 text-muted-foreground">
					<div class="text-center">
						<Calendar size={36} class="mx-auto mb-3 opacity-25" />
						<p class="text-sm">Aucune séance à venir</p>
					</div>
				</div>
			{:else}
				<div class="divide-y divide-border-card">
					{#each data.upcomingBookings as booking}
						<div class="px-6 py-4">
							<div class="flex items-center gap-3 mb-2">
								<div class="flex-1 min-w-0">
									<div class="font-medium text-foreground text-sm truncate">
										{booking.clientName}
									</div>
									<div class="text-xs text-muted-foreground truncate">{booking.productName}</div>
								</div>
								<div class="text-right">
									<div class="text-sm font-semibold text-foreground">
										{formatTime(booking.startTime)}
									</div>
									<div class="text-xs text-muted-foreground">{booking.duration} min</div>
								</div>
							</div>
							<div class="flex items-center gap-1 text-xs text-muted-foreground">
								<span class="inline-block w-2 h-2 rounded-full bg-dark-300"></span>
								{booking.roomName}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>

	<!-- Quick Actions -->
	<div class="flex flex-wrap gap-3">
		<a
			href="/admin/products"
			class="inline-flex items-center gap-2 px-4 py-2 bg-background border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors text-sm font-medium"
		>
			<Package size={15} />
			<span class="whitespace-nowrap">Gérer les produits</span>
		</a>
		<a
			href="/admin/bookings"
			class="inline-flex items-center gap-2 px-4 py-2 bg-background border border-border-input text-foreground-alt rounded-lg hover:bg-muted transition-colors text-sm font-medium"
		>
			<Calendar size={15} />
			<span class="whitespace-nowrap">Voir les réservations</span>
		</a>
	</div>
</div>
