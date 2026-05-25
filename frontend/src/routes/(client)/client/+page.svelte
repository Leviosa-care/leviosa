<script lang="ts">
	import type { PageProps } from './$types';
	import { CalendarDays, Clock, ArrowRight, PlusCircle, CheckCircle } from '@lucide/svelte';

	let { data }: PageProps = $props();

	function formatDate(iso: string): string {
		return new Date(iso).toLocaleDateString('fr-FR', {
			weekday: 'long',
			day: 'numeric',
			month: 'long',
			year: 'numeric',
		});
	}

	function formatTime(iso: string): string {
		return new Date(iso).toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' });
	}

	function formatCents(cents: number): string {
		return new Intl.NumberFormat('fr-FR', { style: 'currency', currency: 'EUR' }).format(cents / 100);
	}

	function statusBadge(status: string): string {
		switch (status) {
			case 'confirmed': return 'bg-blue-100 text-blue-700';
			case 'completed': return 'bg-green-100 text-green-700';
			case 'cancelled': return 'bg-red-100 text-red-700';
			case 'no_show': return 'bg-orange-100 text-orange-700';
			default: return 'bg-gray-100 text-gray-700';
		}
	}

	function statusLabel(status: string): string {
		switch (status) {
			case 'confirmed': return 'À venir';
			case 'completed': return 'Terminé';
			case 'cancelled': return 'Annulé';
			case 'no_show': return 'Absent';
			default: return status;
		}
	}
</script>

<svelte:head>
	<title>Mon espace | Leviosa</title>
</svelte:head>

<div class="space-y-8">
	<!-- Welcome -->
	<div>
		<h1 class="text-2xl lg:text-3xl font-bold text-foreground">Bienvenue</h1>
		<p class="text-muted-foreground mt-1">Votre espace personnel</p>
	</div>

	<!-- Prochaine séance -->
	<div>
		<h2 class="text-lg font-semibold text-foreground mb-3 flex items-center gap-2">
			<CalendarDays size={20} />
			Prochaine séance
		</h2>

		{#if data.nextBooking}
			<div class="bg-card rounded-lg border border-border-card p-5 shadow-card">
				<div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
					<div>
						<div class="flex items-center gap-2 mb-1">
							<span class="px-2.5 py-0.5 rounded-full text-xs font-medium {statusBadge(data.nextBooking.status)}">
								{statusLabel(data.nextBooking.status)}
							</span>
						</div>
						<p class="text-foreground font-medium">{formatDate(data.nextBooking.slot_start_time)}</p>
						<div class="flex items-center gap-1.5 text-sm text-muted-foreground mt-1">
							<Clock size={14} />
							{formatTime(data.nextBooking.slot_start_time)} – {formatTime(data.nextBooking.slot_end_time)}
						</div>
						{#if data.nextBooking.total_price_cents}
							<p class="text-sm text-muted-foreground mt-1">
								{formatCents(data.nextBooking.total_price_cents)}
							</p>
						{/if}
					</div>
					<a
						href="/client/bookings/{data.nextBooking.id}"
						class="inline-flex items-center gap-1.5 text-sm font-medium text-foreground hover:underline"
					>
						Voir le détail <ArrowRight size={14} />
					</a>
				</div>
			</div>
		{:else}
			<div class="bg-card rounded-lg border border-border-card p-8 text-center shadow-card">
				<p class="text-muted-foreground mb-4">Aucune séance à venir</p>
				<a
					href="/book"
					class="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg bg-foreground text-background text-sm font-medium hover:bg-foreground/90 transition-colors"
				>
					<PlusCircle size={16} />
					Réserver une séance
				</a>
			</div>
		{/if}
	</div>

	<!-- CTA: Réserver -->
	{#if data.nextBooking}
		<a
			href="/book"
			class="flex items-center justify-center gap-2 w-full sm:w-auto px-6 py-3 rounded-lg bg-foreground text-background font-medium hover:bg-foreground/90 transition-colors"
		>
			<PlusCircle size={18} />
			Réserver une séance
		</a>
	{/if}

	<!-- Mes réservations récentes -->
	<div>
		<div class="flex items-center justify-between mb-3">
			<h2 class="text-lg font-semibold text-foreground flex items-center gap-2">
				<CheckCircle size={20} />
				Réservations récentes
			</h2>
			<a href="/client/bookings" class="text-sm font-medium text-muted-foreground hover:text-foreground flex items-center gap-1">
				Voir tout <ArrowRight size={14} />
			</a>
		</div>

		{#if data.recentCompleted.length > 0}
			<div class="space-y-3">
				{#each data.recentCompleted as booking (booking.id)}
					<div class="bg-card rounded-lg border border-border-card p-4 shadow-card flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
						<div>
							<p class="text-foreground font-medium text-sm">{formatDate(booking.slot_start_time)}</p>
							<p class="text-sm text-muted-foreground">{formatTime(booking.slot_start_time)} – {formatTime(booking.slot_end_time)}</p>
						</div>
						<div class="flex items-center gap-3">
							<span class="px-2.5 py-0.5 rounded-full text-xs font-medium {statusBadge(booking.status)}">
								{statusLabel(booking.status)}
							</span>
							{#if booking.total_price_cents}
								<span class="text-sm font-medium text-foreground">{formatCents(booking.total_price_cents)}</span>
							{/if}
						</div>
					</div>
				{/each}
			</div>
		{:else}
			<div class="bg-card rounded-lg border border-border-card p-6 text-center text-muted-foreground shadow-card">
				Aucune réservation terminée
			</div>
		{/if}
	</div>
</div>
