<script lang="ts">
	import { enhance } from '$app/forms';
	import { onMount } from 'svelte';

	let formElement: HTMLFormElement;

	onMount(() => {
		formElement.requestSubmit();
	});
</script>

<form method="POST" bind:this={formElement} use:enhance>
	<p class="text-center text-gray-600">Déconnexion en cours…</p>
</form>
