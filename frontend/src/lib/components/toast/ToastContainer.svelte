<script lang="ts">
	import { getContext } from 'svelte';
	import { onMount } from 'svelte';
	import Toast from './Toast.svelte';
	import type { ToastContext } from './types';

	const TOAST_KEY = Symbol('toast');
	const toast = getContext<ToastContext>(TOAST_KEY);
	let mounted = $state(false);

	onMount(() => {
		mounted = true;
	});

	function handleClose(id: number) {
		if (!toast) return;
		const index = toast.toasts.findIndex((t) => t.id === id);
		if (index !== -1) {
			toast.toasts.splice(index, 1);
		}
	}
</script>

{#if mounted && toast?.toasts.length > 0}
	<div class="fixed top-4 right-4 z-[100] flex flex-col gap-2 pointer-events-none">
		{#each toast.toasts as toastItem (toastItem.id)}
			<Toast toast={toastItem} onClose={handleClose} />
		{/each}
	</div>
{/if}
