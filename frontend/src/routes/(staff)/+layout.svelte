<script lang="ts">
	import type { LayoutProps } from "./$types";
	import { setUserContext } from "$lib/context/user.svelte";
	import StaffSidebar from "$lib/navigation/StaffSidebar.svelte";
	import { adminDarkMode } from "$lib/stores/adminDarkMode.svelte";
	import { browser } from "$app/environment";

	let { children, data }: LayoutProps = $props();

	setUserContext(data.user);

	// Apply dark mode class to html element (same pattern as admin)
	$effect(() => {
		if (browser) {
			if (adminDarkMode.get()) {
				document.documentElement.classList.add('dark');
			} else {
				document.documentElement.classList.remove('dark');
			}
		}
	});

	// Clean up: remove dark class when leaving staff routes
	$effect(() => {
		return () => {
			if (browser) {
				document.documentElement.classList.remove('dark');
			}
		};
	});
</script>

<!-- Mobile Layout -->
<div class="lg:hidden flex flex-col min-h-screen bg-muted/25">
	<main class="flex-1">
		<StaffSidebar user={data.user} permissions={data.permissions} unreadCount={data.unreadCount} />
		{@render children()}
	</main>
</div>

<!-- Desktop Layout -->
<div class="hidden lg:flex min-h-screen bg-muted/25">
	<StaffSidebar user={data.user} permissions={data.permissions} unreadCount={data.unreadCount} />
	<main class="flex-1 bg-background overflow-x-hidden">
		{@render children()}
	</main>
</div>
