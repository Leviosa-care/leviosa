<script lang="ts">
	import { ArrowLeft, ChevronLeft } from 'lucide-svelte';

	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	interface Props {
		pathname?: string;
		variant?: 'chevron' | 'arrow';
		color?: 'black' | 'white';
	}

	let { pathname = '', variant = 'chevron', color = 'black' }: Props = $props();
	function goBack() {
		if (history.length > 0) {
			history.back();
		} else if (pathname != '') {
			goto(pathname);
		} else {
			const previousPath = $page.url.pathname.split('/').slice(0, -1).join('/');
			goto(previousPath);
		}
	}
	// TODO: add other variant for that button other than the one with the black stroke.
</script>

<button onclick={goBack} class="back">
	{#if variant === 'arrow'}
		<ArrowLeft {color} />
	{:else}
		<ChevronLeft {color} />
	{/if}
</button>

<style>
	.back {
		background: transparent;
	}
</style>
