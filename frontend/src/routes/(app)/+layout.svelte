<script lang="ts">
	import type { LayoutProps } from "./$types";
	import Navigation from "./navigation.svelte";

	let { children, data }: LayoutProps = $props();
</script>

<Navigation user={data.user} />

<!-- Main content with bottom padding on mobile to avoid overlap with bottom nav -->
<main class="min-h-screen bg-white pb-20 xl:pb-0">
	{@render children()}
</main>
