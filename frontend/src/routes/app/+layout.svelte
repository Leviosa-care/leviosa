<script lang="ts">
	import NavigationBar from '$lib/components/navigation/NavigationBar.svelte';

	interface Props {
		children: import('svelte').Snippet;
		data: import('./$types').PageData;
	}

	let { children, data }: Props = $props();
	const { role } = data;

	// page transition
	import { onNavigate } from '$app/navigation';
	onNavigate((navigation) => {
		const isSameUrl = navigation?.from?.url.href === navigation?.to?.url.href;
		const toPathname = navigation?.to?.url.pathname;
		const fromPathname = navigation?.from?.url.pathname;
		// TODO: handle the nav from app/ to app/... paths
		const isSubRoute = toPathname?.includes(String(fromPathname)) && fromPathname != '/app';
		if (!document.startViewTransition || isSameUrl || !isSubRoute) return;
		return new Promise((resolve) => {
			document.startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});
</script>

<div class="layout">
	<NavigationBar {role} />
	<div class="content">
		{@render children?.()}
	</div>
</div>

<style>
	.layout {
		view-transition-name: pushing;
		position: relative;
		display: flex;
	}
	.content {
		min-height: 100vh;
		flex: 1;
		overflow-y: auto;
	}
	/* do the media queries here */
	@keyframes push-new {
		from {
			transform: translateX(100%);
		}
		to {
			transform: translateX(0%);
		}
	}
	@keyframes push-old {
		to {
			transform: translateX(-30%);
			opacity: 0.3;
		}
	}
	/* TODO: here there is a root, need to see how that affects other pages ? */
	:root::view-transition-old(root) {
		animation: 250ms ease-out both push-old;
	}
	:root::view-transition-new(pushed) {
		animation: 250ms ease-out both push-new;
	}
</style>
