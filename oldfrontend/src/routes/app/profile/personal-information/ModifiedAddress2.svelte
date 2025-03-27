<script lang="ts">
	type AddressProps = { address: string };
	let { address }: AddressProps = $props();
	import TextInput from '$lib/components/forms/atoms/TextInput.svelte';
</script>

<div class="content grid" style="--gap: 1.5rem;">
	<TextInput
		id="address2"
		label="Detail d'une adresse permanente"
		name="address2"
		placeholder="Detail d'une adresse permanente"
		value={address}
	/>
</div>

<style>
	.content {
		margin-top: 1rem;
	}
</style>
