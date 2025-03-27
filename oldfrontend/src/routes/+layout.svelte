<script lang="ts">
	import { onMount } from 'svelte';
	import '../styles/style.css';
	import '../app.css';

	async function detectSWUpdate() {
		const registation = await navigator.serviceWorker.ready;
		registation.addEventListener('updatefound', () => {
			const newSW = registation.installing;
			newSW?.addEventListener('statechange', () => {
				if (newSW.state === 'installed') {
					// include logic, make popup etc.. use the confirm API
					if (confirm('New update available ! Reload to update ?')) {
						newSW.postMessage({ type: 'SKIP WAITING' });
						window.location.reload();
					}
				}
			});
		});
	}
	onMount(() => detectSWUpdate);

	// NOTE: the part about making the transition to the pages (view transition API)
	import { onNavigate } from '$app/navigation';
	interface Props {
		children?: import('svelte').Snippet;
	}

	let { children }: Props = $props();
	onNavigate(() => {
		if (!document.startViewTransition) return;
		return new Promise((resolve) => {
			resolve();
		});
	});
</script>

{@render children?.()}
