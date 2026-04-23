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

<div class="flex flex-col lg:flex-row">
	<AdminSidebar user={data.user} permissions={data.permissions} />
	{@render children()}
</div>
