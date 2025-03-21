<script lang="ts">
	let { cards } = $props();

	import { eventstate } from '$lib/stores/eventbar';
	eventstate.set('Evenements a venir');

	import { Images } from 'lucide-svelte';
	import { redirectTo } from '$lib/scripts/redirect';
	// TODO: is this the thing to do ?
	import { reservationstate } from '$lib/stores/reservationtab';
	reservationstate.set('consultations');

	import EventCard from '$lib/components/ui/EventCard.svelte';
	import ConsultationCard from '$lib/components/ui/ConsultationCard.svelte';
	import NoEvent from './NoEvent.svelte';
	import NoConsultation from './NoConsultation.svelte';
	import EventNavigationBar from '$lib/components/navigation/EventNavigationBar.svelte';
	import ConsultationNavigationBar from '$lib/components/navigation/ConsultationNavigationBar.svelte';
	import Tabs from '$lib/components/Tabs.svelte';

	function handleTab(): void {
		if ($reservationstate === 'consultations') reservationstate.set('events');
		else reservationstate.set('consultations');
	}
	const offers = [
		{ name: 'Consultations', action: () => handleTab() },
		{ name: 'Evenements', action: () => handleTab() }
	];

	const role = 'user';
</script>

<div class="event-header grid">
	<div class="event-header-top container flex">
		<h2 class="page-title">Reservations</h2>
		<button onclick={() => redirectTo('galerie')}>
			<Images strokeWidth={1.5} absoluteStrokeWidth={true} />
		</button>
	</div>
	<div class="container">
		<Tabs {offers} isSecondary={true} />
	</div>
	{#if $reservationstate === 'consultations'}
		<ConsultationNavigationBar {role} />
	{:else}
		<EventNavigationBar {role} />
	{/if}
</div>
{#if $reservationstate === 'events'}
	<div class="content grid" style="padding-inline: 0.75rem;">
		<div class="grid" style="margin-top: 1rem;">
			{#if cards.length > 0}
				<EventCard />
			{:else}
				<NoEvent />
			{/if}
		</div>
	</div>
{:else}
	<div class="content grid" style="padding-inline: 0.75rem;">
		{#if cards.length > 0}
			<ConsultationCard id="" />
		{:else}
			<NoConsultation />
		{/if}
	</div>
{/if}

<style>
	.event-header {
		background-color: hsl(var(--clr-light-primary));
		padding-top: 2rem;
		flex-shrink: 0;
		flex: none;
	}
	.event-header-top {
		justify-content: space-between;
		align-items: center;
	}
	.content {
		padding-top: 1rem;
		padding-bottom: 7rem;
	}
</style>
