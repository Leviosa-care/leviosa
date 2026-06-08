<script lang="ts">
	import type { LayoutProps } from './$types';
	import ClientSidebar from '$lib/navigation/ClientSidebar.svelte';
	import ProfileBanner from './ProfileBanner.svelte';

	let { children, data }: LayoutProps = $props();
</script>

<!-- Mobile Layout -->
<div class="lg:hidden flex flex-col min-h-screen bg-muted/25">
	<ProfileBanner show={data.profileIncomplete} />
	<ClientSidebar user={data.user} unreadCount={data.unreadCount} />
	<main class="flex-1 px-4 py-6">
		{@render children()}
	</main>
</div>

<!-- Desktop Layout -->
<div class="hidden lg:flex min-h-screen bg-muted/25">
	<ClientSidebar user={data.user} unreadCount={data.unreadCount} />
	<div class="flex-1 flex flex-col">
		<ProfileBanner show={data.profileIncomplete} />
		<main class="flex-1 bg-background overflow-x-hidden p-6 lg:p-10">
			{@render children()}
		</main>
	</div>
</div>
