<script lang="ts">
	import { NAV_STATES } from '$lib/types';
	import { navigationState } from '$lib/stores/persisted_stores.svelte';
	navigationState.set(NAV_STATES.Reservations);

	interface Props {
		children?: import('svelte').Snippet;
		data: import('./$types').PageData;
	}

	let { children }: Props = $props();
</script>

<div class="page">
	<div class="slot">
		{@render children?.()}
	</div>
</div>

<style>
	.page {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}
	.slot {
		flex-grow: 1;
		overflow: auto;
		/* background-color: hsl(var(--clr-light-secondary)); */
		background-color: #f7f7f9;
	}
</style>
