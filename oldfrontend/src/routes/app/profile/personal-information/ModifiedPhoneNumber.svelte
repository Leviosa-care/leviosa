<script lang="ts">
	type PhoneNumberProps = { phoneNumber: string };
	let { phoneNumber }: PhoneNumberProps = $props();
	import TextInput from '$lib/components/forms/atoms/TextInput.svelte';
</script>

<div class="content grid" style="--gap: 1.5rem;">
	<TextInput
		id="phoneNumber"
		label="Code postal d'une adresse permanente"
		name="phoneNumber"
		placeholder="Numero de telephone"
		value={phoneNumber}
	/>
</div>

<style>
	.content {
		margin-top: 1rem;
	}
</style>
