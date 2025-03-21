<script lang="ts">
	import ReservationUser from './ReservationUser.svelte';
	import ReservationUserPremium from './ReservationUserPremium.svelte';
	import ReservationAdmin from './ReservationAdmin.svelte';

	interface Props {
		data: import('./$types').PageData;
	}

	let { data }: Props = $props();
	let { cards, role } = data;
</script>

{#if role === 'user'}
	<ReservationUser {cards} />
{:else if role === 'admin'}
	<ReservationAdmin {cards} />
{:else if role === 'userPremium'}
	<ReservationUserPremium {cards} />
{/if}
