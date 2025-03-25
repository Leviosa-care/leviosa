<script lang="ts">
	import ServiceNavigationBarIcon from './atoms/ServiceNavigationBarIcon.svelte';
	import { SERVICE_STATES_ARRAY } from '$lib/types';
	import { serviceState } from '$lib/stores/persisted_stores.svelte';

	type BarType = 'light' | 'dark';
	interface Props {
		type?: BarType;
	}

	let { type = 'dark' }: Props = $props();
</script>

<nav class="navigation-bar flex">
	{#each SERVICE_STATES_ARRAY as name}
		<ServiceNavigationBarIcon {name} active={$serviceState === name} {type} />
	{/each}
</nav>

<style>
	.navigation-bar {
		justify-content: space-between;
		align-items: center;
		border-bottom: 1px solid hsl(var(--clr-stroke));
		text-wrap: nowrap;
	}
</style>
