<script lang="ts">
	import { redirectTo } from '$lib/scripts/redirect';
	import { Calendar, Clock } from 'lucide-svelte';
	interface Props {
		id: string;
	}

	let { id }: Props = $props();
</script>

<button onclick={() => redirectTo(`/app/events/${id}`)} class="card grid" style="--gap: 0.5rem;">
	<div class="grid">
		<div>
			<p class="title">Consultation Leviosa massage</p>
			<div class="label flex" style="--gap:0.5rem;">
				<Clock size={16} />
				<p>Prevue dans moins d'un mois</p>
			</div>
		</div>
		<div>
			<p class="label">Budget</p>
			<p class="price">$5000-7000</p>
		</div>
	</div>
	<div class="separator"></div>
	<div class="grid" style="--gap: 1rem;">
		<div class="row flex">
			<div class="flex" style="--gap:0.5rem;">
				<img class="avatar" src="https://pbs.twimg.com/media/FA9tIzxUUAUy80c.jpg" alt="avatar" />
				<div>
					<p>Livio</p>
					<p>HENRY</p>
				</div>
			</div>
			<div class="row-right flex" style="--gap:0.5rem;">
				<div class="calendar">
					<Calendar size={20} strokeWidth={1.5} absoluteStrokeWidth={true} />
				</div>
				<div>
					<p>Oct 15, 2024</p>
					<p class="label">Deadline</p>
				</div>
			</div>
		</div>
		<div class="flex" style="--gap: 0.5rem;">
			<div class="tag">Massage</div>
			<div class="tag">Bien etre</div>
			<div class="tag">+2</div>
		</div>
	</div>
</button>

<style>
	.card {
		border-radius: 1rem;
		overflow: hidden;
		background: transparent;
		padding: 0.75rem;
		text-align: left;
		font-size: var(--fs--1);
		box-shadow:
			rgba(0, 0, 0, 0.05) 0px 6px 24px 0px,
			rgba(0, 0, 0, 0.08) 0px 0px 0px 1px;
		background-color: hsl(var(--clr-light-primary));
	}
	.title,
	.price {
		font-weight: 500;
	}
	.title {
		color: hsl(var(--clr-grey-700));
		font-size: var(--fs-0);
	}
	.price {
		color: hsl(var(--clr-grey-700));
	}
	.label {
		color: hsl(var(--clr-grey-400));
		align-items: center;
	}
	.row {
		align-items: center;
		justify-content: space-between;
	}
	.avatar {
		--dimension: 2.5rem;
		width: var(--dimension);
		height: var(--dimension);
		background-color: hsl(var(--clr-dark-primary));
		border-radius: 100%;
	}
	.calendar,
	.tag {
		background-color: hsl(var(--clr-light-secondary));
		padding: 0.5rem;
		border-radius: 100%;
		/* TODO: change the color if necessary */
	}
	.row-right {
		align-items: center;
	}
	.calendar {
		height: fit-content;
		display: grid;
		place-content: center;
		flex-grow: 0;
		color: hsl(var(--clr-grey-400));
	}
	.tag {
		padding: 0.4rem 0.75rem;
		color: hsl(var(--clr-grey-400));
		border-radius: 2rem;
	}
</style>
