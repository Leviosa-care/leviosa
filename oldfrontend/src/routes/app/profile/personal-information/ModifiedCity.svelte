<script lang="ts">
	type CityProps = { city: string };
	let { city }: CityProps = $props();
	import TextInput from '$lib/components/forms/atoms/TextInput.svelte';
</script>

<div class="content grid" style="--gap: 1.5rem;">
	<TextInput
		id="city"
		label="Ville"
		name="city"
		placeholder="Nom de la ville d'une adresse permanente"
		value={city}
	/>
</div>

<style>
	.content {
		margin-top: 1rem;
	}
</style>
