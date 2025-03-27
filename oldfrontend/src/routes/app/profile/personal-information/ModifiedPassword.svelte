<script lang="ts">
	import PasswordInput from '$lib/components/forms/atoms/PasswordInput.svelte';
</script>

<div class="content grid" style="--gap: 1.5rem;">
	<PasswordInput id="new-password" label="Nouveau mot de passe" />
	<PasswordInput id="confirm-new-password" label="Confirmation nouveau mot de passe" />
</div>

<style>
	.content {
		margin-top: 1rem;
	}
</style>
