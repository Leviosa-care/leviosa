<script lang="ts">
	// =======================
	// Imports
	// =======================

	// Types
	import type { ServiceState } from '$lib/types';

	// Stores
	import { serviceState } from '$lib/stores/persisted_stores.svelte';

	// =======================
	// Props and Types

	type BarType = 'light' | 'dark';
	interface Props {
		// =======================
		active?: boolean;
		name?: ServiceState;
		type?: BarType;
	}

	let { active = $bindable(false), name = 'A propos', type = 'dark' }: Props = $props();

	// =======================
	// Functions
	// =======================

	function setState(event: MouseEvent): void {
		const targetElement = event.currentTarget as HTMLButtonElement;
		const id = targetElement.id as ServiceState;

		if (!active) active = true;
		serviceState.set(id);
	}
</script>

<button id={name} class:active class={type} onclick={setState}>
	<p class="name">{name}</p>
</button>

<style>
	button {
		border-bottom: 2px solid transparent;
		background: transparent;
		transition: border-color 0.3s ease;
		color: hsl(var(--clr-grey-200));
	}
	button.dark {
		color: hsl(var(--clr-dark-ternary));
	}
	button.light {
		color: hsl(var(--clr-light-secondary));
	}
	.name {
		font-size: 0.9rem;
		font-size: var(--fs-0);
		padding-block: 0.5rem;
	}
	button.active {
		font-weight: 600;
		color: hsl(var(--clr-grey-100));
	}
	button.active.dark {
		border-bottom-color: hsl(var(--clr-dark-primary));
		color: hsl(var(--clr-dark-primary));
	}
	button.active.light {
		border-bottom-color: hsl(var(--clr-grey-100));
		color: hsl(var(--clr-grey-100));
	}
</style>
