<script lang="ts">
	import type { LayoutProps } from "./$types";
	import AdminSidebar from "$lib/navigation/AdminSidebar.svelte";

	let { children, data }: LayoutProps = $props();
	let { user, permissions } = data;
</script>

<div class="flex min-h-screen">
	<!-- Sidebar component handles both desktop sidebar and mobile bottom nav -->
	<AdminSidebar {user} {permissions} />

	<!-- Main content area -->
	<main class="flex-1 overflow-x-hidden">
		<!-- Add bottom padding on mobile to prevent content from being hidden behind bottom nav -->
		<div class="pb-20 lg:pb-0">
			{@render children()}
		</div>
	</main>
</div>
