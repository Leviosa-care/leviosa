<script lang="ts">
	import { adminDarkMode } from "$lib/stores/adminDarkMode.svelte";
	import { Moon, Sun } from "@lucide/svelte";

	const isDark = $derived(adminDarkMode.get());
</script>

<button
	onclick={() => adminDarkMode.toggle()}
	class="p-3 rounded-lg transition-colors text-dark-500 hover:text-dark-900 dark:text-dark-400 dark:hover:text-dark-100 hover:bg-dark-100 dark:hover:bg-dark-800"
	aria-label={isDark ? "Passer en mode clair" : "Passer en mode sombre"}
	type="button"
>
	{#if isDark}
		<Sun size={18} />
	{:else}
		<Moon size={18} />
	{/if}
</button>
