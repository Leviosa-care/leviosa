<script lang="ts">
	import { adminDarkMode } from "$lib/stores/adminDarkMode.svelte";
	import { Moon, Sun } from "@lucide/svelte";

	const isDark = $derived(adminDarkMode.get());
</script>

<button
	onclick={() => adminDarkMode.toggle()}
	class="p-3 rounded-lg transition-colors text-muted-foreground hover:text-foreground hover:bg-surface-hover"
	aria-label={isDark ? "Passer en mode clair" : "Passer en mode sombre"}
	type="button"
>
	{#if isDark}
		<Sun size={18} />
	{:else}
		<Moon size={18} />
	{/if}
</button>
