<script lang="ts">
	// global components
	import Button from '$lib/components/Button.svelte';
	import Figure from './Figure.svelte';
	interface Props {
		// TODO: change the vars that you need for that page
		description: string;
		duration: number;
		clients_count: number;
		positive_responses: number;
	}

	let { description, duration, clients_count, positive_responses }: Props = $props();
</script>

<div class="flow content">
	<div class="container">
		<div class="card grid" style="--gap: 1rem;">
			<h3 class="title">Description</h3>
			<p class="description">{description}</p>
			<div class="card-overlay"></div>
		</div>
	</div>
	<div class="figures container flex" style="--gap: 0rem;">
		<Figure figure={duration} label="minutes" />
		<Figure figure={clients_count} label="satisfait.e.s" />
		<Figure figure="{positive_responses}%" label="avis positifs" />
	</div>
	<div class="container grid">
		<Button onClick={() => console.log('make the reservation brother')}>Reserve ta place</Button>
	</div>
</div>

<style>
	.figures {
		justify-content: space-between;
		align-items: center;
		max-width: 500px;
	}
	.title {
		color: hsl(var(--clr-grey-100));
		font-size: var(--fs-1);
	}
	.card {
		--card-border-radius: 1rem;
		--card-padding: 1rem;
		padding: var(--card-padding);
		border-radius: var(--card-border-radius);
		backdrop-filter: blur(15px) saturate(150%);
		color: hsl(var(--clr-light-ternary));
		box-shadow:
			rgba(60, 64, 67, 0.3) 0px 1px 2px 0px,
			rgba(60, 64, 67, 0.15) 0px 1px 3px 1px;
		position: relative;
		overflow: hidden;
	}
	.card-overlay {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.25);
		z-index: -1;
	}
</style>
