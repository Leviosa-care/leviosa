<script>
	import { Calendar } from 'lucide-svelte';

	import Button from '$lib/components/Button.svelte';
</script>

<div class="card container grid" style="--gap: 2rem;">
	<div class="grid">
		<div class="icon">
			<Calendar size={32} color="black" />
		</div>
		<div class="description center">
			<p class="">Vous n'avez aucune reservations</p>
			<p class="">Vos reservations futures s'afficheront ici</p>
		</div>
	</div>
</div>
<div class="btn grid">
	<Button onClick={() => console.log('go to make a reservation')}>Nouvelle reservation</Button>
</div>

<style>
	.card {
		padding: 2rem;
		place-content: center;
	}
	.description {
		margin-inline: auto;
		text-wrap: balance;
		font-size: var(--fs--1);
	}
	.description p:first-child {
		font-weight: 600;
		color: hsl(var(--clr-grey-700));
		font-size: var(--fs-0);
	}
	.icon {
		margin-inline: auto;
		padding: 0.5rem;
		border-radius: 0.5rem;
	}
	.btn {
		place-content: center;
	}
</style>
