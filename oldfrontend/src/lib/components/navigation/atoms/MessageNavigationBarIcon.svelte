<script lang="ts">
	import type { MessageState } from '$lib/types';
	import { messageState } from '$lib/stores/persisted_stores.svelte';
	interface Props {
		active?: boolean;
		name?: MessageState;
	}

	let { active = $bindable(false), name = 'Conversations' }: Props = $props();

	function setState(event: MouseEvent) {
		let targetElement = event.currentTarget as HTMLButtonElement;
		let id = targetElement.id as MessageState;
		if (!active) active = !active;
		messageState.set(id);
	}
</script>

<button class="name container" id={name} class:active onclick={setState}>
	{name}
</button>

<style>
	button {
		border-bottom: 3px solid transparent;
		background: transparent;
		transition: border-color 0.3s ease;
		color: hsl(var(--clr-grey-400));
		width: 100%;
	}
	.name {
		font-size: var(--fs--1);
		font-weight: 500;
		padding-block: 0.5rem;
		margin-inline: auto;
	}
	button.active {
		border-bottom-color: hsl(var(--clr-grey-700));
		color: hsl(var(--clr-grey-700));
	}
</style>
