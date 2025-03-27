<script lang="ts">
	import { Flower } from 'lucide-svelte';

	interface Props {
		serviceName?: string; // export let icon: string = "Massage";
	}

	let { serviceName = 'Massage' }: Props = $props();

	// NOTE: tuto : https://www.youtube.com/watch?v=3yfswsnD2sw
</script>

<div class="services snaps-inline grid" style="--gap: 1rem;">
	{#each Array(6) as _, index}
		<button class="service grid" style="--gap: 0.25rem;">
			<div class="illustration grid illustration_{index}">
				<Flower
					size={40}
					color="hsl(var(--clr-grey-100))"
					strokeWidth={2}
					absoluteStrokeWidth={true}
				/>
			</div>
			<p class="name">{serviceName}</p>
		</button>
	{/each}
</div>

<style>
	.services {
		grid-auto-flow: column;
		grid-auto-columns: 22%;

		overflow-x: auto;
		scrollbar-width: none;
		overscroll-behavior-inline: contain;
	}
	.snaps-inline {
		scroll-snap-type: inline mandatory;
		scroll-padding-inline: 0.75rem;
	}
	.snaps-inline > * {
		scroll-snap-align: start;
	}
	.service {
		background: transparent;
	}
	.service {
		background: transparent;
	}
	.illustration {
		border-radius: 0.5rem;
		padding: 0.5rem;
		background-color: hsl(var(--clr-grey-600));
		place-content: center;
	}
	.name {
		color: hsl(var(--clr-grey-600));
		font-size: var(--fs--1);
	}
	.illustration_1 {
		background-color: rebeccapurple;
	}
	.illustration_2 {
		background-color: orangered;
	}
	.illustration_3 {
		background-color: hotpink;
	}
	.illustration_4 {
		background-color: hsl(var(--clr-accent));
	}
</style>
