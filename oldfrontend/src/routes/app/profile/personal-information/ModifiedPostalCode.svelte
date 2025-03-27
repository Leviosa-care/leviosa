<script lang="ts">
	type PostalCodeProps = { postalCode: number };
	let { postalCode }: PostalCodeProps = $props();
	import TextInput from '$lib/components/forms/atoms/TextInput.svelte';
</script>

<div class="content grid" style="--gap: 1.5rem;">
	<TextInput
		id="postalCode"
		label="Code postal"
		name="postalCode"
		placeholder="Code postal d'une adresse permanente"
		value={postalCode}
	/>
</div>

<style>
	.content {
		margin-top: 1rem;
	}
</style>
