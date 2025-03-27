<script lang="ts">
	import type { ConsultationState } from '$lib/types';
	import { consultationState } from '$lib/stores/persisted_stores.svelte';
	interface Props {
		active?: boolean;
		href: string;
		name?: ConsultationState;
	}

	let { active = $bindable(false), href, name = 'Consultations a venir' }: Props = $props();

	function setState(event: MouseEvent) {
		let targetElement = event.currentTarget as HTMLButtonElement;
		let id = targetElement.id as ConsultationState;
		if (!active) active = !active;
		consultationState.set(id);
	}
</script>

<button id={name} class:active onclick={setState}>
	<a class="flex" {href}>
		<p class="name">{name}</p>
	</a>
</button>

<style>
	button {
		border-bottom: 3px solid transparent;
		background: transparent;
		transition: border-color 0.3s ease;
		color: hsl(var(--clr-grey-400));
	}
	a {
		color: inherit;
		text-decoration: none;
	}
	.name {
		font-size: var(--fs--1);
		font-weight: 500;
		padding-block: 0.5rem;
	}
	button.active {
		border-bottom-color: hsl(var(--clr-grey-700));
		color: hsl(var(--clr-grey-700));
	}
</style>
