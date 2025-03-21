<script lang="ts">
	import { navstate } from '$lib/stores/navbar';
	navstate.set('reservations'); // just to forget the value stored in localstore when reconecting and I had the page to another link.

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
