<script lang="ts">
	import type { LayoutProps } from "./$types";
	import { setUserContext } from "$lib/context/user.svelte";
	import AdminSidebar from "$lib/navigation/AdminSidebar.svelte";
	import { adminDarkMode } from "$lib/stores/adminDarkMode.svelte";
	import { browser } from "$app/environment";

	let { data, children }: LayoutProps = $props();

	setUserContext(data.user);

	// Apply dark mode class to html element only when in admin routes
	$effect(() => {
		if (browser) {
			if (adminDarkMode.get()) {
				document.documentElement.classList.add('dark');
			} else {
				document.documentElement.classList.remove('dark');
			}
		}
	});

	// Clean up: remove dark class when leaving admin routes
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
		{@render children()}
	</main>
</div>

<!-- Desktop Layout -->
<div class="hidden lg:flex min-h-screen bg-muted/25">
	<AdminSidebar user={data.user} permissions={data.permissions} unreadCount={data.unreadCount} />
	<main class="flex-1 bg-background">
		{@render children()}
	</main>
</div>
