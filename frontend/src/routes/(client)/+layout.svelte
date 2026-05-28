<script lang="ts">
	import type { LayoutProps } from './$types';
	import ClientNav from './ClientNav.svelte';
	import ProfileBanner from './ProfileBanner.svelte';

	let { children, data }: LayoutProps = $props();
</script>

<div class="min-h-screen bg-background">
	<ClientNav user={data.user} />
	<ProfileBanner show={data.profileIncomplete} />
	<main class="max-w-5xl mx-auto px-4 py-6 lg:py-10">
		{@render children()}
	</main>
</div>
