<script lang="ts">
	type FullnameProps = { firstname: string; lastname: string };
	let { firstname, lastname }: FullnameProps = $props();
	import TextInput from '$lib/components/forms/atoms/TextInput.svelte';
</script>

<div class="content grid" style="--gap: 1.5rem;">
	<TextInput
		id="firstname"
		label="Prenom sur la piece d'identite"
		name="firstname"
		placeholder="John"
		value={firstname}
	/>
	<TextInput
		id="lastname"
		label="Nom sur la piece d'identite"
		name="lastname"
		placeholder="DOE"
		value={lastname}
	/>
</div>

<style>
	.content {
		margin-top: 1rem;
	}
</style>
