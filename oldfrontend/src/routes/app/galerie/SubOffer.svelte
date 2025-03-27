<script lang="ts">
	type Offer = {
		name: string;
		action: () => void;
	};

	let action = 0;

	let selectIndex: number = $state(0);
	function buttonAction(index: number) {
		selectIndex = index;
		offers[index].action();
	}
	interface Props {
		offers?: Offer[];
		isSecondary?: boolean;
	}

	let {
		offers = [
			{ name: 'Photos', action: () => (action = 1) },
			{ name: 'Albums', action: () => (action = 2) }
		],
		isSecondary = false
	}: Props = $props();
</script>

<div class="offers grid" style="--gap: 0rem;">
	{#each offers as offer, index}
		<button
			class="offer"
			class:selected={index === selectIndex}
			onclick={() => buttonAction(index)}
			class:secondary={isSecondary}
		>
			{offer.name}
		</button>
	{/each}
</div>

<style>
	.offers {
		--border: 0.3rem;
		background-color: hsl(var(--clr-dark-primary));
		border-radius: calc(3 * var(--border) / 2);
		padding: 0.2rem;

		grid-auto-flow: column;
		grid-auto-columns: minmax(max-content, 1fr);
	}
	.offer {
		--padding: 0.2rem;
		font-size: 0.8rem;
		border-radius: var(--border);
		padding: calc(var(--padding) * 2) calc(var(--padding) * 4);
		background: transparent;
		font-weight: 600;
	}
	.offer:not(.selected) {
		color: hsl(var(--clr-dark-ternary));
	}
	.selected {
		background-color: hsl(var(--clr-light-primary));
		color: hsl(var(--clr-dark-primary));
		box-shadow:
			rgba(0, 0, 0, 0.05) 0px 6px 24px 0px,
			rgba(0, 0, 0, 0.08) 0px 0px 0px 1px;
	}
	.offers:has(:global(.secondary)) {
		background-color: hsl(var(--clr-light-ternary));
		background-color: hsl(var(--clr-light-secondary));
	}
	.secondary {
		color: hsl(var(--clr-dark-primary));
	}
	.secondary.selected {
		background-color: hsl(var(--clr-light-primary));
	}
</style>
