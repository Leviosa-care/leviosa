<script lang="ts">
	// TODO: use the icon for massage and preparation mental in this
	import { Sparkles } from 'lucide-svelte';
	import { redirectTo } from '$lib/scripts/redirect';
	interface Props {
		id?: string;
		eventType?: string;
		title?: string;
		date?: string;
		eventBegin?: string;
		city?: string;
		prestationType?: string;
		prestationStart?: string;
	}

	let {
		id = 'awt9av4wetq24tioet34',
		eventType = 'Meetup',
		title = 'Journee speciale sante mentale',
		date = 'Vendredi 21 Juin',
		eventBegin = '10:00',
		city = 'Ivry-Sur-Seine, Paris',
		prestationType = 'massage',
		prestationStart = '11:00'
	}: Props = $props();
</script>

<button
	onclick={() => redirectTo(`/app/reservations/events/${id}`)}
	class="card grid"
	style="--gap: 0.5rem;"
	role="link"
	aria-label={`Aller a l'evenement avec l'ID ${id}`}
>
	<div class="relative">
		<img
			class="img"
			src="https://i.pinimg.com/originals/dc/36/fc/dc36fcd109235b0f4f591ff32bae1db9.jpg"
			alt="Lieu de de l'evenement"
		/>
		{#if prestationType !== ''}
			<div class="tag flex" style="--gap: 0.25rem;">
				<Sparkles size={16} color="hsl(var(--clr-light-primary))" />
				<p class="capitalize">
					{prestationType} :
					<span class="bold">{prestationStart}</span>
				</p>
			</div>
		{/if}
	</div>
	<div
		class="content grid"
		style="--gap: 1rem; padding-bottom: 1rem;"
		role="region"
		aria-label="Contenu de la carte"
	>
		<p class="title" aria-label="Titre de l'evenement">
			<span class="bold">Leviosa {eventType}</span>
			- {title}
		</p>
		<div class="grid" style="--gap: 0.2rem;">
			<p aria-label="Timestamp" class="date">{date} | {eventBegin}</p>
			<p aria-label="Localisation" class="location">{city}</p>
		</div>
	</div>
</button>

<style>
	.card {
		border-radius: 0.5rem;
		overflow: hidden;
		box-shadow:
			rgba(0, 0, 0, 0.1) 0px 1px 3px 0px,
			rgba(0, 0, 0, 0.06) 0px 1px 2px 0px;
		background-color: hsl(var(--clr-light-primary));
		box-shadow:
			rgba(0, 0, 0, 0.05) 0px 6px 24px 0px,
			rgba(0, 0, 0, 0.08) 0px 0px 0px 1px;
	}
	.img {
		inline-size: 100%;
		aspect-ratio: 16/9;
		object-fit: cover;
	}
	.tag {
		padding: 0.25rem 0.75rem;
		border-radius: 2rem;
		width: fit-content;
		background-color: hsl(var(--clr-accent));
		background-color: hsl(var(--clr-dark-primary));
		color: hsl(var(--clr-light-primary));
		align-items: center;
		position: absolute;
		--pos: 1rem;
		top: var(--pos);
		left: var(--pos);
	}
	.card .title {
		color: hsl(var(--clr-grey-700));
		font-size: var(--fs-1);
		font-weight: 500;
		line-height: 1.25;
	}
	.card .content {
		padding-inline: 0.4rem;
		text-align: left;
	}
	.bold {
		font-weight: 600;
	}
	.date {
		color: hsl(var(--clr-accent));
		font-weight: 500;
	}
	.location {
		color: hsl(var(--clr-grey-500));
	}
</style>
