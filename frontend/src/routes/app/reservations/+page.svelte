<script lang="ts">
	import { ROLES } from '$lib/types/Navigation';

	import ReservationBasic from './ReservationBasic.svelte';
	// import ReservationPremium from './ReservationPremium.svelte';
	// import ReservationAdmin from './ReservationAdmin.svelte';

	interface Props {
		data: import('./$types').PageData;
	}

	let { data }: Props = $props();
	console.log('the data received is:', data);
	let { role, eventCards, consultationCards } = data;
	$inspect({ role, eventCards, consultationCards });
</script>

{#if role === ROLES.Basic}
	<ReservationBasic {eventCards} {consultationCards} />
	<!-- TODO: add the rest of the user after that -->
	<!-- {:else if role === ROLES.Admin} -->
	<!-- 	<ReservationAdmin {cards} /> -->
	<!-- {:else if role === ROLES.Premium} -->
	<!-- 	<ReservationPremium {cards} /> -->
{/if}
