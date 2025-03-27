<script lang="ts">
	import { Calendar } from 'lucide-svelte';
	import { redirectTo } from '$lib/scripts/redirect';
</script>

<div class="card grid" style="--gap: 2rem;">
	<div class="icon stroke">
		<Calendar />
	</div>
	<p class="description center">Tu n'as pas encore de reservations prevues</p>
	<button onclick={() => redirectTo('/app/events')} class="cta">Decouvre nos evenements</button>
</div>

<style>
	.card {
		padding: 1rem;
		border-radius: 0.5rem;
		min-width: 100%;

		background-color: hsl(var(--clr-light-primary));
		box-shadow: rgba(0, 0, 0, 0.05) 0px 0px 0px 1px;
	}
	.description {
		width: 85%;
		margin-inline: auto;
	}
	.icon {
		margin-inline: auto;
		width: fit-content;
		padding: 0.5rem;
		border-radius: 0.5rem;
		background-color: hsl(var(--clr-light-secondary));
	}
	.cta {
		background-color: hsl(var(--clr-dark-primary));
		color: hsl(var(--clr-light-primary));
		border-radius: 0.5rem;
		padding: 0.5rem 1rem;
		text-wrap: pretty;
	}
</style>
