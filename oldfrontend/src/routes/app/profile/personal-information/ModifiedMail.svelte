<script lang="ts">
	type MailProps = { email: string };
	let { email }: MailProps = $props();
	import TextInput from '$lib/components/forms/atoms/TextInput.svelte';
</script>

<div class="content grid" style="--gap: 1.5rem;">
	<TextInput
		id="email"
		label="Votre adresse e-mail"
		name="email"
		placeholder="jeanne.dupont@gmail.com"
		value={email}
	/>
</div>

<style>
	.content {
		margin-top: 1rem;
	}
</style>
