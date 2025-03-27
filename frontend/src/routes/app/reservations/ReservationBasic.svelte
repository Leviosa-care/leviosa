<script lang="ts">
	import { Images } from 'lucide-svelte';
	import { redirectTo } from '$lib/scripts/redirect';

	import EventCard from '$lib/components/ui/EventCard.svelte';
	import ConsultationCard from '$lib/components/ui/ConsultationCard.svelte';
	import NoEvent from './NoEvent.svelte';
	import NoConsultation from './NoConsultation.svelte';
	import EventNavigationBar from '$lib/components/navigation/EventNavigationBar.svelte';
	import ConsultationNavigationBar from '$lib/components/navigation/ConsultationNavigationBar.svelte';
	import Tabs from '$lib/components/Tabs.svelte';

	import { EVENT_STATES, RESERVATION_STATES } from '$lib/types';
	import { eventState, reservationState } from '$lib/stores/persisted_stores.svelte';

	eventState.set(EVENT_STATES.EvenementsAVenir);
	reservationState.set(RESERVATION_STATES.Consultations);

	let { eventCards, consultationCards } = $props();

	// TODO: change that function, this is stupid
	function handleTab(): void {
		if ($reservationState === RESERVATION_STATES.Consultations)
			reservationState.set(RESERVATION_STATES.Events);
		else reservationState.set(RESERVATION_STATES.Consultations);
	}
	const offers = [
		{ name: 'Consultations', action: () => handleTab() },
		{ name: 'Evenements', action: () => handleTab() }
	];

	import { ROLES } from '$lib/types/Navigation';
	const role = ROLES.Basic;
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
	{#if $reservationState === RESERVATION_STATES.Consultations}
		<ConsultationNavigationBar {role} />
	{:else}
		<EventNavigationBar {role} />
	{/if}
</div>
{#if $reservationState === RESERVATION_STATES.Events}
	<div class="content grid px-3">
		<div class="grid" style="margin-top: 1rem;">
			{#if eventCards.length > 0}
				<EventCard />
			{:else}
				<NoEvent />
			{/if}
		</div>
	</div>
{:else}
	<div class="content grid" style="padding-inline: 0.75rem;">
		{#if eventCards.length > 0}
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
