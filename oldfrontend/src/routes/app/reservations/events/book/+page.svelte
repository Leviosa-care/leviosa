<script lang="ts">
	import PrestataireRadio from './PrestataireRadio.svelte';
	import EventPicker from '$lib/components/forms/EventPicker.svelte';
	import PrestatairePicker from '$lib/components/forms/PrestatairePicker.svelte';
	import Button from '$lib/components/Button.svelte';

	import type { PageData } from './$types';
	interface Props {
		data: PageData;
	}

	let { data }: Props = $props();
	let { prestataires, prestations, monthsData } = $derived(data);

	function sendData() {
		console.log('sending the data to the backend brother');
	}
</script>

<form method="POST" class="content grid" style="--gap: 1.5rem;">
	<div class="container grid" style="--gap:1rem;">
		<p class="subtitle">Choisis un type de prestation</p>
		<PrestataireRadio {prestations} />
	</div>
	<div class="container grid" style="--gap:1rem;">
		<p class="subtitle">Choisis une date et un horaire</p>
		<EventPicker {monthsData} />
	</div>
	<div class="container grid" style="--gap:1rem;">
		<p class="subtitle">Choisis un(e) prestataire</p>
		<PrestatairePicker {prestataires} />
	</div>
	<div class="button-container container grid">
		<Button buttonType="submit" onClick={sendData}>Finalise ta reservation</Button>
	</div>
</form>

<style>
	.content {
		padding-bottom: 7rem;
		position: relative;
	}
	.content > * {
		margin-top: 2rem;
	}
	.button-container {
		margin-top: 2rem;
	}
	.subtitle {
		color: hsl(var(--clr-grey-600));
		font-size: var(--fs-0);
		font-weight: 600;
	}
</style>
