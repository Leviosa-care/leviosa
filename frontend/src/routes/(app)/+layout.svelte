<script lang="ts">
	import type { LayoutProps } from "./$types";
	import Navigation from "./navigation.svelte";

	let { children, data }: LayoutProps = $props();
</script>

<Navigation user={data.user} />

<main class="min-h-screen bg-white">
	{@render children()}
</main>
