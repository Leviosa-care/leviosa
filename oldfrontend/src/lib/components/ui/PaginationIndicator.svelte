<script lang="ts">
	interface Props {
		count?: number;
		selectedElement?: number;
	}

	let { count = 1, selectedElement = 0 }: Props = $props();
</script>

<div class="component flex" style="--gap: 0.5rem;">
	{#each Array(count) as _, index}
		<div class:active={index <= selectedElement} class="bar"></div>
	{/each}
</div>

<style>
	.component {
		align-items: center;
	}
	.bar {
		height: 6px;
		width: 25px;
		border-radius: 0.2rem;
		background-color: hsl(var(--clr-light-ternary));
	}
	.bar.active {
		background-color: hsl(var(--clr-dark-primary));
	}
</style>
